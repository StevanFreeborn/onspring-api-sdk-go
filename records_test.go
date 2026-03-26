package onspring_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestRecords(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Records.Get(nilContext, 1, 1)

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

			_, err := client.Records.Get(ctx, 1, 1)

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

			_, err := client.Records.Get(t.Context(), 1, 1)

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

			_, err := invalidClient.Records.Get(t.Context(), 1, 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Records.Get(t.Context(), 1, 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the correct endpoint and return a record", func(t *testing.T) {
			appId := 1
			recordId := 2

			expectedRecord := onspring.Record{
				AppId:    appId,
				RecordId: recordId,
				FieldData: []onspring.RecordFieldValue{
					{
						Type:    "String",
						FieldId: 1,
						Value:   "test value",
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/records/appId/%d/recordId/%d", appId, recordId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				jsonData, _ := json.Marshal(expectedRecord)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			record, err := client.Records.Get(t.Context(), appId, recordId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if record.AppId != expectedRecord.AppId {
				t.Errorf("Expected AppId %d but got %d", expectedRecord.AppId, record.AppId)
			}

			if record.RecordId != expectedRecord.RecordId {
				t.Errorf("Expected RecordId %d but got %d", expectedRecord.RecordId, record.RecordId)
			}

			if len(record.FieldData) != len(expectedRecord.FieldData) {
				t.Errorf("Expected %d field data items but got %d", len(expectedRecord.FieldData), len(record.FieldData))
			}
		})

		t.Run("it should include fieldIds query param when provided", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				fieldIds := r.URL.Query().Get("fieldIds")

				if fieldIds != "1,2,3" {
					t.Errorf("Expected query param fieldIds to be 1,2,3 but got %s", fieldIds)
				}

				jsonData, _ := json.Marshal(onspring.Record{})

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			_, _ = client.Records.Get(t.Context(), 1, 1, onspring.WithFieldIds([]int{1, 2, 3}))
		})

		t.Run("it should include dataFormat query param when provided", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				dataFormat := r.URL.Query().Get("dataFormat")

				if dataFormat != "Formatted" {
					t.Errorf("Expected query param dataFormat to be Formatted but got %s", dataFormat)
				}

				jsonData, _ := json.Marshal(onspring.Record{})

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			_, _ = client.Records.Get(t.Context(), 1, 1, onspring.WithRecordDataFormat("Formatted"))
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Records.List(nilContext, 1)

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

			_, err := client.Records.List(ctx, 1)

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

			_, err := client.Records.List(t.Context(), 1)

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

			_, err := invalidClient.Records.List(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Records.List(t.Context(), 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the correct endpoint and return page of records", func(t *testing.T) {
			appId := 1
			expectedPageNumber := 1
			expectedPageSize := 50

			expectedPage := onspring.Page[onspring.Record]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items: []onspring.Record{
					{
						AppId:     appId,
						RecordId:  1,
						FieldData: []onspring.RecordFieldValue{},
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/records/appId/%d", appId)

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

			page, err := client.Records.List(t.Context(), appId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedPage, page) {
				t.Errorf("Expected %v but got %v", expectedPage, page)
			}
		})

		t.Run("it should include paging and record options when provided", func(t *testing.T) {
			appId := 1

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber != "2" {
					t.Errorf("Expected query param pageNumber to be 2 but got %s", pageNumber)
				}

				pageSize := r.URL.Query().Get("pageSize")

				if pageSize != "10" {
					t.Errorf("Expected query param pageSize to be 10 but got %s", pageSize)
				}

				fieldIds := r.URL.Query().Get("fieldIds")

				if fieldIds != "1,2" {
					t.Errorf("Expected query param fieldIds to be 1,2 but got %s", fieldIds)
				}

				dataFormat := r.URL.Query().Get("dataFormat")

				if dataFormat != "Formatted" {
					t.Errorf("Expected query param dataFormat to be Formatted but got %s", dataFormat)
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"pageNumber":2,"pageSize":10,"totalPages":2,"totalRecords":2,"items":[]}`))
			})

			_, _ = client.Records.List(
				t.Context(),
				appId,
				onspring.WithFieldIds([]int{1, 2}),
				onspring.WithRecordDataFormat("Formatted"),
				onspring.WithPaging(onspring.ForPageNumber(2), onspring.WithPageSize(10)),
			)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("it should return an error if fails to retrieve any pages of records", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			for _, err := range client.Records.ListAll(t.Context(), 1) {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}
		})

		t.Run("it should return all the records from multiple pages", func(t *testing.T) {
			expectedRecords := []onspring.Record{
				{AppId: 1, RecordId: 1, FieldData: []onspring.RecordFieldValue{}},
				{AppId: 1, RecordId: 2, FieldData: []onspring.RecordFieldValue{}},
			}

			pageOne := onspring.Page[onspring.Record]{
				TotalPages: 2, TotalRecords: 2, PageNumber: 1, PageSize: 1,
				Items: []onspring.Record{expectedRecords[0]},
			}

			pageTwo := onspring.Page[onspring.Record]{
				TotalPages: 2, TotalRecords: 2, PageNumber: 2, PageSize: 1,
				Items: []onspring.Record{expectedRecords[1]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

			retrievedRecords := []onspring.Record{}

			for record, _ := range client.Records.ListAll(t.Context(), 1) {
				retrievedRecords = append(retrievedRecords, record)
			}

			if !reflect.DeepEqual(expectedRecords, retrievedRecords) {
				t.Errorf("Expected %v but got %v", expectedRecords, retrievedRecords)
			}
		})

		t.Run("it should return records and errors if some pages fail and some succeed", func(t *testing.T) {
			expectedRecords := []onspring.Record{
				{AppId: 1, RecordId: 1, FieldData: []onspring.RecordFieldValue{}},
			}

			pageOne := onspring.Page[onspring.Record]{
				TotalPages: 2, TotalRecords: 2, PageNumber: 1, PageSize: 1,
				Items: []onspring.Record{expectedRecords[0]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

			retrievedRecords := []onspring.Record{}
			encounteredErrors := []error{}

			for record, err := range client.Records.ListAll(t.Context(), 1) {
				if err != nil {
					encounteredErrors = append(encounteredErrors, err)
				} else {
					retrievedRecords = append(retrievedRecords, record)
				}
			}

			if !reflect.DeepEqual(expectedRecords, retrievedRecords) {
				t.Errorf("Expected %v but got %v", expectedRecords, retrievedRecords)
			}

			if len(encounteredErrors) != 1 {
				t.Errorf("Expected to receive one error, but received %d", len(encounteredErrors))
			}
		})

		t.Run("it should pass record options through to each page request", func(t *testing.T) {
			page := onspring.Page[onspring.Record]{
				TotalPages: 1, TotalRecords: 0, PageNumber: 1, PageSize: 1,
				Items: []onspring.Record{},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				fieldIds := r.URL.Query().Get("fieldIds")

				if fieldIds != "1,2" {
					t.Errorf("Expected query param fieldIds to be 1,2 but got %s", fieldIds)
				}

				jsonData, _ := json.Marshal(page)
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			for range client.Records.ListAll(t.Context(), 1, onspring.WithFieldIds([]int{1, 2})) {
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Records.GetMany(nilContext, onspring.GetManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			_, err := client.Records.GetMany(ctx, onspring.GetManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			_, err := client.Records.GetMany(t.Context(), onspring.GetManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			_, err := invalidClient.Records.GetMany(t.Context(), onspring.GetManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Records.GetMany(t.Context(), onspring.GetManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a POST request to the /records/batch-get endpoint and return a batch of records", func(t *testing.T) {
			request := onspring.GetManyRecordsRequest{
				AppId:      1,
				RecordIds:  []int{1, 2},
				FieldIds:   []int{1},
				DataFormat: "Formatted",
			}

			expectedBatch := onspring.RecordBatch{
				Count: 2,
				Items: []onspring.Record{
					{AppId: 1, RecordId: 1, FieldData: []onspring.RecordFieldValue{}},
					{AppId: 1, RecordId: 2, FieldData: []onspring.RecordFieldValue{}},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/records/batch-get" {
					t.Errorf("Expected /records/batch-get endpoint, got %s", r.URL.Path)
				}

				var body onspring.GetManyRecordsRequest
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !reflect.DeepEqual(request, body) {
					t.Errorf("Expected body to be %v but got %v", request, body)
				}

				jsonData, _ := json.Marshal(expectedBatch)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			batch, err := client.Records.GetMany(t.Context(), request)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedBatch, batch) {
				t.Errorf("Expected %v but got %v", expectedBatch, batch)
			}
		})
	})

	t.Run("Query", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Records.Query(nilContext, onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"})

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

			_, err := client.Records.Query(ctx, onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"})

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

			_, err := client.Records.Query(t.Context(), onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"})

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

			_, err := invalidClient.Records.Query(t.Context(), onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Records.Query(t.Context(), onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a POST request to the /records/query endpoint and return page of records", func(t *testing.T) {
			request := onspring.QueryRecordsRequest{
				AppId:      1,
				Filter:     "field eq 'value'",
				FieldIds:   []int{1},
				DataFormat: "Formatted",
			}

			expectedPage := onspring.Page[onspring.Record]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items: []onspring.Record{
					{AppId: 1, RecordId: 1, FieldData: []onspring.RecordFieldValue{}},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/records/query" {
					t.Errorf("Expected /records/query endpoint, got %s", r.URL.Path)
				}

				var body onspring.QueryRecordsRequest
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !reflect.DeepEqual(request, body) {
					t.Errorf("Expected body to be %v but got %v", request, body)
				}

				jsonData, _ := json.Marshal(expectedPage)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			page, err := client.Records.Query(t.Context(), request)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedPage, page) {
				t.Errorf("Expected %v but got %v", expectedPage, page)
			}
		})

		t.Run("it should include paging query params when provided", func(t *testing.T) {
			expectedPageNumber := 2
			expectedPageSize := 10

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				pageNumber := r.URL.Query().Get("pageNumber")
				pageSize := r.URL.Query().Get("pageSize")

				if pageNumber != strconv.Itoa(expectedPageNumber) {
					t.Errorf("Expected query param pageNumber to be %d but got %s", expectedPageNumber, pageNumber)
				}

				if pageSize != strconv.Itoa(expectedPageSize) {
					t.Errorf("Expected query param pageSize to be %d but got %s", expectedPageSize, pageSize)
				}

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"pageNumber":2,"pageSize":10,"totalPages":2,"totalRecords":2,"items":[]}`))
			})

			_, _ = client.Records.Query(
				t.Context(),
				onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"},
				onspring.ForPageNumber(expectedPageNumber),
				onspring.WithPageSize(expectedPageSize),
			)
		})
	})

	t.Run("QueryAll", func(t *testing.T) {
		t.Run("it should return an error if fails to retrieve any pages", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			for _, err := range client.Records.QueryAll(t.Context(), onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"}) {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}
		})

		t.Run("it should return all the records from multiple pages", func(t *testing.T) {
			expectedRecords := []onspring.Record{
				{AppId: 1, RecordId: 1, FieldData: []onspring.RecordFieldValue{}},
				{AppId: 1, RecordId: 2, FieldData: []onspring.RecordFieldValue{}},
			}

			pageOne := onspring.Page[onspring.Record]{
				TotalPages: 2, TotalRecords: 2, PageNumber: 1, PageSize: 1,
				Items: []onspring.Record{expectedRecords[0]},
			}

			pageTwo := onspring.Page[onspring.Record]{
				TotalPages: 2, TotalRecords: 2, PageNumber: 2, PageSize: 1,
				Items: []onspring.Record{expectedRecords[1]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
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

			retrievedRecords := []onspring.Record{}

			for record, _ := range client.Records.QueryAll(t.Context(), onspring.QueryRecordsRequest{AppId: 1, Filter: "field eq 'value'"}) {
				retrievedRecords = append(retrievedRecords, record)
			}

			if !reflect.DeepEqual(expectedRecords, retrievedRecords) {
				t.Errorf("Expected %v but got %v", expectedRecords, retrievedRecords)
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Records.Save(nilContext, onspring.SaveRecordRequest{AppId: 1, Fields: map[string]any{}})

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

			_, err := client.Records.Save(ctx, onspring.SaveRecordRequest{AppId: 1, Fields: map[string]any{}})

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

			_, err := client.Records.Save(t.Context(), onspring.SaveRecordRequest{AppId: 1, Fields: map[string]any{}})

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

			_, err := invalidClient.Records.Save(t.Context(), onspring.SaveRecordRequest{AppId: 1, Fields: map[string]any{}})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Records.Save(t.Context(), onspring.SaveRecordRequest{AppId: 1, Fields: map[string]any{}})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a PUT request to the /records endpoint and return a response when successful", func(t *testing.T) {
			recordId := 1
			request := onspring.SaveRecordRequest{
				AppId:    1,
				RecordId: &recordId,
				Fields:   map[string]any{"1": "value"},
			}

			expectedResponse := onspring.SaveRecordResponse{
				Id:       1,
				Warnings: []string{"warning"},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("Expected PUT method, got %s", r.Method)
				}

				if r.URL.Path != "/records" {
					t.Errorf("Expected /records endpoint, got %s", r.URL.Path)
				}

				var body onspring.SaveRecordRequest
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !reflect.DeepEqual(request, body) {
					t.Errorf("Expected body to be %v but got %v", request, body)
				}

				jsonData, _ := json.Marshal(expectedResponse)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			response, err := client.Records.Save(t.Context(), request)

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

			err := client.Records.Delete(nilContext, 1, 1)

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

			err := client.Records.Delete(ctx, 1, 1)

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

			err := client.Records.Delete(t.Context(), 1, 1)

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

			err := invalidClient.Records.Delete(t.Context(), 1, 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			err := client.Records.Delete(t.Context(), 1, 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a DELETE request to the correct endpoint and return no error if receives 204 status code", func(t *testing.T) {
			appId := 1
			recordId := 2

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("Expected DELETE method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/records/appId/%d/recordId/%d", appId, recordId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusNoContent)
			})

			err := client.Records.Delete(t.Context(), appId, recordId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	})

	t.Run("DeleteMany", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			var nilContext context.Context = nil

			err := client.Records.DeleteMany(nilContext, onspring.DeleteManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			err := client.Records.DeleteMany(ctx, onspring.DeleteManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			err := client.Records.DeleteMany(t.Context(), onspring.DeleteManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

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

			err := invalidClient.Records.DeleteMany(t.Context(), onspring.DeleteManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			err := client.Records.DeleteMany(t.Context(), onspring.DeleteManyRecordsRequest{AppId: 1, RecordIds: []int{1}})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a POST request to the /records/batch-delete endpoint and return no error if receives 204 status code", func(t *testing.T) {
			request := onspring.DeleteManyRecordsRequest{
				AppId:     1,
				RecordIds: []int{1, 2},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/records/batch-delete" {
					t.Errorf("Expected /records/batch-delete endpoint, got %s", r.URL.Path)
				}

				var body onspring.DeleteManyRecordsRequest
				err := json.NewDecoder(r.Body).Decode(&body)

				if err != nil {
					t.Errorf("Expected to decode request body, but got error: %v", err)
				}

				if !reflect.DeepEqual(request, body) {
					t.Errorf("Expected body to be %v but got %v", request, body)
				}

				w.WriteHeader(http.StatusNoContent)
			})

			err := client.Records.DeleteMany(t.Context(), request)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	})
}
