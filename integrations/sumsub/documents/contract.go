package documents

type Config struct {
	Host   string
	APIKey string
}

const FrontSide = "FRONT_SIDE"
const BackSide = "BACK_SIDE"

type Documents interface {
	UploadDocument(
		applicantID string,
		document Document,
	) (*Metadata, error)
}
type Mock struct {
	UploadDocumentFn func(
		applicantID string,
		document Document,
	) (*Metadata, error)
}

func (mock Mock) UploadDocument(
	applicantID string,
	document Document,
) (*Metadata, error) {
	return mock.UploadDocumentFn(applicantID, document)
}
