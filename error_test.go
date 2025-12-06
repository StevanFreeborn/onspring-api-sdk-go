package onspring_test

import (
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestOnspringAPIError(t *testing.T) {
	t.Run("it should create OnspringAPIError instance", func(t *testing.T) {
		statusCode := 500
		message := "Internal Server Error"

		err := &onspring.OnspringAPIError{
			StatusCode: statusCode,
			Message:    message,
		}

		if err.StatusCode != statusCode {
			t.Errorf("Expected status code %d, got %d", statusCode, err.StatusCode)
		}

		if err.Message != message {
			t.Errorf("Expected message %s, got %s", message, err.Message)
		}
	})

	t.Run("it should return formatted error message", func(t *testing.T) {
		err := &onspring.OnspringAPIError{
			StatusCode: 404,
			Message:    "Not Found",
		}

		expectedMessage := "onspring api error: status=404 message=Not Found"

		if err.Error() != expectedMessage {
			t.Errorf("Expected error message %s, got %s", expectedMessage, err.Error())
		}
	})
}
