package custom

import (
	"testing"
	"time"

	"ivxv.ee/common/collector/cookie"
	"ivxv.ee/common/collector/cryptoutil"
)

var sharedSecret *cookie.C

var msg = "Expected no errors, got %v\n"
var msg2 = "Expected RawBearerRegexError, got %v\n"

func init() {
	var err error
	sharedSecret, err = cookie.New(make([]byte, 16))
	if err != nil {
		panic(err)
	}
}

func TestCorrectBearer(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msg, err)
	}

	bearerToken := NewFromEmptyBuilder().
		WithNonce(nonce).
		WithTimeStamp(time.Now()).
		Build()

	payload, err := bearerToken.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig := sharedSecret.Create([]byte(payload))

	rawComplete, err := NewFromExistingBuilder().
		WithPayload(payload).
		WithSignature(string(sig)).
		Build().
		Marshal()
	if err != nil {
		t.Errorf(msg, err)
	}

	tBearer, err := NewFromRawBuilder().
		WithBearer(rawComplete).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig2, err := tBearer.Signature()
	if err != nil {
		t.Errorf(msg, err)
	}

	data, err := sharedSecret.Open(sig2)
	if err != nil || data == nil {
		t.Errorf(msg, err)
	}

	p, err := tBearer.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	err = tBearer.Verify()
	if err != nil {
		t.Errorf(msg, err)
	}

	if nonce != p {
		//nolint:lll
		msg = "Expected generatedNonce == parsedFromBearerNone, got generatedNonce: %v, parsedFromBearerNone: %v\n"
		t.Errorf(msg, nonce, p)
	}
}

func TestWithoutTimeStampBearer(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msg, err)
	}

	bearerToken := NewFromEmptyBuilder().
		WithNonce(nonce).
		// WithTimeStamp(time.Now()).
		Build()

	payload, err := bearerToken.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig := sharedSecret.Create([]byte(payload))
	if err != nil {
		t.Errorf(msg, err)
	}

	rawComplete, err := NewFromExistingBuilder().
		WithPayload(payload).
		WithSignature(string(sig)).
		Build().
		Marshal()

	if err != nil {
		t.Errorf(msg, err)
	}

	tBearer, err := NewFromRawBuilder().
		WithBearer(rawComplete).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig2, err := tBearer.Signature()
	if err != nil {
		t.Errorf(msg, err)
	}

	data, err := sharedSecret.Open(sig2)
	if err != nil || data == nil {
		t.Errorf(msg, err)
	}

	err = tBearer.Verify()
	if err == nil {
		msg := "Expected error, got %v\n"
		t.Errorf(msg, err)
	}

	if _, ok := err.(ExpiredBearerTokenError); !ok {
		msg := "Expected ExpiredBearerTokenError, got %v\n"
		t.Errorf(msg, err)
	}
}

func TestWithoutNonceBearer(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msg, err)
	}

	bearerToken := NewFromEmptyBuilder().
		// WithNonce(nonce).
		WithTimeStamp(time.Now()).
		Build()

	payload, err := bearerToken.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig := sharedSecret.Create([]byte(payload))
	if err != nil {
		t.Errorf(msg, err)
	}

	rawComplete, err := NewFromExistingBuilder().
		WithPayload(payload).
		WithSignature(string(sig)).
		Build().
		Marshal()

	if err != nil {
		t.Errorf(msg, err)
	}

	tBearer, err := NewFromRawBuilder().
		WithBearer(rawComplete).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig2, err := tBearer.Signature()
	if err != nil {
		t.Errorf(msg, err)
	}

	data, err := sharedSecret.Open(sig2)
	if err != nil || data == nil {
		t.Errorf(msg, err)
	}

	err = tBearer.Verify()
	if err != nil {
		msg := "Expected error, got %v\n"
		t.Errorf(msg, err)
	}

	p, err := tBearer.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	if nonce == p {
		//nolint:lll
		msg = "Expected generatedNonce != parsedFromBearerNone, got generatedNonce: %v, parsedFromBearerNone: %v\n"
		t.Errorf(msg, nonce, p)
	}

	if p != "" {
		msg = `Expected parsedFromBearerNone == "", got parsedFromBearerNone: %v\n`
		t.Errorf(msg, p)
	}
}

func TestWithoutSignatureBearer(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msg, err)
	}

	bearerToken := NewFromEmptyBuilder().
		WithNonce(nonce).
		WithTimeStamp(time.Now()).
		Build()

	payload, err := bearerToken.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	rawComplete, err := NewFromExistingBuilder().
		WithPayload(payload).
		// WithSignature(string(sig)).
		Build().
		Marshal()

	if err != nil {
		t.Errorf(msg, err)
	}

	_, err = NewFromRawBuilder().
		WithBearer(rawComplete).
		Build().
		Unmarshal()
	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(RawBearerRegexError); !ok {
		t.Errorf(msg2, err)
	}
}

func TestWithoutPayloadBearer(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msg, err)
	}

	bearerToken := NewFromEmptyBuilder().
		WithNonce(nonce).
		WithTimeStamp(time.Now()).
		Build()

	payload, err := bearerToken.Payload()
	if err != nil {
		t.Errorf(msg, err)
	}

	sig := sharedSecret.Create([]byte(payload))

	rawComplete, err := NewFromExistingBuilder().
		// WithPayload(payload).
		WithSignature(string(sig)).
		Build().
		Marshal()

	if err != nil {
		t.Errorf(msg, err)
	}

	_, err = NewFromRawBuilder().
		WithBearer(rawComplete).
		Build().
		Unmarshal()

	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(RawBearerRegexError); !ok {
		t.Errorf(msg2, err)
	}
}

func TestWithoutRawBearerBearer(t *testing.T) {
	_, err := NewFromRawBuilder().
		Build().
		Unmarshal()

	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(RawBearerRegexError); !ok {
		t.Errorf(msg2, err)
	}
}
