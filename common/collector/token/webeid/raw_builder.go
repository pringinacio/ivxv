package webeid

import (
	"crypto/x509"

	"ivxv.ee/common/collector/token"
)

type fromRaw struct {
	Token     string            `json:"-"`
	Nonce     []byte            `json:"-"`
	Origin    []byte            `json:"-"`
	Sig       []byte            `json:"-"`
	Format    string            `json:"-"`
	Algorithm string            `json:"-"`
	Cert      *x509.Certificate `json:"-"`
}

type FromRawBuilder struct {
	token  string
	nonce  string
	origin []byte
}

// NewFromRawBuilder is a Builder-pattern constructor, which is used
// to prepare a raw Web eID token to be unmarshalled it to a token.Token.
func NewFromRawBuilder() *FromRawBuilder {
	return new(FromRawBuilder)
}

func (frb *FromRawBuilder) WithToken(t string) *FromRawBuilder {
	frb.token = t
	return frb
}

func (frb *FromRawBuilder) WithNonce(n string) *FromRawBuilder {
	frb.nonce = n
	return frb
}

func (frb *FromRawBuilder) WithOrigin(o []byte) *FromRawBuilder {
	frb.origin = o
	return frb
}

// Build will produce a Web eID token Unmarshaler with only Unmarshal() method
// implemented. Use that method to convert raw Web eID token to a token.Token.
//
// Once unmarshalled you get access to following methods:
//
// a) Verify() -> verifies Web eID token according to the Web eID token specification
//
// b) Signature() -> returns a Signature of a Web eID token
func (frb *FromRawBuilder) Build() token.Unmarshaler {
	return &fromRaw{
		Token:  frb.token,
		Nonce:  []byte(frb.nonce),
		Origin: frb.origin,
	}
}
