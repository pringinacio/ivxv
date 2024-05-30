package custom

import (
	"time"

	"ivxv.ee/common/collector/token"
)

type fromRaw struct {
	Nonce     string    `json:"nonce"`
	CreatedAt time.Time `json:"createdAt"`
	Bearer    string    `json:"-"`
	Sig       []byte    `json:"-"`
}

type FromRawBuilder struct {
	bearer string
}

// NewFromRawBuilder is a Builder-pattern constructor, which is used
// to prepare a raw Bearer token to be unmarshalled it to a token.Token.
func NewFromRawBuilder() *FromRawBuilder {
	return new(FromRawBuilder)
}

func (frb *FromRawBuilder) WithBearer(b string) *FromRawBuilder {
	frb.bearer = b
	return frb
}

// Build will prepare a raw Bearer token to be unmarshalled into a Bearer token.
//
// Unmarshalled Bearer token implements all methods of a token.Token, except
// Header().
//
// It returns a nonce on a Payload() call.
//
// Verify() method will only check expiration time of a Bearer token, which is
// 2 minutes by default.
//
// Signature() will return a signature over a Bearer token, all crypto checks
// are done outside the implementation, since shared secret is used.
func (frb *FromRawBuilder) Build() token.Unmarshaler {
	return &fromRaw{
		Bearer: frb.bearer,
	}
}
