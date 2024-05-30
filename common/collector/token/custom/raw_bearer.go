package custom

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"ivxv.ee/common/collector/token"
)

func (c *fromRaw) Unmarshal() (token.Token, error) {
	// Regex over raw bearer token
	ok := rawBearerRegex(c.Bearer)
	if !ok {
		return nil, RawBearerRegexError{Bearer: c.Bearer}
	}

	// Split raw bearer token <payload>.<signature>
	bJSON64, bSig64 := splitRawBearer(c.Bearer)
	if bJSON64 == "" || bSig64 == "" {
		return nil, SplitRawBearerTokenError{
			Payload:   bJSON64,
			Signature: bSig64,
		}
	}

	// Base64 decode payload
	bJSON, err := base64.StdEncoding.DecodeString(bJSON64)
	if err != nil {
		return nil, Base64DecodePayloadError{Err: err}
	}

	// Base64 decode signature
	bSig, err := base64.StdEncoding.DecodeString(bSig64)
	if err != nil {
		return nil, Base64DecodeSignatureError{Err: err}
	}

	// JSON unmarshal payload
	err = json.Unmarshal(bJSON, c)
	if err != nil {
		return nil, JSONUnmarshalError{Err: err}
	}

	c.Sig = bSig
	return c, nil
}

// Header is not implemented.
func (c *fromRaw) Header() (string, error) {
	return "", nil
}

// Payload returns a nonce.
func (c *fromRaw) Payload() (string, error) {
	return c.Nonce, nil
}

// Signature returns a signature.
func (c *fromRaw) Signature() ([]byte, error) {
	return c.Sig, nil
}

// Verify only verifies expiration time of a token.
func (c *fromRaw) Verify() error {
	now := time.Now()
	expiresAt := c.CreatedAt.Add(ttl)
	if !now.Before(expiresAt) {
		return ExpiredBearerTokenError{
			Now:       now.String(),
			ExpiredAt: expiresAt.String(),
		}
	}
	return nil
}
