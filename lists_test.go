package onspring_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestLists(t *testing.T) {
	t.Run("Save", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			listItem := onspring.SaveListItemRequest{
				Name: "List Value",
			}

			_, err := client.Lists.Save(nilContext, 1, listItem)

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

			listItem := onspring.SaveListItemRequest{
				Name: "List Value",
			}

			_, err := client.Lists.Save(ctx, 1, listItem)

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

			listItem := onspring.SaveListItemRequest{
				Name: "List Value",
			}

			_, err := client.Lists.Save(t.Context(), 1, listItem)

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

			listItem := onspring.SaveListItemRequest{
				Name: "List Value",
			}

			_, err := invalidClient.Lists.Save(t.Context(), 1, listItem)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			listItem := onspring.SaveListItemRequest{
				Name: "List Value",
			}

			_, err := client.Lists.Save(t.Context(), 1, listItem)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a PUT request to the /lists/id/:listId/items endpoint and return a response when successful", func(t *testing.T) {
			listId := 1
			numericValue := 0.0
			color := "#ffffff"

			listItem := onspring.SaveListItemRequest{
				Name:         "List Value",
				NumericValue: &numericValue,
				Color:        &color,
			}

			expectedResponse := onspring.SaveListItemResponse{
				Id: "d4a3c2b1-e5f6-7890-abcd-ef1234567890",
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("Expected PUT method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/lists/id/%d/items", listId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				var body onspring.SaveListItemRequest
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !reflect.DeepEqual(listItem, body) {
					t.Errorf("Expected body to be %v but got %v", listItem, body)
				}

				jsonData, _ := json.Marshal(expectedResponse)

				w.WriteHeader(http.StatusCreated)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			response, err := client.Lists.Save(t.Context(), listId, listItem)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedResponse, response) {
				t.Errorf("Expected %v but got %v", expectedResponse, response)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			var nilContext context.Context = nil

			err := client.Lists.Delete(nilContext, 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

			if err == nil {
				t.Errorf("Expected error for nil context, got nil")
			}
		})

		t.Run("it should return an error if context is canceled", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			ctx, cancel := context.WithCancel(t.Context())

			cancel()

			err := client.Lists.Delete(ctx, 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

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

			err := client.Lists.Delete(t.Context(), 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

			if err == nil {
				t.Errorf("Expected network error, got nil")
			}
		})

		t.Run("it should return an error if create a request fails", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			invalidClient := onspring.NewClient(
				"test-api-key",
				onspring.WithBaseURL("http://[::1]:namedport"),
				onspring.WithHTTPClient(client.HTTPClient()),
			)

			err := invalidClient.Lists.Delete(t.Context(), 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			err := client.Lists.Delete(t.Context(), 1, "d4a3c2b1-e5f6-7890-abcd-ef1234567890")

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a DELETE request to the /lists/id/:listId/itemId/:itemId endpoint and return no error if receives 204 status code", func(t *testing.T) {
			listId := 1
			itemId := "d4a3c2b1-e5f6-7890-abcd-ef1234567890"

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("Expected DELETE method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/lists/id/%d/itemId/%s", listId, itemId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusNoContent)
			})

			err := client.Lists.Delete(t.Context(), listId, itemId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	})
}
