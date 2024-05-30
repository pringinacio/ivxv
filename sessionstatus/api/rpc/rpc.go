// Package rpc defines an API for any implementation services that
// wish to communicate with a session status reporting service (server).
//
// This API is not supposed to be a part of common/collector packages collection,
// since this API is the session status service specific and therefore cannot
// be reused elsewhere.
package rpc

const (
	sessionIDField = "SessionID"
	callerField    = "Caller"
	authField      = "Auth"
	leaseField     = "Lease"
	okField        = "Ok"
)
