package custom

import "ivxv.ee/common/collector/token"

type fromExisting struct {
	Sig  string `json:"-"`
	Data string `json:"-"`
}

type FromExistingBuilder struct {
	payload   string
	signature string
}

// NewFromExistingBuilder is a Builder-pattern constructor, which is used
// to finalize a Bearer token before Marshalling it to a raw Bearer token.
func NewFromExistingBuilder() *FromExistingBuilder {
	return new(FromExistingBuilder)
}

func (feb *FromExistingBuilder) WithPayload(p string) *FromExistingBuilder {
	feb.payload = p
	return feb
}

func (feb *FromExistingBuilder) WithSignature(s string) *FromExistingBuilder {
	feb.signature = s
	return feb
}

// Build will prepare a Bearer token to be marshalled into a raw Bearer token string.
func (feb *FromExistingBuilder) Build() token.Marshaller {
	return &fromExisting{
		Sig:  feb.signature,
		Data: feb.payload,
	}
}
