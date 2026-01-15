package onspring_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestApps(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Apps.Get(nilContext)

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

			_, err := client.Apps.Get(ctx)

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

			_, err := client.Apps.Get(t.Context())

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

			_, err := invalidClient.Apps.Get(t.Context())

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should perform a GET request to the /apps endpoint and return page of apps if receives 200 status code", func(t *testing.T) {
			expectedPageNumber := 1
			expectedPageSize := 50

			expectedPage := onspring.Page[onspring.App]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items: []onspring.App{
					{
						Href: "https://test.com",
						Id:   1,
						Name: "App",
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")
				pageSize := r.URL.Query().Get("pageSize")

				if pageNumber != strconv.Itoa(expectedPageNumber) {
					t.Errorf("Expected query param pageNumber to be %d but got %s", expectedPageNumber, pageNumber)
				}

				if pageSize != strconv.Itoa(expectedPageSize) {
					t.Errorf("Expected query param pageSize to be %d but got %s", expectedPageSize, pageSize)
				}

				jsonData, _ := json.Marshal(expectedPage)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
			})

			page, err := client.Apps.Get(t.Context())

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedPage, page) {
				t.Errorf("Expected %v but got %v", expectedPage, page)
			}
		})

		t.Run("it should perform a GET request to the /apps endpoint with non-default paging information when provided", func(t *testing.T) {
			expectedPageNumber := 2
			expectedPageSize := 1

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")
				pageSize := r.URL.Query().Get("pageSize")

				if pageNumber != strconv.Itoa(expectedPageNumber) {
					t.Errorf("Expected query param pageNumber to be %d but got %s", expectedPageNumber, pageNumber)
				}

				if pageSize != strconv.Itoa(expectedPageSize) {
					t.Errorf("Expected query param pageSize to be %d but got %s", expectedPageSize, pageSize)
				}

				w.WriteHeader(http.StatusOK)
			})

			client.Apps.Get(
				t.Context(),
				onspring.ForPageNumber(expectedPageNumber),
				onspring.WithPageSize(expectedPageSize),
			)
		})

		t.Run("it should return an error if the /apps endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Apps.Get(t.Context())

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	})
}
