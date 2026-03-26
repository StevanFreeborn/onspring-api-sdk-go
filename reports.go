package onspring

import (
	"context"
	"fmt"
	"iter"
	"net/http"
)

const (
	reportsPath = "/reports"
)

// ReportsEndpoint provides access to reports in an Onspring instance.
type ReportsEndpoint struct {
	client *Client
}

// ReportData represents the data returned from an Onspring report.
type ReportData struct {
	Columns []string    `json:"columns"`
	Rows    []ReportRow `json:"rows"`
}

// ReportRow represents a row in a report.
type ReportRow struct {
	RecordId *int  `json:"recordId"`
	Cells    []any `json:"cells"`
}

// Report represents a report associated to an app.
type Report struct {
	AppId       int    `json:"appId"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ReportOption is a functional option for configuring report requests.
type ReportOption func(*reportRequest)

type reportRequest struct {
	DataFormat string
	DataType   string
}

func (r *reportRequest) ToParams() map[string]string {
	params := map[string]string{}

	if r.DataFormat != "" {
		params["apiDataFormat"] = r.DataFormat
	}

	if r.DataType != "" {
		params["dataType"] = r.DataType
	}

	return params
}

func createReportRequest(opts []ReportOption) *reportRequest {
	r := &reportRequest{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// WithDataFormat sets the data format for the report request.
// Valid values are "Raw" and "Formatted".
func WithDataFormat(format string) ReportOption {
	return func(r *reportRequest) {
		r.DataFormat = format
	}
}

// WithDataType sets the data type for the report request.
// Valid values are "ReportData" and "ChartData".
func WithDataType(dataType string) ReportOption {
	return func(r *reportRequest) {
		r.DataType = dataType
	}
}

// Get retrieves a report from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - reportId: The id of the report to retrieve
//   - opts: Optional report configuration functions (e.g., WithDataFormat, WithDataType)
//
// Returns:
//   - ReportData: The data for a report
//   - error: An error if the request fails
func (rp *ReportsEndpoint) Get(ctx context.Context, reportId int, opts ...ReportOption) (ReportData, error) {
	reportReq := createReportRequest(opts)
	path := fmt.Sprintf("%s/id/%d", reportsPath, reportId)
	req, requestCreationErr := rp.client.newRequest(ctx, http.MethodGet, path, reportReq.ToParams(), nil)

	var reportData ReportData

	if requestCreationErr != nil {
		return reportData, requestCreationErr
	}

	responseErr := rp.client.doWithJsonResponse(req, &reportData)

	if responseErr != nil {
		return reportData, responseErr
	}

	return reportData, nil
}

// List retrieves a paginated list of reports for an app from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve reports for
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - Page[Report]: A page of reports with pagination metadata
//   - error: An error if the request fails
func (rp *ReportsEndpoint) List(ctx context.Context, appId int, pagingOpts ...PagingOption) (Page[Report], error) {
	pagingRequest := createPagingRequest(pagingOpts)
	path := fmt.Sprintf("%s/appId/%d", reportsPath, appId)

	req, requestCreationErr := rp.client.newRequest(ctx, http.MethodGet, path, pagingRequest.ToParams(), nil)

	var page Page[Report]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := rp.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}

// ListAll returns an iterator that yields all Reports for an app across all pages.
// It automatically handles pagination by making sequential calls to List
// until all items have been retrieved or the caller stops the iteration.
//
// The iterator yields each Report and any error encountered during fetching.
// If an error occurs during a page request, the error is yielded and
// iteration terminates.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve reports for
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - iter.Seq2[Report, error]: An iterator yielding:
//   - Report: The individual report record.
//   - error: An error if a specific page request fails during iteration.
func (rp *ReportsEndpoint) ListAll(ctx context.Context, appId int, pagingOpts ...PagingOption) iter.Seq2[Report, error] {
	return func(yield func(Report, error) bool) {
		pagingRequest := createPagingRequest(pagingOpts)

		for {
			page, err := rp.List(ctx, appId, ForPageNumber(pagingRequest.PageNumber), WithPageSize(pagingRequest.PageSize))

			if err != nil {
				yield(Report{}, err)
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
