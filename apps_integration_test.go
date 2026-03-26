//go:build integration

package onspring_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestAppsIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		t.Run("should return an app", func(t *testing.T) {
			appId := requireEnvInt(t, "TEST_APP_ID")

			app, err := client.Apps.Get(ctx, appId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if app.Id != appId {
				t.Errorf("Expected app id %d, got %d", appId, app.Id)
			}

			if app.Name == "" {
				t.Error("Expected app name to not be empty")
			}

			if app.Href == "" {
				t.Error("Expected app href to not be empty")
			}
		})

		t.Run("should return a 401 error when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			appId := requireEnvInt(t, "TEST_APP_ID")

			_, err := invalidClient.Apps.Get(ctx, appId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when api key does not have access to the app", func(t *testing.T) {
			appIdNoAccess := requireEnvInt(t, "TEST_APP_ID_NO_ACCESS")

			_, err := client.Apps.Get(ctx, appIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when an app id cannot be found", func(t *testing.T) {
			_, err := client.Apps.Get(ctx, 0)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("should return a paged list of apps", func(t *testing.T) {
			page, err := client.Apps.List(ctx)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if page.PageNumber == 0 {
				t.Error("Expected page number to not be zero")
			}

			if page.PageSize == 0 {
				t.Error("Expected page size to not be zero")
			}

			if page.TotalPages == 0 {
				t.Error("Expected total pages to not be zero")
			}

			if page.TotalRecords == 0 {
				t.Error("Expected total records to not be zero")
			}

			if len(page.Items) == 0 {
				t.Error("Expected items to not be empty")
			}

			for _, app := range page.Items {
				if app.Id == 0 {
					t.Error("Expected app id to not be zero")
				}

				if app.Name == "" {
					t.Error("Expected app name to not be empty")
				}

				if app.Href == "" {
					t.Error("Expected app href to not be empty")
				}
			}
		})

		t.Run("should return a paged list of apps with correct page size and number when passed paging request", func(t *testing.T) {
			page, err := client.Apps.List(ctx, onspring.ForPageNumber(1), onspring.WithPageSize(1))

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
			_, err := client.Apps.List(ctx, onspring.WithPageSize(1001))

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)

			_, err := invalidClient.Apps.List(ctx)

			assertAPIError(t, err, http.StatusUnauthorized)
		})
	})

	t.Run("ListAll", func(t *testing.T) {
		t.Run("should iterate all apps", func(t *testing.T) {
			var apps []onspring.App

			for app, err := range client.Apps.ListAll(ctx, onspring.WithPageSize(1)) {
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				apps = append(apps, app)
			}

			if len(apps) == 0 {
				t.Error("Expected to iterate at least one app")
			}
		})
	})

	t.Run("GetMany", func(t *testing.T) {
		t.Run("should return a collection of apps", func(t *testing.T) {
			appIds := requireEnvIntSlice(t, "TEST_APP_IDS")

			batch, err := client.Apps.GetMany(ctx, appIds)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if batch.Count != len(appIds) {
				t.Errorf("Expected count %d, got %d", len(appIds), batch.Count)
			}

			if len(batch.Items) != len(appIds) {
				t.Errorf("Expected %d items, got %d", len(appIds), len(batch.Items))
			}

			for _, app := range batch.Items {
				if app.Id == 0 {
					t.Error("Expected app id to not be zero")
				}

				if app.Name == "" {
					t.Error("Expected app name to not be empty")
				}

				if app.Href == "" {
					t.Error("Expected app href to not be empty")
				}
			}
		})

		t.Run("should return a 401 error when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			appIds := requireEnvIntSlice(t, "TEST_APP_IDS")

			_, err := invalidClient.Apps.GetMany(ctx, appIds)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when api key does not have access to the apps", func(t *testing.T) {
			appIdsNoAccess := requireEnvIntSlice(t, "TEST_APP_IDS_NO_ACCESS")

			_, err := client.Apps.GetMany(ctx, appIdsNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})
	})
}
