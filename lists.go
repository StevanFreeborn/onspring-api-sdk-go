package onspring

import (
	"context"
	"fmt"
	"net/http"
)

const (
	listsPath = "/lists"
)

// ListsEndpoint provides access to lists in an Onspring instance.
type ListsEndpoint struct {
	client *Client
}

// SaveListItemRequest represents a request to save a list item.
type SaveListItemRequest struct {
	Id           string   `json:"id,omitempty"`
	Name         string   `json:"name"`
	NumericValue *float64 `json:"numericValue,omitempty"`
	Color        *string  `json:"color,omitempty"`
}

// SaveListItemResponse represents the response for saving a list item.
type SaveListItemResponse struct {
	Id string `json:"id"`
}

// Save creates or updates a list item in the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - listId: The id of the list to save the item to
//   - item: The list item to save
//
// Returns:
//   - SaveListItemResponse: The response containing the saved list item's id
//   - error: An error if the request fails
func (l *ListsEndpoint) Save(ctx context.Context, listId int, item SaveListItemRequest) (SaveListItemResponse, error) {
	path := fmt.Sprintf("%s/id/%d/items", listsPath, listId)
	req, requestCreationErr := l.client.newRequest(ctx, http.MethodPut, path, nil, item)

	var response SaveListItemResponse

	if requestCreationErr != nil {
		return response, requestCreationErr
	}

	responseErr := l.client.doWithJsonResponse(req, &response)

	if responseErr != nil {
		return response, responseErr
	}

	return response, nil
}

// Delete removes a list item from the Onspring API.
//
// Parameters:
//   - ctx: The context for the request
//   - listId: The id of the list
//   - itemId: The id of the list item to delete
//
// Returns:
//   - error: An error if the request fails
func (l *ListsEndpoint) Delete(ctx context.Context, listId int, itemId string) error {
	path := fmt.Sprintf("%s/id/%d/itemId/%s", listsPath, listId, itemId)
	req, requestCreationErr := l.client.newRequest(ctx, http.MethodDelete, path, nil, nil)

	if requestCreationErr != nil {
		return requestCreationErr
	}

	return l.client.do(req)
}
