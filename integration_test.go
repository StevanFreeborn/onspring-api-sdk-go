//go:build integration

package onspring_test

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func loadEnvFile(t *testing.T) {
	t.Helper()

	file, err := os.Open(".env")
	if err != nil {
		t.Fatalf("Failed to open .env file: %v", err)
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")

		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		t.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		t.Fatalf("Failed to read .env file: %v", err)
	}
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()

	value := os.Getenv(key)

	if value == "" {
		t.Fatalf("Required environment variable %s is not set", key)
	}

	return value
}

func requireEnvInt(t *testing.T, key string) int {
	t.Helper()

	value := requireEnv(t, key)
	intValue, err := strconv.Atoi(value)

	if err != nil {
		t.Fatalf("Environment variable %s is not a valid integer: %v", key, err)
	}

	return intValue
}

func requireEnvIntSlice(t *testing.T, key string) []int {
	t.Helper()

	value := requireEnv(t, key)
	parts := strings.Split(value, ",")
	result := make([]int, len(parts))

	for i, part := range parts {
		intValue, err := strconv.Atoi(strings.TrimSpace(part))

		if err != nil {
			t.Fatalf("Environment variable %s contains invalid integer '%s': %v", key, part, err)
		}

		result[i] = intValue
	}

	return result
}

func createClient(t *testing.T) *onspring.Client {
	t.Helper()

	return onspring.NewClient(
		requireEnv(t, "SANDBOX_API_KEY"),
		onspring.WithBaseURL(requireEnv(t, "API_BASE_URL")),
	)
}

func createInvalidClient(t *testing.T) *onspring.Client {
	t.Helper()

	return onspring.NewClient(
		"invalid-api-key",
		onspring.WithBaseURL(requireEnv(t, "API_BASE_URL")),
	)
}

func assertAPIError(t *testing.T, err error, expectedStatus int) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	var apiErr *onspring.OnspringAPIError

	if !errors.As(err, &apiErr) {
		t.Fatalf("Expected OnspringAPIError, got %T: %v", err, err)
	}

	if apiErr.StatusCode != expectedStatus {
		t.Errorf("Expected status code %d, got %d", expectedStatus, apiErr.StatusCode)
	}
}

func addTestRecord(t *testing.T, client *onspring.Client, appId, textFieldId int) int {
	t.Helper()

	ctx := context.Background()

	response, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
		AppId: appId,
		Fields: map[string]any{
			fmt.Sprintf("%d", textFieldId): "test",
		},
	})

	if err != nil {
		t.Fatalf("Failed to create test record: %v", err)
	}

	t.Cleanup(func() {
		_ = client.Records.Delete(context.Background(), appId, response.Id)
	})

	return response.Id
}
