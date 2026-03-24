// Package main provides the web server entry point for Chain Lens.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"chainlens/internal/analyzer"
	"chainlens/internal/models"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()

	// Health check endpoint (required)
	mux.HandleFunc("GET /api/health", handleHealth)

	// Transaction analysis endpoint
	mux.HandleFunc("POST /api/analyze", handleAnalyze)

	// Block analysis endpoint
	mux.HandleFunc("POST /api/analyze/block", handleAnalyzeBlock)

	// Serve static files for the web UI
	staticFS := http.FileServer(http.Dir("web/static"))
	mux.Handle("/", staticFS)

	// Print the URL to stdout (required by spec)
	url := fmt.Sprintf("http://127.0.0.1:%s", port)
	fmt.Println(url)

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		server.Close()
	}()

	// Start server (blocks until shutdown)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// handleHealth handles GET /api/health
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// AnalyzeRequest is the request body for /api/analyze
type AnalyzeRequest struct {
	Network  string        `json:"network"`
	RawTx    string        `json:"raw_tx"`
	Prevouts []PrevoutData `json:"prevouts"`
}

// PrevoutData matches the fixture prevout format
type PrevoutData struct {
	Txid            string `json:"txid"`
	Vout            uint32 `json:"vout"`
	ValueSats       uint64 `json:"value_sats"`
	ScriptPubkeyHex string `json:"script_pubkey_hex"`
}

// handleAnalyze handles POST /api/analyze for transaction analysis
func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_JSON", err.Error()))
		return
	}

	// Build fixture from request
	fixture := &models.Fixture{
		Network: req.Network,
		RawTx:   req.RawTx,
	}
	for _, p := range req.Prevouts {
		fixture.Prevouts = append(fixture.Prevouts, models.FixturePrevout{
			Txid:            p.Txid,
			Vout:            p.Vout,
			ValueSats:       p.ValueSats,
			ScriptPubkeyHex: p.ScriptPubkeyHex,
		})
	}

	// Validate fixture
	if err := fixture.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if analysisErr, ok := err.(*models.AnalysisError); ok {
			json.NewEncoder(w).Encode(analysisErr.ToErrorResult())
		} else {
			json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_FIXTURE", err.Error()))
		}
		return
	}

	// Analyze transaction
	result, err := analyzer.AnalyzeTransaction(fixture)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if analysisErr, ok := err.(*models.AnalysisError); ok {
			json.NewEncoder(w).Encode(analysisErr.ToErrorResult())
		} else {
			json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_TX", err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// BlockAnalyzeRequest is the request body for /api/analyze/block
type BlockAnalyzeRequest struct {
	BlkData []byte `json:"blk_data"` // Base64 encoded blk*.dat content
	RevData []byte `json:"rev_data"` // Base64 encoded rev*.dat content
	XorKey  []byte `json:"xor_key"`  // Base64 encoded xor.dat content (8 bytes)
}

// handleAnalyzeBlock handles POST /api/analyze/block for block analysis
func handleAnalyzeBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req BlockAnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_JSON", err.Error()))
		return
	}

	// Validate XOR key
	if len(req.XorKey) != 8 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_XOR_KEY",
			fmt.Sprintf("XOR key must be 8 bytes, got %d", len(req.XorKey))))
		return
	}

	// XOR decrypt block data
	blkData := make([]byte, len(req.BlkData))
	copy(blkData, req.BlkData)
	for i := range blkData {
		blkData[i] ^= req.XorKey[i%8]
	}

	// XOR decrypt undo data
	revData := make([]byte, len(req.RevData))
	copy(revData, req.RevData)
	for i := range revData {
		revData[i] ^= req.XorKey[i%8]
	}

	// Analyze all blocks from the decrypted data
	results, err := analyzer.AnalyzeBlockFromData(blkData, revData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if analysisErr, ok := err.(*models.AnalysisError); ok {
			json.NewEncoder(w).Encode(analysisErr.ToErrorResult())
		} else {
			json.NewEncoder(w).Encode(models.NewErrorResult("INVALID_BLOCK_HEADER", err.Error()))
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}
