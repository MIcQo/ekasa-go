package ekasa

import (
	"net/http"
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	publicKey := "pub"
	privateKey := "priv"
	tenantID := "tenant"
	baseURL := "https://test.ekasa.sk"
	timeout := 15 * time.Second

	client := NewClient(publicKey, privateKey,
		WithBaseURL(baseURL),
		WithTenantID(tenantID),
		WithTimeout(timeout),
	)

	if client.opts.baseURL != baseURL {
		t.Errorf("Expected baseURL %s, got %s", baseURL, client.opts.baseURL)
	}

	if client.opts.tenantID != tenantID {
		t.Errorf("Expected tenantID %s, got %s", tenantID, client.opts.tenantID)
	}

	if client.opts.timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.opts.timeout)
	}

	if client.opts.httpClient.Timeout != timeout {
		t.Errorf("Expected httpClient.Timeout %v, got %v", timeout, client.opts.httpClient.Timeout)
	}
}

func TestWithHttpClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Second}
	client := NewClient("pub", "priv", WithHttpClient(customClient))

	if client.opts.httpClient != customClient {
		t.Errorf("Expected custom http client to be set")
	}
}
