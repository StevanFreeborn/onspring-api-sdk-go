package onspring

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	filesPath = "/files"
)

// FilesEndpoint provides access to files in an Onspring instance.
type FilesEndpoint struct {
	client *Client
}

// FileInfo represents the metadata of a file in Onspring.
type FileInfo struct {
	Type         string `json:"type"`
	ContentType  string `json:"contentType"`
	Name         string `json:"name"`
	CreatedDate  string `json:"createdDate"`
	ModifiedDate string `json:"modifiedDate"`
	Owner        string `json:"owner"`
	Notes        string `json:"notes"`
	FileHref     string `json:"fileHref"`
}

// FileContent represents the content of a downloaded file.
type FileContent struct {
	FileName    string
	ContentType string
	Data        []byte
}

// SaveFileRequest represents a request to upload a file.
type SaveFileRequest struct {
	RecordId     int
	FieldId      int
	Notes        string
	ModifiedDate string
	FileName     string
	FileContents io.Reader
}

// SaveFileResponse represents the response for saving a file.
type SaveFileResponse struct {
	Id int `json:"id"`
}

func fileBasePath(recordId, fieldId, fileId int) string {
	return fmt.Sprintf("%s/recordId/%d/fieldId/%d/fileId/%d", filesPath, recordId, fieldId, fileId)
}

// GetInfo retrieves a file's metadata from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - recordId: The id of the record
//   - fieldId: The id of the field
//   - fileId: The id of the file
//
// Returns:
//   - FileInfo: The file's metadata
//   - error: An error if the request fails
func (f *FilesEndpoint) GetInfo(ctx context.Context, recordId, fieldId, fileId int) (FileInfo, error) {
	path := fileBasePath(recordId, fieldId, fileId)
	req, requestCreationErr := f.client.newRequest(ctx, http.MethodGet, path, nil, nil)

	var fileInfo FileInfo

	if requestCreationErr != nil {
		return fileInfo, requestCreationErr
	}

	responseErr := f.client.doWithJsonResponse(req, &fileInfo)

	if responseErr != nil {
		return fileInfo, responseErr
	}

	return fileInfo, nil
}

// GetContent retrieves a file's content from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - recordId: The id of the record
//   - fieldId: The id of the field
//   - fileId: The id of the file
//
// Returns:
//   - FileContent: The file's content including name, content type, and data
//   - error: An error if the request fails
func (f *FilesEndpoint) GetContent(ctx context.Context, recordId, fieldId, fileId int) (FileContent, error) {
	path := fmt.Sprintf("%s/file", fileBasePath(recordId, fieldId, fileId))
	req, requestCreationErr := f.client.newRequest(ctx, http.MethodGet, path, nil, nil)

	var fileContent FileContent

	if requestCreationErr != nil {
		return fileContent, requestCreationErr
	}

	data, headers, responseErr := f.client.doWithBytesResponse(req)

	if responseErr != nil {
		return fileContent, responseErr
	}

	fileContent.Data = data
	fileContent.ContentType = headers.Get("Content-Type")

	contentDisposition := headers.Get("Content-Disposition")

	if contentDisposition != "" {
		_, params, err := mime.ParseMediaType(contentDisposition)

		if err == nil {
			fileContent.FileName = params["filename"]
		}
	}

	return fileContent, nil
}

// Delete removes a file from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - recordId: The id of the record
//   - fieldId: The id of the field
//   - fileId: The id of the file
//
// Returns:
//   - error: An error if the request fails
func (f *FilesEndpoint) Delete(ctx context.Context, recordId, fieldId, fileId int) error {
	path := fileBasePath(recordId, fieldId, fileId)
	req, requestCreationErr := f.client.newRequest(ctx, http.MethodDelete, path, nil, nil)

	if requestCreationErr != nil {
		return requestCreationErr
	}

	return f.client.do(req)
}

// Save uploads a file to the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - saveReq: The file upload request containing record id, field id, file name, contents, and optional metadata
//
// Returns:
//   - SaveFileResponse: The response containing the saved file's id
//   - error: An error if the request fails
func (f *FilesEndpoint) Save(ctx context.Context, saveReq SaveFileRequest) (SaveFileResponse, error) {
	var response SaveFileResponse

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("RecordId", strconv.Itoa(saveReq.RecordId))
	_ = writer.WriteField("FieldId", strconv.Itoa(saveReq.FieldId))

	if saveReq.Notes != "" {
		_ = writer.WriteField("Notes", saveReq.Notes)
	}

	if saveReq.ModifiedDate != "" {
		_ = writer.WriteField("ModifiedDate", saveReq.ModifiedDate)
	}

	filePart, err := writer.CreateFormFile("File", saveReq.FileName)

	if err != nil {
		return response, fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(filePart, saveReq.FileContents)

	if err != nil {
		return response, fmt.Errorf("failed to copy file contents: %w", err)
	}

	err = writer.Close()

	if err != nil {
		return response, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, requestCreationErr := f.client.newMultipartRequest(ctx, filesPath, body, writer.FormDataContentType())

	if requestCreationErr != nil {
		return response, requestCreationErr
	}

	responseErr := f.client.doWithJsonResponse(req, &response)

	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}
