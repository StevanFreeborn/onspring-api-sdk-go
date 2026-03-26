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

func TestReports(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Reports.Get(nilContext, 1)

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

			_, err := client.Reports.Get(ctx, 1)

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

			_, err := client.Reports.Get(t.Context(), 1)

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

			_, err := invalidClient.Reports.Get(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Reports.Get(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the /reports/id/:reportId endpoint and return a report", func(t *testing.T) {
			reportId := 1
			recordId := 100

			expectedReport := onspring.ReportData{
				Columns: []string{"Column1", "Column2"},
				Rows: []onspring.ReportRow{
					{
						RecordId: &recordId,
						Cells:    []any{"value1", "value2"},
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/reports/id/%d", reportId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				jsonData, _ := json.Marshal(expectedReport)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			report, err := client.Reports.Get(t.Context(), reportId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedReport.Columns, report.Columns) {
				t.Errorf("Expected columns %v but got %v", expectedReport.Columns, report.Columns)
			}

			if len(report.Rows) != len(expectedReport.Rows) {
				t.Errorf("Expected %d rows but got %d", len(expectedReport.Rows), len(report.Rows))
			}
		})

		t.Run("it should include apiDataFormat query param when provided", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/reports/id/%d", 1)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				apiDataFormat := r.URL.Query().Get("apiDataFormat")

				if apiDataFormat != "Formatted" {
					t.Errorf("Expected query param apiDataFormat to be Formatted but got %s", apiDataFormat)
				}

				expectedReport := onspring.ReportData{
					Columns: []string{"Column1"},
					Rows:    []onspring.ReportRow{},
				}

				jsonData, _ := json.Marshal(expectedReport)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			_, _ = client.Reports.Get(t.Context(), 1, onspring.WithDataFormat("Formatted"))
		})

		t.Run("it should include dataType query param when provided", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/reports/id/%d", 1)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				dataType := r.URL.Query().Get("dataType")

				if dataType != "ChartData" {
					t.Errorf("Expected query param dataType to be ChartData but got %s", dataType)
				}

				expectedReport := onspring.ReportData{
					Columns: []string{"Column1"},
					Rows:    []onspring.ReportRow{},
				}

				jsonData, _ := json.Marshal(expectedReport)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			_, _ = client.Reports.Get(t.Context(), 1, onspring.WithDataType("ChartData"))
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Reports.List(nilContext, 1)

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

			_, err := client.Reports.List(ctx, 1)

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

			_, err := client.Reports.List(t.Context(), 1)

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

			_, err := invalidClient.Reports.List(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Reports.List(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the /reports/appId/:appId endpoint and return page of reports", func(t *testing.T) {
			appId := 1
			expectedPageNumber := 1
			expectedPageSize := 50

			expectedPage := onspring.Page[onspring.Report]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items: []onspring.Report{
					{
						AppId:       1,
						Id:          1,
						Name:        "Report",
						Description: "A test report",
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/reports/appId/%d", appId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
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

			page, err := client.Reports.List(t.Context(), appId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedPage, page) {
				t.Errorf("Expected %v but got %v", expectedPage, page)
			}
		})

		t.Run("it should perform a GET request to the /reports/appId/:appId endpoint with non-default paging information when provided", func(t *testing.T) {
			appId := 1
			expectedPageNumber := 2
			expectedPageSize := 1

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/reports/appId/%d", appId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
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

			_, _ = client.Reports.List(
				t.Context(),
				appId,
				onspring.ForPageNumber(expectedPageNumber),
				onspring.WithPageSize(expectedPageSize),
			)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("it should return an error if fails to retrieve any pages of reports", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			for _, err := range client.Reports.ListAll(t.Context(), 1) {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}
		})

		t.Run("it should return all the reports from multiple pages", func(t *testing.T) {
			expectedReports := []onspring.Report{
				{
					AppId:       1,
					Id:          1,
					Name:        "Report 1",
					Description: "First report",
				},
				{
					AppId:       1,
					Id:          2,
					Name:        "Report 2",
					Description: "Second report",
				},
			}

			pageOne := onspring.Page[onspring.Report]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        []onspring.Report{expectedReports[0]},
			}

			pageTwo := onspring.Page[onspring.Report]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        []onspring.Report{expectedReports[1]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
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

			retrievedReports := []onspring.Report{}

			for report, _ := range client.Reports.ListAll(t.Context(), 1) {
				retrievedReports = append(retrievedReports, report)
			}

			if !slices.Equal(expectedReports, retrievedReports) {
				t.Errorf("Expected %v but got %v", expectedReports, retrievedReports)
			}
		})

		t.Run("it should return reports and errors if some pages fail and some succeed", func(t *testing.T) {
			expectedReports := []onspring.Report{
				{
					AppId:       1,
					Id:          1,
					Name:        "Report 1",
					Description: "First report",
				},
			}

			pageOne := onspring.Page[onspring.Report]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        []onspring.Report{expectedReports[0]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
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

			retrievedReports := []onspring.Report{}
			encounteredErrors := []error{}

			for report, err := range client.Reports.ListAll(t.Context(), 1) {
				if err != nil {
					encounteredErrors = append(encounteredErrors, err)
				} else {
					retrievedReports = append(retrievedReports, report)
				}
			}

			if !slices.Equal(expectedReports, retrievedReports) {
				t.Errorf("Expected %v but got %v", expectedReports, retrievedReports)
			}

			if len(encounteredErrors) != 1 {
				t.Errorf("Expected to receive one error, but received %d", len(encounteredErrors))
			}
		})

		t.Run("it should start paging from specified page number when given", func(t *testing.T) {
			expectedReports := []onspring.Report{
				{
					AppId:       1,
					Id:          2,
					Name:        "Report 2",
					Description: "Second report",
				},
			}

			pageTwo := onspring.Page[onspring.Report]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        []onspring.Report{expectedReports[0]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "2" {
					jsonData, _ := json.Marshal(pageTwo)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedReports := []onspring.Report{}

			for report, _ := range client.Reports.ListAll(t.Context(), 1, onspring.ForPageNumber(2)) {
				retrievedReports = append(retrievedReports, report)
			}

			if !slices.Equal(expectedReports, retrievedReports) {
				t.Errorf("Expected %v but got %v", expectedReports, retrievedReports)
			}
		})

		t.Run("it should retrieve pages using specified page size when given", func(t *testing.T) {
			expectedReports := []onspring.Report{
				{
					AppId:       1,
					Id:          1,
					Name:        "Report 1",
					Description: "First report",
				},
				{
					AppId:       1,
					Id:          2,
					Name:        "Report 2",
					Description: "Second report",
				},
			}

			page := onspring.Page[onspring.Report]{
				TotalPages:   1,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     2,
				Items:        expectedReports,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				pageSize := r.URL.Query().Get("pageSize")

				if pageSize == "2" {
					jsonData, _ := json.Marshal(page)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedReports := []onspring.Report{}

			for report, _ := range client.Reports.ListAll(t.Context(), 1, onspring.WithPageSize(2)) {
				retrievedReports = append(retrievedReports, report)
			}

			if !slices.Equal(expectedReports, retrievedReports) {
				t.Errorf("Expected %v but got %v", expectedReports, retrievedReports)
			}
		})
	})
}
