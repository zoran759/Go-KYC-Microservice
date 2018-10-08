package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"modulus/kyc/http"
)

type service struct {
	host   string
	apiKey string
}

// NewService constructs a new verification service object.
func NewService(config Config) Verificator {
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
