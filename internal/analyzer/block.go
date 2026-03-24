// Package analyzer provides block analysis functions.
package analyzer

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"os"
	"unicode/utf8"

	"chainlens/internal/models"
	"chainlens/internal/parser"
)

// utf8Valid is a helper that wraps utf8.Valid for use in block processing.
func utf8Valid(b []byte) bool {
	return utf8.Valid(b)
}

// AnalyzeBlock parses and analyzes a block from blk*.dat, rev*.dat, and xor.dat files.
func AnalyzeBlock(blkPath, revPath, xorPath string) ([]*models.BlockResult, error) {
	// Read XOR key
	xorKey, err := readXORKey(xorPath)
	if err != nil {
		return nil, err
	}

	// Read and decode blk*.dat
	blkData, err := readAndXOR(blkPath, xorKey)
	if err != nil {
		return nil, err
	}

	// Read and decode rev*.dat
	revData, err := readAndXOR(revPath, xorKey)
	if err != nil {
		return nil, err
	}

	return analyzeAllBlocks(blkData, revData)
}

// AnalyzeBlockFromData parses and analyzes a block from already-decrypted blk and rev data.
// Used by the web API which receives raw (already XOR-decrypted) bytes.
func AnalyzeBlockFromData(blkData, revData []byte) ([]*models.BlockResult, error) {
	return analyzeAllBlocks(blkData, revData)
}

// analyzeAllBlocks is the shared implementation for both file-based and data-based analysis.
// It parses all blocks from the blk data and matches each to its corresponding undo entry by index.
func analyzeAllBlocks(blkData, revData []byte) ([]*models.BlockResult, error) {
	// Parse block file
	blockEntries, err := parser.ParseBlockFile(blkData)
	if err != nil {
		return nil, err
	}

	// Parse undo file - get raw undo entries
	undoEntries, err := parseUndoFileRaw(revData)
	if err != nil {
		return nil, err
	}

	if len(blockEntries) == 0 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidBlockHeader,
			"no blocks found in blk data")
	}
	if len(undoEntries) == 0 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeUndoMismatch,
			"no undo entries found in rev data")
	}

	// Only analyze the first block, but we need all undo entries for matching.
	entry := blockEntries[0]
	txs, err := parser.ParseBlockTransactions(entry.RawData)
	if err != nil {
		return nil, err
	}
	expected := len(txs) - 1 // non-coinbase tx count

	// Match the first block to its undo entry by non-coinbase tx count.
	// Block and undo entries may not be in the same order in the files.
	var matchingUndo []byte
	for _, undoData := range undoEntries {
		count, _, err := parser.ReadVarInt(undoData, 0)
		if err != nil {
			continue
		}
		if int(count) == expected {
			matchingUndo = undoData
			break
		}
	}

	if matchingUndo == nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeUndoMismatch,
			"no matching undo entry found for block with %d non-coinbase txs", expected)
	}

	result, err := analyzeBlockEntry(entry, matchingUndo)
	if err != nil {
		return nil, err
	}

	return []*models.BlockResult{result}, nil
}

// readXORKey reads the 8-byte XOR key from xor.dat.
func readXORKey(path string) ([]byte, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to read XOR key: %v", err)
	}
	if len(key) != 8 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidXORKey,
			"XOR key must be 8 bytes, got %d", len(key))
	}
	return key, nil
}

// readAndXOR reads a file and XORs it with the key.
func readAndXOR(path string, key []byte) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to read file %s: %v", path, err)
	}

	// XOR each byte with the repeating 8-byte key
	for i := range data {
		data[i] ^= key[i%8]
	}

	return data, nil
}

