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
func NewService(config Config) Verification {
	return service{
		host:   config.Host,
		apiKey: config.APIKey,
	}
}

func (service service) StartVerification(applicantID string) (bool, *int, error) {
	//
	// TODO: WTF this does here???!!!
	//
	// https://developers.sumsub.com/#requesting-an-applicant-re-check
	//
	// You CAN programatically ask us to re-check an applicant
	//
	// IN CASE YOU OR YOUR USER BELIEVE THAT OUR SYSTEM MADE A MISTAKE,
	//
	// or you want to request some additional checks AGREED WITH US IN ADVANCE.
	//
	_, responseBytes, err := http.Post(fmt.Sprintf("%s/resources/applicants/%s/status/pending?key=%s",
		service.host,
		applicantID,
		service.apiKey,
	),
		http.Headers{},
		make([]byte, 0),
	)
	if err != nil {
		return false, nil, err
	}

	response := new(StartVerificationResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return false, nil, err
	}
	if response.Error.Description != nil {
		return false, response.Error.Code, errors.New(*response.Error.Description)
	}

	return response.OK == 1, nil, nil
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
