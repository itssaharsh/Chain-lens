// Package xor provides XOR decoding for Bitcoin Core obfuscated files.
package xor

import (
	"os"

	"chainlens/internal/models"
)

// =============================================================================
// BITCOIN CORE XOR OBFUSCATION
// =============================================================================
//
// Since Bitcoin Core v0.13, block and undo files (blk*.dat, rev*.dat) are
// XOR-obfuscated with a key stored in the chainstate LevelDB database.
//
// The xor.dat file contains this key (typically 8 bytes).
//
// To decode:
//   decoded[i] = encoded[i] ^ key[i % len(key)]
//
// If the key is all zeros, no transformation is needed.
// =============================================================================

// Key represents an XOR obfuscation key.
type Key []byte

// LoadKey reads the XOR key from a file.
func LoadKey(path string) (Key, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, models.NewAnalysisErrorf(models.ErrCodeFileReadError,
			"failed to read XOR key file: %v", err)
	}

	if len(data) == 0 {
		return nil, models.NewAnalysisErrorf(models.ErrCodeInvalidXORKey,
			"XOR key file is empty")
	}

	return Key(data), nil
}

// IsZero returns true if the key is all zeros (no obfuscation needed).
func (k Key) IsZero() bool {
	for _, b := range k {
		if b != 0 {
			return false
		}
	}
	return true
}

// Decode decodes XOR-obfuscated data in place.
// If the key is all zeros, this is a no-op.
//
// The XOR operation is applied cyclically:
//   data[i] ^= key[i % len(key)]
func (k Key) Decode(data []byte) {
	if k.IsZero() || len(k) == 0 {
		return // No transformation needed
	}

	keyLen := len(k)
	for i := range data {
		data[i] ^= k[i%keyLen]
	}
}

// DecodeRange decodes a portion of data starting at the given file offset.
// This is important because the XOR key position depends on the file offset,
// not the buffer offset.
//
// Example: If reading bytes 100-199 from a file, we need:
//   data[0] ^= key[100 % keyLen]
//   data[1] ^= key[101 % keyLen]
//   etc.
func (k Key) DecodeRange(data []byte, fileOffset int64) {
	if k.IsZero() || len(k) == 0 {
		return
	}

	keyLen := int64(len(k))
	for i := range data {
		keyIdx := (fileOffset + int64(i)) % keyLen
		data[i] ^= k[keyIdx]
	}
}

// DecodeCopy returns a decoded copy of the data (does not modify original).
func (k Key) DecodeCopy(data []byte) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	k.Decode(result)
	return result
}

// DecodeCopyRange returns a decoded copy with file offset consideration.
func (k Key) DecodeCopyRange(data []byte, fileOffset int64) []byte {
	result := make([]byte, len(data))
	copy(result, data)
	k.DecodeRange(result, fileOffset)
	return result
}