// parseUndoFileRaw parses undo file and returns raw undo data for each block.
// The undo file structure is: magic(4) + size(4) + data(size) + checksum(32) for each entry.
func parseUndoFileRaw(data []byte) ([][]byte, error) {
	var undos [][]byte
	offset := 0

	for offset < len(data) {
		// Check minimum size for magic + size
		if offset+8 > len(data) {
			break
		}

		// Read magic (4 bytes)
		magic := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4

		// Validate magic
		if magic != parser.MainnetMagic && magic != parser.TestnetMagic && magic != parser.RegtestMagic {
			// Try to find next valid magic
			found := false
			for offset < len(data)-4 {
				testMagic := binary.LittleEndian.Uint32(data[offset : offset+4])
				if testMagic == parser.MainnetMagic || testMagic == parser.TestnetMagic || testMagic == parser.RegtestMagic {
					found = true
					break
				}
				offset++
			}
			if !found {
				break
			}
			continue
		}

		// Read undo size (4 bytes)
		if offset+4 > len(data) {
			break
		}
		undoSize := binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4

		// Validate size (including 32-byte checksum after data)
		if uint32(offset)+undoSize+32 > uint32(len(data)) {
			return nil, models.NewAnalysisErrorf(models.ErrCodeTruncatedUndo,
				"undo size %d exceeds remaining data at offset %d", undoSize, offset)
		}

		// Store raw undo data
		undoData := make([]byte, undoSize)
		copy(undoData, data[offset:offset+int(undoSize)])
		undos = append(undos, undoData)
		offset += int(undoSize)

		// Skip 32-byte checksum after undo data
		offset += 32
	}

	return undos, nil
}

// analyzeBlockEntry analyzes a single block with its undo data.
func analyzeBlockEntry(entry *parser.BlockFileEntry, undoData []byte) (*models.BlockResult, error) {
	// Parse transactions from block
	txs, err := parser.ParseBlockTransactions(entry.RawData)
	if err != nil {
		return nil, err
	}

	// Get input counts for non-coinbase transactions
	inputCounts := parser.GetInputCounts(txs)

	// Parse undo data
	var blockUndo *parser.BlockUndo
	if len(inputCounts) > 0 {
		blockUndo, err = parser.ParseBlockUndoData(undoData, inputCounts)
		if err != nil {
			return nil, err
		}
	}

	// Build prevout map from undo data
	prevoutMap := make(map[string]models.FixturePrevout)
	if blockUndo != nil {
		for txIdx, txUndo := range blockUndo.TxUndos {
			tx := txs[txIdx+1] // +1 because txIdx 0 in undo corresponds to tx 1 (skip coinbase)
			for inIdx, undoEntry := range txUndo {
				input := tx.Inputs[inIdx]
				txidHex := hex.EncodeToString(reversedBytes(input.PrevTxid[:]))
				key := outpointKey(txidHex, input.PrevVout)
				prevoutMap[key] = models.FixturePrevout{
					Txid:            txidHex,
					Vout:            input.PrevVout,
					ValueSats:       undoEntry.Value,
					ScriptPubkeyHex: undoEntry.ScriptPubKeyHex,
				}
			}
		}
	}

	// Determine network from magic
	network := "mainnet"
	switch entry.Magic {
	case parser.TestnetMagic:
		network = "testnet"
	case parser.RegtestMagic:
		network = "regtest"
	}

	// Verify merkle root
	txids := make([]string, len(txs))
	for i, tx := range txs {
		txids[i] = tx.Txid()
	}
	computedRoot, err := parser.ComputeMerkleRootFromTxids(txids)
	if err != nil {
		return nil, err
	}
	merkleValid := computedRoot == entry.Header.MerkleRoot

	// Error on invalid merkle root
	if !merkleValid {
		return nil, models.NewAnalysisErrorf(models.ErrCodeMerkleRootMismatch,
			"computed merkle root does not match block header merkle root")
	}

	// Extract BIP34 height from coinbase
	bip34Height := extractBIP34Height(txs[0])

	// Analyze each transaction
	txResults := make([]models.TransactionResult, len(txs))
	var totalFees uint64
	var totalWeight int

	for i, tx := range txs {
		if i == 0 {
			// Coinbase transaction
			txResults[i] = analyzeCoinbase(tx, network)
		} else {
			// Regular transaction
			result := analyzeBlockTx(tx, prevoutMap, network, i)
			txResults[i] = result
			totalFees += result.FeeSats
		}
		totalWeight += tx.Weight()
	}

	// Calculate average fee rate
	avgFeeRate := 0.0
	if totalWeight > 0 {
		totalVbytes := (totalWeight + 3) / 4
		avgFeeRate = float64(totalFees) / float64(totalVbytes)
		avgFeeRate = math.Round(avgFeeRate*100) / 100
	}

	// Build script type summary
	scriptSummary := buildScriptTypeSummary(txResults)

	// Calculate total coinbase output
	coinbaseOutputTotal := uint64(0)
	for _, out := range txs[0].Outputs {
		coinbaseOutputTotal += out.Value
	}

	result := &models.BlockResult{
		OK:   true,
		Mode: models.ModeBlock,
		BlockHeader: models.BlockHeader{
			Version:         entry.Header.Version,
			PrevBlockHash:   entry.Header.PrevBlockHashHex(),
			MerkleRoot:      entry.Header.MerkleRootHex(),
			MerkleRootValid: merkleValid,
			Timestamp:       entry.Header.Timestamp,
			Bits:            entry.Header.BitsHex(),
			Nonce:           entry.Header.Nonce,
			BlockHash:       entry.Header.BlockHash(),
		},
		TxCount: len(txs),
		Coinbase: models.CoinbaseInfo{
			BIP34Height:       bip34Height,
			CoinbaseScriptHex: hex.EncodeToString(txs[0].Inputs[0].ScriptSig),
			TotalOutputSats:   coinbaseOutputTotal,
		},
		Transactions: txResults,
		BlockStats: models.BlockStats{
			TotalFeesSats:     totalFees,
			TotalWeight:       totalWeight,
			AvgFeeRateSatVB:   avgFeeRate,
			ScriptTypeSummary: scriptSummary,
		},
	}

	return result, nil
}

