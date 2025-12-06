package onspring

import "net/http"

func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) BaseURL() string {
	return c.baseURL
}

func (c *Client) APIKey() string {
	return c.apiKey
}

func (c *Client) APIVersion() string {
	return c.apiVersion
}
