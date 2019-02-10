package common

// List of available values of DocumentType.
const (
	IDCardType   DocumentType = "idcard"
	PassportType DocumentType = "passport"
)

// DocumentType defines a document type.
type DocumentType string

// Document represents a document.
type Document struct {
	Type          DocumentType
	Number        string
	CountryAlpha2 string
	Image         *DocumentFile
}
