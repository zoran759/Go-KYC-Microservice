package documents

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNewService(t *testing.T) {
	testService := service{
		host:   "test_host",
		apiKey: "api_key",
	}

	service := NewService(Config{
		Host:   "test_host",
		APIKey: "api_key",
	})

	assert.Equal(t, testService, service)
}

func Test_service_UploadDocument(t *testing.T) {
	documentsService := service{
		host:   "https://test_host.subsub.com",
		apiKey: "api_key",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants/test_applicant_id/info/idDoc?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				request.FormValue("metadata"),
			), nil
		},
	)

	testMetadata := Metadata{
		DocumentType: "PASSPORT",
		Country:      "ALB",
	}

	response, err := documentsService.UploadDocument(
		"test_applicant_id",
		Document{
			Metadata: testMetadata,
			File: File{
				Data: []byte{10, 11, 12, 123, 21, 12, 32, 23, 54, 35, 234, 53, 23, 43, 234, 234, 12, 12, 12, 12, 34, 23, 43, 52},
			},
		},
	)

	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, testMetadata, *response)
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants/test_applicant_id/info/idDoc?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			// Return metadata that is not matching request metadata to test that we return metadata from response to consumer
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"idDocType":"PASPORT","country":"RUS"}`,
			), nil
		},
	)

	response, err = documentsService.UploadDocument(
		"test_applicant_id",
		Document{
			Metadata: testMetadata,
			File: File{
				Data: []byte{10, 11, 12, 123, 21, 12, 32, 23, 54, 35, 234, 53, 23, 43, 234, 234, 12, 12, 12, 12, 34, 23, 43, 52},
			},
		},
	)

	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, Metadata{
			DocumentType: "PASPORT",
			Country:      "RUS",
		}, *response)
	}
}

func Test_service_UploadDocumentError(t *testing.T) {
	documentsService := service{
		host:   "https://test_host.subsub.com",
		apiKey: "api_key",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants/test_applicant_id/info/idDoc?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"description": "Cannot read a metadata object from the body","code": 400}`,
			), nil
		},
	)

	response, err := documentsService.UploadDocument("test_applicant_id", Document{})
	if assert.Error(t, err) && assert.Nil(t, response) {
		assert.Equal(t, "Cannot read a metadata object from the body", err.Error())
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants/test_applicant_id/info/idDoc?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`,
			), nil
		},
	)

	response, err = documentsService.UploadDocument("test_applicant_id", Document{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants/test_applicant_id/info/idDoc?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = documentsService.UploadDocument("test_applicant_id", Document{})
	assert.Error(t, err)
	assert.Nil(t, response)
}
