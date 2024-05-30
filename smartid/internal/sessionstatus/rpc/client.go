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
	// This should be a StatusReadResp.Caller value when calling RPC.Authenticate
	Empty = ""

	// This should be a StatusReadResp.Caller value when calling RPC.AuthenticateStatus
	Authenticate = "RPC.Authenticate"

	// This should be a StatusReadResp.Caller value when calling RPC.VoterChoices
	AuthenticateStatus = "RPC.AuthenticateStatus"

	// This should be a StatusReadResp.Caller value when calling RPC.GetCertificate
	VoterChoices = "RPC.VoterChoices"

	// This should be a StatusReadResp.Caller value when calling RPC.GetCertificateStatus
	GetCertificate = "RPC.GetCertificate"

	// This should be a StatusReadResp.Caller value when calling RPC.Sign
	GetCertificateStatus = "RPC.GetCertificateStatus"

	// This should be a StatusReadResp.Caller value when calling RPC.SignStatus
	Sign = "RPC.Sign"

	// This should be a StatusReadResp.Caller value when calling RPC.Vote
	SignStatus = "RPC.SignStatus"
)

const exitCodeOK = 0

type RPC struct {
	authTTL int64
	voteTTL int64
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
		voteTTL: c.Conf.Technical.Status.Session.VoteTTL,
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

	// NB! Most important part, that prevents any attack on SessionID
	var ok bool
	var ttl string
	switch serviceMethod {
	case Authenticate:
		ok, err = verifyStatusReadResp(&respRead, authenticateHandler)
		ttl = strconv.FormatInt(r.authTTL, 10)
	case AuthenticateStatus:
		ok, err = verifyStatusReadResp(&respRead, authenticateStatusHandler)
		ttl = strconv.FormatInt(r.authTTL, 10)
	case GetCertificate:
		ok, err = verifyStatusReadResp(&respRead, getCertificateHandler)
		ttl = strconv.FormatInt(r.voteTTL, 10)
		respRead.Lease = ""
	case GetCertificateStatus:
		ok, err = verifyStatusReadResp(&respRead, getCertificateStatusHandler)
		ttl = strconv.FormatInt(r.voteTTL, 10)
	case Sign:
		ok, err = verifyStatusReadResp(&respRead, signHandler)
		ttl = strconv.FormatInt(r.voteTTL, 10)
	case SignStatus:
		ok, err = verifyStatusReadResp(&respRead, signStatusHandler)
		ttl = strconv.FormatInt(r.voteTTL, 10)
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

// authenticateHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.Authenticate request.
func authenticateHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.Authenticate is the very first client request to IVXV,
	// so IVXV requires no previous interactions
	firstTime := r.Caller == Empty && r.Auth == client.NoAuth

	if !(firstTime) {
		return false, AuthenticateInvalidCallerOrAuthForSessionID{
			Method: Authenticate,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}

	r.Auth = client.SmartIDAuth
	return true, nil
}

// authenticateStatusHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.AuthenticateStatus request.
func authenticateStatusHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.AuthenticateStatus is the second client request to IVXV,
	// so IVXV requires RPC.Authenticate previously interacted
	secondTime := r.Caller == Authenticate && r.Auth == client.SmartIDAuth

	// If RPC.AuthenticateStatus is still processing Smart-ID request
	// then voting app will still send RPC.AuthenticateStatus request
	// to finish Smart-ID authentication. So client can send
	// RPC.AuthenticateStatus queries as many as wants
	nTime := r.Caller == AuthenticateStatus && r.Auth == client.SmartIDAuth

	if !(secondTime) && !(nTime) {
		return false, AuthenticateStatusInvalidCallerOrAuthForSessionID{
			Method: AuthenticateStatus,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}

// getCertificateHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.GetCertificate request.
func getCertificateHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.GetCertificate is the fourth client request to IVXV,
	// third interaction should be done with RPC.VoterChoices
	fourthTime := r.Caller == VoterChoices && r.Auth == client.SmartIDAuth

	if !(fourthTime) {
		return false, GetCertificateInvalidCallerOrAuthForSessionID{
			Method: GetCertificate,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}

// getCertificateStatusHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.GetCertificateStatus request.
func getCertificateStatusHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.GetCertificateStatus is the fifth client request to IVXV,
	// fourth interaction should be done with RPC.GetCertificate
	fifthTime := r.Caller == GetCertificate && r.Auth == client.SmartIDAuth

	// If RPC.GetCertificateStatus is still processing Smart-ID request
	// then voting app will still send RPC.GetCertificateStatus request
	// to receive Smart-ID signing certificate. So client can send
	// RPC.GetCertificateStatus queries as many as wants
	nTime := r.Caller == GetCertificateStatus && r.Auth == client.SmartIDAuth

	if !(fifthTime) && !(nTime) {
		return false, GetCertificateStatusInvalidCallerOrAuthForSessionID{
			Method: GetCertificateStatus,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}

// signHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.Sign request.
func signHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.Sign is the sixth client request to IVXV,
	// so IVXV requires RPC.GetCertificateStatus previously interacted
	sixthTime := r.Caller == GetCertificateStatus && r.Auth == client.SmartIDAuth

	if !(sixthTime) {
		return false, SignInvalidCallerOrAuthForSessionID{
			Method: Sign,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}

// signStatusHandler performs filter operation on StatusReadResp r to
// detect invalid SessionID in a client RPC.SignStatus request.
func signStatusHandler(r *api.StatusReadResp) (bool, error) {
	// RPC.Sign is the seventh client request to IVXV,
	// so IVXV requires RPC.Sign previously interacted
	seventhTime := r.Caller == Sign && r.Auth == client.SmartIDAuth

	// If RPC.SignStatus is still processing Smart-ID request
	// then voting app will still send RPC.SignStatus request
	// to check whether document is signed with Smart-ID. So
	// client can send RPC.SignStatus queries as many as wants
	nTime := r.Caller == SignStatus && r.Auth == client.SmartIDAuth

	if !(seventhTime) && !(nTime) {
		return false, SignStatusInvalidCallerOrAuthForSessionID{
			Method: SignStatus,
			Caller: r.Caller,
			Auth:   r.Auth,
		}
	}
	return true, nil
}
