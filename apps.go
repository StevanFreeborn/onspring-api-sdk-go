package onspring

import (
	"context"
	"net/http"
)

const (
	appsPath = "/apps"
)

type AppsEndpoint struct {
	client *Client
}

type App struct {
	Href string `json:"href"`
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (p *AppsEndpoint) Get(ctx context.Context, pagingOpts ...PagingOption) (Page[App], error) {
	pagingRequest := &PagingRequest{
		pageSize:   1,
		pageNumber: 50,
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
