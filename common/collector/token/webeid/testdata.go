package webeid

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"os"
)

// getCertAndSig will generate test eID user's certificate and a signature.
func getCertAndSig(keyPath, certPath string, originURL, nonce []byte) ([]byte, []byte) {
	// Read authentication private key of eID user
	keyPem, err := os.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	// Read authentication certificate of eID user
	certPem, err := os.ReadFile(certPath)
	if err != nil {
		panic(err)
	}

	// Private key is in PEM format, so decode it to DER
	keyDer, _ := pem.Decode(keyPem)

	// Certificate is in PEM format, so decode it to DER
	certDer, _ := pem.Decode(certPem)

	// Parse private key in DER form to PKCS8
	pkcs8Key, err := x509.ParsePKCS8PrivateKey(keyDer.Bytes)
	if err != nil {
		panic(err)
	}

	// Extract private key as Go struct from PKCS8
	ecdsaKey, ok := pkcs8Key.(*ecdsa.PrivateKey)
	if !ok {
		panic("Cannot cast PKCS8 private key to ECDSA private key")
	}

	// SHA256(origin), it is hardcoded here, because certDer has
	// SignatureAlgorithm == ECDSA-SHA256
	h1 := sha256.New()
	h1.Write(originURL)
	originHash := h1.Sum(nil)

	// SHA256(nonce)
	h2 := sha256.New()
	h2.Write(nonce)
	nonceHash := h2.Sum(nil)

	// concatHash = SHA384(SHA256(origin), SHA256(nonce))
	originHash = append(originHash, nonceHash...)
	h3 := sha512.New384()
	h3.Write(originHash)
	concatHash := h3.Sum(nil)

	// Sign result with a eID user's private key
	sig, err := ecdsa.SignASN1(rand.Reader, ecdsaKey, concatHash)
	if err != nil {
		panic(err)
	}

	// Return eID user's authentication certificate and signature
	return certDer.Bytes, sig
}

// GenerateTestToken will generate test Web eID authentication token
// that is used in RPC TokenReq.
func GenerateTestToken(keyPath, certPath string, originURL, nonce []byte) string {
	// Generate eID user's auth cert and signature
	cert, sig := getCertAndSig(keyPath, certPath, originURL, nonce)

	// Build Web eID auth token up
	t := &webEidAuthToken{
		UnverifiedCertificate: base64.StdEncoding.EncodeToString(cert),
		Algorithm:             "ES256", // OID 1.2.840.10045.4.3.2
		Signature:             base64.StdEncoding.EncodeToString(sig),
		Format:                "web-eid:1.0",
		AppVersion:            "https://web-eid.eu/web-eid-app/releases/v2.0.0",
	}

	// Convert from Go struct to []byte
	token, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	// Return Web eID auth token []byte as a string
	return string(token)
}

// GenerateBadUnverifiedCertificateToken will generate Web eID token,
// where incorrect eID user's auth cert is used.
func GenerateBadUnverifiedCertificateToken(keyPath, certPath string, originURL, nonce []byte) string {
	_, sig := getCertAndSig(keyPath, certPath, originURL, nonce)

	t := &webEidAuthToken{
		// This is bad certificate
		UnverifiedCertificate: base64.StdEncoding.EncodeToString([]byte("Reji")),
		Algorithm:             "ES256",
		Signature:             base64.StdEncoding.EncodeToString(sig),
		Format:                "web-eid:1.0",
		AppVersion:            "https://web-eid.eu/web-eid-app/releases/v2.0.0",
	}

	token, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(token)
}

// GenerateBadAlgorithmToken will generate Web eID auth token which
// declares one Algorithm field, but uses different one for the actual
// signing.
func GenerateBadAlgorithmToken(keyPath, certPath string, originURL, nonce []byte) string {
	cert, sig := getCertAndSig(keyPath, certPath, originURL, nonce)

	t := &webEidAuthToken{
		UnverifiedCertificate: base64.StdEncoding.EncodeToString(cert),
		Algorithm:             "ES512", // should be ES256
		Signature:             base64.StdEncoding.EncodeToString(sig),
		Format:                "web-eid:1.0",
		AppVersion:            "https://web-eid.eu/web-eid-app/releases/v2.0.0",
	}

	token, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(token)
}

// GenerateBadSignatureToken will generate Web eID with the incorrect Signature.
func GenerateBadSignatureToken(keyPath, certPath string, originURL, nonce []byte) string {
	cert, _ := getCertAndSig(keyPath, certPath, originURL, nonce)

	t := &webEidAuthToken{
		UnverifiedCertificate: base64.StdEncoding.EncodeToString(cert),
		Algorithm:             "ES256",
		// This is incorrect signature
		Signature:  base64.StdEncoding.EncodeToString([]byte("Defi")),
		Format:     "web-eid:1.0",
		AppVersion: "https://web-eid.eu/web-eid-app/releases/v2.0.0",
	}

	token, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(token)
}

// GenerateBadFormatToken will generate Web eID with incorrect Format field.
func GenerateBadFormatToken(keyPath, certPath string, originURL, nonce []byte) string {
	cert, sig := getCertAndSig(keyPath, certPath, originURL, nonce)

	t := &webEidAuthToken{
		UnverifiedCertificate: base64.StdEncoding.EncodeToString(cert),
		Algorithm:             "ES256",
		Signature:             base64.StdEncoding.EncodeToString(sig),
		Format:                "web-eid1.0", // correct is web-eid:1.0
		AppVersion:            "https://web-eid.eu/web-eid-app/releases/v2.0.0",
	}

	token, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	return string(token)
}
