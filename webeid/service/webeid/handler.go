package main

import (
	"crypto/x509"
	"time"

	"ivxv.ee/common/collector/auth"
	"ivxv.ee/common/collector/cryptoutil"
	"ivxv.ee/common/collector/log"
	"ivxv.ee/common/collector/server"
	status "ivxv.ee/common/collector/status/client/rpc"
	"ivxv.ee/common/collector/token"
	"ivxv.ee/common/collector/token/custom"
	"ivxv.ee/common/collector/token/webeid"
	internal "ivxv.ee/webeid/internal/sessionstatus/rpc"
)

// ChallengeReq is a client RPC call to retrieve a cryptographic nonce.
type ChallengeReq struct {
	server.Header
}

// ChallengeResp is a server RPC response to the ChallengeReq.
type ChallengeResp struct {
	server.Header
	Challenge []byte
	Bearer    string
}

// Challenge is an RPC endpoint to provide a cryptographic nonce.
func (r *RPC) Challenge(req ChallengeReq, resp *ChallengeResp) (err error) {
	log.Log(req.Ctx, ChallengeRequest{})

	// Check that election period is not ended
	if !time.Now().Before(r.authEnd) {
		log.Log(req.Ctx, ChallengeVotingEnded{})
		return server.ErrVotingEnd
	}

	// Build up VerifyReq for session status service
	verifyReq := status.NewVerifyReqBuilder().
		WithServiceMethod(internal.Challenge).
		WithRequest(req.Header).
		Build()

	// SessionID security check
	ok, err := r.status.Verify(&verifyReq)
	if err != nil {
		log.Error(req.Ctx, ChallengeVerifySessionIDError{Err: err})
		return server.ErrBadRequest
	}
	if !ok {
		log.Error(req.Ctx, ChallengeUpdateSessionIDError{})
		return server.ErrBadRequest
	}

	// Generate 44 byte nonce
	nonce, err := cryptoutil.Nonce44Bytes()
	if err != nil {
		log.Error(req.Ctx, ChallengeGenerateNonceError{Err: err})
		return server.ErrInternal
	}

	// Generate brand-new bearer token
	bearerToken := custom.NewFromEmptyBuilder().
		WithNonce(nonce).
		WithTimeStamp(time.Now()).
		Build()

	// Payload is a serialized bearerToken
	payload, err := bearerToken.Payload()
	if err != nil {
		log.Error(req.Ctx, ChallengeExtractPayloadFromBearerError{Err: err})
		return server.ErrInternal
	}

	// Create a signature over a payload using shared secret
	sig := r.cookie.Create([]byte(payload))

	// Create new bearer token with a payload and a signature,
	// and serialize it to base64 string
	rawComplete, err := custom.NewFromExistingBuilder().
		WithPayload(payload).
		WithSignature(string(sig)).
		Build().
		Marshal()
	if err != nil {
		log.Error(req.Ctx, ChallengeMarshalBearerError{Err: err})
		return server.ErrInternal
	}

	// Return raw bearer token to client
	resp.Bearer = rawComplete

	// Return nonce to client
	resp.Challenge = []byte(nonce)

	log.Log(req.Ctx, ChallengeResponse{
		Challenge: resp.Challenge,
		Bearer:    resp.Bearer,
	})

	return nil
}

// TokenReq is a client RPC call to validate a Web eID authentication token.
type TokenReq struct {
	server.Header
	// Token has a maximum size of 16000 bytes.
	// Token size is not defined by a backend and is completely up to the user.
	Token string `size:"16000"`
	// Bearer has a maximum size of 400 bytes.
	// Bearer size is fully defined by a backend.
	Bearer string `size:"400"`
}

// TokenResp is a server RPC response to the TokenReq.
type TokenResp struct {
	server.Header
	Status       string
	GivenName    string
	Surname      string
	PersonalCode string
	AuthToken    []byte
}

