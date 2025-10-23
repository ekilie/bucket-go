package model

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
