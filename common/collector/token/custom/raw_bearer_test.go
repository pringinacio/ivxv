package custom

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"
)

func TestRawBearerUnmarshalBadRawBearer(t *testing.T) {
	b64Payload := base64.StdEncoding.EncodeToString([]byte("Hello"))
	b64Sig := base64.StdEncoding.EncodeToString([]byte("World"))

	rawBearer := NewFromRawBuilder().
		WithBearer(b64Payload + "." + b64Sig).
		Build()

	msg := "Expected JSONUnmarshalError, got %v\n"
	_, err := rawBearer.Unmarshal()
	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(JSONUnmarshalError); !ok {
		t.Errorf(msg, err)
	}
}

func TestRawBearerBadBase64Signature(t *testing.T) {
	b64Payload := base64.StdEncoding.EncodeToString([]byte("Hello"))
	b64Sig := "World"

	rawBearer := NewFromRawBuilder().
		WithBearer(b64Payload + "." + b64Sig).
		Build()

	msg := "Expected Base64DecodeSignature, got %v\n"
	_, err := rawBearer.Unmarshal()
	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(Base64DecodeSignatureError); !ok {
		t.Errorf(msg, err)
	}
}

func TestRawBearerBadBase64Payload(t *testing.T) {
	b64Payload := "Hello"
	b64Sig := base64.StdEncoding.EncodeToString([]byte("World"))

	rawBearer := NewFromRawBuilder().
		WithBearer(b64Payload + "." + b64Sig).
		Build()

	msg := "Expected Base64DecodePayload, got %v\n"
	_, err := rawBearer.Unmarshal()
	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(Base64DecodePayloadError); !ok {
		t.Errorf(msg, err)
	}
}

func TestRawBearerBadSplitRawBearerToken(t *testing.T) {
	b64Payload := base64.StdEncoding.EncodeToString([]byte("Hello"))
	b64Sig := base64.StdEncoding.EncodeToString([]byte("World"))

	rawBearer := NewFromRawBuilder().
		WithBearer(b64Payload + "." + b64Sig + ".").
		Build()

	msg := "Expected SplitRawBearerTokenError, got %v\n"
	_, err := rawBearer.Unmarshal()
	if err == nil {
		t.Errorf(msg, err)
	}
	if _, ok := err.(SplitRawBearerTokenError); !ok {
		t.Errorf(msg, err)
	}
}

func TestRawBearerDoesntImplementHeader(t *testing.T) {
	rawBearer := fromRaw{Nonce: "123", CreatedAt: time.Now()}

	marshalled, err := json.Marshal(rawBearer)
	if err != nil {
		if err == nil {
			t.Errorf(msg, err)
		}
	}

	b64Payload := base64.StdEncoding.EncodeToString(marshalled)
	b64Sig := base64.StdEncoding.EncodeToString([]byte("World"))

	toUnmarshalBearer := NewFromRawBuilder().
		WithBearer(b64Payload + "." + b64Sig).
		Build()

	bearer, err := toUnmarshalBearer.Unmarshal()
	if err != nil {
		t.Errorf(msg, err)
	}

	header, err := bearer.Header()
	if header != "" && err != nil {
		t.Errorf(msg, err)
	}
}