// Token is an RPC endpoint to validate a Web eID authentication token.
func (r *RPC) Token(req TokenReq, resp *TokenResp) (err error) {
	log.Log(req.Ctx, TokenRequest{})

	// Check that election period is not ended
	if !time.Now().Before(r.authEnd) {
		log.Log(req.Ctx, TokenVotingEnded{})
		return server.ErrVotingEnd
	}

	// Build up VerifyReq for session status service
	verifyReq := status.NewVerifyReqBuilder().
		WithServiceMethod(internal.Token).
		WithRequest(req.Header).
		Build()

	// SessionID security check
	ok, err := r.status.Verify(&verifyReq)
	if err != nil {
		log.Error(req.Ctx, TokenVerifySessionIDError{Err: err})
		return server.ErrBadRequest
	}
	if !ok {
		log.Error(req.Ctx, TokenUpdateSessionIDError{})
		return server.ErrBadRequest
	}

	// Generate bearer token from a raw bearer token (base64 string)
	bearerToken, err := custom.NewFromRawBuilder().
		WithBearer(req.Bearer).
		Build().
		Unmarshal()
	if err != nil {
		log.Error(req.Ctx, TokenUnmarshalBearerError{Err: err})
		return server.ErrBadRequest
	}

	// Get bearer token signature
	sig, err := bearerToken.Signature()
	if err != nil {
		log.Error(req.Ctx, TokenExtractSignatureFromBearerError{Err: err})
		return server.ErrBadRequest
	}

	// Verify bearer token signature
	data, err := r.cookie.Open(sig)
	if err != nil {
		log.Error(req.Ctx, TokenVerifyBearerSignatureError{Err: err})
		return server.ErrBadRequest
	}
	if data == nil {
		log.Error(req.Ctx, TokenBearerSignatureIsEmptyAfterDecryptError{Err: err})
		return server.ErrBadRequest
	}

	// Payload returns nonce
	nonce, err := bearerToken.Payload()
	if err != nil {
		log.Error(req.Ctx, TokenExtractPayloadFromBearerError{Err: err})
		return server.ErrBadRequest
	}

	// Generate Web eID token from a raw Web eID token (json string)
	wToken, err := webeid.NewFromRawBuilder().
		WithToken(req.Token).
		WithNonce(nonce).
		WithOrigin(r.origin).
		Build().
		Unmarshal()
	if err != nil {
		log.Error(req.Ctx, TokenUnmarshalWebeidError{Err: err})
		return server.ErrBadRequest
	}

	// Parse auth cert from Web eID token
	cert := wToken.(token.Certifier).Certify()

	// Add auth cert to the req.Ctx context, so that auth.tls verifier
	// could perform TLS certificate validation
	req.Ctx = server.TLSClientKey(req.Ctx, []*x509.Certificate{cert})

	// Validate auth cert from the req.Ctx and perform OCSP check
	authCert, _, err := r.auther.Auth.Verify(req.Ctx, auth.TLS, nil)
	if err != nil {
		log.Error(req.Ctx, TokenVerifyWebeidAuthCertCAnOCSPError{Err: err})
		return server.ErrCertificate
	}

	// Verify Web eID auth token signature
	err = wToken.Verify()
	if err != nil {
		log.Error(req.Ctx, TokenVerifyWebeidError{Err: err})
		return server.ErrBadRequest
	}

	// Retrieve specific data from client certificate from Web eID auth token
	resp.Status = StatusOK
	resp.GivenName = findName(authCert, oid["givenName"])
	resp.Surname = findName(authCert, oid["surname"])
	resp.PersonalCode = personalCode(authCert)

	// Create client authentication cookie (ticket)
	if resp.AuthToken, err = r.ticket.Create(cert.Subject); err != nil {
		log.Error(req.Ctx, TokenAuthenticationTicketError{Err: err})
		return server.ErrInternal
	}

	log.Log(req.Ctx, TokenResponse{
		Status:       resp.Status,
		GivenName:    resp.GivenName,
		Surname:      resp.Surname,
		PersonalCode: resp.PersonalCode,
		AuthToken:    log.Sensitive(resp.AuthToken),
	})

	return nil
}
