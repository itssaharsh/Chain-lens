// Package analyzer provides transaction and block analysis functions.
package analyzer

import (
	"encoding/hex"
	"fmt"
	"math"
	"unicode/utf8"

	"chainlens/internal/address"
	"chainlens/internal/models"
	"chainlens/internal/parser"
)

// AnalyzeTransaction parses and analyzes a transaction from a fixture.
func AnalyzeTransaction(fixture *models.Fixture) (*models.TransactionResult, error) {
	// Parse the transaction (ParseTransaction takes hex string)
	tx, err := parser.ParseTransaction(fixture.RawTx)
	if err != nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidTx, "failed to parse transaction: %v", err)
	}

	// Build prevout map for lookup
	prevoutMap, err := buildPrevoutMap(fixture.Prevouts, tx)
	if err != nil {
		return nil, err
	}

	// Calculate totals
	totalInputSats := uint64(0)
	for _, prevout := range prevoutMap {
		totalInputSats += prevout.ValueSats
	}

	totalOutputSats := uint64(0)
	for _, output := range tx.Outputs {
		totalOutputSats += output.Value
	}

	if totalInputSats < totalOutputSats {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidTx,
			"total inputs (%d) less than total outputs (%d)", totalInputSats, totalOutputSats)
	}
	feeSats := totalInputSats - totalOutputSats

	// Calculate fee rate
	vbytes := tx.Vbytes()
	feeRate := float64(feeSats) / float64(vbytes)
	feeRate = math.Round(feeRate*100) / 100 // Round to 2 decimals

	// Detect RBF signaling
	rbfSignaling := detectRBFSignaling(tx)

	// Classify locktime
	locktimeType, locktimeValue := classifyLocktime(tx.Locktime)

	// Build result
	result := &models.TransactionResult{
		OK:              true,
		Network:         fixture.Network,
		Segwit:          tx.IsSegwit,
		Txid:            tx.Txid(),
		Version:         tx.Version,
		Locktime:        tx.Locktime,
		SizeBytes:       tx.TotalSize,
		Weight:          tx.Weight(),
		Vbytes:          vbytes,
		TotalInputSats:  totalInputSats,
		TotalOutputSats: totalOutputSats,
		FeeSats:         feeSats,
		FeeRateSatVB:    feeRate,
		RBFSignaling:    rbfSignaling,
		LocktimeType:    locktimeType,
		LocktimeValue:   locktimeValue,
	}

	// Set wtxid for SegWit transactions
	if tx.IsSegwit {
		result.Wtxid = tx.Wtxid()
	}

	// Calculate SegWit savings
	if tx.IsSegwit {
		result.SegwitSavings = calculateSegwitSavings(tx)
	}

	// Process inputs
	result.Vin = make([]models.Vin, len(tx.Inputs))
	for i, input := range tx.Inputs {
		txidHex := hex.EncodeToString(reversedBytes(input.PrevTxid[:]))
		prevout := prevoutMap[outpointKey(txidHex, input.PrevVout)]
		var witness []parser.RawWitnessItem
		if i < len(tx.Witnesses) {
			witness = tx.Witnesses[i]
		}
		vin, err := processInput(input, witness, prevout)
		if err != nil {
			return nil, err
		}
		result.Vin[i] = *vin
	}

	// Process outputs
	result.Vout = make([]models.Vout, len(tx.Outputs))
	for i, output := range tx.Outputs {
		vout, err := processOutput(uint32(i), output)
		if err != nil {
			return nil, err
		}
		result.Vout[i] = *vout
	}

	// Generate warnings
	result.Warnings = generateWarnings(result)

	return result, nil
}

// buildPrevoutMap builds a map of outpoints to prevout data and validates.
func buildPrevoutMap(prevouts []models.FixturePrevout, tx *parser.RawTransaction) (map[string]models.FixturePrevout, error) {
	// Build required outpoints set
	required := make(map[string]bool)
	for _, input := range tx.Inputs {
		key := outpointKey(hex.EncodeToString(reversedBytes(input.PrevTxid[:])), input.PrevVout)
		if required[key] {
			return nil, models.NewAnalysisErrorf(models.ErrCodeDuplicatePrevout,
				"duplicate input outpoint: %s", key)
		}
		required[key] = true
	}

	// Build prevout map
	provided := make(map[string]models.FixturePrevout)
	for _, p := range prevouts {
		key := outpointKey(p.Txid, p.Vout)
		if _, exists := provided[key]; exists {
		return nil, models.NewAnalysisErrorf(models.ErrCodeDuplicatePrevout,
				"duplicate prevout in fixture: %s", key)
		}
		provided[key] = p
	}

	// Verify all required prevouts are provided
	for key := range required {
		if _, exists := provided[key]; !exists {
			return nil, models.NewAnalysisErrorf(models.ErrCodeMissingPrevout,
				"missing prevout: %s", key)
		}
	}

	// Verify no extra prevouts
	for key := range provided {
		if !required[key] {
			return nil, models.NewAnalysisErrorf(models.ErrCodeExtraPrevout,
				"extra prevout not referenced by any input: %s", key)
		}
	}

	return provided, nil
}

