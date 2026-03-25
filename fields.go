package onspring

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
)

const (
	fieldsPath = "/fields"
)

// FieldsEndpoint provides access to fields in an Onspring instance.
type FieldsEndpoint struct {
	client *Client
}

// Field represents an Onspring field
type Field struct {
	Id         int    `json:"id"`
	AppId      int    `json:"appId"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Status     string `json:"status"`
	IsRequired bool   `json:"isRequired"`
	IsUnique   bool   `json:"isUnique"`
	TypeData   any    `json:"-"`
}

// FormulaField represents a formula field type in Onspring
type FormulaField struct {
	OutputType string   `json:"outputType"`
	Values     []string `json:"values"`
}

// ReferenceField represents a reference field type in Onspring
type ReferenceField struct {
	Multiplicity   string `json:"multiplicity"`
	ReferenceAppId string `json:"referenceAppId"`
}

// ListField represents a list field type in Onspring
type ListField struct {
	Multiplicity string   `json:"multiplicity"`
	Values       []string `json:"values"`
	ListId       int      `json:"listId"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for Field.
// This allows for custom deserialization logic based on the 'Type' field.
func (f *Field) UnmarshalJSON(data []byte) error {
	type Alias Field

	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(f),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	switch f.Type {
	case "Formula":
		var formulaField FormulaField

		if err := json.Unmarshal(data, &formulaField); err != nil {
			return err
		}

		f.TypeData = formulaField
	case "Reference":
		var referenceField ReferenceField

		if err := json.Unmarshal(data, &referenceField); err != nil {
			return err
		}

		f.TypeData = referenceField
	case "List":
		var listField ListField

		if err := json.Unmarshal(data, &listField); err != nil {
			return err
		}

		f.TypeData = listField
	default:
	}

	return nil
}

// FieldBatch represents a batch of Onspring fields
type FieldBatch struct {
	Count int     `json:"count"`
	Items []Field `json:"items"`
}

// Get retrieves a field from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - id: The id of the field to retrieve
//
// Returns:
//   - Field: A field
//   - error: An error if the request fails
func (f *FieldsEndpoint) Get(ctx context.Context, id int) (Field, error) {
	path := fmt.Sprintf("%s/id/%d", fieldsPath, id)
	req, requestCreationErr := f.client.newRequest(ctx, http.MethodGet, path, nil, nil)

	var field Field

	if requestCreationErr != nil {
		return field, requestCreationErr
	}

	responseErr := f.client.doWithJsonResponse(req, &field)

	if responseErr != nil {
		return field, responseErr
	}

	return field, nil
}

// GetMany retrieves a batch of fields from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - ids: The ids of the fields to retrieve
//
// Returns:
//   - FieldBatch: A batch of fields
//   - error: An error if the request fails
func (f *FieldsEndpoint) GetMany(ctx context.Context, ids []int) (FieldBatch, error) {
	path := fmt.Sprintf("%s/batch-get", fieldsPath)
	req, requestCreationErr := f.client.newRequest(ctx, http.MethodPost, path, nil, ids)

	var fieldBatch FieldBatch

	if requestCreationErr != nil {
		return fieldBatch, requestCreationErr
	}

	responseErr := f.client.doWithJsonResponse(req, &fieldBatch)

	if responseErr != nil {
		return fieldBatch, responseErr
	}

	return fieldBatch, nil
}

// List retrieves a paginated list of fields from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve fields for
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - Page[Field]: A page of fields with pagination metadata
//   - error: An error if the request fails
func (f *FieldsEndpoint) List(ctx context.Context, appId int, pagingOpts ...PagingOption) (Page[Field], error) {
	pagingRequest := createPagingRequest(pagingOpts)
	path := fmt.Sprintf("%s/appId/%d", fieldsPath, appId)

	req, requestCreationErr := f.client.newRequest(ctx, http.MethodGet, path, pagingRequest.ToParams(), nil)

	var page Page[Field]

	if requestCreationErr != nil {
		return page, requestCreationErr
	}

	responseErr := f.client.doWithJsonResponse(req, &page)

	if responseErr != nil {
		return page, responseErr
	}

	return page, nil
}

// ListAll returns an iterator that yields all Fields for an app across all pages.
// It automatically handles pagination by making sequential calls to List
// until all items have been retrieved or the caller stops the iteration.
//
// The iterator yields each Field and any error encountered during fetching.
// If an error occurs during a page request, the error is yielded and
// iteration terminates.
//
// Parameters:
//   - ctx: The context for the request
//   - appId: The id of the app to retrieve fields for
//   - pagingOpts: Optional paging configuration functions (e.g., ForPageNumber, WithPageSize)
//
// Returns:
//   - iter.Seq2[Field, error]: An iterator yielding:
//   - Field: The individual field record.
//   - error: An error if a specific page request fails during iteration.
func (f *FieldsEndpoint) ListAll(ctx context.Context, appId int, pagingOpts ...PagingOption) iter.Seq2[Field, error] {
	return func(yield func(Field, error) bool) {
		pagingRequest := createPagingRequest(pagingOpts)

		for {
			page, err := f.List(ctx, appId, ForPageNumber(pagingRequest.PageNumber), WithPageSize(pagingRequest.PageSize))

			if err != nil {
				yield(Field{}, err)
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
