//go:build integration

package onspring_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestRecordsIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		t.Run("should get a record", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			record, err := client.Records.Get(ctx, surveyId, recordId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if record.AppId != surveyId {
				t.Errorf("Expected appId %d, got %d", surveyId, record.AppId)
			}

			if record.RecordId != recordId {
				t.Errorf("Expected recordId %d, got %d", recordId, record.RecordId)
			}

			if len(record.FieldData) == 0 {
				t.Error("Expected field data to not be empty")
			}

			for _, field := range record.FieldData {
				if field.FieldId == 0 {
					t.Error("Expected field id to not be zero")
				}

				if field.Type == "" {
					t.Error("Expected field type to not be empty")
				}
			}
		})

		t.Run("should get a record when fieldIds and data format are passed as parameters", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			record, err := client.Records.Get(
				ctx,
				surveyId,
				recordId,
				onspring.WithFieldIds([]int{textFieldId}),
				onspring.WithRecordDataFormat("Formatted"),
			)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if record.AppId != surveyId {
				t.Errorf("Expected appId %d, got %d", surveyId, record.AppId)
			}

			if record.RecordId != recordId {
				t.Errorf("Expected recordId %d, got %d", recordId, record.RecordId)
			}

			if len(record.FieldData) == 0 {
				t.Error("Expected field data to not be empty")
			}
		})

		t.Run("should return a 401 error when an invalid API key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			_, err := invalidClient.Records.Get(ctx, surveyId, recordId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 404 error when an invalid record id is used", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			_, err := client.Records.Get(ctx, surveyId, 0)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("should get records", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			page, err := client.Records.List(ctx, surveyId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber == 0 {
				t.Error("Expected page number to not be zero")
			}

			if page.TotalRecords == 0 {
				t.Error("Expected total records to not be zero")
			}

			if len(page.Items) == 0 {
				t.Error("Expected items to not be empty")
			}

			for _, record := range page.Items {
				if record.AppId != surveyId {
					t.Errorf("Expected appId %d, got %d", surveyId, record.AppId)
				}

				if record.RecordId == 0 {
					t.Error("Expected recordId to not be zero")
				}
			}
		})

		t.Run("should get records when fieldIds, paging information, and data format are passed as parameters", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			page, err := client.Records.List(
				ctx,
				surveyId,
				onspring.WithFieldIds([]int{textFieldId}),
				onspring.WithRecordDataFormat("Formatted"),
				onspring.WithPaging(onspring.ForPageNumber(1), onspring.WithPageSize(1)),
			)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber != 1 {
				t.Errorf("Expected page number 1, got %d", page.PageNumber)
			}

			if page.PageSize != 1 {
				t.Errorf("Expected page size 1, got %d", page.PageSize)
			}

			if len(page.Items) != 1 {
				t.Errorf("Expected 1 item, got %d", len(page.Items))
			}
		})

		t.Run("should return a 401 error when an invalid API key is passed", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			_, err := invalidClient.Records.List(ctx, surveyId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when the api key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")

			_, err := client.Records.List(ctx, appIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("should iterate all records for an app", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			var records []onspring.Record

			for record, err := range client.Records.ListAll(ctx, surveyId, onspring.WithPaging(onspring.WithPageSize(1))) {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				records = append(records, record)
			}

			if len(records) == 0 {
				t.Error("Expected to iterate at least one record")
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("should get records", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			batch, err := client.Records.GetMany(ctx, onspring.GetManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: []int{recordId},
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if batch.Count == 0 {
				t.Error("Expected count to not be zero")
			}

			if len(batch.Items) == 0 {
				t.Error("Expected items to not be empty")
			}

			for _, record := range batch.Items {
				if record.AppId != surveyId {
					t.Errorf("Expected appId %d, got %d", surveyId, record.AppId)
				}

				if record.RecordId == 0 {
					t.Error("Expected recordId to not be zero")
				}
			}
		})

		t.Run("should get records when field ids and data format are passed as parameters", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			batch, err := client.Records.GetMany(ctx, onspring.GetManyRecordsRequest{
				AppId:      surveyId,
				RecordIds:  []int{recordId},
				FieldIds:   []int{textFieldId},
				DataFormat: "Formatted",
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if batch.Count == 0 {
				t.Error("Expected count to not be zero")
			}

			if len(batch.Items) == 0 {
				t.Error("Expected items to not be empty")
			}
		})

		t.Run("should return a 400 error if too many record ids are passed", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			ids := make([]int, 101)
			for i := range ids {
				ids[i] = i + 1
			}

			_, err := client.Records.GetMany(ctx, onspring.GetManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: ids,
			})

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error if the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			_, err := invalidClient.Records.GetMany(ctx, onspring.GetManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: []int{recordId},
			})

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error if the user does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			_, err := client.Records.GetMany(ctx, onspring.GetManyRecordsRequest{
				AppId:     appIdNoAccess,
				RecordIds: []int{recordId},
			})

			assertAPIError(t, err, http.StatusForbidden)
		})
	})

	t.Run("Query", func(t *testing.T) {
		t.Run("should return records", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			page, err := client.Records.Query(ctx, onspring.QueryRecordsRequest{
				AppId:  surveyId,
				Filter: filter,
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber == 0 {
				t.Error("Expected page number to not be zero")
			}

			if page.TotalRecords == 0 {
				t.Error("Expected total records to not be zero")
			}

			if len(page.Items) == 0 {
				t.Error("Expected items to not be empty")
			}

			for _, record := range page.Items {
				if record.AppId != surveyId {
					t.Errorf("Expected appId %d, got %d", surveyId, record.AppId)
				}

				if record.RecordId == 0 {
					t.Error("Expected recordId to not be zero")
				}
			}
		})

		t.Run("should return records when data format, paging information, and fields are passed as parameters", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			page, err := client.Records.Query(
				ctx,
				onspring.QueryRecordsRequest{
					AppId:      surveyId,
					Filter:     filter,
					FieldIds:   []int{textFieldId},
					DataFormat: "Formatted",
				},
				onspring.ForPageNumber(1),
				onspring.WithPageSize(1),
			)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber != 1 {
				t.Errorf("Expected page number 1, got %d", page.PageNumber)
			}

			if page.PageSize != 1 {
				t.Errorf("Expected page size 1, got %d", page.PageSize)
			}

			if len(page.Items) != 1 {
				t.Errorf("Expected 1 item, got %d", len(page.Items))
			}
		})

		t.Run("should return a 400 error if page size is invalid", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			_, err := client.Records.Query(
				ctx,
				onspring.QueryRecordsRequest{
					AppId:  surveyId,
					Filter: filter,
				},
				onspring.WithPageSize(1001),
			)

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error if the API key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			_, err := invalidClient.Records.Query(ctx, onspring.QueryRecordsRequest{
				AppId:  surveyId,
				Filter: filter,
			})

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error if the API key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			_, err := client.Records.Query(ctx, onspring.QueryRecordsRequest{
				AppId:  appIdNoAccess,
				Filter: filter,
			})

			assertAPIError(t, err, http.StatusForbidden)
		})
	})

	t.Run("QueryAll", func(t *testing.T) {
		t.Run("should iterate all records matching a query", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			autoNumberField := requireEnvInt(t, "TEST_SURVEY_AUTO_NUMBER_FIELD")
			filter := fmt.Sprintf("%d gt 0", autoNumberField)

			var records []onspring.Record

			for record, err := range client.Records.QueryAll(
				ctx,
				onspring.QueryRecordsRequest{
					AppId:  surveyId,
					Filter: filter,
				},
				onspring.WithPageSize(1),
			) {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				records = append(records, record)
			}

			if len(records) == 0 {
				t.Error("Expected to iterate at least one record")
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("should add a record when no record id is passed", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			response, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId: surveyId,
				Fields: map[string]any{
					fmt.Sprintf("%d", textFieldId): "integration test record",
				},
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.Id == 0 {
				t.Error("Expected record id to not be zero")
			}

			if response.Warnings == nil {
				t.Error("Expected warnings to not be nil")
			}

			t.Cleanup(func() {
				_ = client.Records.Delete(context.Background(), surveyId, response.Id)
			})
		})

		t.Run("should update a record when a record id is passed", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			recordId := addTestRecord(t, client, surveyId, textFieldId)

			response, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId:    surveyId,
				RecordId: &recordId,
				Fields: map[string]any{
					fmt.Sprintf("%d", textFieldId): "updated integration test record",
				},
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.Id == 0 {
				t.Error("Expected record id to not be zero")
			}
		})

		t.Run("should return a 400 error when field data is empty", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			_, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId:  surveyId,
				Fields: map[string]any{},
			})

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error when the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			_, err := invalidClient.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId: surveyId,
				Fields: map[string]any{
					fmt.Sprintf("%d", textFieldId): "test",
				},
			})

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when the api key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			_, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId: appIdNoAccess,
				Fields: map[string]any{
					fmt.Sprintf("%d", textFieldId): "test",
				},
			})

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when the record id is not found", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			recordId := 0

			_, err := client.Records.Save(ctx, onspring.SaveRecordRequest{
				AppId:    surveyId,
				RecordId: &recordId,
				Fields: map[string]any{
					fmt.Sprintf("%d", textFieldId): "test",
				},
			})

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete a record", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			recordId := addTestRecord(t, client, surveyId, textFieldId)

			err := client.Records.Delete(ctx, surveyId, recordId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		t.Run("should return a 401 error when the API key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			err := invalidClient.Records.Delete(ctx, surveyId, recordId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when the API key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")
			recordId := requireEnvInt(t, "TEST_SURVEY_RECORD_ID")

			err := client.Records.Delete(ctx, appIdNoAccess, recordId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when the record does not exist", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			err := client.Records.Delete(ctx, surveyId, 0)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("DeleteMany", func(t *testing.T) {
		t.Run("should delete records", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			recordId1 := addTestRecord(t, client, surveyId, textFieldId)
			recordId2 := addTestRecord(t, client, surveyId, textFieldId)

			err := client.Records.DeleteMany(ctx, onspring.DeleteManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: []int{recordId1, recordId2},
			})

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		t.Run("should return a 400 error when no record ids are provided", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			err := client.Records.DeleteMany(ctx, onspring.DeleteManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: []int{},
			})

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 error when an invalid API key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			err := invalidClient.Records.DeleteMany(ctx, onspring.DeleteManyRecordsRequest{
				AppId:     surveyId,
				RecordIds: []int{1},
			})

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when the API key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")

			err := client.Records.DeleteMany(ctx, onspring.DeleteManyRecordsRequest{
				AppId:     appIdNoAccess,
				RecordIds: []int{1},
			})

			assertAPIError(t, err, http.StatusForbidden)
		})
	})
}
