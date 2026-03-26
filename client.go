// Package onspring provides a Go SDK for interacting with the Onspring API.
// It offers a type-safe, idiomatic Go interface for making API requests
// to the Onspring platform.
package onspring

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// defaultBaseURL is the default base URL for the Onspring API.
	defaultBaseURL = "https://api.onspring.com"
	// defaultTimeout is the default HTTP client timeout duration.
	defaultTimeout = 120 * time.Second
	// defaultAPIKeyHeader is the HTTP header name used for API key authentication.
	defaultAPIKeyHeader = "x-api-key"
	// defaultAPIVersionHeader is the HTTP header name used to specify the API version.
	defaultAPIVersionHeader = "x-api-version"
	// defaultAPIVersion is the default Onspring API version to use.
	defaultAPIVersion = "2.0"
)

// Client is the main client for interacting with the Onspring API.
// It manages HTTP communication, authentication, and API versioning.
// All API endpoints are accessed through this client.
type Client struct {
	// httpClient is the underlying HTTP client used to make requests.
	httpClient *http.Client
	// baseURL is the base URL for the Onspring API.
	baseURL string
	// apiKey is the API key used for authentication.
	apiKey string
	// apiVersion is the API version to use for requests.
	apiVersion string

	// Ping provides access to the ping endpoint for health checks.
	Ping *PingEndpoint
	// Apps provides access to the apps within an Onspring instance.
	Apps *AppsEndpoint
	// Fields provides access to the fields within an Onspring instance.
	Fields *FieldsEndpoint
	// Lists provides access to the lists within an Onspring instance.
	Lists *ListsEndpoint
	// Reports provides access to the reports within an Onspring instance.
	Reports *ReportsEndpoint
	// Files provides access to the files within an Onspring instance.
	Files *FilesEndpoint
}

// NewClient creates a new Onspring API client with the provided API key.
// It initializes the client with default settings including a 100-second timeout,
// the production API base URL, and API version 2.0.
//
// Optional configuration can be provided using Option functions such as
// WithHTTPClient, WithBaseURL, and WithAPIVersion.
//
// Parameters:
//   - apiKey: The API key for authenticating with the Onspring API
//   - opts: Optional configuration functions to customize the client
//
// Returns:
//   - *Client: A configured Onspring API client ready to make requests
//
// Example:
//
//	client := onspring.NewClient("your-api-key")
//	client := onspring.NewClient("your-api-key", onspring.WithHTTPClient(customHTTPClient))
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    defaultBaseURL,
		apiKey:     apiKey,
		apiVersion: defaultAPIVersion,
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Ping = &PingEndpoint{client: c}
	c.Apps = &AppsEndpoint{client: c}
	c.Fields = &FieldsEndpoint{client: c}
	c.Lists = &ListsEndpoint{client: c}
	c.Reports = &ReportsEndpoint{client: c}
	c.Files = &FilesEndpoint{client: c}

	return c
}

// do executes an HTTP request and handles the response.
// It performs the actual HTTP call using the configured HTTP client,
// checks the response status code, and handles any API errors.
//
// Parameters:
//   - req: The HTTP request to execute
//
// Returns:
//   - error: nil if the request succeeds, or an error if the request fails
//     or returns a non-2xx status code
func (c *Client) do(req *http.Request) error {
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleAPIError(resp)
	}

	return nil
}

// doWithJsonResponse executes an HTTP request and decodes the JSON response.
// It performs the HTTP call, checks the response status code, and decodes
// the JSON response body into the provided variable.
//
// Parameters:
//   - req: The HTTP request to execute
//   - v: A pointer to the variable where the decoded JSON response will be stored
//
// Returns:
//   - error: nil if the request and decoding succeed, or an error if the request fails,
//     returns a non-2xx status code, or if JSON decoding fails
func (c *Client) doWithJsonResponse(req *http.Request, v any) error {
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleAPIError(resp)
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

// doWithBytesResponse executes an HTTP request and returns the raw response bytes and headers.
// It performs the HTTP call, checks the response status code, and reads
// the response body into a byte slice.
//
// Parameters:
//   - req: The HTTP request to execute
//
// Returns:
//   - []byte: The raw response body bytes
//   - http.Header: The response headers
//   - error: nil if the request succeeds, or an error if the request fails
//     or returns a non-2xx status code
func (c *Client) doWithBytesResponse(req *http.Request) ([]byte, http.Header, error) {
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, nil, c.handleAPIError(resp)
	}

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, resp.Header, nil
}

// newMultipartRequest creates a new multipart/form-data HTTP request for the Onspring API.
// It constructs the full URL, sets required authentication headers,
// and prepares the multipart form data with the provided context.
//
// Parameters:
//   - ctx: The context for the request
//   - path: The API endpoint path
//   - body: The request body as a reader (should be a multipart form body)
//   - contentType: The content type header value (including boundary)
//
// Returns:
//   - *http.Request: The prepared HTTP request
//   - error: An error if the context is nil or request creation fails
func (c *Client) newMultipartRequest(ctx context.Context, path string, body io.Reader, contentType string) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}

	fullURL := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), strings.TrimLeft(path, "/"))
	validUrl, urlParsingError := url.Parse(fullURL)

	if urlParsingError != nil {
		return nil, fmt.Errorf("failed to parse the request url: %w", urlParsingError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, validUrl.String(), body)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(defaultAPIKeyHeader, c.apiKey)
	req.Header.Set(defaultAPIVersionHeader, c.apiVersion)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// handleAPIError processes error responses from the Onspring API.
// It attempts to decode the error message from the response body.
// If decoding fails, it falls back to using the HTTP status text.
//
// Parameters:
//   - resp: The HTTP response containing the error
//
// Returns:
//   - error: An OnspringAPIError with the status code and error message
func (c *Client) handleAPIError(resp *http.Response) error {
	var errBody struct {
		Message string `json:"message"`
	}

	decodeErr := json.NewDecoder(resp.Body).Decode(&errBody)

	if decodeErr != nil {
		errBody.Message = http.StatusText(resp.StatusCode)
	}

	return &OnspringAPIError{
		StatusCode: resp.StatusCode,
		Message:    errBody.Message,
	}
}

// newRequest creates a new HTTP request for the Onspring API.
// It constructs the full URL, sets required authentication headers,
// and prepares the request with the provided context.
//
// Parameters:
//   - ctx: The context for the request
//   - method: The HTTP method
//   - path: The API endpoint path
//   - queryParams: The query parameters for the request
//   - body: The request body
//
// Returns:
//   - *http.Request: The prepared HTTP request
//   - error: An error if the context is nil or request creation fails
func (c *Client) newRequest(ctx context.Context, method, path string, queryParams map[string]string, body any) (*http.Request, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context must not be nil")
	}

	fullURL := fmt.Sprintf("%s/%s", strings.TrimRight(c.baseURL, "/"), strings.TrimLeft(path, "/"))
	validUrl, urlParsingError := url.Parse(fullURL)

	if urlParsingError != nil {
		return nil, fmt.Errorf("failed to parse the request url: %w", urlParsingError)
	}

	q := validUrl.Query()

	for key, value := range queryParams {
		q.Add(key, value)
	}

	validUrl.RawQuery = q.Encode()

	var bodyReader io.Reader

	if body != nil {
		jsonData, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}

		bodyReader = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, validUrl.String(), bodyReader)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(defaultAPIKeyHeader, c.apiKey)
	req.Header.Set(defaultAPIVersionHeader, c.apiVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}
