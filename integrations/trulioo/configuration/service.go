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
		errorCode = &code
		err := &Error{}
		if len(responseBytes) == 0 {
			err.Message = "Unknown error"
		} else {
			err1 := json.Unmarshal(responseBytes, err)
			if err1 != nil {
				err.Message = string(responseBytes)
			}
		}
		return nil, errorCode, err
	}

	consents := make(Consents, 0)

	return consents, nil, json.Unmarshal(responseBytes, &consents)
}
