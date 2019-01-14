package model

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
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"gif":  true,
	"bmp":  true,
	"svg":  true,
	"psd":  true,
	"tif":  true,
	"tiff": true,
	"webp": true,
	"pdf":  true,
}

// IsAcceptedFileExt tests the file extension for acceptance by the API.
func IsAcceptedFileExt(ext string) bool {
	return acceptedExt[ext]
}
