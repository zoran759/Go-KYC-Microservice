package documents

type Document struct {
	Metadata
	File
}

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

type File struct {
	Data        []byte
	Filename    string
	ContentType string
}

type Error struct {
	Description *string `json:"description"`
}

type UploadDocumentResponse struct {
	Metadata
	Error
}