// extractBIP34Height extracts the block height from coinbase scriptSig (BIP34).
func extractBIP34Height(coinbase *parser.RawTransaction) int64 {
	if len(coinbase.Inputs) == 0 {
		return -1
	}

	scriptSig := coinbase.Inputs[0].ScriptSig
	if len(scriptSig) < 1 {
		return -1
	}

	// First byte is the push length for the height
	pushLen := int(scriptSig[0])

	// BIP34 height is pushed using minimally-encoded push
	// 0x01-0x4b means push that many bytes directly
	if pushLen >= 1 && pushLen <= 8 && len(scriptSig) >= 1+pushLen {
		heightBytes := scriptSig[1 : 1+pushLen]
		var height int64
		for i := 0; i < len(heightBytes); i++ {
			height |= int64(heightBytes[i]) << (8 * i)
		}
		return height
	}

	return -1
}

// analyzeCoinbase creates a TransactionResult for the coinbase transaction.
func analyzeCoinbase(tx *parser.RawTransaction, network string) models.TransactionResult {
	// Coinbase has no inputs to spend, so no fee calculation
	vbytes := tx.Vbytes()

	result := models.TransactionResult{
		OK:              true,
		Network:         network,
		Segwit:          tx.IsSegwit,
		Txid:            tx.Txid(),
		Version:         tx.Version,
		Locktime:        tx.Locktime,
		SizeBytes:       tx.TotalSize,
		Weight:          tx.Weight(),
		Vbytes:          vbytes,
		TotalInputSats:  0,
		TotalOutputSats: 0,
		FeeSats:         0,
		FeeRateSatVB:    0,
		RBFSignaling:    false,
		LocktimeType:    models.LocktimeNone,
		LocktimeValue:   0,
	}

	if tx.IsSegwit {
		result.Wtxid = tx.Wtxid()
	}

	// Process coinbase input
	result.Vin = make([]models.Vin, 1)
	input := tx.Inputs[0]
	scriptAsm, _ := parser.DisassembleScript(input.ScriptSig)

	witnessHex := make([]string, len(tx.Witnesses[0]))
	for i, w := range tx.Witnesses[0] {
		witnessHex[i] = hex.EncodeToString(w)
	}

	result.Vin[0] = models.Vin{
		Txid:         models.CoinbaseTxid,
		Vout:         models.CoinbaseVout,
		Sequence:     input.Sequence,
		ScriptSigHex: hex.EncodeToString(input.ScriptSig),
		ScriptAsm:    scriptAsm,
		Witness:      witnessHex,
		ScriptType:   "coinbase",
		Address:      nil,
		Prevout: models.VinPrevout{
			ValueSats:       0,
			ScriptPubkeyHex: "",
		},
		RelativeTimelock: models.RelativeTimelock{Enabled: false},
	}

	// Process outputs
	result.Vout = make([]models.Vout, len(tx.Outputs))
	for i, output := range tx.Outputs {
		result.TotalOutputSats += output.Value

		scriptType := parser.ClassifyScript(output.ScriptPubkey)
		addr := getAddressFromScript(output.ScriptPubkey)
		scriptAsm, _ := parser.DisassembleScript(output.ScriptPubkey)

		result.Vout[i] = models.Vout{
			N:               uint32(i),
			ValueSats:       output.Value,
			ScriptPubkeyHex: hex.EncodeToString(output.ScriptPubkey),
			ScriptAsm:       scriptAsm,
			ScriptType:      scriptTypeToOutputString(scriptType),
			Address:         addr,
		}
	}

	return result
}

