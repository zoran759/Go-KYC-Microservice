package configuration

import (
	"gitlab.com/lambospeed/kyc/http"
	"encoding/json"
	"github.com/pkg/errors"
	stdhttp "net/http"
)

type service struct {
	config Config
}

func NewService(config Config) Configuration {
	return service{
		config: config,
	}
}

func (service service) Consents(countryAlpha2 string) (Consents, error) {
	if countryAlpha2 == "" {
		return nil, errors.New("No country code provided")
	}
	code, responseBytes, err := http.Get(
		service.config.Host+"/consents/Identity Verification/"+countryAlpha2,
		http.Headers{
			"Authorization": "Basic " + service.config.Token,
		})

	if err != nil {
		return nil, err
	}

	if code != stdhttp.StatusOK {
		errResp := new(Error)
		if err := json.Unmarshal(responseBytes, errResp); err != nil {
			return nil, err
		}

		if errResp.Message != "" {
			return nil, errResp
		}

		return nil, errors.New("Unknown error")
	}

	consents := make(Consents, 0)

	return consents, json.Unmarshal(responseBytes, &consents)
}
