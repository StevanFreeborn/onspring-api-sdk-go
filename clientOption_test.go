package onspring_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestWithHTTPClient(t *testing.T) {
	t.Run("it should set a custom HTTP client on the Onsring client", func(t *testing.T) {
		customHTTPClient := &http.Client{Timeout: 10 * time.Second}

		clientWithCustomHTTP := onspring.NewClient(
			"test-api-key",
			onspring.WithHTTPClient(customHTTPClient),
		)

		if clientWithCustomHTTP.HTTPClient() != customHTTPClient {
			t.Errorf("Expected custom HTTP client to be set")
		}
	})
}

func TestWithBaseURL(t *testing.T) {
	t.Run("it should set a custom base URL on the Onsring client", func(t *testing.T) {
		customBaseURL := "https://custom.onspring.com"

		clientWithCustomURL := onspring.NewClient(
			"test-api-key",
			onspring.WithBaseURL(customBaseURL),
		)

		if clientWithCustomURL.BaseURL() != customBaseURL {
			t.Errorf("Expected base URL to be %s, got %s", customBaseURL, clientWithCustomURL.BaseURL())
		}
	})
}

func TestWithAPIVersion(t *testing.T) {
	t.Run("it should set a custom API version on the Onsring client", func(t *testing.T) {
		customVersion := "3.0"

		clientWithCustomVersion := onspring.NewClient(
			"test-api-key",
			onspring.WithAPIVersion(customVersion),
		)

		if clientWithCustomVersion.APIVersion() != customVersion {
			t.Errorf("Expected API version to be %s, got %s", customVersion, clientWithCustomVersion.APIVersion())
		}
	})
}
