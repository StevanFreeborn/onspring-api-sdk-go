package onspring

import (
	"context"
	"net/http"
)

const (
	// pingPath is the API endpoint path for the ping health check.
	pingPath = "/ping"
)

// PingEndpoint provides access to the Onspring API ping endpoint.
// It can be used to verify API connectivity and authentication.
type PingEndpoint struct {
	// client is the parent Client used to make API requests.
	client *Client
}

// Get performs a ping request to verify API connectivity.
// This is a simple health check that can be used to test if the API is
// reachable.
//
// Parameters:
//   - ctx: The context for the request
//
// Returns:
//   - error: nil if the ping succeeds, or an error if the request fails
//     or authentication is invalid
//
// Example:
//
//	client := onspring.NewClient("your-api-key")
//	err := client.Ping.Get(context.Background())
//	if err != nil {
//	    log.Fatal("Ping failed:", err)
//	}
func (p *PingEndpoint) Get(ctx context.Context) error {
	req, err := p.client.newRequest(ctx, http.MethodGet, pingPath, nil, nil)

	if err != nil {
		return err
	}

	return p.client.do(req)
}
