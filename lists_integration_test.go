//go:build integration

package onspring_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestListsIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("Save", func(t *testing.T) {
		t.Run("should add a list item", func(t *testing.T) {
			listId := requireEnvInt(t, "TEST_LIST_ID")

			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item",
			}

			response, err := client.Lists.Save(ctx, listId, item)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.Id == "" {
				t.Error("Expected list item id to not be empty")
			}

			t.Cleanup(func() {
				_ = client.Lists.Delete(context.Background(), listId, response.Id)
			})
		})

		t.Run("should update a list item", func(t *testing.T) {
			listId := requireEnvInt(t, "TEST_LIST_ID")

			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item",
			}

			addResponse, err := client.Lists.Save(ctx, listId, item)

			if err != nil {
				t.Fatalf("Failed to add list item: %v", err)
			}

			t.Cleanup(func() {
				_ = client.Lists.Delete(context.Background(), listId, addResponse.Id)
			})

			updateItem := onspring.SaveListItemRequest{
				Id:   addResponse.Id,
				Name: "Updated Integration Test Item",
			}

			updateResponse, err := client.Lists.Save(ctx, listId, updateItem)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if updateResponse.Id == "" {
				t.Error("Expected list item id to not be empty")
			}
		})

		t.Run("should return a 401 error when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			listId := requireEnvInt(t, "TEST_LIST_ID")

			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item",
			}

			_, err := invalidClient.Lists.Save(ctx, listId, item)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when api key does not have access to the list", func(t *testing.T) {
			listIdNoAccess := requireEnvInt(t, "TEST_LIST_ID_NO_ACCESS")

			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item",
			}

			_, err := client.Lists.Save(ctx, listIdNoAccess, item)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when the list does not exist", func(t *testing.T) {
			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item",
			}

			_, err := client.Lists.Save(ctx, 0, item)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete a list item", func(t *testing.T) {
			listId := requireEnvInt(t, "TEST_LIST_ID")

			item := onspring.SaveListItemRequest{
				Name: "Integration Test Item To Delete",
			}

			response, err := client.Lists.Save(ctx, listId, item)

			if err != nil {
				t.Fatalf("Failed to add list item: %v", err)
			}

			err = client.Lists.Delete(ctx, listId, response.Id)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		t.Run("should return a 401 error when an invalid api key is used", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			listId := requireEnvInt(t, "TEST_LIST_ID")

			err := invalidClient.Lists.Delete(ctx, listId, "some-item-id")

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 error when api key does not have access to the list", func(t *testing.T) {
			listIdNoAccess := requireEnvInt(t, "TEST_LIST_ID_NO_ACCESS")
			listItemIdNoAccess := requireEnv(t, "TEST_LIST_ITEM_ID_NO_ACCESS")

			err := client.Lists.Delete(ctx, listIdNoAccess, listItemIdNoAccess)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 error when the list does not exist", func(t *testing.T) {
			err := client.Lists.Delete(ctx, 0, "a57e3d33-9195-4039-9ac0-c180013b043e")

			assertAPIError(t, err, http.StatusNotFound)
		})
	})
}
