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
	"github.com/ekilie/bucket-go/client"
	"github.com/ekilie/bucket-go/model"
	"github.com/ekilie/bucket-go/util"
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

// If baseURL is empty, it defaults to BaseURL.
// ...existing code...

// UploadFile uploads a file to the API and returns the response.
func (c *Client) UploadFile(filePath string) (*UploadResponse, error) {
	// ...existing code...

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