package applicants

import (
	"encoding/json"
	"fmt"
	"modulus/kyc/http"

	"github.com/pkg/errors"
)

type service struct {
	host   string
	apiKey string
}

func NewService(config Config) Applicants {
	return service{
		host:   config.Host,
		apiKey: config.APIKey,
	}
}

func (service service) CreateApplicant(email string, applicant ApplicantInfo) (*CreateApplicantResponse, error) {
	requestBytes, err := json.Marshal(CreateApplicantRequest{
		Email: email,
		Info:  applicant,
	})
	if err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(
		fmt.Sprintf("%s/resources/applicants?key=%s",
			service.host,
			service.apiKey,
		),
		http.Headers{
			"Content-Type": "application/json",
		},
		requestBytes,
	)
	if err != nil {
		return nil, err
	}

	response := new(CreateApplicantResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	if response.Error.Description != nil {
		return response, errors.New(*response.Error.Description)
	}

	return response, nil
}
