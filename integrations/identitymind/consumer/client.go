package consumer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"time"

	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/http"
)

const (
	contentType = "application/json"

	consumerEndpoint       = "/account/consumer"
	filesUploadEndpoint    = "/account/consumer/%s/files"
	stateRetrievalEndpoint = "/account/consumer/v2/%s"
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
func (c *Client) CheckCustomer(customer *common.UserData) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	result = common.Error

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

	response, err := c.sendRequest(body)
	if err != nil {
		err = fmt.Errorf("during sending request: %s", err)
		return
	}

	// FIXME: I can't understand from the available docs how documents verification is doing and how to implement it.

	if response.State == UnderReview {
		response, err = c.pollApplicationState(response.KYCTxID)
		if err != nil {
			err = fmt.Errorf("during retrieving current KYC state: %s", err)
			return
		}
	}

	result, details, err = response.toResult()

	return
}

// sendRequest sends a vefirication request into the API.
// It returns a response from the API or the error if occured.
func (c *Client) sendRequest(body []byte) (response *ApplicationResponseData, err error) {
	headers := http.Headers{
		"Content-Type":  contentType,
		"Authorization": c.credentials,
	}

	_, resp, err := http.Post(c.host+consumerEndpoint, headers, body)
	if err != nil {
		return
	}

	response = &ApplicationResponseData{}
	err = json.Unmarshal(resp, response)

	return
}

// pollApplicationState polls IDM API for the current state of a consumer KYC so long
// as the returned state is "Under Review" during one hour.
// If the application is not found then an error message is provided in the response.
func (c *Client) pollApplicationState(mtid string) (response *ApplicationResponseData, err error) {
	headers := http.Headers{
		"Authorization": c.credentials,
	}
	deadline := time.Now().Add(time.Hour)
	endpoint := fmt.Sprintf(c.host+stateRetrievalEndpoint, mtid)

	for {
		var (
			status int
			resp   []byte
		)

		status, resp, err = http.Get(endpoint, headers)
		if err != nil {
			return
		}

		if status == stdhttp.StatusOK {
			response = &ApplicationResponseData{}
			if err = json.Unmarshal(resp, response); err != nil {
				return
			}
			if len(response.ErrorMessage) > 0 {
				err = errors.New(response.ErrorMessage)
				return
			}
			if response.State != UnderReview {
				return
			}
		}

		if time.Now().After(deadline) {
			err = errors.New("polling for the KYC status was aborted after an hour")
			return
		}

		time.Sleep(5 * time.Minute)
	}
}

// Ensure implementation conformance to the interface.
var _ common.CustomerChecker = (*Client)(nil)
