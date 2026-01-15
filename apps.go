package onspring

import (
	"context"
	"net/http"
	"strconv"
)

const (
	appsPath = "/apps"
)

type AppsEndpoint struct {
	client *Client
}

type PagingRequest struct {
	pageSize   int
	pageNumber int
}

func (pr *PagingRequest) ToParams() map[string]string {
	return map[string]string{
		"pageSize":   strconv.Itoa(pr.pageNumber),
		"pageNumber": strconv.Itoa(pr.pageSize),
	}
}

type PagingOption func(*PagingRequest)

type Page[T any] struct {
	PageNumber   int `json:"pageNumber"`
	PageSize     int `json:"pageSize"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"totalRecords"`
	Items        []T `json:"items"`
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
