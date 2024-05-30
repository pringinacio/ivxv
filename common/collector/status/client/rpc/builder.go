package rpc

type StatusReqBuilder struct {
	serviceMethod string
	request       any
}

// NewStatusReqBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusReq.
func NewStatusReqBuilder() *StatusReqBuilder {
	return new(StatusReqBuilder)
}

func (srb *StatusReqBuilder) WithServiceMethod(s string) *StatusReqBuilder {
	srb.serviceMethod = s
	return srb
}

func (srb *StatusReqBuilder) WithRequest(r any) *StatusReqBuilder {
	srb.request = r
	return srb
}

// Build returns a StatusReq.
func (srb *StatusReqBuilder) Build() StatusReq {
	return StatusReq{
		ServiceMethod: srb.serviceMethod,
		Request:       srb.request,
	}
}

type StatusRespBuilder struct {
	response map[string]any
}

// NewStatusRespBuilder is a Builder-pattern constructor, which is used
// to prepare a StatusResp.
func NewStatusRespBuilder() *StatusRespBuilder {
	return new(StatusRespBuilder)
}

func (srb *StatusRespBuilder) WithResponse(r any) *StatusRespBuilder {
	resp, ok := r.(*StatusResp)
	if ok {
		// If r is a *StatusResp, then add a response, otherwise keep it nil
		srb.response = resp.Response
	}
	return srb
}

// Build returns a StatusResp.
func (srb *StatusRespBuilder) Build() StatusResp {
	return StatusResp{
		Response: srb.response,
	}
}

type VerifyReqBuilder struct {
	serviceMethod string
	request       any
}

// NewVerifyReqBuilder is a Builder-pattern constructor, which is used
// to prepare a VerifyReq.
func NewVerifyReqBuilder() *VerifyReqBuilder {
	return new(VerifyReqBuilder)
}

func (vrb *VerifyReqBuilder) WithServiceMethod(s string) *VerifyReqBuilder {
	vrb.serviceMethod = s
	return vrb
}

func (vrb *VerifyReqBuilder) WithRequest(r any) *VerifyReqBuilder {
	vrb.request = r
	return vrb
}

// Build returns a VerifyReq.
func (vrb *VerifyReqBuilder) Build() VerifyReq {
	return VerifyReq{
		ServiceMethod: vrb.serviceMethod,
		Request:       vrb.request,
	}
}
