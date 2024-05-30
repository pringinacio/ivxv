package rpc

import (
	"strconv"

	"ivxv.ee/common/collector/command"
	"ivxv.ee/common/collector/server"
	"ivxv.ee/common/collector/status/client"
	status "ivxv.ee/common/collector/status/client/rpc"
	api "ivxv.ee/sessionstatus/api/rpc"
)

const (
	// This should be a StatusReadResp.Caller value when calling RPC.Challenge
	Empty = ""

	// This should be a StatusReadResp.Caller value when calling RPC.Token
	Challenge = "RPC.Challenge"

	// This should be a StatusReadResp.Caller value when calling RPC.VoterChoices
	Token = "RPC.Token"
)

const exitCodeOK = 0

type RPC struct {
	authTTL int64
	client  client.TLSDialer
}

// NewClient initializes session status server client.
func NewClient(c *command.C) (client.Verifier, int) {
	// Initialize RPC TLS session status client
	tlsDialer, errCode := api.NewClient(c)
	if errCode != exitCodeOK {
		return nil, errCode
	}

	return &RPC{
		client:  tlsDialer,
		authTTL: c.Conf.Technical.Status.Session.AuthTTL,
	}, exitCodeOK
}

func (r *RPC) Verify(dto interface{}) (bool, error) {
	// dto should cast to *status.VerifyReq
	verifyReq, err := status.CastAnyToVerifyReq(dto)
	if err != nil {
		return false, CastAnyToVerifyReqError{Err: err}
	}

	// verifyReq.Request should cast to server.Header
	header, err := api.CastVerifyRequestToServerHeader(verifyReq)
	if err != nil {
		return false, CastVerifyRequestToServerHeaderError{Err: err}
	}

	// Send request to session status server and verify response
	ok, err := r.verifyAndUpdateSessionStatus(verifyReq.ServiceMethod, *header)
	if err != nil {
		return false, VerifyAndUpdateSessionStatusError{Err: err}
	}

	return ok, nil
}

// verifyAndUpdateSessionStatus will first check h.Header.SessionID
// record against the underlying storage, and if everything is correct,
// then will update h.Header.SessionID record in the underlying storage
// by marking session status Caller as serviceMethod.
//
// Note, that here serviceMethod is the RPC method that calls this function.
func (r *RPC) verifyAndUpdateSessionStatus(serviceMethod string, h server.Header) (bool, error) {
	// Create new session read status request
	reqRead := api.NewSessionStatusReadReqBuilder().
		WithHeader(h).
		Build()

	// Create new RPC request to status server, embeds session status request
	reqReadRPC := status.NewStatusReqBuilder().
		WithServiceMethod(api.Endpoint.SessionStatusRead).
		WithRequest(reqRead).
		Build()

	// RPC call to .WithServiceMethod(...)
	respReadRaw, err := r.client.TLSDial(&reqReadRPC)
	if err != nil {
		return false, SessionReadReqTLSDialError{Err: err}
	}

	// Process raw RPC response, doesn't care about the embedded status type
	respReadRPC := status.NewStatusRespBuilder().
		WithResponse(respReadRaw).
		Build()

	// Process session read status response
	respRead := api.NewSessionStatusReadRespBuilder().
		WithResponse(respReadRPC.Response).
		Build()

	// NB! Most important part, that prevents any attacks on SessionID
	var ok bool
	var ttl string
	switch serviceMethod {
	case Challenge:
		ok, err = verifyStatusReadResp(&respRead, challengeHandler)
		ttl = strconv.FormatInt(r.authTTL, 10)
	case Token:
		ok, err = verifyStatusReadResp(&respRead, tokenHandler)
		ttl = strconv.FormatInt(r.authTTL, 10)
	}
	if !ok || err != nil {
		return ok, VerifyStatusReadRespError{Err: err}
	}

	// Create new session update status request
	reqUpdate := api.NewSessionStatusUpdateReqBuilder().
		WithHeader(h).
		WithCaller(serviceMethod).
		WithAuth(respRead.Auth).
		WithLease(respRead.Lease).
		WithTTL(ttl).
		Build()

	// Create new RPC request to status server, embeds session status request
	reqUpdateRPC := status.NewStatusReqBuilder().
		WithServiceMethod(api.Endpoint.SessionStatusUpdate).
		WithRequest(reqUpdate).
		Build()

	// RPC call to .WithServiceMethod(...)
	respUpdateRaw, err := r.client.TLSDial(&reqUpdateRPC)
	if err != nil {
		return false, SessionUpdateReqTLSDialError{Err: err}
	}

	// Process raw RPC response, doesn't care about the embedded status type
	respUpdateRPC := status.NewStatusRespBuilder().
		WithResponse(respUpdateRaw).
		Build()

	// Process session update status response
	respUpdate := api.NewSessionStatusUpdateRespBuilder().
		WithResponse(respUpdateRPC.Response).
		Build()

	// If true, then status has been successfully updated
	ok = respUpdate.Ok
	if !ok {
		return false, SessionStatusUpdateError{
			Caller: reqUpdate.Caller,
			Auth:   respRead.Auth,
		}
	}

	return true, nil
}

// verifyStatusReadResp r by applying an appropriate handler h.
func verifyStatusReadResp(r *api.StatusReadResp,
	h func(*api.StatusReadResp) (bool, error)) (bool, error) {
	return h(r)
}

// challengeHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.Challenge request.
func challengeHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.Challenge is the very first client request to IVXV,
	// so IVXV requires no previous interactions
	firstTime := r.Caller == Empty && r.Auth == client.NoAuth

	if !(firstTime) {
		return false, ChallengeInvalidCallerOrAuthForSessionID{
			Method: Challenge,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}

	r.Auth = client.WebeIDAuth
	return true, nil
}

// tokenHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.Token request.
func tokenHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.Token is the second client request to IVXV,
	// so IVXV requires RPC.Challenge previously interacted
	secondTime := r.Caller == Challenge && r.Auth == client.WebeIDAuth

	if !(secondTime) {
		return false, TokenInvalidCallerOrAuthForSessionID{
			Method: Token,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}
