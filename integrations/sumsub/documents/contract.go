package documents

// Config represents the configuration of the service.
type Config struct {
	Host   string
	APIKey string
}

// The document subtype values.
const (
	FrontSide = "FRONT_SIDE"
	BackSide  = "BACK_SIDE"
)

// Documents represents the documents uploading interface.
type Documents interface {
	UploadDocument(
		applicantID string,
		document Document,
	) (*Metadata, *int, error)
}

// Mock represents the mock of the service for the testing.
type Mock struct {
	UploadDocumentFn func(
		applicantID string,
		document Document,
	) (*Metadata, *int, error)
}

// UploadDocument implements the Documents interface for Mock.
func (mock Mock) UploadDocument(
	applicantID string,
	document Document,
) (*Metadata, *int, error) {
	return mock.UploadDocumentFn(applicantID, document)
}
