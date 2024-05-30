package custom

import (
	"bytes"
	"testing"
)

func TestEmptyBearerDoesntImplementHeader(t *testing.T) {
	emptyBearer := NewFromEmptyBuilder().Build()
	header, err := emptyBearer.Header()

	if header != "" && err != nil {
		t.Errorf(msg, err)
	}
}

func TestEmptyBearerDoesntImplementSignature(t *testing.T) {
	emptyBearer := NewFromEmptyBuilder().Build()
	sig, err := emptyBearer.Signature()

	if !bytes.Equal(sig, nil) || err != nil {
		msg := "Expected sig: nil, err: nil, got sig %v, err: %v\n"
		t.Errorf(msg, sig, err)
	}
}

func TestEmptyBearerDoesntImplementVerify(t *testing.T) {
	emptyBearer := NewFromEmptyBuilder().Build()
	err := emptyBearer.Verify()

	if err != nil {
		msg := "Expected err == nil, got %v\n"
		t.Errorf(msg, err)
	}
}