// outpointKey creates a unique key for an outpoint.
func outpointKey(txid string, vout uint32) string {
	return fmt.Sprintf("%s:%d", txid, vout)
}

// detectRBFSignaling checks if any input signals RBF (BIP125).
func detectRBFSignaling(tx *parser.RawTransaction) bool {
	for _, input := range tx.Inputs {
		// RBF signaled if sequence < 0xFFFFFFFE
		if input.Sequence < 0xFFFFFFFE {
			return true
		}
	}
	return false
}

// classifyLocktime determines the locktime type.
func classifyLocktime(locktime uint32) (string, uint32) {
	if locktime == 0 {
		return models.LocktimeNone, 0
	}
	if locktime < models.LocktimeThreshold {
		return models.LocktimeBlock, locktime
	}
	return models.LocktimeTimestamp, locktime
}

// calculateSegwitSavings calculates the SegWit weight savings.
func calculateSegwitSavings(tx *parser.RawTransaction) *models.SegwitSavings {
	witnessBytes := tx.WitnessSize
	nonWitnessBytes := tx.NonWitnessSize
	totalBytes := tx.TotalSize

	weightActual := tx.Weight()
	// If legacy: all bytes would count as non-witness (4x weight)
	weightIfLegacy := totalBytes * 4

	savingsPct := float64(weightIfLegacy-weightActual) / float64(weightIfLegacy) * 100
	savingsPct = math.Round(savingsPct*100) / 100 // Round to 2 decimals

	return &models.SegwitSavings{
		WitnessBytes:    witnessBytes,
		NonWitnessBytes: nonWitnessBytes,
		TotalBytes:      totalBytes,
		WeightActual:    weightActual,
		WeightIfLegacy:  weightIfLegacy,
		SavingsPct:      savingsPct,
	}
}

// processInput processes a single input.
func processInput(input parser.RawInput, witness []parser.RawWitnessItem, prevout models.FixturePrevout) (*models.Vin, error) {
	// Decode prevout scriptPubKey
	prevoutScript, err := hex.DecodeString(prevout.ScriptPubkeyHex)
	if err != nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidTx, "invalid prevout script hex")
	}

	// Convert witness to [][]byte for classification
	witnessData := make([][]byte, len(witness))
	witnessHexStrs := make([]string, len(witness))
	for i, w := range witness {
		witnessData[i] = []byte(w)
		witnessHexStrs[i] = hex.EncodeToString(w)
	}

	// Classify input
	inputType := parser.ClassifyInput(prevoutScript, input.ScriptSig, witnessData)

	// Get address from prevout script
	addr := getAddressFromScript(prevoutScript)

	// Disassemble scriptSig
	scriptAsm, _ := parser.DisassembleScript(input.ScriptSig)

	// Detect relative timelock (BIP68)
	relTimelock := parseRelativeTimelock(input.Sequence)

	vin := &models.Vin{
		Txid:         hex.EncodeToString(reversedBytes(input.PrevTxid[:])),
		Vout:         input.PrevVout,
		Sequence:     input.Sequence,
		ScriptSigHex: hex.EncodeToString(input.ScriptSig),
		ScriptAsm:    scriptAsm,
		Witness:      witnessHexStrs,
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

	return vin, nil
}

