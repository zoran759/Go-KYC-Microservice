package configuration

import (
	"encoding/json"
	"errors"
	"modulus/kyc/http"
	stdhttp "net/http"
)

type service struct {
	config Config
}

// NewService constructs new configuration service object.
func NewService(config Config) Configuration {
	return service{
		config: config,
	}
}

func (service service) Consents(countryAlpha2 string) (Consents, *int, error) {
	if countryAlpha2 == "" {
		return nil, nil, errors.New("No country code provided")
	}
	code, responseBytes, err := http.Get(
		service.config.Host+"/consents/Identity Verification/"+countryAlpha2,
		http.Headers{
			"Authorization": "Basic " + service.config.Token,
		})

	if err != nil {
		return nil, nil, err
	}

	var errorCode *int
	if code != stdhttp.StatusOK {
		if code != 0 {
			errorCode = &code
		}
		errResp := new(Error)
		if err := json.Unmarshal(responseBytes, errResp); err != nil {
			return nil, errorCode, err
		}

		if errResp.Message != "" {
			return nil, errorCode, errResp
		}

		return nil, errorCode, errors.New("Unknown error")
	}

	consents := make(Consents, 0)

	return consents, nil, json.Unmarshal(responseBytes, &consents)
}
