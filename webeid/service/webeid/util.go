package main

import (
	"crypto/x509/pkix"
	"encoding/asn1"
	"strings"
)

const (
	StatusOK = "OK"
	PNOEE    = "PNOEE-"
)

// oid is a map of asn1.ObjectIdentifier that is used to extract particular
// data from a x509.Certificate.Subject.Names.
var oid = map[string][]int{
	"personalCode": asn1.ObjectIdentifier{2, 5, 4, 5},
	"givenName":    asn1.ObjectIdentifier{2, 5, 4, 42},
	"surname":      asn1.ObjectIdentifier{2, 5, 4, 4},
}

// findName searches name for oID and returns the value for that oID or an
// empty string. Panics if the value for the oID is not a string.
func findName(name *pkix.Name, oID asn1.ObjectIdentifier) string {
	for _, n := range name.Names {
		if n.Type.Equal(oID) {
			return n.Value.(string)
		}
	}
	return ""
}

// personalCode reads personal code (isikukood) from a cert.
func personalCode(cert *pkix.Name) string {
	code := findName(cert, oid["personalCode"])
	return strings.TrimSuffix(code, PNOEE)
}
