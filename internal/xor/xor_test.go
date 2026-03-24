package xor

import (
	"bytes"
	"testing"
)

func TestXORDecode(t *testing.T) {
	// Test basic XOR decoding
	key := Key{0xAB, 0xCD, 0xEF}

	// Test data
	original := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	encoded := make([]byte, len(original))
	copy(encoded, original)

	// Encode (XOR with key)
	for i := range encoded {
		encoded[i] ^= key[i%len(key)]
	}

	// Verify encoding changed the data
	if bytes.Equal(encoded, original) {
		t.Error("Encoding should change the data")
	}

	// Decode
	key.Decode(encoded)

	// Verify decoding restored the original
	if !bytes.Equal(encoded, original) {
		t.Errorf("Decoding failed:\nexpected: %x\ngot:      %x", original, encoded)
	}
}

func TestXORDecodeRange(t *testing.T) {
	// Test XOR decoding with file offset
	key := Key{0xAB, 0xCD, 0xEF, 0x12}

	// Simulate reading bytes 100-109 from a file
	fileOffset := int64(100)
	data := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A}
	original := make([]byte, len(data))
	copy(original, data)

	// Encode with offset consideration
	encoded := make([]byte, len(data))
	for i := range data {
		keyIdx := (fileOffset + int64(i)) % int64(len(key))
		encoded[i] = data[i] ^ key[keyIdx]
	}

	// Decode with DecodeRange
	key.DecodeRange(encoded, fileOffset)

	if !bytes.Equal(encoded, original) {
		t.Errorf("DecodeRange failed:\nexpected: %x\ngot:      %x", original, encoded)
	}
}

func TestZeroKeyIsNoOp(t *testing.T) {
	// Zero key should not modify data
	key := Key{0x00, 0x00, 0x00, 0x00}

	if !key.IsZero() {
		t.Error("All-zero key should be detected as zero")
	}

	data := []byte{0x01, 0x02, 0x03, 0x04}
	original := make([]byte, len(data))
	copy(original, data)

	key.Decode(data)

	if !bytes.Equal(data, original) {
		t.Error("Zero key should not modify data")
	}
}

func TestDecodeCopy(t *testing.T) {
	key := Key{0xAB, 0xCD}
	data := []byte{0x01, 0x02, 0x03, 0x04}
	original := make([]byte, len(data))
	copy(original, data)

	decoded := key.DecodeCopy(data)

	// Original should be unchanged
	if !bytes.Equal(data, original) {
		t.Error("DecodeCopy should not modify original")
	}

	// Decoded should be different
	if bytes.Equal(decoded, original) {
		t.Error("DecodeCopy should produce decoded data")
	}

	// Decoding the decoded should give back original
	key.Decode(decoded)
	if !bytes.Equal(decoded, original) {
		t.Error("Double decode should restore original")
	}
}
