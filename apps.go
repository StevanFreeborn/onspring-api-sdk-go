package onspring

import (
	"context"
	"iter"
	"net/http"
)

const (
	appsPath      = "/apps"
	appsBatchPath = "/apps/batch-get"
)

// AppsEndpoint provides access to apps in an Onspring instance.
type AppsEndpoint struct {
	client *Client
}

// App represents an Onspring app
type App struct {
	Href string `json:"href"`
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// AppBatch represents a batch of Onspring apps
type AppBatch struct {
	Count int   `json:"count"`
	Items []App `json:"items"`
}

// Get retrieves a paginated list of apps from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - Page[App]: A page of apps with pagination metadata
//   - error: An error if the request fails
func (a *AppsEndpoint) Get(ctx context.Context, pagingOpts ...PagingOption) (Page[App], error) {
	pagingRequest := createPagingRequest(pagingOpts)

	req, requestCreationErr := a.client.newRequest(ctx, http.MethodGet, appsPath, pagingRequest.ToParams(), nil)

	var page Page[App]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := a.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}

// GetAll returns an iterator that yields all Apps across all pages.
// It automatically handles pagination by making sequential calls to Get
// until all items have been retrieved or the caller stops the iteration.
//
// The iterator yields each *App and any error encountered during fetching.
// If an error occurs during a page request, the error is yielded and
// iteration terminates.
//
// Parameters:
//   - ctx: The context for the request
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - iter.Seq2[App, error]: An iterator yielding:
//   - App: The individual application record.
//   - error: An error if a specific page request fails during iteration.
func (a *AppsEndpoint) GetAll(ctx context.Context, pagingOpts ...PagingOption) iter.Seq2[App, error] {
	return func(yield func(App, error) bool) {
		pagingRequest := createPagingRequest(pagingOpts)

		for {
			page, err := a.Get(ctx, ForPageNumber(pagingRequest.PageNumber), WithPageSize(pagingRequest.PageSize))

			if err != nil {
				yield(App{}, err)
				return
			}

			for _, item := range page.Items {
				if !yield(item, nil) {
					return
				}
			}

			if page.TotalPages == page.PageNumber {
				break
			}

			pagingRequest.PageNumber++
		}
	}
}

func (a *AppsEndpoint) GetBatch(ctx context.Context, appIds []int) (AppBatch, error) {
	req, requestCreationErr := a.client.newRequest(ctx, http.MethodPost, appsBatchPath, nil, appIds)

	var appBatch AppBatch

	if requestCreationErr != nil {
		return appBatch, requestCreationErr
	}

	responseErr := a.client.doWithJsonResponse(req, &appBatch)

	if responseErr != nil {
		return appBatch, responseErr
	}

	return appBatch, nil
}

func createPagingRequest(pagingOpts []PagingOption) *PagingRequest {
	pagingRequest := &PagingRequest{
		PageNumber: 1,
		PageSize:   50,
	}

	for _, opt := range pagingOpts {
		opt(pagingRequest)
	}

	return pagingRequest
}
