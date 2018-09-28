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

func NewService(config Config) Verification {
	return service{
		host:   config.Host,
		apiKey: config.APIKey,
	}
}

func (service service) StartVerification(applicantID string) (bool, error) {
	_, responseBytes, err := http.Post(fmt.Sprintf("%s/resources/applicants/%s/status/pending?key=%s",
		service.host,
		applicantID,
		service.apiKey,
	),
		http.Headers{},
		make([]byte, 0),
	)
	if err != nil {
		return false, err
	}

	response := new(StartVerificationResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return false, err
	}
	if response.Error.Description != nil {
		return false, errors.New(*response.Error.Description)
	}

	return response.OK == 1, nil
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
		return "", nil, errors.New(*response.Error.Description)
	}

	return response.ReviewStatus, &response.ReviewResult, nil
}
