//go:build integration

package onspring_test

import (
	"context"
	"testing"
)

func TestPingIntegration(t *testing.T) {
	loadEnvFile(t)

	t.Run("Get should be able to connect to the API", func(t *testing.T) {
		client := createClient(t)

		err := client.Ping.Get(context.Background())

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
