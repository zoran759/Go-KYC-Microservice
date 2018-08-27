package documents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/lambospeed/kyc/http"
	"mime/multipart"
)

type service struct {
	host   string
	apiKey string
}

func NewService(config Config) Documents {
	return service{
		host:   config.Host,
		apiKey: config.APIKey,
	}
}

func (service service) UploadDocument(
	applicantID string,
	document Document,
) (*Metadata, error) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := http.CreateFormFile(writer, "content", document.File.Filename, document.File.ContentType)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(document.File.Data); err != nil {
		return nil, err
	}

	metadataBytes, err := json.Marshal(document.Metadata)
	if err != nil {
		return nil, err
	}

	if err := writer.WriteField("metadata", string(metadataBytes)); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(fmt.Sprintf("%s/resources/applicants/%s/info/idDoc?key=%s",
		service.host,
		applicantID,
		service.apiKey,
	), http.Headers{
		"Content-Type": writer.FormDataContentType(),
	}, body.Bytes())
	if err != nil {
		return nil, err
	}

	response := new(UploadDocumentResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}
	if response.Error.Description != nil {
		return nil, errors.New(*response.Error.Description)
	}

	return &response.Metadata, nil
}
