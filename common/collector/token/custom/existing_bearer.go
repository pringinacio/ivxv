package custom

import "encoding/base64"

func (b *fromExisting) Marshal() (string, error) {
	// Bearer token == <base64>.<base64>
	return b.Data + "." + base64.StdEncoding.EncodeToString([]byte(b.Sig)), nil
}
