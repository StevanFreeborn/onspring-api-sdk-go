package onspring

import "strconv"

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
