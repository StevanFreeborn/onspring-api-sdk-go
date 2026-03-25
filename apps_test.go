package onspring_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"strconv"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestApps(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Apps.List(nilContext)

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

			_, err := client.Apps.List(ctx)

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

			_, err := client.Apps.List(t.Context())

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

			_, err := invalidClient.Apps.List(t.Context())

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
				_, _ = w.Write(jsonData)
			})

			page, err := client.Apps.List(t.Context())

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

			_, _ = client.Apps.List(
				t.Context(),
				onspring.ForPageNumber(expectedPageNumber),
				onspring.WithPageSize(expectedPageSize),
			)
		})

		t.Run("it should return an error if the /apps endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Apps.List(t.Context())

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("it should return an error if fails to retrieve any pages of apps", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			for _, err := range client.Apps.ListAll(t.Context()) {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}
		})

		t.Run("it should return all the apps from multiple pages", func(t *testing.T) {
			expectedApps := []onspring.App{
				{
					Href: "https://test.com",
					Id:   1,
					Name: "App",
				},
				{
					Href: "https://test.com",
					Id:   2,
					Name: "App",
				},
			}

			pageOne := onspring.Page[onspring.App]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        []onspring.App{expectedApps[0]},
			}

			pageTwo := onspring.Page[onspring.App]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        []onspring.App{expectedApps[1]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "1" {
					jsonData, _ := json.Marshal(pageOne)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}

				if pageNumber == "2" {
					jsonData, _ := json.Marshal(pageTwo)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedApps := []onspring.App{}

			for app, _ := range client.Apps.ListAll(t.Context()) {
				retrievedApps = append(retrievedApps, app)
			}

			if !slices.Equal(expectedApps, retrievedApps) {
				t.Errorf("Expected %v but got %v", expectedApps, retrievedApps)
			}
		})

		t.Run("it should return apps and errors if some pages fail and some succeed", func(t *testing.T) {
			expectedApps := []onspring.App{
				{
					Href: "https://test.com",
					Id:   1,
					Name: "App",
				},
			}

			pageOne := onspring.Page[onspring.App]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        []onspring.App{expectedApps[0]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "1" {
					jsonData, _ := json.Marshal(pageOne)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}

				if pageNumber == "2" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			retrievedApps := []onspring.App{}
			encounteredErrors := []error{}

			for app, err := range client.Apps.ListAll(t.Context()) {
				if err != nil {
					encounteredErrors = append(encounteredErrors, err)
				} else {
					retrievedApps = append(retrievedApps, app)
				}
			}

			if !slices.Equal(expectedApps, retrievedApps) {
				t.Errorf("Expected %v but got %v", expectedApps, retrievedApps)
			}

			if len(encounteredErrors) != 1 {
				t.Errorf("Expected to receive one error, but received %d", len(encounteredErrors))
			}
		})

		t.Run("it should start paging from specified page number when given", func(t *testing.T) {
			expectedApps := []onspring.App{
				{
					Href: "https://test.com",
					Id:   2,
					Name: "App",
				},
			}

			pageTwo := onspring.Page[onspring.App]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        []onspring.App{expectedApps[0]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "2" {
					jsonData, _ := json.Marshal(pageTwo)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedApps := []onspring.App{}

			for app, _ := range client.Apps.ListAll(t.Context(), onspring.ForPageNumber(2)) {
				retrievedApps = append(retrievedApps, app)
			}

			if !slices.Equal(expectedApps, retrievedApps) {
				t.Errorf("Expected %v but got %v", expectedApps, retrievedApps)
			}
		})

		t.Run("it should retrieve pages using specified page size when given", func(t *testing.T) {
			expectedApps := []onspring.App{
				{
					Href: "https://test.com",
					Id:   1,
					Name: "App",
				},
				{
					Href: "https://test.com",
					Id:   2,
					Name: "App",
				},
			}

			page := onspring.Page[onspring.App]{
				TotalPages:   1,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     2,
				Items:        expectedApps,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				if r.URL.Path != "/apps" {
					t.Errorf("Expected /apps endpoint, got %s", r.URL.Path)
				}

				pageSize := r.URL.Query().Get("pageSize")

				if pageSize == "2" {
					jsonData, _ := json.Marshal(page)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedApps := []onspring.App{}

			for app, _ := range client.Apps.ListAll(t.Context(), onspring.WithPageSize(2)) {
				retrievedApps = append(retrievedApps, app)
			}

			if !slices.Equal(expectedApps, retrievedApps) {
				t.Errorf("Expected %v but got %v", expectedApps, retrievedApps)
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Apps.GetMany(nilContext, []int{})

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

			_, err := client.Apps.GetMany(ctx, []int{})

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

			_, err := client.Apps.GetMany(t.Context(), []int{})

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

			_, err := invalidClient.Apps.GetMany(t.Context(), []int{})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the /apps/batch-get endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Apps.GetMany(t.Context(), []int{})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should return a batch of apps when the /apps/batch-get endpoint returns a 200 status code", func(t *testing.T) {
			apps := []onspring.App{
				{
					Href: "https://test.com",
					Id:   1,
					Name: "App",
				},
			}

			expectedBatch := onspring.AppBatch{
				Count: len(apps),
				Items: apps,
			}

			expectedIds := []int{apps[0].Id}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/apps/batch-get" {
					t.Errorf("Expected /apps/batch-get endpoint, got %s", r.URL.Path)
				}

				var ids []int
				err := json.NewDecoder(r.Body).Decode(&ids)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !slices.Equal(expectedIds, ids) {
					t.Errorf("Expected body to be %v but got %v", expectedIds, ids)
				}

				jsonData, _ := json.Marshal(expectedBatch)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			batch, _ := client.Apps.GetMany(t.Context(), expectedIds)

			if !reflect.DeepEqual(expectedBatch, batch) {
				t.Errorf("Expected %v but got %v", expectedBatch, batch)
			}
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Apps.Get(nilContext, 0)

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

			_, err := client.Apps.Get(ctx, 0)

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

			_, err := client.Apps.Get(t.Context(), 0)

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

			_, err := invalidClient.Apps.Get(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the /apps/id/:id endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Apps.Get(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should return an app if the /apps/id/:id endpoint returns a 200 status code", func(t *testing.T) {
			expectedApp := onspring.App{
				Href: "https://test.com",
				Id:   1,
				Name: "App",
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/apps/id/%d", expectedApp.Id)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				jsonData, _ := json.Marshal(expectedApp)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			app, _ := client.Apps.Get(t.Context(), expectedApp.Id)

			if !reflect.DeepEqual(expectedApp, app) {
				t.Errorf("Expected %v but got %v", expectedApp, app)
			}
		})
	})
}
