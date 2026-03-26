//go:build integration

package onspring_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestFieldsIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		t.Run("should return a field", func(t *testing.T) {
			fieldId := requireEnvInt(t, "TEST_FIELD_ID")

			field, err := client.Fields.Get(ctx, fieldId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if field.Id != fieldId {
				t.Errorf("Expected field id %d, got %d", fieldId, field.Id)
			}

			if field.Name == "" {
				t.Error("Expected field name to not be empty")
			}

			if field.AppId == 0 {
				t.Error("Expected field appId to not be zero")
			}

			if field.Type == "" {
				t.Error("Expected field type to not be empty")
			}

			if field.Status == "" {
				t.Error("Expected field status to not be empty")
			}
		})

		t.Run("should return a 401 error when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			fieldId := requireEnvInt(t, "TEST_FIELD_ID")

			_, err := invalidClient.Fields.Get(ctx, fieldId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when api key does not have access to the field", func(t *testing.T) {
			fieldIdNoAccess := requireEnvInt(t, "TEST_FIELD_ID_NO_ACCESS")

			_, err := client.Fields.Get(ctx, fieldIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when the field does not exist", func(t *testing.T) {
			_, err := client.Fields.Get(ctx, 0)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("should return a paged list of fields", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			page, err := client.Fields.List(ctx, surveyId)

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

			for _, field := range page.Items {
				if field.Id == 0 {
					t.Error("Expected field id to not be zero")
				}

				if field.Name == "" {
					t.Error("Expected field name to not be empty")
				}

				if field.AppId == 0 {
					t.Error("Expected field appId to not be zero")
				}

				if field.Type == "" {
					t.Error("Expected field type to not be empty")
				}

				if field.Status == "" {
					t.Error("Expected field status to not be empty")
				}
			}
		})

		t.Run("should return a paged list of fields with correct page size and number when passed paging request", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			page, err := client.Fields.List(ctx, surveyId, onspring.ForPageNumber(1), onspring.WithPageSize(1))

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

		t.Run("should return a 400 response when an invalid page size is used", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			_, err := client.Fields.List(ctx, surveyId, onspring.WithPageSize(1001))

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			_, err := invalidClient.Fields.List(ctx, surveyId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when api key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")

			_, err := client.Fields.List(ctx, appIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("should iterate all fields for an app", func(t *testing.T) {
			surveyId := requireEnvInt(t, "TEST_SURVEY_ID")

			var fields []onspring.Field

			for field, err := range client.Fields.ListAll(ctx, surveyId, onspring.WithPageSize(1)) {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				fields = append(fields, field)
			}

			if len(fields) == 0 {
				t.Error("Expected to iterate at least one field")
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("should return a collection of fields", func(t *testing.T) {
			fieldIds := requireEnvIntSlice(t, "TEST_FIELD_IDS")

			batch, err := client.Fields.GetMany(ctx, fieldIds)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if batch.Count != len(fieldIds) {
				t.Errorf("Expected count %d, got %d", len(fieldIds), batch.Count)
			}

			if len(batch.Items) != len(fieldIds) {
				t.Errorf("Expected %d items, got %d", len(fieldIds), len(batch.Items))
			}

			for _, field := range batch.Items {
				if field.Id == 0 {
					t.Error("Expected field id to not be zero")
				}

				if field.Name == "" {
					t.Error("Expected field name to not be empty")
				}
			}
		})

		t.Run("should return a 401 response when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			fieldIds := requireEnvIntSlice(t, "TEST_FIELD_IDS")

			_, err := invalidClient.Fields.GetMany(ctx, fieldIds)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when the api key does not have access to any of the requested fields", func(t *testing.T) {
			fieldIdsNoAccess := requireEnvIntSlice(t, "TEST_FIELD_IDS_NO_ACCESS")

			_, err := client.Fields.GetMany(ctx, fieldIdsNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})
	})
}
