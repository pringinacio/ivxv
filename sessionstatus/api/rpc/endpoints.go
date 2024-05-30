package rpc

// Endpoint is a collection of available session status server RPC endpoints.
var Endpoint = struct {
	SessionStatusRead   string
	SessionStatusUpdate string
	SessionStatusDelete string
}{
	"RPC.SessionStatusRead",
	"RPC.SessionStatusUpdate",
	"RPC.SessionStatusDelete",
}
