package server

import (
	"net/url"
)

// VerifyHTTPSOrigin verifies origin to be compliant with https://<URL>:443 schema,
// e.g. valid URL is https://mydomain.ee:443
func VerifyHTTPSOrigin(origin string) bool {
	r, err := url.ParseRequestURI(origin)
	if err != nil {
		return false
	}

	// `?` not allowed, e.g. https://mydomain.ee:443?
	if r.ForceQuery {
		return false
	}

	// Hostname must present
	if r.Hostname() == "" {
		return false
	}

	// Scheme is only `https` lowercase
	if r.Scheme != "https" {
		return false
	}

	// Port only 443
	if r.Port() != "443" {
		return false
	}

	// No query paths allowed, e.g. https://mydomain.ee:443/api/v2,
	// even https://mydomain.ee:443/
	if r.Path != "" {
		return false
	}

	return true
}