// analyzeBlockTx analyzes a non-coinbase transaction from a block.
func analyzeBlockTx(tx *parser.RawTransaction, prevoutMap map[string]models.FixturePrevout, network string, txIndex int) models.TransactionResult {
	vbytes := tx.Vbytes()

	// Calculate totals
	totalInputSats := uint64(0)
	totalOutputSats := uint64(0)

	for _, output := range tx.Outputs {
		totalOutputSats += output.Value
	}

	result := models.TransactionResult{
		OK:              true,
		Network:         network,
		Segwit:          tx.IsSegwit,
		Txid:            tx.Txid(),
		Version:         tx.Version,
		Locktime:        tx.Locktime,
		SizeBytes:       tx.TotalSize,
		Weight:          tx.Weight(),
		Vbytes:          vbytes,
		TotalOutputSats: totalOutputSats,
		RBFSignaling:    detectRBFSignaling(tx),
	}

	if tx.IsSegwit {
		result.Wtxid = tx.Wtxid()
	}

	// Classify locktime
	result.LocktimeType, result.LocktimeValue = classifyLocktime(tx.Locktime)

	// Process inputs
	result.Vin = make([]models.Vin, len(tx.Inputs))
	for i, input := range tx.Inputs {
		txidHex := hex.EncodeToString(reversedBytes(input.PrevTxid[:]))
		key := outpointKey(txidHex, input.PrevVout)
		prevout, hasPrevout := prevoutMap[key]

		var witness []parser.RawWitnessItem
		if i < len(tx.Witnesses) {
			witness = tx.Witnesses[i]
		}

		if hasPrevout {
			totalInputSats += prevout.ValueSats
		}

		vin := processBlockInput(input, witness, prevout, hasPrevout)
		result.Vin[i] = vin
	}

	result.TotalInputSats = totalInputSats
	if totalInputSats >= totalOutputSats {
		result.FeeSats = totalInputSats - totalOutputSats
		if vbytes > 0 {
			result.FeeRateSatVB = float64(result.FeeSats) / float64(vbytes)
			result.FeeRateSatVB = math.Round(result.FeeRateSatVB*100) / 100
		}
	}

	// Calculate SegWit savings
	if tx.IsSegwit {
		result.SegwitSavings = calculateSegwitSavings(tx)
	}

	// Process outputs
	result.Vout = make([]models.Vout, len(tx.Outputs))
	for i, output := range tx.Outputs {
		vout := processBlockOutput(uint32(i), output)
		result.Vout[i] = vout
	}

	// Generate warnings
	result.Warnings = generateWarnings(&result)

	return result
}

