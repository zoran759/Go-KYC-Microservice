package verification

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/http"

	"github.com/google/uuid"
)

const (
	endpointUsers = "users"
	endpointOAuth = "oauth"
	appLanguage   = "en"
)

type service struct {
	config Config
}

// NewService constructs the new verification service object.
func NewService(config Config) Verification {
	config.fingerprint = fmt.Sprintf("%x", sha256.Sum256([]byte(config.ClientID+config.ClientSecret)))

	return service{
		config: config,
	}
}

// TODO: resend on fail, process errors, PROCESS RESPONSE PERMISSIONS!!!
func (service service) CreateUser(user User) (resp *Response, code *string, err error) {
	body, err := json.Marshal(user)
	if err != nil {
		return
	}

	headers := service.composeHeaders(true, "")
	endpoint := service.config.Host + endpointUsers

	status, response, err := http.Post(endpoint, headers, body)
	if err != nil {
		return
	}
	if status != stdhttp.StatusOK {
		code, err = MapErrorResponse(response)
		return
	}

	resp = &Response{}

	if err = json.Unmarshal(response, resp); err != nil {
		resp = nil
	}

	return
}

func (service service) AddPhysicalDocs(userID, rtoken string, docs PhysicalDocs) (code *string, err error) {
	key, err := service.getOAuthKey(userID, rtoken)
	if err != nil {
		return
	}

	headers := service.composeHeaders(false, key)
	endpoint := service.config.Host + endpointUsers + "/" + userID

	for _, doc := range docs.Documents {
		body, err1 := json.Marshal(PhysicalDocs{
			Documents: []Document{doc},
		})
		if err1 != nil {
			return nil, err1
		}

		status, response, err1 := http.Patch(endpoint, headers, body)
		if err1 != nil {
			return nil, err1
		}
		if status != stdhttp.StatusOK {
			code, err = MapErrorResponse(response)
			return
		}
	}

	return
}

func (service service) GetUser(userID string) (resp *Response, code *string, err error) {
	headers := service.composeHeaders(true, "")
	endpoint := service.config.Host + endpointUsers + "/" + userID

	status, response, err := http.Get(endpoint, headers)
	if err != nil {
		return
	}
	if status != stdhttp.StatusOK {
		code, err = MapErrorResponse(response)
		return
	}

	resp = &Response{}

	if err = json.Unmarshal(response, resp); err != nil {
		resp = nil
	}

	return
}

func (service service) getOAuthKey(userID, rtoken string) (key string, err error) {
	req := OAuthRequest{
		RefreshToken: rtoken,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	headers := service.composeHeaders(false, "")
	endpoint := service.config.Host + endpointOAuth + "/" + userID

	status, resp, err := http.Post(endpoint, headers, body)
	if err != nil {
		return
	}
	if status != stdhttp.StatusOK {
		_, err = MapErrorResponse(resp)
		return
	}

	response := &OAuthResponse{}
	if err = json.Unmarshal(resp, response); err != nil {
		return
	}

	key = response.OAuthKey

	return
}

func (service service) composeHeaders(useIdempotency bool, oauthKey string) http.Headers {
	headers := http.Headers{
		// required
		"X-SP-GATEWAY": service.config.ClientID + "|" + service.config.ClientSecret,
		// required
		"X-SP-USER": oauthKey + "|" + service.config.fingerprint,
		// required
		"X-SP-USER-IP": "127.0.0.1",
		// required
		"Content-Type": "application/json",
	}

	if useIdempotency {
		// optional
		headers["X-SP-IDEMPOTENCY-KEY"] = newIdempotencyKey()
	}

	return headers
}

// newIdempotencyKey returns new idempotency key.
func newIdempotencyKey() string {
	return uuid.New().String()
}
