package common

// List of available values of DocumentType.
const (
	IDCardType        DocumentType = "idcard"
	PassportType      DocumentType = "passport"
	DriverLicenseType DocumentType = "drivers"
	CreditCardType    DocumentType = "credit_card"
	DebitCardType     DocumentType = "debit_card"
)

// DocumentType defines a document type.
type DocumentType string

// Document represents a document.
type Document struct {
	Type          DocumentType
	Number        string
	CountryAlpha2 string
	IssuedDate    Time
	ValidUntil    Time
	Image         *DocumentFile
}
