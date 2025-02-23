package cryptoutil

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"ivxv.ee/common/collector/errors"

	_ "crypto/sha256"
	_ "crypto/sha512"
)

func TestDigestInfo(t *testing.T) {
	tests := []struct {
		name string
		hash crypto.Hash
		want []byte
	}{
		{"SHA-224", crypto.SHA224, []byte{
			0x30, 0x2d, 0x30, 0x0d,
			0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x04,
			0x04, 0x1c,
			0xd1, 0x4a, 0x02, 0x8c, 0x2a, 0x3a, 0x2b, 0xc9,
			0x47, 0x61, 0x02, 0xbb, 0x28, 0x82, 0x34, 0xc4,
			0x15, 0xa2, 0xb0, 0x1f, 0x82, 0x8e, 0xa6, 0x2a,
			0xc5, 0xb3, 0xe4, 0x2f,
		}},
		{"SHA-256", crypto.SHA256, []byte{
			0x30, 0x31, 0x30, 0x0d,
			0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x01,
			0x04, 0x20,
			0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14,
			0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24,
			0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c,
			0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55,
		}},
		{"SHA-384", crypto.SHA384, []byte{
			0x30, 0x41, 0x30, 0x0d,
			0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x02,
			0x04, 0x30,
			0x38, 0xb0, 0x60, 0xa7, 0x51, 0xac, 0x96, 0x38,
			0x4c, 0xd9, 0x32, 0x7e, 0xb1, 0xb1, 0xe3, 0x6a,
			0x21, 0xfd, 0xb7, 0x11, 0x14, 0xbe, 0x07, 0x43,
			0x4c, 0x0c, 0xc7, 0xbf, 0x63, 0xf6, 0xe1, 0xda,
			0x27, 0x4e, 0xde, 0xbf, 0xe7, 0x6f, 0x65, 0xfb,
			0xd5, 0x1a, 0xd2, 0xf1, 0x48, 0x98, 0xb9, 0x5b,
		}},
		{"SHA-512", crypto.SHA512, []byte{
			0x30, 0x51, 0x30, 0x0d,
			0x06, 0x09, 0x60, 0x86, 0x48, 0x01, 0x65, 0x03, 0x04, 0x02, 0x03,
			0x04, 0x40,
			0xcf, 0x83, 0xe1, 0x35, 0x7e, 0xef, 0xb8, 0xbd,
			0xf1, 0x54, 0x28, 0x50, 0xd6, 0x6d, 0x80, 0x07,
			0xd6, 0x20, 0xe4, 0x05, 0x0b, 0x57, 0x15, 0xdc,
			0x83, 0xf4, 0xa9, 0x21, 0xd3, 0x6c, 0xe9, 0xce,
			0x47, 0xd0, 0xd1, 0x3c, 0x5d, 0x85, 0xf2, 0xb0,
			0xff, 0x83, 0x18, 0xd2, 0x87, 0x7e, 0xec, 0x2f,
			0x63, 0xb9, 0x31, 0xbd, 0x47, 0x41, 0x7a, 0x81,
			0xa5, 0x38, 0x32, 0x7a, 0xf9, 0x27, 0xda, 0x3e,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := DigestInfo(test.hash, nil); !bytes.Equal(got, test.want) {
				t.Errorf("unexpected results: got %x, want %x", got, test.want)
			}
		})
	}
}

func TestPEMDecode(t *testing.T) {
	tests := []struct {
		name  string
		file  string
		cause error
	}{
		{"PEM", "certificate.pem", nil},
		{"DER", "certificate.der", new(NotPEMEncodingError)},
		{"trailing data", "certificate-trailing-data.pem", new(PEMTrailingDataError)},
		{"wrong type", "certificate-wrong-type.pem", new(PEMBlockTypeError)},
		{"with headers", "certificate-with-header.pem", new(PEMHeadersError)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", test.file))
			if err != nil {
				t.Fatalf("failed to read %s: %v\n", test.file, err)
			}
			_, err = PEMDecode(string(data), "CERTIFICATE")
			if err != test.cause && errors.CausedBy(err, test.cause) == nil {
				t.Errorf("unexpected error: got %v, want cause %T", err, test.cause)
			}
		})
	}
}

func TestNonce44Bytes(t *testing.T) {
	// Generate a nonce
	nonce, err := Nonce44Bytes()
	if err != nil {
		msg := "Expected no errors, got %v\n"
		t.Errorf(msg, err)
	}

	// Test that nonce is base64 encoded string
	_, err = base64.StdEncoding.DecodeString(nonce)
	if err != nil {
		msg := "Expected nonce to be base64, but got %v"
		t.Errorf(msg, err)
	}

	// Test that nonce is base64 encoded string with a length of 44 bytes
	if len(nonce) != 44 {
		msg := "Expected nonce to be %d bytes, but got %d"
		t.Errorf(msg, 44, len(nonce))
	}

	// Test that base64 decoded nonce is 32 bytes long
	b, _ := base64.StdEncoding.DecodeString(nonce)
	if len(b) != 32 {
		msg := "Expected base64 decoded nonce to be %d bytes, but got %d"
		t.Errorf(msg, 32, len(b))
	}
}
