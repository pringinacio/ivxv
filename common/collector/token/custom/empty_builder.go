package custom

import (
	"time"

	"ivxv.ee/common/collector/token"
)

type fromEmpty struct {
	Nonce     string    `json:"nonce"`
	CreatedAt time.Time `json:"createdAt"`
	Sig       string    `json:"-"`
	Data      string    `json:"-"`
}

type FromEmptyBuilder struct {
	nonce     string
	timestamp time.Time
}

// NewFromEmptyBuilder is a Builder-pattern constructor, which is used
// to create brand-new Bearer token, which is in an incomplete form.
//
// Incomplete form means that Bearer token doesn't have a signature
// and a payload yet. This particular Bearer token implementation generates
// payload on a Payload() call, but doesn't generate signature on a Signature()
// call, use NewFromExistingBuilder to create a Bearer token with a signature
// and don't forget to add a payload.
func NewFromEmptyBuilder() *FromEmptyBuilder {
	return new(FromEmptyBuilder)
}

func (feb *FromEmptyBuilder) WithNonce(n string) *FromEmptyBuilder {
	feb.nonce = n
	return feb
}

func (feb *FromEmptyBuilder) WithTimeStamp(t time.Time) *FromEmptyBuilder {
	feb.timestamp = t
	return feb
}

// Build returns a Bearer token with only Payload() method implemented.
//
// Payload() method will serialize nonce and timestamp into []byte.
func (feb *FromEmptyBuilder) Build() token.Token {
	return &fromEmpty{
		Nonce:     feb.nonce,
		CreatedAt: feb.timestamp,
	}
}
