package webeid

import (
	"testing"

	"ivxv.ee/common/collector/cryptoutil"
)

const (
	testKeyPath  = "testdata/voter.auth.key"
	testCertPath = "testdata/voter.auth.pem"
)

// RPC.TokenReq.Token
var tokenReqToken string

var originURL = []byte("https://ivxv1.test.ivxv.ee:443")
var challenge []byte

func TestCorrectWebeidToken(t *testing.T) {
	// Generate nonce
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	// Generate test Web eID token just like "Token" you send to RPC.TokenReq
	tokenReqToken = GenerateTestToken(
		testKeyPath, testCertPath, originURL, []byte(nonce))
	challenge = []byte(nonce)

	// Unmarshal raw Web eID token to token.Token
	wToken, err := NewFromRawBuilder().
		WithToken(tokenReqToken).
		WithNonce(string(challenge)).
		WithOrigin(originURL).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	// Verify Web eID token according to Web eID specification
	err = wToken.Verify()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}
}

func TestBadUnverifiedCertificateWebeidToken(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	tokenReqToken = GenerateBadUnverifiedCertificateToken(
		testKeyPath, testCertPath, originURL, []byte(nonce))
	challenge = []byte(nonce)

	// Unmarshalling should be unsuccessful
	_, err = NewFromRawBuilder().
		WithToken(tokenReqToken).
		WithNonce(string(challenge)).
		WithOrigin(originURL).
		Build().
		Unmarshal()
	if _, ok := err.(Base64DecodeUnverifiedCertificateError); !ok {
		msg := "Expected Base64DecodeUnverifiedCertificateError, got %v\n"
		t.Errorf(msg, err)
	}
}

func TestBadAlgorithmWebeidToken(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	tokenReqToken = GenerateBadAlgorithmToken(
		testKeyPath, testCertPath, originURL, []byte(nonce))
	challenge = []byte(nonce)

	wToken, err := NewFromRawBuilder().
		WithToken(tokenReqToken).
		WithNonce(string(challenge)).
		WithOrigin(originURL).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	err = wToken.Verify()
	if _, ok := err.(CheckSignatureError); !ok {
		msg := "Expected CheckSignatureError, got %v"
		t.Errorf(msg, err)
	}
}

func TestBadSignatureWebeidToken(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	tokenReqToken = GenerateBadSignatureToken(
		testKeyPath, testCertPath, originURL, []byte(nonce))
	challenge = []byte(nonce)

	wToken, err := NewFromRawBuilder().
		WithToken(tokenReqToken).
		WithNonce(string(challenge)).
		WithOrigin(originURL).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	err = wToken.Verify()
	if _, ok := err.(CheckSignatureError); !ok {
		msg := "Expected CheckSignatureError, got %v"
		t.Errorf(msg, err)
	}
}

func TestBadFormatWebeidToken(t *testing.T) {
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	tokenReqToken = GenerateBadFormatToken(
		testKeyPath, testCertPath, originURL, []byte(nonce))
	challenge = []byte(nonce)

	wToken, err := NewFromRawBuilder().
		WithToken(tokenReqToken).
		WithNonce(string(challenge)).
		WithOrigin(originURL).
		Build().
		Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	err = wToken.Verify()
	if _, ok := err.(VerifyFormatError); !ok {
		msg := "Expected VerifyFormatError, got %v"
		t.Errorf(msg, err)
	}
}
