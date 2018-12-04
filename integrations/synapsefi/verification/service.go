package verification

import (
	"encoding/json"
	"log"

	"modulus/kyc/http"
)

const (
	endpointUsers = "users"
	endpointOauth = "oauth"
	appLanguage   = "en"
)

type service struct {
	config Config
}

// NewService constructs the new verification service object.
func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

// TODO: resend on fail, process errors, PROCESS RESPONSE PERMISSIONS!!!
func (service service) CreateUser(request CreateUserRequest) (*UserResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	headers := service.composeHeaders(true, "")
	host := service.config.Host + endpointUsers

	log.Printf("Request user: %+v\n\nHeaders: %+v\n\nHost: %+v\n\n", request.Logins, headers, host)

	responseStatus, responseBytes, err := http.Post(host, headers, requestBytes)

	if err != nil {
		return nil, err
	}

	if responseStatus != 200 {
		err, _ := MapResponseError(responseBytes)
		return nil, err
	}

	response := &UserResponse{}

	if err := json.Unmarshal(responseBytes, response); err != nil {
		log.Printf("Error decoding SynapseFi response: %v", err)
		log.Printf("Response: %q", response)
		return nil, err
	}

	return response, nil
}

func (service service) AddDocument(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	log.Println("Adding documents...")

	headers := service.composeHeaders(false, userOAuth)
	host := service.config.Host + endpointUsers + "/" + userID

	log.Printf("Request: %+v\n\nHeaders: %+v\n\nHost: %+v\n\n", request, headers, host)

	responseStatus, responseBytes, err := http.Request("PATCH", host, headers, requestBytes)

	if err != nil {
		return nil, err
	}

	if responseStatus != 200 {
		err, _ := MapResponseError(responseBytes)
		return nil, err
	}

	response := &UserResponse{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		log.Printf("Error decoding SynapseFi response: %v", err)
		log.Printf("Response: %q", response)
		return nil, err
	}

	return response, nil
}

func (service service) GetUser(userID string) (*UserResponse, error) {

	headers := service.composeHeaders(true, "")
	_, responseBytes, err := http.Get(
		service.config.Host+endpointUsers+"/"+userID,
		headers,
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

func (service service) GetOauthKey(userID string, request CreateOauthRequest) (*OauthResponse, error) {
	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	log.Println("Get OAuth key...")

	headers := service.composeHeaders(false, "")
	host := service.config.Host + endpointOauth + "/" + userID

	log.Printf("Request: %+v\n\nHeaders: %+v\n\nHost: %+v\n\n", request, headers, host)

	responseStatus, responseBytes, err := http.Post(host, headers, requestBytes)

	if err != nil {
		return nil, err
	}

	if responseStatus != 200 {
		err, _ := MapResponseError(responseBytes)
		return nil, err
	}

	response := &OauthResponse{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		log.Printf("Error decoding SynapseFi response: %v", err)
		log.Printf("Response: %q", response)
		return nil, err
	}

	return response, nil
}

func (service service) composeHeaders(isIdempodent bool, oauthKey string) http.Headers {
	headers := http.Headers{}

	// required
	headers["X-SP-GATEWAY"] = service.config.ClientID + "|" + service.config.ClientSecret
	// required
	headers["X-SP-USER"] = oauthKey + "|e83cf6ddcf778e37bfe3d48fc78a6502062fc"
	// required
	headers["X-SP-USER-IP"] = "127.0.0.1"
	// required
	headers["Content-Type"] = "application/json"

	if isIdempodent {
		// optional
		headers["X-SP-IDEMPOTENCY-KEY"] = generateIdempodencyKey()
	}

	return headers
}
