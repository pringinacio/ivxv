package conf

import (
	"context"
	"testing"

	"ivxv.ee/common/collector/log"

	_ "ivxv.ee/common/collector/container/dummy"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		trust     string
		election  string
		technical string
	}{
		{"dummy", "testdata/trust.dummy", "testdata/election.dummy", "testdata/technical.dummy"},
		{"prefixed", "testdata/prefix.trust.dummy", "", ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, _, err := New(log.TestContext(context.Background()),
				test.trust, test.election, test.technical); err != nil {

				t.Fatal("New failed:", err)
			}
		})
	}
}
