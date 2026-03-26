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

func TestFields(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Fields.Get(nilContext, 0)

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

			_, err := client.Fields.List(ctx, 0)

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

			_, err := client.Fields.Get(t.Context(), 0)

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

			_, err := invalidClient.Fields.Get(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the /fields/id/:id endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Fields.Get(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should return a field if the /fields/id/:id endpoint returns a 200 status code", func(t *testing.T) {
			expectedField := onspring.Field{
				Id:         1,
				AppId:      1,
				Name:       "Field",
				Type:       "Text",
				Status:     "Enabled",
				IsRequired: true,
				IsUnique:   false,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/id/%d", expectedField.Id)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				jsonData, _ := json.Marshal(expectedField)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			field, _ := client.Fields.Get(t.Context(), expectedField.Id)

			if !reflect.DeepEqual(expectedField, field) {
				t.Errorf("Expected %v but got %v", expectedField, field)
			}
		})

	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Run("it should unmarshal a Formula field correctly", func(t *testing.T) {
			jsonStr := `{
				"id": 1,
				"appId": 10,
				"name": "Calc Field",
				"type": "Formula",
				"status": "Enabled",
				"isRequired": true,
				"isUnique": false,
				"outputType": "Number",
				"values": [
					{"id": "a1", "name": "Value1", "sortOrder": 1, "numericValue": 1, "color": "#fff"},
					{"id": "b2", "name": "Value2", "sortOrder": 2, "numericValue": 2, "color": "#000"}
				]
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if field.Type != "Formula" {
				t.Errorf("Expected Type 'Formula', got %s", field.Type)
			}

			formulaField, ok := field.TypeData.(onspring.FormulaField)

			if !ok {
				t.Fatalf("Expected TypeData to be FormulaField, got %T", field.TypeData)
			}

			expectedFormulaField := onspring.FormulaField{
				OutputType: "Number",
				Values: []onspring.ListValue{
					{Id: "a1", Name: "Value1", SortOrder: 1, NumericValue: 1, Color: "#fff"},
					{Id: "b2", Name: "Value2", SortOrder: 2, NumericValue: 2, Color: "#000"},
				},
			}

			if !reflect.DeepEqual(formulaField, expectedFormulaField) {
				t.Errorf("Expected FormulaField %+v, got %+v", expectedFormulaField, formulaField)
			}
		})

		t.Run("it should unmarshal a Reference field correctly", func(t *testing.T) {
			jsonStr := `{
				"id": 2,
				"appId": 20,
				"name": "Ref Field",
				"type": "Reference",
				"status": "Enabled",
				"isRequired": false,
				"isUnique": true,
				"multiplicity": "OneToOne",
				"referenceAppId": "123"
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if field.Type != "Reference" {
				t.Errorf("Expected Type 'Reference', got %s", field.Type)
			}

			referenceField, ok := field.TypeData.(onspring.ReferenceField)

			if !ok {
				t.Fatalf("Expected TypeData to be ReferenceField, got %T", field.TypeData)
			}

			expectedReferenceField := onspring.ReferenceField{
				Multiplicity:   "OneToOne",
				ReferenceAppId: "123",
			}

			if !reflect.DeepEqual(referenceField, expectedReferenceField) {
				t.Errorf("Expected ReferenceField %+v, got %+v", expectedReferenceField, referenceField)
			}
		})

		t.Run("it should unmarshal a List field correctly", func(t *testing.T) {
			jsonStr := `{
				"id": 3,
				"appId": 30,
				"name": "List Field",
				"type": "List",
				"status": "Enabled",
				"isRequired": false,
				"isUnique": false,
				"multiplicity": "MultiSelect",
				"values": [
					{"id": "aaa", "name": "OptionA", "sortOrder": 1, "numericValue": 0, "color": "#ffffff"},
					{"id": "bbb", "name": "OptionB", "sortOrder": 2, "numericValue": 0, "color": "#000000"}
				],
				"listId": 456
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if field.Type != "List" {
				t.Errorf("Expected Type 'List', got %s", field.Type)
			}

			listField, ok := field.TypeData.(onspring.ListField)

			if !ok {
				t.Fatalf("Expected TypeData to be ListField, got %T", field.TypeData)
			}

			expectedListField := onspring.ListField{
				Multiplicity: "MultiSelect",
				Values: []onspring.ListValue{
					{Id: "aaa", Name: "OptionA", SortOrder: 1, NumericValue: 0, Color: "#ffffff"},
					{Id: "bbb", Name: "OptionB", SortOrder: 2, NumericValue: 0, Color: "#000000"},
				},
				ListId: 456,
			}

			if !reflect.DeepEqual(listField, expectedListField) {
				t.Errorf("Expected ListField %+v, got %+v", expectedListField, listField)
			}
		})

		t.Run("it should handle unknown field types without TypeData", func(t *testing.T) {
			jsonStr := `{
				"id": 4,
				"appId": 40,
				"name": "Text Field",
				"type": "Text",
				"status": "Enabled",
				"isRequired": true,
				"isUnique": false
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err != nil {
				t.Fatalf("UnmarshalJSON failed: %v", err)
			}

			if field.Type != "Text" {
				t.Errorf("Expected Type 'Text', got %s", field.Type)
			}

			if field.TypeData != nil {
				t.Errorf("Expected TypeData to be nil for 'Text' type, got %+v", field.TypeData)
			}
		})

		t.Run("it should return an error for malformed JSON for main field struct", func(t *testing.T) {
			jsonStr := `{
				"id": "invalid",
				"appId": 10,
				"name": "Calc Field",
				"type": "Formula"
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err == nil {
				t.Fatalf("Expected UnmarshalJSON to fail for malformed main struct JSON, got nil")
			}
		})

		t.Run("it should return an error for malformed JSON for TypeData (Formula)", func(t *testing.T) {
			jsonStr := `{
				"id": 1,
				"appId": 10,
				"name": "Calc Field",
				"type": "Formula",
				"outputType": 123,
				"values": ["1", "2", "3"]
			}`

			var field onspring.Field

			err := json.Unmarshal([]byte(jsonStr), &field)

			if err == nil {
				t.Fatalf("Expected UnmarshalJSON to fail for malformed Formula TypeData JSON, got nil")
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Fields.GetMany(nilContext, []int{})

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

			_, err := client.Fields.GetMany(ctx, []int{})

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

			_, err := client.Fields.GetMany(t.Context(), []int{})

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

			_, err := invalidClient.Fields.GetMany(t.Context(), []int{})

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the /fields/batch-get endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Fields.GetMany(t.Context(), []int{})

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should return a batch of fields when the /fields/batch-get endpoint returns a 200 status code", func(t *testing.T) {
			fields := []onspring.Field{
				{
					Id:         1,
					AppId:      1,
					Name:       "Field 1",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
			}

			expectedBatch := onspring.FieldBatch{
				Count: len(fields),
				Items: fields,
			}

			expectedIds := []int{fields[0].Id}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/fields/batch-get" {
					t.Errorf("Expected /fields/batch-get endpoint, got %s", r.URL.Path)
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

			batch, _ := client.Fields.GetMany(t.Context(), expectedIds)

			if !reflect.DeepEqual(expectedBatch, batch) {
				t.Errorf("Expected %v but got %v", expectedBatch, batch)
			}
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Fields.List(nilContext, 0)

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

			_, err := client.Fields.List(ctx, 0)

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

			_, err := client.Fields.List(t.Context(), 0)

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

			_, err := invalidClient.Fields.List(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return a page of fields if the /fields/appId/:appId endpoint returns a 200 status code", func(t *testing.T) {
			appId := 1
			expectedPageNumber := 1
			expectedPageSize := 50

			expectedPage := onspring.Page[onspring.Field]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     50,
				Items: []onspring.Field{
					{
						Id:         1,
						AppId:      appId,
						Name:       "Field",
						Type:       "Text",
						Status:     "Enabled",
						IsRequired: true,
						IsUnique:   false,
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)

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

			page, err := client.Fields.List(t.Context(), appId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedPage, page) {
				t.Errorf("Expected %v but got %v", expectedPage, page)
			}
		})

		t.Run("it should perform a GET request to the /fields/appId/:appId endpoint with non-default paging information when provided", func(t *testing.T) {
			appId := 1
			expectedPageNumber := 2
			expectedPageSize := 1

			expectedPage := onspring.Page[onspring.Field]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items: []onspring.Field{
					{
						Id:         1,
						AppId:      appId,
						Name:       "Field",
						Type:       "Text",
						Status:     "Enabled",
						IsRequired: true,
						IsUnique:   false,
					},
				},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)

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

			_, _ = client.Fields.List(
				t.Context(),
				appId,
				onspring.ForPageNumber(expectedPageNumber),
				onspring.WithPageSize(expectedPageSize),
			)
		})

		t.Run("it should return an error if the /fields/appId/:appId endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Fields.List(t.Context(), 0)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("it should return an error if fails to retrieve any pages of fields", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			for _, err := range client.Fields.ListAll(t.Context(), 0) {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			}
		})

		t.Run("it should return all the fields from multiple pages", func(t *testing.T) {
			appId := 1
			expectedFields := []onspring.Field{
				{
					Id:         1,
					AppId:      appId,
					Name:       "Field 1",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
				{
					Id:         2,
					AppId:      appId,
					Name:       "Field 2",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
			}

			pageOne := onspring.Page[onspring.Field]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        []onspring.Field{expectedFields[0]},
			}

			pageTwo := onspring.Page[onspring.Field]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        []onspring.Field{expectedFields[1]},
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
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

			retrievedFields := []onspring.Field{}

			for field, _ := range client.Fields.ListAll(t.Context(), appId) {
				retrievedFields = append(retrievedFields, field)
			}

			if !slices.Equal(expectedFields, retrievedFields) {
				t.Errorf("Expected %v but got %v", expectedFields, retrievedFields)
			}
		})

		t.Run("it should return fields and errors if some pages fail and some succeed", func(t *testing.T) {
			appId := 1
			expectedFields := []onspring.Field{
				{
					Id:         1,
					AppId:      appId,
					Name:       "Field 1",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
				{
					Id:         2,
					AppId:      appId,
					Name:       "Field 2",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
			}

			page := onspring.Page[onspring.Field]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   1,
				PageSize:     1,
				Items:        expectedFields,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "1" {
					jsonData, _ := json.Marshal(page)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}

				if pageNumber == "2" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			})

			retrievedFields := []onspring.Field{}
			encounteredErrors := []error{}

			for field, err := range client.Fields.ListAll(t.Context(), appId) {
				if err != nil {
					encounteredErrors = append(encounteredErrors, err)
				} else {
					retrievedFields = append(retrievedFields, field)
				}
			}

			if !slices.Equal(expectedFields, retrievedFields) {
				t.Errorf("Expected %v but got %v", expectedFields, retrievedFields)
			}

			if len(encounteredErrors) != 1 {
				t.Errorf("Expected to receive one error, but received %d", len(encounteredErrors))
			}
		})

		t.Run("it should start paging from specified page number when given", func(t *testing.T) {
			appId := 1
			expectedFields := []onspring.Field{
				{
					Id:         1,
					AppId:      appId,
					Name:       "Field 1",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
				{
					Id:         2,
					AppId:      appId,
					Name:       "Field 2",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
			}

			pageTwo := onspring.Page[onspring.Field]{
				TotalPages:   2,
				TotalRecords: 2,
				PageNumber:   2,
				PageSize:     1,
				Items:        expectedFields,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				pageNumber := r.URL.Query().Get("pageNumber")

				if pageNumber == "2" {
					jsonData, _ := json.Marshal(pageTwo)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedFields := []onspring.Field{}

			for field, _ := range client.Fields.ListAll(t.Context(), appId, onspring.ForPageNumber(2)) {
				retrievedFields = append(retrievedFields, field)
			}

			if !slices.Equal(expectedFields, retrievedFields) {
				t.Errorf("Expected %v but got %v", expectedFields, retrievedFields)
			}
		})

		t.Run("it should retrieve pages using specified page size when given", func(t *testing.T) {
			appId := 1
			expectedFields := []onspring.Field{
				{
					Id:         1,
					AppId:      appId,
					Name:       "Field 1",
					Type:       "Text",
					Status:     "Enabled",
					IsRequired: true,
					IsUnique:   false,
				},
			}

			page := onspring.Page[onspring.Field]{
				TotalPages:   1,
				TotalRecords: 1,
				PageNumber:   1,
				PageSize:     1,
				Items:        expectedFields,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/fields/appId/%d", appId)
				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				pageSize := r.URL.Query().Get("pageSize")

				if pageSize == "1" {
					jsonData, _ := json.Marshal(page)

					w.WriteHeader(http.StatusOK)
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(jsonData)
				}
			})

			retrievedFields := []onspring.Field{}

			for field, _ := range client.Fields.ListAll(t.Context(), appId, onspring.WithPageSize(1)) {
				retrievedFields = append(retrievedFields, field)
			}

			if !slices.Equal(expectedFields, retrievedFields) {
				t.Errorf("Expected %v but got %v", expectedFields, retrievedFields)
			}
		})
	})
}
