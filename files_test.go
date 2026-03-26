package onspring_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/StevanFreeborn/onspring-api-sdk-go"
)

func TestFiles(t *testing.T) {
	t.Run("GetInfo", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Files.GetInfo(nilContext, 1, 1, 1)

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

			_, err := client.Files.GetInfo(ctx, 1, 1, 1)

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

			_, err := client.Files.GetInfo(t.Context(), 1, 1, 1)

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

			_, err := invalidClient.Files.GetInfo(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Files.GetInfo(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the correct endpoint and return file info", func(t *testing.T) {
			recordId := 1
			fieldId := 2
			fileId := 3

			expectedFileInfo := onspring.FileInfo{
				Type:         "Attachment",
				ContentType:  "application/pdf",
				Name:         "test.pdf",
				CreatedDate:  "2024-01-01T00:00:00Z",
				ModifiedDate: "2024-01-02T00:00:00Z",
				Owner:        "admin",
				Notes:        "Test notes",
				FileHref:     "https://api.onspring.com/files/1/2/3/file",
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/files/recordId/%d/fieldId/%d/fileId/%d", recordId, fieldId, fileId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				jsonData, _ := json.Marshal(expectedFileInfo)

				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			fileInfo, err := client.Files.GetInfo(t.Context(), recordId, fieldId, fileId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedFileInfo, fileInfo) {
				t.Errorf("Expected %v but got %v", expectedFileInfo, fileInfo)
			}
		})
	})

	t.Run("GetContent", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			var nilContext context.Context = nil

			_, err := client.Files.GetContent(nilContext, 1, 1, 1)

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

			_, err := client.Files.GetContent(ctx, 1, 1, 1)

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

			_, err := client.Files.GetContent(t.Context(), 1, 1, 1)

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

			_, err := invalidClient.Files.GetContent(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-200 status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.Files.GetContent(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a GET request to the correct endpoint and return file content", func(t *testing.T) {
			recordId := 1
			fieldId := 2
			fileId := 3
			expectedData := []byte("file content here")
			expectedContentType := "application/pdf"
			expectedFileName := "test.pdf"

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("Expected GET method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/files/recordId/%d/fieldId/%d/fileId/%d/file", recordId, fieldId, fileId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", expectedContentType)
				w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, expectedFileName))
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(expectedData)
			})

			fileContent, err := client.Files.GetContent(t.Context(), recordId, fieldId, fileId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if fileContent.FileName != expectedFileName {
				t.Errorf("Expected file name %s but got %s", expectedFileName, fileContent.FileName)
			}

			if fileContent.ContentType != expectedContentType {
				t.Errorf("Expected content type %s but got %s", expectedContentType, fileContent.ContentType)
			}

			if !bytes.Equal(expectedData, fileContent.Data) {
				t.Errorf("Expected data %v but got %v", expectedData, fileContent.Data)
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			})

			var nilContext context.Context = nil

			err := client.Files.Delete(nilContext, 1, 1, 1)

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

			err := client.Files.Delete(ctx, 1, 1, 1)

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

			err := client.Files.Delete(t.Context(), 1, 1, 1)

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

			err := invalidClient.Files.Delete(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			err := client.Files.Delete(t.Context(), 1, 1, 1)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a DELETE request to the correct endpoint and return no error if receives 204 status code", func(t *testing.T) {
			recordId := 1
			fieldId := 2
			fileId := 3

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("Expected DELETE method, got %s", r.Method)
				}

				expectedPath := fmt.Sprintf("/files/recordId/%d/fieldId/%d/fileId/%d", recordId, fieldId, fileId)

				if r.URL.Path != expectedPath {
					t.Errorf("Expected %s endpoint, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusNoContent)
			})

			err := client.Files.Delete(t.Context(), recordId, fieldId, fileId)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Run("it should return an error if context is nil", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			})

			var nilContext context.Context = nil

			saveReq := onspring.SaveFileRequest{
				RecordId:     1,
				FieldId:      1,
				FileName:     "test.pdf",
				FileContents: strings.NewReader("file content"),
			}

			_, err := client.Files.Save(nilContext, saveReq)

			if err == nil {
				t.Errorf("Expected error for nil context, got nil")
			}
		})

		t.Run("it should return an error if context is canceled", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			})

			ctx, cancel := context.WithCancel(t.Context())

			cancel()

			saveReq := onspring.SaveFileRequest{
				RecordId:     1,
				FieldId:      1,
				FileName:     "test.pdf",
				FileContents: strings.NewReader("file content"),
			}

			_, err := client.Files.Save(ctx, saveReq)

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

			saveReq := onspring.SaveFileRequest{
				RecordId:     1,
				FieldId:      1,
				FileName:     "test.pdf",
				FileContents: strings.NewReader("file content"),
			}

			_, err := client.Files.Save(t.Context(), saveReq)

			if err == nil {
				t.Errorf("Expected network error, got nil")
			}
		})

		t.Run("it should return an error if create a request fails", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
			})

			invalidClient := onspring.NewClient(
				"test-api-key",
				onspring.WithBaseURL("http://[::1]:namedport"),
				onspring.WithHTTPClient(client.HTTPClient()),
			)

			saveReq := onspring.SaveFileRequest{
				RecordId:     1,
				FieldId:      1,
				FileName:     "test.pdf",
				FileContents: strings.NewReader("file content"),
			}

			_, err := invalidClient.Files.Save(t.Context(), saveReq)

			if err == nil {
				t.Errorf("Expected request creation error, got nil")
			}
		})

		t.Run("it should return an error if the endpoint returns a non-2xx status code", func(t *testing.T) {
			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			})

			saveReq := onspring.SaveFileRequest{
				RecordId:     1,
				FieldId:      1,
				FileName:     "test.pdf",
				FileContents: strings.NewReader("file content"),
			}

			_, err := client.Files.Save(t.Context(), saveReq)

			if err == nil {
				t.Errorf("Expected error, got nil")
			}
		})

		t.Run("it should perform a POST request to the /files endpoint with multipart form data and return a response when successful", func(t *testing.T) {
			expectedRecordId := 1
			expectedFieldId := 2
			expectedNotes := "Test notes"
			expectedModifiedDate := "2024-01-01T00:00:00Z"
			expectedFileName := "test.pdf"
			expectedFileContent := "file content here"

			expectedResponse := onspring.SaveFileResponse{
				Id: 42,
			}

			_, client := setupMockServer(t, func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}

				if r.URL.Path != "/files" {
					t.Errorf("Expected /files endpoint, got %s", r.URL.Path)
				}

				contentType := r.Header.Get("Content-Type")

				if !strings.HasPrefix(contentType, "multipart/form-data") {
					t.Errorf("Expected multipart/form-data content type, got %s", contentType)
				}

				err := r.ParseMultipartForm(10 << 20)

				if err != nil {
					t.Errorf("Expected to parse multipart form, but got error: %v", err)
				}

				recordId := r.FormValue("RecordId")

				if recordId != fmt.Sprintf("%d", expectedRecordId) {
					t.Errorf("Expected RecordId %d but got %s", expectedRecordId, recordId)
				}

				fieldId := r.FormValue("FieldId")

				if fieldId != fmt.Sprintf("%d", expectedFieldId) {
					t.Errorf("Expected FieldId %d but got %s", expectedFieldId, fieldId)
				}

				notes := r.FormValue("Notes")

				if notes != expectedNotes {
					t.Errorf("Expected Notes %s but got %s", expectedNotes, notes)
				}

				modifiedDate := r.FormValue("ModifiedDate")

				if modifiedDate != expectedModifiedDate {
					t.Errorf("Expected ModifiedDate %s but got %s", expectedModifiedDate, modifiedDate)
				}

				file, header, err := r.FormFile("File")

				if err != nil {
					t.Errorf("Expected to get file from form, but got error: %v", err)
				}

				defer func() {
					_ = file.Close()
				}()

				if header.Filename != expectedFileName {
					t.Errorf("Expected file name %s but got %s", expectedFileName, header.Filename)
				}

				fileBytes, _ := io.ReadAll(file)

				if string(fileBytes) != expectedFileContent {
					t.Errorf("Expected file content %s but got %s", expectedFileContent, string(fileBytes))
				}

				jsonData, _ := json.Marshal(expectedResponse)

				w.WriteHeader(http.StatusCreated)
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write(jsonData)
			})

			saveReq := onspring.SaveFileRequest{
				RecordId:     expectedRecordId,
				FieldId:      expectedFieldId,
				Notes:        expectedNotes,
				ModifiedDate: expectedModifiedDate,
				FileName:     expectedFileName,
				FileContents: strings.NewReader(expectedFileContent),
			}

			response, err := client.Files.Save(t.Context(), saveReq)

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !reflect.DeepEqual(expectedResponse, response) {
				t.Errorf("Expected %v but got %v", expectedResponse, response)
			}
		})
	})
}
