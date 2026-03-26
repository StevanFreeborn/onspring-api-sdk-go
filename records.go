package onspring

import (
	"context"
	"fmt"
	"iter"
	"net/http"
	"strconv"
	"strings"
)

const (
	recordsPath = "/records"
)

// RecordsEndpoint provides access to records in an Onspring instance.
type RecordsEndpoint struct {
	client *Client
}

// Record represents an Onspring record.
type Record struct {
	AppId     int                `json:"appId"`
	RecordId  int                `json:"recordId"`
	FieldData []RecordFieldValue `json:"fieldData"`
}

// RecordFieldValue represents a field value within a record.
type RecordFieldValue struct {
	Type    string `json:"type"`
	FieldId int    `json:"fieldId"`
	Value   any    `json:"value"`
}

// RecordBatch represents a batch of Onspring records.
type RecordBatch struct {
	Count int      `json:"count"`
	Items []Record `json:"items"`
}

// RecordOption is a functional option for configuring record requests.
type RecordOption func(*recordRequest)

type recordRequest struct {
	FieldIds      []int
	DataFormat    string
	PagingRequest PagingRequest
}

func (r *recordRequest) ToParams() map[string]string {
	params := r.PagingRequest.ToParams()

	if len(r.FieldIds) > 0 {
		ids := make([]string, len(r.FieldIds))

		for i, id := range r.FieldIds {
			ids[i] = strconv.Itoa(id)
		}

		params["fieldIds"] = strings.Join(ids, ",")
	}

	if r.DataFormat != "" {
		params["dataFormat"] = r.DataFormat
	}

	return params
}

func (r *recordRequest) ToQueryParams() map[string]string {
	params := map[string]string{}

	if len(r.FieldIds) > 0 {
		ids := make([]string, len(r.FieldIds))

		for i, id := range r.FieldIds {
			ids[i] = strconv.Itoa(id)
		}

		params["fieldIds"] = strings.Join(ids, ",")
	}

	if r.DataFormat != "" {
		params["dataFormat"] = r.DataFormat
	}

	return params
}

func createRecordRequest(opts []RecordOption) *recordRequest {
	r := &recordRequest{
		PagingRequest: PagingRequest{PageNumber: 1, PageSize: 50},
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// WithFieldIds sets the field identifiers to include in the record response.
func WithFieldIds(ids []int) RecordOption {
	return func(r *recordRequest) {
		r.FieldIds = ids
	}
}

// WithRecordDataFormat sets the data format for the record response.
// Valid values are "Raw" and "Formatted".
func WithRecordDataFormat(format string) RecordOption {
	return func(r *recordRequest) {
		r.DataFormat = format
	}
}

// WithPaging applies paging options to a record request.
func WithPaging(opts ...PagingOption) RecordOption {
	return func(r *recordRequest) {
		for _, opt := range opts {
			opt(&r.PagingRequest)
		}
	}
}

// GetManyRecordsRequest represents a request to get a batch of records.
type GetManyRecordsRequest struct {
	AppId      int    `json:"appId"`
	RecordIds  []int  `json:"recordIds"`
	FieldIds   []int  `json:"fieldIds,omitempty"`
	DataFormat string `json:"dataFormat,omitempty"`
}

// QueryRecordsRequest represents a request to query records.
type QueryRecordsRequest struct {
	AppId      int    `json:"appId"`
	Filter     string `json:"filter"`
	FieldIds   []int  `json:"fieldIds,omitempty"`
	DataFormat string `json:"dataFormat,omitempty"`
}

// SaveRecordRequest represents a request to save a record.
type SaveRecordRequest struct {
	AppId    int            `json:"appId"`
	RecordId *int           `json:"recordId,omitempty"`
	Fields   map[string]any `json:"fields"`
}

// SaveRecordResponse represents the response for saving a record.
type SaveRecordResponse struct {
	Id       int      `json:"id"`
	Warnings []string `json:"warnings"`
}

// DeleteManyRecordsRequest represents a request to delete a batch of records.
type DeleteManyRecordsRequest struct {
	AppId     int   `json:"appId"`
	RecordIds []int `json:"recordIds"`
}

// Get retrieves a record from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app
//   - recordId: The id of the record to retrieve
//   - opts: Optional record configuration functions (e.g., WithFieldIds, WithRecordDataFormat)
//
// Returns:
//   - Record: A record
//   - error: An error if the request fails
func (rc *RecordsEndpoint) Get(ctx context.Context, appId, recordId int, opts ...RecordOption) (Record, error) {
	recordReq := createRecordRequest(opts)
	path := fmt.Sprintf("%s/appId/%d/recordId/%d", recordsPath, appId, recordId)
	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodGet, path, recordReq.ToQueryParams(), nil)

	var record Record

	if requestCreationErr != nil {
		return record, requestCreationErr
	}

	responseErr := rc.client.doWithJsonResponse(req, &record)

	if responseErr != nil {
		return record, responseErr
	}

	return record, nil
}

// List retrieves a paginated list of records for an app from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve records for
//   - opts: Optional record configuration functions (e.g., WithFieldIds, WithRecordDataFormat, WithPaging)
//
// Returns:
//   - Page[Record]: A page of records with pagination metadata
//   - error: An error if the request fails
func (rc *RecordsEndpoint) List(ctx context.Context, appId int, opts ...RecordOption) (Page[Record], error) {
	recordReq := createRecordRequest(opts)
	path := fmt.Sprintf("%s/appId/%d", recordsPath, appId)

	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodGet, path, recordReq.ToParams(), nil)

	var page Page[Record]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := rc.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}

