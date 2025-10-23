package util

// BaseURL is the default base URL for the Ekilie Bucket API.
const BaseURL = "https://bucket.ekilie.com"

// Endpoint is the API endpoint for file uploads.
const Endpoint = "/api/store/v1/index.php"

// MaxFileSize is the maximum allowed file size in bytes (100MB).
const MaxFileSize = 100 * 1024 * 1024

// AllowedExtensions is the set of allowed file extensions (lowercase).
var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
	".pdf":  true,
	".txt":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".zip":  true,
	".rar":  true,
	".tar":  true,
	".gz":   true,
	".json": true,
	".xml":  true,
}
