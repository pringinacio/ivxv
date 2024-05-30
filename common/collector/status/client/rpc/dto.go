package rpc

// VerifyReq is a DTO to pass as a req to a Verify() interface method.
type VerifyReq struct {
	// ServiceMethod is an RPC endpoint of a status server, e.g. RPC.SessionStatusRead.
	ServiceMethod string

	// Request is a data to pass to the RPC endpoint.
	Request any
}

// StatusReq is a DTO to pass as a req to a TLSDial() interface method.
type StatusReq struct {
	// ServiceMethod is an RPC endpoint of a status server, e.g. RPC.SessionStatusRead.
	ServiceMethod string

	// Request is a data to pass to the RPC endpoint.
	Request any
}

// StatusResp is a DTO that TLSDial() interface method returns as a resp.
type StatusResp struct {
	// Response is a data that RPC endpoint returns to a caller.
	Response map[string]any
}
