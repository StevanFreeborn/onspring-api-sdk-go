package onspring

import "strconv"

// PagingRequest contains pagination parameters for API requests.
type PagingRequest struct {
	// The page number to retrieve
	PageNumber int
	// The size of pages to retrieve
	PageSize int
}

// ToParams converts the paging request to a map of query parameters.
func (pr *PagingRequest) ToParams() map[string]string {
	return map[string]string{
		"pageNumber": strconv.Itoa(pr.PageNumber),
		"pageSize":   strconv.Itoa(pr.PageSize),
	}
}

// PagingOption is a function that modifies a PagingRequest.
type PagingOption func(*PagingRequest)

// ForPageNumber sets the page number for a paging request.
func ForPageNumber(pageNumber int) PagingOption {
	return func(pr *PagingRequest) {
		pr.PageNumber = pageNumber
	}
}

// WithPageSize sets the page size for a paging request.
func WithPageSize(pageSize int) PagingOption {
	return func(pr *PagingRequest) {
		pr.PageSize = pageSize
	}
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
