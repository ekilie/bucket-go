# Ekilie Bucket Go Client

This package provides a reusable Go client for uploading files to the Ekilie Bucket API. It handles file validation, multipart uploads, and response parsing.

## Features

- Simple client for Ekilie Bucket API
- File type and size validation
- Multipart file upload
- Structured response and error handling

## Installation

Clone the repository and use Go modules to import the packages:

```
git clone https://github.com/ekilie/bucket-go.git
cd bucket-go
go mod tidy
```

## Usage

### 1. Create a client

```go
import "github.com/ekilie/bucket-go/client"

apiKey := "your-api-key"
c := client.NewClient(apiKey)
```

### 2. Upload a file

```go
import (
	"github.com/ekilie/bucket-go/store"
	"github.com/ekilie/bucket-go/client"
)

resp, err := store.UploadFile(c, "/path/to/file.jpg")
if err != nil {
	// handle error
}
fmt.Println("File URL:", resp.URL)
fmt.Printf("Metadata: %+v\n", resp.Metadata)
```

### 3. Full Example

See `main.go` for a complete demo:

```go
package main

import (
	"fmt"
	"os"
	"github.com/ekilie/bucket-go/client"
	"github.com/ekilie/bucket-go/store"
)

func main() {
	apiKey := "your-api-key"
	filePath := "sample.jpg"

	c := client.NewClient(apiKey)
	resp, err := store.UploadFile(c, filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("File uploaded successfully! URL: %s\n", resp.URL)
	fmt.Printf("Metadata: %+v\n", resp.Metadata)
}
```

## API Reference

- `client.NewClient(apiKey string, baseURL ...string) *Client` - Create a new API client
- `store.UploadFile(c *client.Client, filePath string) (*model.UploadResponse, error)` - Upload a file

## File Validation

- Maximum file size: 100MB
- Allowed extensions: jpg, jpeg, png, gif, webp, svg, pdf, txt, doc, docx, xls, xlsx, ppt, pptx, zip, rar, tar, gz, json, xml

## License

MIT

# bucket-go

ekilie bucket go client
