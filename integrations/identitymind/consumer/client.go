package consumer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

const (
	contentType = "application/json"

	consumerEndpoint       = "/account/consumer"
	stateRetrievalEndpoint = "/account/consumer/v2/"
)

// Client defines the client for IdentityMind API.
// It shouldn't be instantiated directly.
// Use NewClient() constructor instead.
type Client struct {
	host        string
	credentials string
}

// NewClient constructs new client object.
func NewClient(config Config) *Client {
	return &Client{
		host:        config.Host,
		credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Username+":"+config.Password)),
	}
}

// CheckCustomer implements customer verification using IdentityMind API.
func (c *Client) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	req := &KYCRequestData{}
	if err = req.populateFields(customer); err != nil {
		err = fmt.Errorf("invalid customer data: %s", err)
		return
	}

	body, err := req.createRequestBody()
	if err != nil {
		err = fmt.Errorf("during creating request body: %s", err)
		return
	}

	response, errorCode, err := c.sendRequest(body)
	if err != nil {
		if errorCode != nil {
			result.ErrorCode = fmt.Sprintf("%d", *errorCode)
		}
		err = fmt.Errorf("during sending request: %s", err)
		return
	}

	result, err = response.toResult()

	return
}

// sendRequest sends a vefirication request into the API.
// It returns a response from the API or the error if occured.
func (c *Client) sendRequest(body []byte) (response *ApplicationResponseData, errorCode *int, err error) {
	headers := http.Headers{
		"Content-Type":  contentType,
		"Authorization": c.credentials,
	}

	status, resp, err := http.Post(c.host+consumerEndpoint, headers, body)
	if err != nil {
		return
	}
	if status != stdhttp.StatusOK {
		errorCode = &status
		err = errors.New("http error")
		return
	}

	response = &ApplicationResponseData{}
	err = json.Unmarshal(resp, response)

	return
}

// CheckStatus queries IDM API for the current state of a consumer KYC.
// If the application is not found then an error message is provided in the response.
func (c *Client) CheckStatus(referenceID string) (result common.KYCResult, err error) {
	headers := http.Headers{
		"Authorization": c.credentials,
	}

	status, resp, err := http.Get(c.host+stateRetrievalEndpoint+referenceID, headers)
	if err != nil {
		err = fmt.Errorf("during sending request: %s", err)
		return
	}
	if status != stdhttp.StatusOK {
		result.ErrorCode = fmt.Sprintf("%d", status)
		err = errors.New("http error")
		return
	}

	response := &ApplicationResponseData{}
	if err = json.Unmarshal(resp, response); err != nil {
		return
	}
	if len(response.ErrorMessage) > 0 {
		err = errors.New(response.ErrorMessage)
		return
	}

	result, err = response.toResult()

	return
}
