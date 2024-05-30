package rpc

import "ivxv.ee/common/collector/server"

type StatusReadReqBuilder struct {
	header server.Header
}

// NewSessionStatusReadReqBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusReadReq.
func NewSessionStatusReadReqBuilder() *StatusReadReqBuilder {
	return &StatusReadReqBuilder{}
}

func (srrb *StatusReadReqBuilder) WithHeader(h server.Header) *StatusReadReqBuilder {
	srrb.header = h
	return srrb
}

// Build returns StatusReadReq.
func (srrb *StatusReadReqBuilder) Build() StatusReadReq {
	return StatusReadReq{Header: srrb.header}
}

type StatusReadRespBuilder struct {
	response map[string]any
}

// NewSessionStatusReadRespBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusReadResp.
func NewSessionStatusReadRespBuilder() *StatusReadRespBuilder {
	return &StatusReadRespBuilder{}
}

func (srrb *StatusReadRespBuilder) WithResponse(r map[string]any) *StatusReadRespBuilder {
	srrb.response = r
	return srrb
}

// Build returns StatusReadResp.
func (srrb *StatusReadRespBuilder) Build() StatusReadResp {
	return StatusReadResp{
		Header: server.Header{
			SessionID: srrb.response[sessionIDField].(string),
		},
		Caller: srrb.response[callerField].(string),
		Auth:   srrb.response[authField].(string),
		Lease:  srrb.response[leaseField].(string),
	}
}

type StatusUpdateReqBuilder struct {
	header server.Header
	caller string
	auth   string
	lease  string
	ttl    string
}

// NewSessionStatusUpdateReqBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusUpdateReq.
func NewSessionStatusUpdateReqBuilder() *StatusUpdateReqBuilder {
	return &StatusUpdateReqBuilder{}
}

func (surb *StatusUpdateReqBuilder) WithHeader(h server.Header) *StatusUpdateReqBuilder {
	surb.header = h
	return surb
}

func (surb *StatusUpdateReqBuilder) WithCaller(c string) *StatusUpdateReqBuilder {
	surb.caller = c
	return surb
}

func (surb *StatusUpdateReqBuilder) WithAuth(a string) *StatusUpdateReqBuilder {
	surb.auth = a
	return surb
}

func (surb *StatusUpdateReqBuilder) WithLease(l string) *StatusUpdateReqBuilder {
	surb.lease = l
	return surb
}

func (surb *StatusUpdateReqBuilder) WithTTL(t string) *StatusUpdateReqBuilder {
	surb.ttl = t
	return surb
}

// Build returns StatusUpdateReq.
func (surb *StatusUpdateReqBuilder) Build() StatusUpdateReq {
	return StatusUpdateReq{
		Header: surb.header,
		Caller: surb.caller,
		Auth:   surb.auth,
		Lease:  surb.lease,
		TTL:    surb.ttl,
	}
}

type StatusUpdateRespBuilder struct {
	response map[string]any
}

// NewSessionStatusUpdateRespBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusUpdateResp.
func NewSessionStatusUpdateRespBuilder() *StatusUpdateRespBuilder {
	return &StatusUpdateRespBuilder{}
}

func (surb *StatusUpdateRespBuilder) WithResponse(r map[string]any) *StatusUpdateRespBuilder {
	surb.response = r
	return surb
}

// Build returns StatusUpdateResp.
func (surb *StatusUpdateRespBuilder) Build() StatusUpdateResp {
	return StatusUpdateResp{
		Ok: surb.response[okField].(bool),
	}
}

type StatusDeleteReqBuilder struct {
	header server.Header
}

// NewSessionStatusDeleteReqBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusDeleteReq.
func NewSessionStatusDeleteReqBuilder() *StatusDeleteReqBuilder {
	return new(StatusDeleteReqBuilder)
}

func (sdrb *StatusDeleteReqBuilder) WithHeader(h server.Header) *StatusDeleteReqBuilder {
	sdrb.header = h
	return sdrb
}

// Build returns StatusDeleteReq.
func (sdrb *StatusDeleteReqBuilder) Build() StatusDeleteReq {
	return StatusDeleteReq{Header: sdrb.header}
}

type StatusDeleteRespBuilder struct {
	response map[string]any
}

// NewSessionStatusDeleteRespBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusDeleteResp.
func NewSessionStatusDeleteRespBuilder() *StatusDeleteRespBuilder {
	return new(StatusDeleteRespBuilder)
}

func (sdrb *StatusDeleteRespBuilder) WithResponse(r map[string]any) *StatusDeleteRespBuilder {
	sdrb.response = r
	return sdrb
}

// Build returns StatusDeleteResp.
func (sdrb *StatusDeleteRespBuilder) Build() StatusDeleteResp {
	return StatusDeleteResp{
		Ok: sdrb.response[okField].(bool),
	}
}
