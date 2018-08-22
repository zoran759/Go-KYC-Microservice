package documents

type Document struct {
	Metadata
	File
}

type Metadata struct {
	DocumentType    string  `json:"idDocType"`
	DocumentSubType *string `json:"idDocSubType"`
	Country         string  `json:"country"`
	FirstName       *string `json:"firstName"`
	MiddleName      *string `json:"middleName"`
	LastName        *string `json:"lastName"`
	DateIssued      *string `json:"issuedDate"`
	ValidUntil      *string `json:"validUntil"`
	Number          *string `json:"number"`
	DateOfBirth     *string `json:"dob"`
	PlaceOfBirth    *string `json:"placeOfBirth"`
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
