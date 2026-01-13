package onspring_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestNewClient(t *testing.T) {
	t.Run("it should create client with default settings", func(t *testing.T) {
		apiKey := "test-api-key"
		client := onspring.NewClient(apiKey)

		if client.APIKey() != apiKey {
			t.Errorf("Expected apiKey %s, got %s", apiKey, client.APIKey())
		}

		if client.BaseURL() != "https://api.onspring.com" {
			t.Errorf("Expected default baseURL, got %s", client.BaseURL())
		}

		if client.HTTPClient().Timeout != 120*time.Second {
			t.Errorf("Expected default timeout, got %v", client.HTTPClient().Timeout)
		}

	})

	t.Run("it should create client with custom settings", func(t *testing.T) {
		apiKey := "test-api-key"
		customURL := "https://custom.onspring.com"
		customHTTPClient := &http.Client{Timeout: 50 * time.Second}

		client := onspring.NewClient(
			apiKey,
			onspring.WithBaseURL(customURL),
			onspring.WithHTTPClient(customHTTPClient),
		)

		if client.APIKey() != apiKey {
			t.Errorf("Expected apiKey %s, got %s", apiKey, client.APIKey())
		}

		if client.BaseURL() != customURL {
			t.Errorf("Expected baseURL %s, got %s", customURL, client.BaseURL())
		}

		if client.HTTPClient() != customHTTPClient {
			t.Errorf("Expected custom HTTP client, got %v", client.HTTPClient())
		}
	})
}
