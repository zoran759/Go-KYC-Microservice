package verification

import (
	"encoding/json"
	"gitlab.com/lambospeed/kyc/http"
)

type service struct {
	config Config
}

func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

func (service service) CreateUser(request CreateUserRequest) (*UserResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	_, responseBytes, err := http.Post(
		service.config.Host,
		http.Headers{
			"X-SP-GATEWAY": service.config.ClientID + "|" + service.config.ClientSecret,
			"X-SP-USER":    "|",
			"X-SP-USER-ID": "0.0.0.0",
		},
		requestBytes,
	)
	if err != nil {
		return nil, err
	}

	response := new(UserResponse)

	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (service service) GetUser(userID string) (*UserResponse, error) {
	_, responseBytes, err := http.Get(
		service.config.Host+"/"+userID,
		http.Headers{
			"X-SP-GATEWAY": service.config.ClientID + "|" + service.config.ClientSecret,
			"X-SP-USER":    "|",
			"X-SP-USER-ID": "0.0.0.0",
		},
	)
	if err != nil {
		return nil, err
	}

	response := new(UserResponse)

	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	return response, nil
}
