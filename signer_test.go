package ekasa

import (
	"testing"
	"time"
)

func TestSigner_Sign(t *testing.T) {
	publicKey := "8248d4058e36840bea23ebbe3e602f6e"
	privateKey := "388e98f5c56728266991583a6f8fcd1b9279cdc00b5c371bffc0ea402b14d954"
	tenantID := "39fd0386-1c7b-5fb0-201f-36725cbfcacc"

	// 2021-06-24 15:47:57.213 UTC
	mockTime := time.Date(2021, 6, 24, 15, 47, 57, 213000000, time.UTC)

	s := newSigner(publicKey, privateKey)
	s.timeProvider = func() time.Time { return mockTime }

	headers := map[string]string{
		"__tenant": tenantID,
	}
	rawURL := "http://localhost:5000/api/v1/registrations/receipts/processable/first?cashRegisterCode=88812345678900001"

	err := s.sign("GET", rawURL, headers, nil)
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}

	expectedAuth := "NWS4-HMAC-SHA256 Credential%3D8248d4058e36840bea23ebbe3e602f6e%2CSignedHeaders%3Dhost%253Bx-nd-content-sha256%253Bx-nd-date%253B__tenant%2CTimestamp%3D2021-06-24T15%253A47%253A57.213Z%2CSignature%3D510cfbf80d0f27e78a51a0e06407cdc6f01e2758dd6aea655fb64387cf252e6f"

	if headers["Authorization"] != expectedAuth {
		t.Errorf("Expected Authorization:\n%s\nGot:\n%s", expectedAuth, headers["Authorization"])
	}

	if headers["x-nd-date"] != "2021-06-24T15:47:57.213Z" {
		t.Errorf("Expected x-nd-date: 2021-06-24T15:47:57.213Z, Got: %s", headers["x-nd-date"])
	}

	if headers["x-nd-content-sha256"] != emptyBodyHash {
		t.Errorf("Expected x-nd-content-sha256: %s, Got: %s", emptyBodyHash, headers["x-nd-content-sha256"])
	}

	if headers["host"] != "localhost:5000" {
		t.Errorf("Expected host: localhost:5000, Got: %s", headers["host"])
	}
}
