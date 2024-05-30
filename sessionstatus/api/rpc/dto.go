package rpc

import "ivxv.ee/common/collector/server"

// StatusReadReq is a data to pass to RPC.SessionStatusRead.
type StatusReadReq struct {
	server.Header
}

// StatusReadResp is a data that RPC.SessionStatusRead responds back.
type StatusReadResp struct {
	server.Header

	// Caller is an RPC method which was previously called.
	Caller string

	// Auth is an authentication method being used.
	Auth string

	// Lease is an optional field that indicates an ID of TTL timer.
	// It may present or may not, depending on implementation.
	Lease string
}

// StatusUpdateReq is a data to pass to RPC.SessionStatusUpdate.
//
// For each field description see StatusReadResp.
type StatusUpdateReq struct {
	server.Header
	// Caller is an RPC method which was currently called.
	Caller string

	// Auth is an authentication method being used.
	Auth string

	// Lease is an optional field that indicates an ID of TTL timer.
	Lease string

	// Lease is an optional field that indicates a TTL in seconds.
	TTL string
}

// StatusUpdateResp is a data that RPC.SessionStatusUpdate responds back.
type StatusUpdateResp struct {
	server.Header

	// Ok is true if updating was successful.
	Ok bool
}

// StatusDeleteReq is a data to pass to RPC.SessionStatusDelete.
type StatusDeleteReq struct {
	server.Header
}

// StatusDeleteResp is a data that RPC.SessionStatusDelete responds back.
type StatusDeleteResp struct {
	server.Header

	// Ok is true if updating was successful.
	Ok bool
}
