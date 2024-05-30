package custom

import (
	"regexp"
	"strings"
	"time"
)

const (
	// <at least 1 char>.<at least 1 char>
	pattern = `(.+)\.(.+)`
	ttl     = 2 * time.Minute
)

// rawBearerRegex checks regex over a raw Bearer token b.
func rawBearerRegex(b string) bool {
	return regexp.MustCompile(pattern).MatchString(b)
}

// splitRaw splits raw Bearer token b into 2 parts (payload, signature).
func splitRawBearer(b string) (string, string) {
	rawList := strings.Split(b, ".")
	if len(rawList) != 2 {
		return "", ""
	}
	return rawList[0], rawList[1]
}
