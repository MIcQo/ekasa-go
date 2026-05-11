package ekasa

import (
	"net/http"
	"time"
)

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*clientOptions)

type clientOptions struct {
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
	tenantID   string
	proxyURL   string
}

func defaultOptions() clientOptions {
	return clientOptions{
		baseURL: "https://ekasa-cloud.ninedigit.sk/api",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		timeout: 30 * time.Second,
	}
}

// WithBaseURL sets the base URL for the API.
func WithBaseURL(url string) ClientOption {
	return func(o *clientOptions) {
		o.baseURL = url
	}
}

// WithHttpClient sets a custom http.Client.
func WithHttpClient(c *http.Client) ClientOption {
	return func(o *clientOptions) {
		o.httpClient = c
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(o *clientOptions) {
		o.timeout = d
		if o.httpClient != nil {
			o.httpClient.Timeout = d
		}
	}
}

// WithTenantID sets the tenant ID header for requests.
func WithTenantID(id string) ClientOption {
	return func(o *clientOptions) {
		o.tenantID = id
	}
}

// WithProxy sets the proxy URL.
func WithProxy(url string) ClientOption {
	return func(o *clientOptions) {
		o.proxyURL = url
		// In Go, proxy is usually configured on the Transport, not directly on the client options
		// but we store it here to apply if needed when building the client.
	}
}

const (
	// ProductionEnvironment is the production API URL.
	ProductionEnvironment = "https://ekasa-cloud.ninedigit.sk/api"
	// PlaygroundEnvironment is the integration/testing API URL.
	PlaygroundEnvironment = "https://ekasa-cloud-int.ninedigit.sk/api"
)