// processOutput processes a single output.
func processOutput(n uint32, output parser.RawOutput) (*models.Vout, error) {
	// Classify output script
	scriptType := parser.ClassifyScript(output.ScriptPubkey)

	// Get address if recognized
	addr := getAddressFromScript(output.ScriptPubkey)

	// Disassemble script
	scriptAsm, _ := parser.DisassembleScript(output.ScriptPubkey)

	// Convert script type to output format (lowercase)
	scriptTypeStr := scriptTypeToOutputString(scriptType)

	vout := &models.Vout{
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
			if len(payloadBytes) > 0 && utf8.Valid(payloadBytes) {
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

	return vout, nil
}

// getAddressFromScript derives an address from a script.
func getAddressFromScript(script []byte) *string {
	scriptType := parser.ClassifyScript(script)

	var addr string
	var err error

	switch scriptType {
	case parser.ScriptTypeP2PKH:
		hash, _ := parser.ExtractPubKeyHash(script)
		addr, err = address.EncodeP2PKH(hash)
	case parser.ScriptTypeP2SH:
		hash, _ := parser.ExtractScriptHash(script)
		addr, err = address.EncodeP2SH(hash)
	case parser.ScriptTypeP2WPKH:
		hash, _ := parser.ExtractPubKeyHash(script)
		addr, err = address.EncodeP2WPKH(hash)
	case parser.ScriptTypeP2WSH:
		hash, _ := parser.ExtractScriptHash(script)
		addr, err = address.EncodeP2WSH(hash)
	case parser.ScriptTypeP2TR:
		key, _ := parser.ExtractTaprootKey(script)
		addr, err = address.EncodeP2TR(key)
	default:
		return nil
	}

	if err != nil {
		return nil
	}
	return &addr
}

// scriptTypeToOutputString converts ScriptType to output JSON string format.
func scriptTypeToOutputString(st parser.ScriptType) string {
	switch st {
	case parser.ScriptTypeP2PKH:
		return "p2pkh"
	case parser.ScriptTypeP2SH:
		return "p2sh"
	case parser.ScriptTypeP2WPKH:
		return "p2wpkh"
	case parser.ScriptTypeP2WSH:
		return "p2wsh"
	case parser.ScriptTypeP2TR:
		return "p2tr"
	case parser.ScriptTypeOPReturn:
		return "op_return"
	default:
		return "unknown"
	}
}

// parseRelativeTimelock parses BIP68 relative timelock from sequence.
func parseRelativeTimelock(sequence uint32) models.RelativeTimelock {
	// Bit 31 disables relative timelock
	if sequence&(1<<31) != 0 {
		return models.RelativeTimelock{Enabled: false}
	}

	// If sequence >= 0xFFFFFFFE, no relative timelock
	if sequence >= 0xFFFFFFFE {
		return models.RelativeTimelock{Enabled: false}
	}

	// Bit 22 determines type: 0 = blocks, 1 = time
	isTime := (sequence & (1 << 22)) != 0

	// Lower 16 bits are the value
	value := sequence & 0xFFFF

	if isTime {
		// Time-based: value is in 512-second units
		seconds := value * 512
		lockType := "time"
		return models.RelativeTimelock{
			Enabled: true,
			Type:    &lockType,
			Value:   &seconds,
		}
	} else {
		// Block-based
		lockType := "blocks"
		valueU32 := uint32(value)
		return models.RelativeTimelock{
			Enabled: true,
			Type:    &lockType,
			Value:   &valueU32,
		}
	}
}

// generateWarnings generates warning codes based on transaction analysis.
func generateWarnings(result *models.TransactionResult) []models.Warning {
	// Initialize as empty slice (not nil) so JSON serializes as [] not null
	warnings := make([]models.Warning, 0)

	// HIGH_FEE: fee > 1,000,000 sats OR fee rate > 200 sat/vB
	if result.FeeSats > 1_000_000 || result.FeeRateSatVB > 200 {
		warnings = append(warnings, models.Warning{Code: models.WarningHighFee})
	}

	// DUST_OUTPUT: any non-op_return output < 546 sats
	for _, vout := range result.Vout {
		if vout.ScriptType != "op_return" && vout.ValueSats < models.DustThreshold {
			warnings = append(warnings, models.Warning{Code: models.WarningDustOutput})
			break
		}
	}

	// UNKNOWN_OUTPUT_SCRIPT: any output with unknown script type
	for _, vout := range result.Vout {
		if vout.ScriptType == "unknown" {
			warnings = append(warnings, models.Warning{Code: models.WarningUnknownOutputScript})
			break
		}
	}

	// RBF_SIGNALING
	if result.RBFSignaling {
		warnings = append(warnings, models.Warning{Code: models.WarningRBFSignaling})
	}

	return warnings
}

// reversedBytes returns a copy of the byte slice with bytes reversed.
func reversedBytes(b []byte) []byte {
	result := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		result[i] = b[len(b)-1-i]
	}
	return result
}
