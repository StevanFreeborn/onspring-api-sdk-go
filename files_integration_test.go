//go:build integration

package onspring_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestFilesIntegration(t *testing.T) {
	loadEnvFile(t)
	client := createClient(t)
	ctx := context.Background()

	t.Run("GetInfo", func(t *testing.T) {
		t.Run("should return information about a file in an attachment field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			info, err := client.Files.GetInfo(ctx, recordId, fieldId, fileId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if info.Name == "" {
				t.Error("Expected name to not be empty")
			}

			if info.ContentType == "" {
				t.Error("Expected content type to not be empty")
			}

			if info.CreatedDate == "" {
				t.Error("Expected created date to not be empty")
			}

			if info.ModifiedDate == "" {
				t.Error("Expected modified date to not be empty")
			}

			if info.Owner == "" {
				t.Error("Expected owner to not be empty")
			}

			if info.Type == "" {
				t.Error("Expected type to not be empty")
			}

			if info.FileHref == "" {
				t.Error("Expected file href to not be empty")
			}
		})

		t.Run("should return information about a file in an image field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_IMAGE_FIELD")
			fileId := requireEnvInt(t, "TEST_IMAGE")

			info, err := client.Files.GetInfo(ctx, recordId, fieldId, fileId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if info.Name == "" {
				t.Error("Expected name to not be empty")
			}

			if info.ContentType == "" {
				t.Error("Expected content type to not be empty")
			}

			if info.Type == "" {
				t.Error("Expected type to not be empty")
			}
		})

		t.Run("should return a 400 response when fieldId is not for a file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetInfo(ctx, recordId, textFieldId, fileId)

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := invalidClient.Files.GetInfo(ctx, recordId, fieldId, fileId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when the api key does not have access to the file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoAccess := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetInfo(ctx, recordId, fieldIdNoAccess, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 403 response when the api key does not have access to the app", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoApp := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_APP")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetInfo(ctx, recordId, fieldIdNoApp, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 response when the file field cannot be found", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetInfo(ctx, recordId, 0, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})

		t.Run("should return a 404 response when the file record cannot be found", func(t *testing.T) {
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetInfo(ctx, 0, fieldId, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("GetContent", func(t *testing.T) {
		t.Run("should return a file in an attachment field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			content, err := client.Files.GetContent(ctx, recordId, fieldId, fileId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(content.Data) == 0 {
				t.Error("Expected data to not be empty")
			}

			if content.ContentType == "" {
				t.Error("Expected content type to not be empty")
			}

			if content.FileName == "" {
				t.Error("Expected file name to not be empty")
			}
		})

		t.Run("should return a file in an image field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_IMAGE_FIELD")
			fileId := requireEnvInt(t, "TEST_IMAGE")

			content, err := client.Files.GetContent(ctx, recordId, fieldId, fileId)

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if len(content.Data) == 0 {
				t.Error("Expected data to not be empty")
			}

			if content.ContentType == "" {
				t.Error("Expected content type to not be empty")
			}

			if content.FileName == "" {
				t.Error("Expected file name to not be empty")
			}
		})

		t.Run("should return a 400 response when fieldId is not for a file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetContent(ctx, recordId, textFieldId, fileId)

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := invalidClient.Files.GetContent(ctx, recordId, fieldId, fileId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when the api key does not have access to the file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoAccess := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetContent(ctx, recordId, fieldIdNoAccess, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 403 response when the api key does not have access to the app", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoApp := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_APP")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetContent(ctx, recordId, fieldIdNoApp, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 response when the file field cannot be found", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetContent(ctx, recordId, 0, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})

		t.Run("should return a 404 response when the file record cannot be found", func(t *testing.T) {
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			_, err := client.Files.GetContent(ctx, 0, fieldId, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("should save a file into an attachment field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			response, err := client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldId,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.Id == 0 {
				t.Error("Expected file id to not be zero")
			}

			t.Cleanup(func() {
				_ = client.Files.Delete(context.Background(), recordId, fieldId, response.Id)
			})
		})

		t.Run("should save a file into an image field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_IMAGE_FIELD")

			file, err := os.Open("testdata/test-image.jpeg")

			if err != nil {
				t.Fatalf("Failed to open test image: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			response, err := client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldId,
				FileName:     "test-image.jpeg",
				FileContents: file,
			})

			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if response.Id == 0 {
				t.Error("Expected file id to not be zero")
			}

			t.Cleanup(func() {
				_ = client.Files.Delete(context.Background(), recordId, fieldId, response.Id)
			})
		})

		t.Run("should return a 400 response when fieldId is not for a file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      textFieldId,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = invalidClient.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldId,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when the api key does not have access to the field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoAccess := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldIdNoAccess,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 403 response when the api key does not have access to the app", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoApp := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_APP")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldIdNoApp,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 response when the file field cannot be found", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      0,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusNotFound)
		})

		t.Run("should return a 404 response when the file record cannot be found", func(t *testing.T) {
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			_, err = client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     0,
				FieldId:      fieldId,
				FileName:     "test-attachment.txt",
				FileContents: file,
			})

			assertAPIError(t, err, http.StatusNotFound)
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("should delete a file from an attachment field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")

			file, err := os.Open("testdata/test-attachment.txt")

			if err != nil {
				t.Fatalf("Failed to open test attachment: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			saveResponse, err := client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldId,
				FileName:     "test-attachment-to-delete.txt",
				FileContents: file,
			})

			if err != nil {
				t.Fatalf("Failed to save file: %v", err)
			}

			err = client.Files.Delete(ctx, recordId, fieldId, saveResponse.Id)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		t.Run("should delete a file from an image field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_IMAGE_FIELD")

			file, err := os.Open("testdata/test-image.jpeg")

			if err != nil {
				t.Fatalf("Failed to open test image: %v", err)
			}

			defer func() {
				_ = file.Close()
			}()

			saveResponse, err := client.Files.Save(ctx, onspring.SaveFileRequest{
				RecordId:     recordId,
				FieldId:      fieldId,
				FileName:     "test-image-to-delete.jpeg",
				FileContents: file,
			})

			if err != nil {
				t.Fatalf("Failed to save file: %v", err)
			}

			err = client.Files.Delete(ctx, recordId, fieldId, saveResponse.Id)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})

		t.Run("should return a 400 response when fieldId is not for a file field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			textFieldId := requireEnvInt(t, "TEST_TEXT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := client.Files.Delete(ctx, recordId, textFieldId, fileId)

			assertAPIError(t, err, http.StatusBadRequest)
		})

		t.Run("should return a 401 response when the api key is invalid", func(t *testing.T) {
			invalidClient := createInvalidClient(t)
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := invalidClient.Files.Delete(ctx, recordId, fieldId, fileId)

			assertAPIError(t, err, http.StatusUnauthorized)
		})

		t.Run("should return a 403 response when the api key does not have access to the field", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoAccess := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := client.Files.Delete(ctx, recordId, fieldIdNoAccess, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 403 response when the api key does not have access to the app", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fieldIdNoApp := requireEnvInt(t, "TEST_ATTACHMENT_FIELD_NO_ACCESS_APP")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := client.Files.Delete(ctx, recordId, fieldIdNoApp, fileId)

			assertAPIError(t, err, http.StatusForbidden)
		})

		t.Run("should return a 404 response when the file field cannot be found", func(t *testing.T) {
			recordId := requireEnvInt(t, "TEST_RECORD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := client.Files.Delete(ctx, recordId, 0, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})

		t.Run("should return a 404 response when the file record cannot be found", func(t *testing.T) {
			fieldId := requireEnvInt(t, "TEST_ATTACHMENT_FIELD")
			fileId := requireEnvInt(t, "TEST_ATTACHMENT")

			err := client.Files.Delete(ctx, 0, fieldId, fileId)

			assertAPIError(t, err, http.StatusNotFound)
		})
	})
}
