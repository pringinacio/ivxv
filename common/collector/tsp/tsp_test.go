package tsp

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"ivxv.ee/common/collector/log"

	_ "crypto/sha256"
)

var client *Client

func init() {
	pem, err := os.ReadFile("testdata/demo_tsu_ecc_2023.pem.crt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	client, err = New(&Conf{
		URL:       "http://demo.sk.ee/tsaecc/",
		Signers:   []string{string(pem)},
		DelayTime: 1,
		Retry:     2,
		MaxSkew:   1,
		MaxAge:    2,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func TestCreate(t *testing.T) {
	if testing.Short() {
		t.Skip("Short mode on, skipping TSP test against live responder")
	}

	ctx := log.TestContext(context.Background())

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	_, err := client.Create(ctx, data, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheck(t *testing.T) {
	bytes, err := os.ReadFile("testdata/test_response_08-11-2023")
	if err != nil {
		t.Fatal(err)
	}

	data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	nonce := []byte{0xc8, 0xac, 0x0e, 0x31, 0x01, 0x26, 0x36, 0xd5, 0x1b,
		0x87, 0x28, 0xb6, 0x18, 0x4f, 0x4e, 0x66, 0xc7, 0x65, 0x2e, 0xcd}
	genTime, err := client.Check(bytes, data, nonce)
	if err != nil {
		t.Fatal(err)
	}
	timestamp := time.Date(2023, time.November, 8, 13, 59, 49, 0, time.UTC)
	if !genTime.Equal(timestamp) {
		t.Errorf("reported genTime %s does not match expected %s", genTime, timestamp)
	}
}
