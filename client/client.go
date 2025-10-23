package client

import (
	"net/http"
)

// Client is the API client for uploading files.
type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

// NewClient creates a new Client instance.
// If baseURL is empty, it defaults to BaseURL.
func NewClient(apiKey string, baseURL ...string) *Client {
	url := "https://bucket.ekilie.com"
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	}
	return &Client{
		APIKey:  apiKey,
		BaseURL: url,
		HTTP:    &http.Client{},
	}
}
