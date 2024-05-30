package rpc

import (
	"context"
	"reflect"
	"testing"

	"ivxv.ee/common/collector/server"
	status "ivxv.ee/common/collector/status/client"
)

// If any new RPC endpoint is about to appear in a IVXV online workflow
// then it should be added here, in case you wish that endpoint to
// receive status checks.
const (
	// Empty RPC
	Empty = ""

	// Mobile-ID or Smart-ID authentication
	Authenticate       = "RPC.Authenticate"
	AuthenticateStatus = "RPC.AuthenticateStatus"

	// Web eID authentication
	Challenge = "RPC.Challenge"
	Token     = "RPC.Token"

	// Choices list
	VoterChoices = "RPC.VoterChoices"

	// Mobile-ID or Smart-ID signing certificate REST query
	GetCertificate = "RPC.GetCertificate"
	// Smart-ID signing certificate received status check
	GetCertificateStatus = "RPC.GetCertificateStatus"

	// Mobile-ID or Smart-ID signing
	Sign       = "RPC.Sign"
	SignStatus = "RPC.SignStatus"

	// Voting, i.e. storing a vote in a database
	Vote = "RPC.Vote"
)

var commonErrTxt3 = "Two structs are not equal. Expected %v, got %v\n"

// Example of StatusReadReq
var headers = []server.Header{
	// Example of any RPC request Header against IVXV backend. That kind of
	// Header will see any RPC method
	{
		Ctx:        context.Background(),
		SessionID:  "0101e9342abab1577b8b2844d6a1d317",
		OS:         "Ubuntu Jammy 22.04 LTS",
		AuthMethod: "",
		AuthToken:  nil,
		DataToken:  nil,
	},
	// Example of SessionID tamper attack to bypass backend SessionID generation.
	// Can be bypassed if SessionID is len(SessionID) != 0 and valid hex
	{
		Ctx:        context.Background(),
		SessionID:  "0",
		OS:         "Ubuntu Jammy 22.04 LTS",
		AuthMethod: "",
		AuthToken:  nil,
		DataToken:  nil,
	},
	// Example of OS JavaScript injection attack which could impact IVXV
	// logmonitor, that will display results on Web
	{
		Ctx:        context.Background(),
		SessionID:  "",
		OS:         `<script>alert("Hello World")</script>`,
		AuthMethod: "",
		AuthToken:  nil,
		DataToken:  nil,
	},
}

var sessionStatusReadResps = []map[string]any{
	// Example of RPC.Authenticate request with Mobile-ID
	// For given SessionID there SHOULD NOT be any session status
	// records in a database
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    Empty,
		"Auth":      status.NoAuth,
		"Lease":     "0",
	},
	// Example of RPC.AuthenticateStatus request with Mobile-ID
	// This time we have updated the session status record in a database
	// and Caller is the previous RPC method that has been invoked and
	// Auth is determined to be "mid". Lease ID also has a value
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    Authenticate,
		"Auth":      status.MobileIDAuth,
		// int64 = 3367438597345620447
		"Lease": "2ebb8b9015b9b1df",
	},
	// Example of Web eID RPC.Token request
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    Challenge,
		"Auth":      status.WebeIDAuth,
		// int64 = 1936286082718970961
		"Lease": "1adf11eeefff3451",
	},
	// Example of ID-card RPC.VoterChoices. In case of ID card it is a
	// first interaction with session status database
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    Empty,
		"Auth":      status.NoAuth,
		"Lease":     "",
	},
	// Example of ID-card RPC.Vote
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    VoterChoices,
		"Auth":      status.IDcardAuth,
		// int64 = 768699984047453265 aaaf8900fff3451
		"Lease": "aaaf8900fff3451",
	},
	// Example of Web eID RPC.Challenge
	{
		"SessionID": "0101e9342abab1577b8b2844d6a1d317",
		"Caller":    status.NoAuth,
		"Auth":      status.WebeIDAuth,
		"Lease":     "",
	},
}

