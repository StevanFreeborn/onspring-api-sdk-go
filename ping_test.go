package onspring_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestGet(t *testing.T) {
	t.Run("it should return an error if context is nil", func(t *testing.T) {
		_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		var nilContext context.Context = nil

		err := client.Ping.Get(nilContext)

		if err == nil {
			t.Errorf("Expected error for nil context, got nil")
		}
	})

	t.Run("it should return an error if context is canceled", func(t *testing.T) {
		_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		ctx, cancel := context.WithCancel(t.Context())

		cancel()

		err := client.Ping.Get(ctx)

		if err == nil {
			t.Errorf("Expected error for canceled context, got nil")
		}
	})

	t.Run("it should return an error if encounters a network error", func(t *testing.T) {
		client := onspring.NewClient(
			"test-api-key",
			onspring.WithBaseURL("http://invalid-url"),
			onspring.WithHTTPClient(&http.Client{Transport: &ErrorTransport{}}),
		)

		err := client.Ping.Get(t.Context())

		if err == nil {
			t.Errorf("Expected network error, got nil")
		}
	})

	t.Run("it should return an error if create a request fails", func(t *testing.T) {
		_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		invalidClient := onspring.NewClient(
			"test-api-key",
			onspring.WithBaseURL("http://[::1]:namedport"),
			onspring.WithHTTPClient(client.HTTPClient()),
		)

		err := invalidClient.Ping.Get(t.Context())

		if err == nil {
			t.Errorf("Expected request creation error, got nil")
		}
	})

	t.Run("it should perform a GET request to the /ping endpoint and return no error if receives 200 status code", func(t *testing.T) {
		_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("Expected GET method, got %s", r.Method)
			}

			if r.URL.Path != "/ping" {
				t.Errorf("Expected /ping endpoint, got %s", r.URL.Path)
			}

			w.WriteHeader(http.StatusOK)
		})

		err := client.Ping.Get(t.Context())

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("it should return an error if the /ping endpoint returns a non-200 status code", func(t *testing.T) {
		_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})

		err := client.Ping.Get(t.Context())

		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})
}
