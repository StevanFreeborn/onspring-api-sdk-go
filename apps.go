package onspring

import (
	"context"
	"net/http"
)

const (
	appsPath = "/apps"
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

// Get retrieves a paginated list of apps from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - Page[App]: A page of apps with pagination metadata
//   - error: An error if the request fails
func (p *AppsEndpoint) Get(ctx context.Context, pagingOpts ...PagingOption) (Page[App], error) {
	pagingRequest := &PagingRequest{
		pageNumber: 1,
		pageSize:   50,
	}

	for _, opt := range pagingOpts {
		opt(pagingRequest)
	}

	req, requestCreationErr := p.client.newRequest(ctx, http.MethodGet, appsPath, pagingRequest.ToParams(), nil)

	var page Page[App]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := p.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}