var sessionStatusUpdateReqs = []map[string]any{
	// Example of RPC.Authenticate Update request with Mobile-ID
	{
		"Header": headers[0],
		"Caller": Authenticate,
		"Auth":   status.MobileIDAuth,
		"Lease":  "0",
	},
	// Example of RPC.AuthenticateStatus request with Mobile-ID
	{
		"Header": headers[0],
		"Caller": AuthenticateStatus,
		"Auth":   status.MobileIDAuth,
		// int64 = 3367438597345620447
		"Lease": "2ebb8b9015b9b1df",
	},
	// Example of Web eID RPC.Token request
	{
		"Header": headers[0],
		"Caller": Token,
		"Auth":   status.WebeIDAuth,
		// int64 = 1936286082718970961
		"Lease": "1adf11eeefff3451",
	},
	// Example of ID-card RPC.VoterChoices
	{
		"Header": headers[0],
		"Caller": VoterChoices,
		"Auth":   status.IDcardAuth,
		"Lease":  "0",
	},
	// Example of ID-card RPC.Vote
	{
		"Header": headers[0],
		"Caller": Vote,
		"Auth":   status.IDcardAuth,
		// int64 = 768699984047453265
		"Lease": "aaaf8900fff34511",
	},
	// Example of Web eID RPC.Challenge
	{
		"Header": headers[0],
		"Caller": Challenge,
		"Auth":   status.WebeIDAuth,
		"Lease":  "0",
	},
}

var sessionStatusUpdateResps = []map[string]any{
	{
		"Header": headers[0],
		"Ok":     true,
	},
	{
		"Header": headers[0],
		"Ok":     false,
	},
}

// Responses for Update and Delete are literally same
var sessionStatusDeleteResps = sessionStatusUpdateResps

// Proves that SessionStatusReadReqBuilder can produce correct StatusReadReq.
func TestSessionStatusReadReqBuilder(t *testing.T) {
	for _, h := range headers {
		expected := StatusReadReq{Header: h}

		got := NewSessionStatusReadReqBuilder().
			WithHeader(h).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}

// Proves that SessionStatusReadRespBuilder can produce correct StatusReadResp.
func TestSessionStatusReadRespBuilder(t *testing.T) {
	for _, sessionStatusReadResp := range sessionStatusReadResps {
		expected := StatusReadResp{
			Header: server.Header{
				SessionID: sessionStatusReadResp["SessionID"].(string),
			},
			Caller: sessionStatusReadResp["Caller"].(string),
			Auth:   sessionStatusReadResp["Auth"].(string),
			Lease:  sessionStatusReadResp["Lease"].(string),
		}

		got := NewSessionStatusReadRespBuilder().
			WithResponse(sessionStatusReadResp).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}

// Proves that SessionStatusUpdateReqBuilder can produce correct StatusUpdateReq.
func TestSessionStatusUpdateReqBuilder(t *testing.T) {
	for _, sessionStatusUpdateReq := range sessionStatusUpdateReqs {
		expected := StatusUpdateReq{
			Header: sessionStatusUpdateReq["Header"].(server.Header),
			Caller: sessionStatusUpdateReq["Caller"].(string),
			Auth:   sessionStatusUpdateReq["Auth"].(string),
			Lease:  sessionStatusUpdateReq["Lease"].(string),
		}

		got := NewSessionStatusUpdateReqBuilder().
			WithHeader(sessionStatusUpdateReq["Header"].(server.Header)).
			WithCaller(sessionStatusUpdateReq["Caller"].(string)).
			WithAuth(sessionStatusUpdateReq["Auth"].(string)).
			WithLease(sessionStatusUpdateReq["Lease"].(string)).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}

// Proves that SessionStatusUpdateRespBuilder can produce correct StatusUpdateResp.
func TestNewSessionStatusUpdateRespBuilder(t *testing.T) {
	for _, sessionStatusUpdateResp := range sessionStatusUpdateResps {
		expected := StatusUpdateResp{
			Ok: sessionStatusUpdateResp["Ok"].(bool),
		}

		got := NewSessionStatusUpdateRespBuilder().
			WithResponse(sessionStatusUpdateResp).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}

// Proves that StatusDeleteReqBuilder can produce correct StatusDeleteReq.
func TestSessionStatusDeleteReqBuilder(t *testing.T) {
	for _, h := range headers {
		expected := StatusDeleteReq{Header: h}

		got := NewSessionStatusDeleteReqBuilder().
			WithHeader(h).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}

// Proves that SessionStatusDeleteRespBuilder can produce correct StatusDeleteResp.
func TestNewSessionStatusDeleteRespBuilder(t *testing.T) {
	for _, sessionStatusDeleteResp := range sessionStatusDeleteResps {
		expected := StatusDeleteResp{
			Ok: sessionStatusDeleteResp["Ok"].(bool),
		}

		got := NewSessionStatusDeleteRespBuilder().
			WithResponse(sessionStatusDeleteResp).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(commonErrTxt3, expected, got)
		}
	}
}
