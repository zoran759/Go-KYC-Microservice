package model

import "strings"

// List of FileType values.
const (
	FileID       FileType = "file_id"
	FileFunds    FileType = "file_funds"
	FileBoard    FileType = "file_board"
	FileAddress  FileType = "file_address"
	FileRegister FileType = "file_register"
	FileActivity FileType = "file_activity"
	FileIncome   FileType = "file_income"
)

// FileType represents an accepted file type.
type FileType string

// File represents a file intended to insert to KYC process.
type File struct {
	Type       FileType `json:"type"`
	Extension  string   `json:"extension"`
	DataBase64 string   `json:"data_base64"`
}

// acceptedExt holds all accepted file extensions.
var acceptedExt = map[string]bool{
	"jpg": true,
	"png": true,
	"pdf": true,
}

// IsAcceptedFileExt tests the file extension for acceptance by the API.
func IsAcceptedFileExt(ext string) bool {
	return acceptedExt[ext]
}

// NormalizeFileExt normalizes the file extension.
func NormalizeFileExt(ext string) string {
	ext = strings.ToLower(ext)
	if ext == "jpeg" {
		ext = "jpg"
	}
	return ext
}
