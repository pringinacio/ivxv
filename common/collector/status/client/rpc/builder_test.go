package rpc

import (
	"context"
	"reflect"
	"testing"

	"ivxv.ee/common/collector/server"
	status "ivxv.ee/common/collector/status/client"
)

const (
	msgTwoStructsNotEqual   = "Two structs are not equal. Expected %v, got %v\n"
	someServiceMethod       = "RPC.SomeServiceMethod"
	someCallerServiceMethod = "RPC.SomeCallerServiceMethod"
)

var header = server.Header{
	Ctx:        context.Background(),
	SessionID:  "0101e9342abab1577b8b2844d6a1d317",
	OS:         "Ubuntu Jammy 22.04 LTS",
	AuthMethod: "",
	AuthToken:  nil,
	DataToken:  nil,
}

var statusReqs = []map[string]any{
	{
		"ServiceMethod": someServiceMethod,
		"Request": struct {
			Header server.Header
		}{
			Header: header,
		},
	},
	{
		"ServiceMethod": someServiceMethod,
		"Request": struct {
			Header server.Header
			Caller string
			Auth   string
			Lease  string
			TTL    string
		}{
			Header: header,
			Caller: someCallerServiceMethod,
			Auth:   status.NoAuth,
			// int64 = 0
			Lease: "0",
			TTL:   "0",
		},
	},
	{
		"ServiceMethod": someServiceMethod,
		"Request": struct {
			Header server.Header
			Caller string
			Auth   string
			Lease  string
			TTL    string
		}{
			Header: header,
			Caller: someCallerServiceMethod,
			Auth:   status.SmartIDAuth,
			// int64 = 768699984047453265
			Lease: "aaaf8900fff3451",
			TTL:   "200",
		},
	},
	{
		"ServiceMethod": someServiceMethod,
		"Request":       nil,
	},
}

var statusResps = []*StatusResp{
	{
		Response: map[string]any{
			"Header": header,
			"Caller": someCallerServiceMethod,
			"Auth":   status.WebeIDAuth,
			// int64 = 0
			"Lease": "",
			"TTL":   "",
		},
	},
	{
		Response: map[string]any{
			"Header": header,
			"Caller": someCallerServiceMethod,
			"Auth":   status.WebeIDAuth,
			// int64 = 768699984047453265
			"Lease": "aaaf8900fff3451",
			"TTL":   "120",
		},
	},
	{
		Response: map[string]any{
			"Header": header,
			"Ok":     true,
		},
	},
	{
		Response: map[string]any{
			"Header": header,
			"Ok":     false,
		},
	},
}

var verifyReqs = []map[string]any{
	{
		"ServiceMethod": someServiceMethod,
		"Request": struct {
			Header server.Header
		}{
			Header: header,
		},
	},
	{
		"ServiceMethod": someServiceMethod,
		"Request":       nil,
	},
}

func TestStatusReqBuilder(t *testing.T) {
	for _, statusReq := range statusReqs {
		expected := StatusReq{
			ServiceMethod: statusReq["ServiceMethod"].(string),
			Request:       statusReq["Request"],
		}

		got := NewStatusReqBuilder().
			WithServiceMethod(statusReq["ServiceMethod"].(string)).
			WithRequest(statusReq["Request"]).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(msgTwoStructsNotEqual, expected, got)
		}
	}
}

func TestStatusRespBuilder(t *testing.T) {
	for _, statusResp := range statusResps {
		expected := StatusResp{
			Response: statusResp.Response,
		}

		got := NewStatusRespBuilder().
			WithResponse(statusResp).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(msgTwoStructsNotEqual, expected, got)
		}
	}
}

func TestVerifyReqBuilder(t *testing.T) {
	for _, verifyReq := range verifyReqs {
		expected := VerifyReq{
			ServiceMethod: verifyReq["ServiceMethod"].(string),
			Request:       verifyReq["Request"],
		}

		got := NewVerifyReqBuilder().
			WithServiceMethod(verifyReq["ServiceMethod"].(string)).
			WithRequest(verifyReq["Request"]).
			Build()

		if !reflect.DeepEqual(expected, got) {
			t.Errorf(msgTwoStructsNotEqual, expected, got)
		}
	}
}
