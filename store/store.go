// Package store provides a reusable client for uploading files to the Ekilie Bucket API.
// It handles file validation, multipart uploads, and response parsing.
// Usage:
//   client := store.NewClient("your-api-key")
//   resp, err := client.UploadFile("/path/to/file.jpg")
//   if err != nil {
//       // handle error
//   }
//   fmt.Println(resp.URL)

package store

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
	"time"
)

// BaseURL is the default base URL for the Ekilie Bucket API.
const BaseURL = "https://bucket.ekilie.com"

// Endpoint is the API endpoint for file uploads.
const Endpoint = "/api/store/v1/index.php"

// MaxFileSize is the maximum allowed file size in bytes (100MB).
const MaxFileSize = 100 * 1024 * 1024

// AllowedExtensions is the set of allowed file extensions (lowercase).
var AllowedExtensions = map[string]bool{
	".jpg":   true,
	".jpeg":  true,
	".png":   true,
	".gif":   true,
	".webp":  true,
	".svg":   true,
	".pdf":   true,
	".txt":   true,
	".doc":   true,
	".docx":  true,
	".xls":   true,
	".xlsx":  true,
	".ppt":   true,
	".pptx":  true,
	".zip":   true,
	".rar":   true,
	".tar":   true,
	".gz":    true,
	".json":  true,
	".xml":   true,
}

// Client is the API client for uploading files.
type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

// NewClient creates a new Client instance.
// If baseURL is empty, it defaults to BaseURL.
func NewClient(apiKey string, baseURL ...string) *Client {
	url := BaseURL
	if len(baseURL) > 0 && baseURL[0] != "" {
		url = baseURL[0]
	}
	return &Client{
		APIKey:  apiKey,
		BaseURL: url,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

// UploadResponse represents the successful API response.
type UploadResponse struct {
	Status   string   `json:"status"`
	URL      string   `json:"url"`
	Metadata Metadata `json:"metadata"`
}

// Metadata contains file details.
type Metadata struct {
	OriginalName string `json:"original_name"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`
	UploadTime   string `json:"upload_time"`
}

// ErrorResponse represents the error API response.
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// UploadFile uploads a file to the API and returns the response.
func (c *Client) UploadFile(filePath string) (*UploadResponse, error) {
	// Validate file existence and size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}
	if fileInfo.Size() > MaxFileSize {
		return nil, errors.New("file exceeds maximum size of 100MB")
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(filePath))
	if !AllowedExtensions[ext] {
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
	url := c.BaseURL + Endpoint
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

	if baseResp.Status == "success" {
		var uploadResp UploadResponse
		if err := json.Unmarshal(respBody, &uploadResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal success response: %w", err)
		}
		return &uploadResp, nil
	} else if baseResp.Status == "error" {
		var errResp ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error response: %w", err)
		}
		return nil, errors.New(errResp.Message)
	}

	return nil, errors.New("unknown response status")
}