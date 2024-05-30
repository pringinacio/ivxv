// Package client defines an API for any implementation services that
// wish to communicate with a status reporting service (server).
package client

// All possible user authentication methods against IVXV backend.
const (
	IDcardAuth   = "id"
	MobileIDAuth = "mid"
	SmartIDAuth  = "sid"
	WebeIDAuth   = "wid"
	NoAuth       = ""
)

// Verifier is the interface that each implementation should use in order
// to embed business logic for status service response validation.
type Verifier interface {
	// Verify returns ok as true if request req is successfully verified.
	Verify(req interface{}) (ok bool, err error)
}

// TLSDialer interface governs the rules of interaction with a
// status service. Data, that is passed in req/resp as an
// interface{} (DTO), is used by implementations to build
// up and then perform requests against a status service.
//
// This means, that implementations should provide DTOs for any
// req/resp.
type TLSDialer interface {
	// TLSDial sends a request req over the TLS connection to a
	// status server, which responds with a resp and an error err.
	TLSDial(req interface{}) (resp interface{}, err error)
}
