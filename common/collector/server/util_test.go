package server

import (
	"testing"
)

var validOrigins = []string{
	"https://ehs-logimonitor-01.prod.riaint.ee:443",
	"https://ehs-ivxv-01.demo.riaint.ee:443",
	"https://ivxv1.kov.ivxv.ee:443",
	"https://ivxv1.ep.ivxv.ee:443",
	"https://irfg3r3rg5g3ree:443",
	"https://.ep.ivxv.ee:443",
	"https://.ep.ivxv.ee:443",
	"https://.e:443",
	"https://.:443",
}

var invalidOrigins = []string{
	"https://234ehs-rp-01demo.riaint.ee:443/api/v2",
	"https://ehs-logimonitor-01.prod.riaint.ee:443/",
	"https://ehs-logimonitor-01.prod.riaint.ee:443#",
	"ws://ehs-logimonitor--01--.prod.riaint.ee:443/",
	"https://234ehs-rp-01demo.riaint.ee:443?",
	"https://.ep.ivxv.ee:443/do?v=IWuSwMC-GjE",
	"https://ehs-rp-01-demo.riaint.ee:8080",
	"https://.ep.ivxv.ee:44#1234E",
	"https://.ep.ivxv.ee:4a43",
	"https://.ep.ivxv.ee:44?",
	"https://:443",
	"http://:443",
	"htt://:443",
	"htt:/:443a",
	"htt::443a",
	"ivxv1",
	":443",
	"443",
	":",
	"",
}

func TestVerifyOriginValid(t *testing.T) {
	msg := "Expected origin %v to be valid"
	for _, origin := range validOrigins {
		if !VerifyHTTPSOrigin(origin) {
			t.Errorf(msg, origin)
		}
	}
}

func TestVerifyOriginInValid(t *testing.T) {
	msg := "Expected origin %v to be invalid"
	for _, origin := range invalidOrigins {
		if VerifyHTTPSOrigin(origin) {
			t.Errorf(msg, origin)
		}
	}
}
