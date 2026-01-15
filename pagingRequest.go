package onspring

import "strconv"

// PagingRequest contains pagination parameters for API requests.
type PagingRequest struct {
	pageNumber int
	pageSize   int
}

// ToParams converts the paging request to a map of query parameters.
func (pr *PagingRequest) ToParams() map[string]string {
	return map[string]string{
		"pageNumber": strconv.Itoa(pr.pageNumber),
		"pageSize":   strconv.Itoa(pr.pageSize),
	}
}

// PagingOption is a function that modifies a PagingRequest.
type PagingOption func(*PagingRequest)

// ForPageNumber sets the page number for a paging request.
func ForPageNumber(pageNumber int) PagingOption {
	return func(pr *PagingRequest) {
		pr.pageNumber = pageNumber
	}
}

// WithPageSize sets the page size for a paging request.
func WithPageSize(pageSize int) PagingOption {
	return func(pr *PagingRequest) {
		pr.pageSize = pageSize
	}
}
