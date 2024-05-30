package custom

import "testing"

func TestRawBearerRegex(t *testing.T) {
	goodRawBearerRegexes := []string{
		"Hello.World",
		"World.Hello",
		`<script>alert("Hello World!")</script>.1`,
		"1.1",
	}

	for _, goodRawBearerRegex := range goodRawBearerRegexes {
		if ok := rawBearerRegex(goodRawBearerRegex); !ok {
			msg := "Expected true for %v, got %v\n"
			t.Errorf(msg, goodRawBearerRegex, ok)
		}
	}

	badRawBearerRegexes := []string{
		"HelloWorld",
		".HelloWorld",
		"HelloWorld.",
		"",
		".",
	}

	for _, badRawBearerRegex := range badRawBearerRegexes {
		if ok := rawBearerRegex(badRawBearerRegex); ok {
			msg := "Expected false for %v, got %v\n"
			t.Errorf(msg, badRawBearerRegex, ok)
		}
	}
}

func TestSplitRawBearer(t *testing.T) {
	msg := "Expected payload: %v and signature %v, got payload: %v and signature %v\n"

	goodRawBearers := map[string][]string{
		"Hello.World": {"Hello", "World"},
		"World.Hello": {"World", "Hello"},
		`<script>alert("Hello World!")</script>.1`: {`<script>alert("Hello World!")</script>`, "1"},
		"1.1": {"1", "1"},
	}

	for k, v := range goodRawBearers {
		payload, signature := splitRawBearer(k)
		if v[0] != payload && v[1] != signature {
			t.Errorf(msg, v[0], v[1], payload, signature)
		}
	}

	badRawBearers := map[string][]string{
		"HelloWorld":  {"", ""},
		".HelloWorld": {"", ""},
		"HelloWorld.": {"", ""},
		"":            {"", ""},
		".":           {"", ""},
	}

	for k, v := range badRawBearers {
		payload, signature := splitRawBearer(k)
		if v[0] != payload && v[1] != signature {
			t.Errorf(msg, v[0], v[1], payload, signature)
		}
	}
}
