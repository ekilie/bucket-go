package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ekilie/bucket-go/model"
	"github.com/ekilie/bucket-go/util"
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

// UploadFile uploads a file to the API and returns the response.
func (c *Client) UploadFile(filePath string) (*model.UploadResponse, error) {
	// Validate file existence and size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	if fileInfo.Size() > util.MaxFileSize {
		return nil, errors.New("file exceeds maximum size of 100MB")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if !util.AllowedExtensions[ext] {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add apikey field
	if err := writer.WriteField("apikey", c.APIKey); err != nil {
		return nil, fmt.Errorf("failed to write apikey: %w", err)
	}

	// Add file field
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	// Create request
	url := c.BaseURL + util.Endpoint
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-OK status: %d - %s", resp.StatusCode, string(respBody))
	}

	// Parse response
	var baseResp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(respBody, &baseResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	switch baseResp.Status {
	case "success":
		var uploadResp model.UploadResponse
		if err := json.Unmarshal(respBody, &uploadResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
		}
		return &uploadResp, nil
	case "error":
		var errResp model.ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		return nil, errors.New(errResp.Message)
	}

	return nil, errors.New("unknown response status")
}
