package webeid

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"reflect"
	"testing"
)

var msgExpectNoErrors = "Expected no errors, got %v\n"

func TestFormatRegex(t *testing.T) {
	goodRegexFormats := []string{
		"web-eid:2.0",
		"web-eid:1.0",
		"web-eid:1.0",
		"web-eid:1.3",
		"web-eid:1.999999999999",
		"web-eid:0.999999999999",
	}

	for _, goodRegexFormat := range goodRegexFormats {
		_, err := formatRegex(goodRegexFormat)
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}
	}

	badRegexFormats := []string{
		"Hello World!",
		"web-eid:1.",
		"web-eid:.9",
		"web-eid:.",
		"web-eid:",
		"web-eid",
		"",
	}

	for _, badRegexFormat := range badRegexFormats {
		_, err := formatRegex(badRegexFormat)
		_, ok := err.(FormatRegexError)
		if err == nil || !ok {
			msg := "Expected FormatRegexError, got no error at %v\n"
			t.Errorf(msg, badRegexFormat)
		}
	}
}

func TestFormatMajorRelease(t *testing.T) {
	goodFormatMajorReleases := []string{
		"web-eid:1.0",
		"web-eid:1.3",
	}

	for _, goodFormatMajorRelease := range goodFormatMajorReleases {
		release, err := formatRegex(goodFormatMajorRelease)
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}

		err = formatMajorRelease(release)
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}
	}

	badFormatMajorReleases := []string{
		"web-eid:2.0",
		"web-eid:0.999999999999",
	}

	for _, badFormatMajorRelease := range badFormatMajorReleases {
		release, err := formatRegex(badFormatMajorRelease)
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}

		err = formatMajorRelease(release)
		_, ok := err.(FormatMajorReleaseError)
		if err == nil || !ok {
			msg := "Expected FormatMajorReleaseError, got no error at %v\n"
			t.Errorf(msg, badFormatMajorRelease)
		}
	}
}

func TestAlgorithmRegex(t *testing.T) {
	onlySupportedAlgos := map[string]string{
		"ES256": "256",
		"ES384": "384",
		"ES512": "512",
		"PS256": "256",
		"PS384": "384",
		"PS512": "512",
		"RS256": "256",
		"RS384": "384",
		"RS512": "512",
	}

	for onlySupportedAlgo, algo1 := range onlySupportedAlgos {
		algo2, err := algorithmRegex(onlySupportedAlgo)
		if algo1 != algo2 {
			msg := "Expected both algorithms equal, got a1: %v, a2: %v\n"
			t.Errorf(msg, algo1, algo2)
		}
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}
	}

	badAlgos := map[string]string{
		"ES257":        "256",
		"EZ384":        "384",
		"512":          "512",
		"PS":           "256",
		"":             "384",
		"PS512PS512":   "512",
		"256RS":        "256",
		" ":            "384",
		"Hello World!": "512",
	}

	for badAlgo, algo1 := range badAlgos {
		algo2, err := algorithmRegex(badAlgo)
		if algo1 == algo2 {
			msg := "Expected both algorithms not equal, got a1: %v, a2: %v\n"
			t.Errorf(msg, algo1, algo2)
		}
		_, ok := err.(AlgorithmRegexError)
		if err == nil || !ok {
			msg := "Expected AlgorithmRegexError, got %v\n"
			t.Errorf(msg, err)
		}
	}
}

func TestHashIt(t *testing.T) {
	goodAlgoAndHash := map[string][]byte{
		"ES256": []byte(""),
		"ES384": []byte("384"),
		"ES512": nil,
	}

	for algo, data := range goodAlgoAndHash {
		_, err := hashIt(data, algo)
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}
	}

	badAlgoAndHash := map[string][]byte{
		"ES257": []byte("256"),
		"EZ384": nil,
		"":      []byte(""),
	}

	for algo, data := range badAlgoAndHash {
		_, err := hashIt(data, algo)
		_, ok := err.(HashDataError)
		if err == nil || !ok {
			msg := "Expected HashDataError, got %v\n"
			t.Errorf(msg, err)
		}
	}
}

func TestHashAlgorithm(t *testing.T) {
	goodAlgos := map[string]hash.Hash{
		"ES256": sha256.New(),
		"ES384": sha512.New384(),
		"RS512": sha512.New(),
	}

	for algo1, algo2 := range goodAlgos {
		hashAlgo, err := hashAlgorithm(algo1)
		if hashAlgo.Size() != algo2.Size() {
			msg := "Expected equal hashes, got h1: %v, h2: %v\n"
			t.Errorf(msg, hashAlgo, algo2)
		}
		if err != nil {
			t.Errorf(msgExpectNoErrors, err)
		}
	}

	badAlgos := map[string]hash.Hash{
		"":    sha512.New384(),
		"512": sha512.New(),
		"RS":  sha256.New(),
	}

	for algo1, algo2 := range badAlgos {
		hashAlgo, err := hashAlgorithm(algo1)
		if reflect.DeepEqual(hashAlgo, algo2) {
			msg := "Expected no hashes, got h1: %v, h2: %v\n"
			t.Errorf(msg, hashAlgo, algo2)
		}
		_, ok := err.(MalformedAlgorithmError)
		if err == nil || !ok {
			msg := "Expected MalformedAlgorithmError, got %v\n"
			t.Errorf(msg, err)
		}
	}
}
