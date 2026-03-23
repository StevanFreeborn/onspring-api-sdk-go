package onspring

import "net/http"

// ClientOption is a functional option for configuring a Client.
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client for the Onspring client.
// If not provided, the client will use http.DefaultClient.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithBaseURL sets a custom base URL for the Onspring API.
// This is useful for testing or when using a different Onspring environment.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithAPIVersion sets a custom API version for the Onspring client.
// If not provided, the client will use the default API version.
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) {
		c.apiVersion = version
	}
}
