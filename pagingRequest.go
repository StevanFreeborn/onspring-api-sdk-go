package onspring

import "strconv"

type PagingRequest struct {
	pageNumber int
	pageSize   int
}

func (pr *PagingRequest) ToParams() map[string]string {
	return map[string]string{
		"pageNumber": strconv.Itoa(pr.pageNumber),
		"pageSize":   strconv.Itoa(pr.pageSize),
	}
}

type PagingOption func(*PagingRequest)

func ForPageNumber(pageNumber int) PagingOption {
	return func(pr *PagingRequest) {
		pr.pageNumber = pageNumber
	}
}

func WithPageSize(pageSize int) PagingOption {
	return func(pr *PagingRequest) {
		pr.pageSize = pageSize
	}
}
