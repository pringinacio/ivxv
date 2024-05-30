package webeid

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"ivxv.ee/common/collector/cryptoutil"
	tkn "ivxv.ee/common/collector/token"
)

func TestUnmarshalWithoutRawToken(t *testing.T) {
	// Build Web eID token without raw Web eID token
	token := NewFromRawBuilder().Build()

	// Unmarshal literally empty raw Web eID token
	_, err := token.Unmarshal()

	if _, ok := err.(JSONUnmarshalError); !ok {
		msg := "Expected JSONUnmarshalError, got %v\n"
		t.Errorf(msg, err)
	}
}

//nolint:lll
func TestUnmarshalWithBadBase64Signature(t *testing.T) {
	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "",
		Signature:             "Hello World",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().WithToken(string(rawToken)).Build()

	_, err = token.Unmarshal()
	if _, ok := err.(Base64DecodeSignatureError); !ok {
		msg := "Expected Base64DecodeSignatureError, got %v\n"
		t.Errorf(msg, err)
	}
}

//nolint:lll
func TestVerifyWithBadHashAlgoForOrigin(t *testing.T) {
	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm",
		Format:                "web-eid:1.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	err = unmarshalled.Verify()
	if _, ok := err.(OriginHashError); !ok {
		msg := "Expected OriginHashError, got %v\n"
		t.Errorf(msg, err)
	}
}

//nolint:lll
func TestVerifyWithBadMajorRelease(t *testing.T) {
	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm",
		Format:                "web-eid:2.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().
		WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	err = unmarshalled.Verify()
	if _, ok := err.(VerifyMajorReleaseError); !ok {
		msg := "Expected VerifyMajorReleaseError, got %v\n"
		t.Errorf(msg, err)
	}
}

//nolint:lll
func TestRawWebeidTokenDoesntImplementHeaderMethod(t *testing.T) {
	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm",
		Format:                "web-eid:1.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().
		WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	header, err := unmarshalled.Header()
	if header != "" && err != nil {
		msg := `Expected header: "" and err: nil, got header: %v, err: %v\n`
		t.Errorf(msg, header, err)
	}
}

//nolint:lll
func TestRawWebeidTokenDoesntImplementPayload(t *testing.T) {
	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm",
		Format:                "web-eid:1.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().
		WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	header, err := unmarshalled.Payload()
	if header != "" && err != nil {
		msg := `Expected payload: "" and err: nil, got payload: %v, err: %v\n`
		t.Errorf(msg, header, err)
	}
}

//nolint:lll
func TestCertifyReturnsUserCert(t *testing.T) {
	certPEM := "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw"

	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm",
		Format:                "web-eid:1.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().
		WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	cert := unmarshalled.(tkn.Certifier).Certify()

	certFromString, err := cryptoutil.Base64Certificate(certPEM)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	if !reflect.DeepEqual(cert, certFromString) {
		msg := "Expected two x509 certificates are equal, but they aren't"
		t.Errorf(msg)
	}
}

//nolint:lll
func TestRawWebeidTokenImplementsSignature(t *testing.T) {
	sig := "MGQCMG4DQut8Zi3EU0eySnxCvHA2pW42gj3OqZrHnIRs3a70DnXq9JrANPtsi13FxhppOgIwfvvkhyiM8p1iQg42rScFxYZ+8hecBl8Sd/4GzuFnyysmi7Rk4g9P6Z5QQBQrTDAm"

	template := webEidAuthToken{
		UnverifiedCertificate: "MIICnDCCAf6gAwIBAgIDBwqKMAoGCCqGSM49BAMCMGMxCzAJBgNVBAYTAkVFMRIwEAYDVQQKDAlTQ0NFSVYgT1kxHzAdBgNVBAsMFklWWFYgVGVzdCBDZXJ0aWZpY2F0ZXMxHzAdBgNVBAMMFlBlcnNvbiBDQSBJbnRlcm1lZGlhdGUwIBcNMjEwNDI4MTQyNzEwWhgPMjEyMTA0MDQxNDI3MTBaMHExCzAJBgNVBAYTAkVFMSMwIQYDVQQDDBpLT0JSQVMsRVVST09QQSw0MTYwMjI5MDA3ODEPMA0GA1UEBAwGS09CUkFTMRAwDgYDVQQqDAdFVVJPT1BBMRowGAYDVQQFExFQTk9FRS00MTYwMjI5MDA3ODB2MBAGByqGSM49AgEGBSuBBAAiA2IABBrsfaycHJPzAzLT9Eob5TWRidsIzIf9MGjYriLFP7vCtEuVXZrRUVJlqTvZ7KIJi2nfpxYJgM/iomiewVd5NIjnkZaKY/vvClNM/3a3R2COX1C/9C/bifek1Pc11Zb6FKN0MHIwHQYDVR0OBBYEFK2G9S7BGci2AicpZ281+9mnq2rMMB8GA1UdIwQYMBaAFIH/R+CRqJF+ti34SbD0J4rYyaY2MA4GA1UdDwEB/wQEAwIDiDAgBgNVHSUBAf8EFjAUBggrBgEFBQcDAgYIKwYBBQUHAwQwCgYIKoZIzj0EAwIDgYsAMIGHAkF+VHMSTvhUd90uEB2fXMIK2ZHcuBeUy/5qx4HQQMlpGy+tUiDlTtpzx3ca7FHlEUy77uNv+nFU/9hu2tI6D5bxgAJCAPoBI1Eilu5BoV3l8RC6ngBSEGTuBRWCP+yPqqNb/I+pnOck3TC76W8nUSOs2l8Q9uuFt0IBlJnde04kBntWPaOw",
		Algorithm:             "NO111",
		Signature:             sig,
		Format:                "web-eid:1.0",
	}

	rawToken, err := json.Marshal(&template)
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	token := NewFromRawBuilder().
		WithToken(string(rawToken)).Build()

	unmarshalled, err := token.Unmarshal()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	signature, err := unmarshalled.Signature()
	if err != nil {
		t.Errorf(msgExpectNoErrors, err)
	}

	b64Sig := base64.StdEncoding.EncodeToString(signature)
	if b64Sig != sig {
		msg := "Expected sig == sig from Web eID, got sig: %v, sig from Web eID: %v\n"
		t.Errorf(msg, sig, b64Sig)
	}
}