// ListAll returns an iterator that yields all Records for an app across all pages.
// It automatically handles pagination by making sequential calls to List
// until all items have been retrieved or the caller stops the iteration.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve records for
//   - opts: Optional record configuration functions (e.g., WithFieldIds, WithRecordDataFormat, WithPaging)
//
// Returns:
//   - iter.Seq2[Record, error]: An iterator yielding:
//   - Record: The individual record.
//   - error: An error if a specific page request fails during iteration.
func (rc *RecordsEndpoint) ListAll(ctx context.Context, appId int, opts ...RecordOption) iter.Seq2[Record, error] {
	return func(yield func(Record, error) bool) {
		recordReq := createRecordRequest(opts)

		for {
			page, err := rc.List(
				ctx,
				appId,
				WithFieldIds(recordReq.FieldIds),
				WithRecordDataFormat(recordReq.DataFormat),
				WithPaging(ForPageNumber(recordReq.PagingRequest.PageNumber), WithPageSize(recordReq.PagingRequest.PageSize)),
			)

			if err != nil {
				yield(Record{}, err)
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if page.PageNumber >= page.TotalPages {
				break
			}

			recordReq.PagingRequest.PageNumber++
		}
	}
}

// GetMany retrieves a batch of records from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - request: The batch get request containing app id, record ids, and optional field ids/data format
//
// Returns:
//   - RecordBatch: A batch of records
//   - error: An error if the request fails
func (rc *RecordsEndpoint) GetMany(ctx context.Context, request GetManyRecordsRequest) (RecordBatch, error) {
	path := fmt.Sprintf("%s/batch-get", recordsPath)
	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodPost, path, nil, request)

	var batch RecordBatch

	if requestCreationErr != nil {
		return batch, requestCreationErr
	}

	responseErr := rc.client.doWithJsonResponse(req, &batch)

	if responseErr != nil {
		return batch, responseErr
	}

	return batch, nil
}

// Query queries records from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - request: The query request containing app id, filter, and optional field ids/data format
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - Page[Record]: A page of records with pagination metadata
//   - error: An error if the request fails
func (rc *RecordsEndpoint) Query(ctx context.Context, request QueryRecordsRequest, pagingOpts ...PagingOption) (Page[Record], error) {
	pagingRequest := createPagingRequest(pagingOpts)
	path := fmt.Sprintf("%s/query", recordsPath)

	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodPost, path, pagingRequest.ToParams(), request)

	var page Page[Record]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := rc.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}

// QueryAll returns an iterator that yields all Records matching a query across all pages.
// It automatically handles pagination by making sequential calls to Query
// until all items have been retrieved or the caller stops the iteration.
//
// Parameters:
//   - ctx: The context for the request
//   - request: The query request containing app id, filter, and optional field ids/data format
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - iter.Seq2[Record, error]: An iterator yielding:
//   - Record: The individual record.
//   - error: An error if a specific page request fails during iteration.
func (rc *RecordsEndpoint) QueryAll(ctx context.Context, request QueryRecordsRequest, pagingOpts ...PagingOption) iter.Seq2[Record, error] {
	return func(yield func(Record, error) bool) {
		pagingRequest := createPagingRequest(pagingOpts)

		for {
			page, err := rc.Query(ctx, request, ForPageNumber(pagingRequest.PageNumber), WithPageSize(pagingRequest.PageSize))

			if err != nil {
				yield(Record{}, err)
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if page.PageNumber >= page.TotalPages {
				break
			}

			pagingRequest.PageNumber++
		}
	}
}

// Save creates or updates a record in the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - request: The save request containing app id, optional record id, and field values
//
// Returns:
//   - SaveRecordResponse: The response containing the saved record's id and any warnings
//   - error: An error if the request fails
func (rc *RecordsEndpoint) Save(ctx context.Context, request SaveRecordRequest) (SaveRecordResponse, error) {
	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodPut, recordsPath, nil, request)

	var response SaveRecordResponse

	if requestCreationErr != nil {
		return response, requestCreationErr
	}

	responseErr := rc.client.doWithJsonResponse(req, &response)

	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

// Delete removes a record from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app
//   - recordId: The id of the record to delete
//
// Returns:
//   - error: An error if the request fails
func (rc *RecordsEndpoint) Delete(ctx context.Context, appId, recordId int) error {
	path := fmt.Sprintf("%s/appId/%d/recordId/%d", recordsPath, appId, recordId)
	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodDelete, path, nil, nil)

	if requestCreationErr != nil {
		return requestCreationErr
	}

	return rc.client.do(req)
}

// DeleteMany removes a batch of records from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - request: The batch delete request containing app id and record ids
//
// Returns:
//   - error: An error if the request fails
func (rc *RecordsEndpoint) DeleteMany(ctx context.Context, request DeleteManyRecordsRequest) error {
	path := fmt.Sprintf("%s/batch-delete", recordsPath)
	req, requestCreationErr := rc.client.newRequest(ctx, http.MethodPost, path, nil, request)

	if requestCreationErr != nil {
		return requestCreationErr
	}

	return rc.client.do(req)
}
