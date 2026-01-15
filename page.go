package onspring

// Page represents a paginated response from the Onspring API.
type Page[T any] struct {
	PageNumber   int `json:"pageNumber"`
	PageSize     int `json:"pageSize"`
	TotalPages   int `json:"totalPages"`
	TotalRecords int `json:"totalRecords"`
	Items        []T `json:"items"`
}
