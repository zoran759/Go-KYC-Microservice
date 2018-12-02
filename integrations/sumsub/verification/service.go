package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/http"
)

type service struct {
	host   string
	apiKey string
}

// NewService constructs a new verification service object.
func NewService(config Config) Verification {
	return service{
		host:   config.Host,
		apiKey: config.APIKey,
	}
}

func (service service) CheckApplicantStatus(applicantID string) (string, *ReviewResult, error) {
	_, responseBytes, err := http.Get(fmt.Sprintf("%s/resources/applicants/%s/status?key=%s",
		service.host,
		applicantID,
		service.apiKey,
	),
		http.Headers{},
	)
	if err != nil {
		return "", nil, err
	}

	response := new(ApplicantStatusResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return "", nil, err
	}
	if response.Error.Description != nil {
		return "", &ReviewResult{ErrorCode: *response.Error.Code}, errors.New(*response.Error.Description)
	}

	return response.ReviewStatus, &response.ReviewResult, nil
}

func (service service) RequestApplicantCheck(applicantID string) (err error) {
	code, responseBytes, err := http.Post(fmt.Sprintf("%s/resources/applicants/%s/status/pending?reason=docs_sent&key=%s",
		service.host, applicantID, service.apiKey), http.Headers{}, nil)
	if err != nil {
		return
	}

	if code == stdhttp.StatusOK {
		return
	}

	response := &ApplicantStatusResponse{}
	if err = json.Unmarshal(responseBytes, response); err != nil {
		return
	}

	return response.Error
}