// processBlockInput processes a single input from a block transaction.
func processBlockInput(input parser.RawInput, witness []parser.RawWitnessItem, prevout models.FixturePrevout, hasPrevout bool) models.Vin {
	var prevoutScript []byte
	if hasPrevout && prevout.ScriptPubkeyHex != "" {
		prevoutScript, _ = hex.DecodeString(prevout.ScriptPubkeyHex)
	}

	// Convert witness
	witnessData := make([][]byte, len(witness))
	witnessHex := make([]string, len(witness))
	for i, w := range witness {
		witnessData[i] = []byte(w)
		witnessHex[i] = hex.EncodeToString(w)
	}

	// Classify input
	var inputType parser.InputScriptType = parser.InputTypeUnknown
	if hasPrevout && len(prevoutScript) > 0 {
		inputType = parser.ClassifyInput(prevoutScript, input.ScriptSig, witnessData)
	}

	// Get address
	var addr *string
	if hasPrevout && len(prevoutScript) > 0 {
		addr = getAddressFromScript(prevoutScript)
	}

	// Disassemble scriptSig
	scriptAsm, _ := parser.DisassembleScript(input.ScriptSig)

	// Parse relative timelock
	relTimelock := parseRelativeTimelock(input.Sequence)

	vin := models.Vin{
		Txid:         hex.EncodeToString(reversedBytes(input.PrevTxid[:])),
		Vout:         input.PrevVout,
		Sequence:     input.Sequence,
		ScriptSigHex: hex.EncodeToString(input.ScriptSig),
		ScriptAsm:    scriptAsm,
		Witness:      witnessHex,
		ScriptType:   string(inputType),
		Address:      addr,
		Prevout: models.VinPrevout{
			ValueSats:       prevout.ValueSats,
			ScriptPubkeyHex: prevout.ScriptPubkeyHex,
		},
		RelativeTimelock: relTimelock,
	}

	// Add witnessScript disassembly for P2WSH/P2SH-P2WSH
	if inputType == parser.InputTypeP2WSH || inputType == parser.InputTypeP2SH_P2WSH {
		witnessScript := parser.GetWitnessScript(inputType, witnessData)
		if witnessScript != nil {
			wsAsm, _ := parser.DisassembleScript(witnessScript)
			vin.WitnessScriptAsm = &wsAsm
		}
	}

	return vin
}

// processBlockOutput processes a single output from a block transaction.
func processBlockOutput(n uint32, output parser.RawOutput) models.Vout {
	scriptType := parser.ClassifyScript(output.ScriptPubkey)
	addr := getAddressFromScript(output.ScriptPubkey)
	scriptAsm, _ := parser.DisassembleScript(output.ScriptPubkey)
	scriptTypeStr := scriptTypeToOutputString(scriptType)

	vout := models.Vout{
		N:               n,
		ValueSats:       output.Value,
		ScriptPubkeyHex: hex.EncodeToString(output.ScriptPubkey),
		ScriptAsm:       scriptAsm,
		ScriptType:      scriptTypeStr,
		Address:         addr,
	}

	// Handle OP_RETURN special fields
	if scriptType == parser.ScriptTypeOPReturn {
		opData, err := parser.ExtractOPReturnData(output.ScriptPubkey)
		if err == nil {
			vout.OPReturnDataHex = &opData.Payload

			// Check if valid UTF-8
			payloadBytes, _ := hex.DecodeString(opData.Payload)
			if len(payloadBytes) > 0 && utf8Valid(payloadBytes) {
				text := string(payloadBytes)
				vout.OPReturnDataUTF8 = &text
			}

			// Set protocol
			protocol := string(opData.Protocol)
			if protocol == "text" {
				protocol = "unknown" // "text" is not a spec protocol
			}
			vout.OPReturnProtocol = &protocol
		}
	}

	return vout
}

// buildScriptTypeSummary counts output script types across all transactions.
func buildScriptTypeSummary(txResults []models.TransactionResult) models.ScriptTypeSummary {
	summary := models.ScriptTypeSummary{}

	for _, tx := range txResults {
		for _, vout := range tx.Vout {
			switch vout.ScriptType {
			case "p2wpkh":
				summary.P2WPKH++
			case "p2tr":
				summary.P2TR++
			case "p2sh":
				summary.P2SH++
			case "p2pkh":
				summary.P2PKH++
			case "p2wsh":
				summary.P2WSH++
			case "op_return":
				summary.OPReturn++
			default:
				summary.Unknown++
			}
		}
	}

	return summary
}
