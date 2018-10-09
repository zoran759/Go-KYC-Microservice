package documents

// Document represents the document for the verification.
type Document struct {
	Metadata
	File
}

// Metadata represents the part of the Document model.
type Metadata struct {
	DocumentType    string `json:"idDocType"`
	DocumentSubType string `json:"idDocSubType,omitempty"`
	Country         string `json:"country"`
	FirstName       string `json:"firstName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	DateIssued      string `json:"issuedDate,omitempty"`
	ValidUntil      string `json:"validUntil,omitempty"`
	Number          string `json:"number,omitempty"`
	DateOfBirth     string `json:"dob,omitempty"`
	PlaceOfBirth    string `json:"placeOfBirth,omitempty"`
}

// File represents the part of the Document model.
type File struct {
	Data        []byte
	Filename    string
	ContentType string
}

// Error represents the error from the API.
type Error struct {
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}

// UploadDocumentResponse represents the response on documents uploading.
type UploadDocumentResponse struct {
	Metadata
	Error
}
