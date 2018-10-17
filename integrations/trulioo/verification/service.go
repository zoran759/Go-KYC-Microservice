package verification

import (
	"encoding/json"
	"modulus/kyc/http"
	"modulus/kyc/integrations/trulioo/configuration"
	stdhttp "net/http"
)

type service struct {
	config Config
}

// NewService constructs a new verification service object.
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

	code, responseBytes, err := http.Post(
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
	if code != stdhttp.StatusOK && code != 0 {
		response.ErrorCode = &code
	}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return response, err
	}
	return response, err
}
