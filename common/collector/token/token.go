package token

import (
	"crypto/x509"
)

// Certifier is the interface that allows any implementation to define
// its own way to retrieve a x509 certificate.
type Certifier interface {
	// Certify may return a x509 certificate to the caller, if any exists.
	Certify() *x509.Certificate
}

// Token is the interface that defines basic operations on a Web token.
type Token interface {
	// Header may return a header part of a token, if any exists.
	// However, if you have some implementation that defines its own
	// meaning for a header, then it is also OK and should be documented
	// for that use case.
	Header() (string, error)

	// Payload may return a payload part of a token, if any exists.
	// However, if you have some implementation that defines its own
	// meaning for a payload, then it is also OK and should be documented
	// for that use case.
	Payload() (string, error)

	// Signature may return a signature part of a token, if any exists.
	Signature() ([]byte, error)

	// Verify should verify a token. Note, that here is not specified
	// what has to be verified exactly, this is up to the implementation.
	Verify() error
}

// Unmarshaler is the interface used to deserialize raw token into a Token.
type Unmarshaler interface {

	// Unmarshal will deserialize raw token into a Token.
	Unmarshal() (Token, error)
}

// Marshaller is the interface used to serialize Token into a raw token.
type Marshaller interface {

	// Marshal will serialize Token into a raw token.
	Marshal() (string, error)
}
