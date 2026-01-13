package onspring_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func setupMockServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *onspring.Client) {
	t.Helper()

	server := httptest.NewServer(handler)

	t.Cleanup(func() {
		server.Close()
	})

	client := onspring.NewClient(
		"test-key",
		onspring.WithBaseURL(server.URL),
		onspring.WithHTTPClient(server.Client()),
	)

	return server, client
}

type ErrorTransport struct{}

func (t *ErrorTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return nil, errors.New("simulated network connection failure")
}
