package verification

import (
	"gitlab.com/lambospeed/kyc/http"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"encoding/json"
)

type service struct {
	config Config
}

func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

func (service service) Verify(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*VerificationResponse, error) {
	request := StartVerificationRequest{
		AcceptTruliooTermsAndConditions: true,
		ConfigurationName:               "Identity Verification",
		ConsentForDataSources:           consents,
		CountryCode:                     countryAlpha2,
		DataFields:                      fields,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(
		service.config.Host+"/verify",
		http.Headers{
			"Authorization": "Basic " + service.config.Token,
			"Content-Type":  "application/json",
		},
		requestBytes,
	)
	if err != nil {
		return nil, err
	}

	response := new(VerificationResponse)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}
	return response, err
}
