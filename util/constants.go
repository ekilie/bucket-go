package util

// BaseURL is the default base URL for the Ekilie Bucket API.
const BaseURL = "https://bucket.ekilie.com"

// Endpoint is the API endpoint for file uploads.
const Endpoint = "/api/store/v1/index.php"

// MaxFileSize is the maximum allowed file size in bytes (100MB).
const MaxFileSize = 100 * 1024 * 1024

// AllowedExtensions is the set of allowed file extensions (lowercase).
var AllowedExtensions = map[string]bool{
	// Images
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,

	// Documents
	".pdf":  true,
	".txt":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,

	// Archives
	".zip":  true,
	".rar":  true,
	".tar":  true,
	".gz":   true,

	// Data
	".json": true,
	".xml":  true,

	// Audio
	".mp3":  true,
	".wav":  true,
	".m4a":  true,
	".aac":  true,
	".ogg":  true,
	".oga":  true,
	".flac": true,
	".opus": true,
	".webm": true,
}
