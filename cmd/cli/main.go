// Package main provides the CLI entry point for Chain Lens.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"chainlens/internal/analyzer"
	"chainlens/internal/models"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)

// Command represents the parsed CLI arguments.
type Command struct {
	Mode        string   // "transaction" or "block"
	FixturePath string   // path to fixture JSON (transaction mode)
	BlkPath     string   // path to blk*.dat (block mode)
	RevPath     string   // path to rev*.dat (block mode)
	XorPath     string   // path to xor.dat (block mode)
}

func main() {
	cmd, err := parseArgs(os.Args[1:])
	if err != nil {
		writeErrorAndExit(err)
	}

	if cmd.Mode == "block" {
		runBlockMode(cmd)
	} else {
		runTransactionMode(cmd)
	}
}

// parseArgs parses command line arguments and returns a Command.
func parseArgs(args []string) (*Command, error) {
	if len(args) == 0 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidFixture,
			"usage: chainlens <fixture.json> OR chainlens --block <blk.dat> <rev.dat> <xor.dat>")
	}

	// Check for block mode
	if args[0] == "--block" {
		if len(args) != 4 {
			return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidFixture,
				"block mode usage: chainlens --block <blk.dat> <rev.dat> <xor.dat>")
		}
		return &Command{
			Mode:    "block",
			BlkPath: args[1],
			RevPath: args[2],
			XorPath: args[3],
		}, nil
	}

	// Transaction mode: single fixture path
	if len(args) != 1 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidFixture,
			"usage: chainlens <fixture.json> OR chainlens --block <blk.dat> <rev.dat> <xor.dat>")
	}

	return &Command{
		Mode:        "transaction",
		FixturePath: args[0],
	}, nil
}

// runTransactionMode handles single transaction analysis.
func runTransactionMode(cmd *Command) {
	// Read fixture file
	fixtureData, err := os.ReadFile(cmd.FixturePath)
	if err != nil {
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to read fixture file: %v", err))
	}

	// Parse fixture JSON
	var fixture models.Fixture
	if err := json.Unmarshal(fixtureData, &fixture); err != nil {
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeInvalidJSON,
			"failed to parse fixture JSON: %v", err))
	}

	// Validate fixture
	if err := fixture.Validate(); err != nil {
		writeErrorAndExit(err)
	}

	// TODO: Actually parse and analyze the transaction
	// For now, this is a placeholder that shows the structure
	result, err := analyzeTransaction(&fixture)
	if err != nil {
		if analysisErr, ok := err.(*models.AnalysisError); ok {
			writeErrorAndExit(analysisErr)
		}
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeInvalidTx, "%v", err))
	}

	// Ensure output directory exists
	if err := os.MkdirAll("out", 0755); err != nil {
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to create output directory: %v", err))
	}

	// Write result to file
	outputPath := filepath.Join("out", result.Txid+".json")
	outputData, _ := json.MarshalIndent(result, "", "  ")
	if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to write output file: %v", err))
	}

	// Also print to stdout for single transaction mode
	fmt.Println(string(outputData))
	os.Exit(ExitSuccess)
}

// runBlockMode handles block file analysis.
func runBlockMode(cmd *Command) {
	// Validate file existence
	for _, path := range []string{cmd.BlkPath, cmd.RevPath, cmd.XorPath} {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileNotFound,
				"file not found: %s", path))
		}
	}

	// TODO: Implement block parsing
	// 1. Read XOR key
	// 2. Read and decode blk*.dat
	// 3. Read and decode rev*.dat
	// 4. Parse block header(s)
	// 5. Parse transactions
	// 6. Match undo data to inputs
	// 7. Verify merkle root
	// 8. Generate output

	results, err := analyzeBlock(cmd.BlkPath, cmd.RevPath, cmd.XorPath)
	if err != nil {
		if analysisErr, ok := err.(*models.AnalysisError); ok {
			writeErrorAndExit(analysisErr)
		}
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeInvalidBlockHeader, "%v", err))
	}

	// Ensure output directory exists
	if err := os.MkdirAll("out", 0755); err != nil {
		writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to create output directory: %v", err))
	}

	// Write each block result to a separate file
	for _, result := range results {
		outputPath := filepath.Join("out", result.BlockHeader.BlockHash+".json")
		outputData, _ := json.MarshalIndent(result, "", "  ")
		if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
			writeErrorAndExit(models.NewAnalysisErrorf(models.ErrCodeFileReadError,
				"failed to write output file: %v", err))
		}
	}

	// Block mode does not print to stdout
	os.Exit(ExitSuccess)
}

// analyzeTransaction parses and analyzes a single transaction.
func analyzeTransaction(fixture *models.Fixture) (*models.TransactionResult, error) {
	return analyzer.AnalyzeTransaction(fixture)
}

// analyzeBlock parses and analyzes a block file.
func analyzeBlock(blkPath, revPath, xorPath string) ([]*models.BlockResult, error) {
	return analyzer.AnalyzeBlock(blkPath, revPath, xorPath)
}

// writeErrorAndExit writes an error response to stdout and exits with error code.
func writeErrorAndExit(err error) {
	var result models.ErrorResult
	if analysisErr, ok := err.(*models.AnalysisError); ok {
		result = analysisErr.ToErrorResult()
	} else {
		result = models.NewErrorResult(models.ErrCodeInvalidTx, err.Error())
	}

	output, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(output))
	os.Exit(ExitError)
}
