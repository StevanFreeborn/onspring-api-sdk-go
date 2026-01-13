package onspring

import (
	"fmt"
)

// OnspringAPIError represents an error returned by the Onspring API.
// It contains the HTTP status code and the error message from the API response.
// This error type implements the error interface.
type OnspringAPIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Message is the error message returned by the API or the HTTP status text.
	Message string
}

// Error returns a formatted error string containing the status code and message.
// It implements the error interface for OnspringAPIError.
//
// Returns:
//   - string: A formatted error message in the format "onspring api error: status={code} message={message}"
func (e *OnspringAPIError) Error() string {
	return fmt.Sprintf("onspring api error: status=%d message=%s", e.StatusCode, e.Message)
}
