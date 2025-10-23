package main

import (
	"fmt"
	"os"

	"github.com/ekilie/bucket-go/client"
	"github.com/ekilie/bucket-go/store"
)

func main() {
	apiKey := "your-api-key"
	filePath := "README.md"

	c := client.NewClient(apiKey)
	resp, err := store.UploadFile(c, filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Upload failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("File uploaded successfully! URL: %s\n", resp.URL)
	fmt.Printf("Metadata: %+v\n", resp.Metadata)
}
