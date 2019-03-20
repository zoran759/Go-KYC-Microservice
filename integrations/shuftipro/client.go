package shuftipro

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

const statusEndpoint = "status"

// Client represents the client of the Shufti Pro API.
// It shouldn't initialized directly, use New() constructor instead.
type Client struct {
	host        string
	headers     http.Headers
	callbackURL string
}

// NewClient constructs new Client object.
func NewClient(config Config) Client {
	return Client{
		host: config.Host,
		headers: http.Headers{
			"Content-Type":  "application/json",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.ClientID+":"+config.SecretKey)),
		},
		callbackURL: config.CallbackURL,
	}
}

// CheckCustomer implements the KYCPlatform interface for the Client.
func (c Client) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	req, err := c.NewRequest(customer)
	if err != nil {
		return
	}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	code, resp, err := http.Post(c.host, c.headers, body)
	if err != nil {
		return
	}
	if code != stdhttp.StatusOK {
		res.ErrorCode = strconv.Itoa(code)
	}

	response := Response{}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if code != stdhttp.StatusOK {
		if _, ok := response.Error.(map[string]interface{}); !ok {
			err = fmt.Errorf("%scheck the error code in the result", event2description[response.Event])
			return
		}
		err = errorFromResponse(resp)
		return
	}

	res = response.ToKYCResult()

	return
}

// CheckStatus implements the KYCPlatform interface for the Client.
func (c Client) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	req := StatusRequest{
		Reference: referenceID,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}

	code, resp, err := http.Post(c.host+statusEndpoint, c.headers, body)
	if err != nil {
		return
	}
	if code != stdhttp.StatusOK {
		res.ErrorCode = strconv.Itoa(code)
	}

	response := Response{}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return
	}

	if code != stdhttp.StatusOK {
		if _, ok := response.Error.(map[string]interface{}); !ok {
			err = fmt.Errorf("%scheck the error code in the result", event2description[response.Event])
			return
		}
		err = errorFromResponse(resp)
		return
	}

	res = response.ToKYCResult()

	return
}

// errorFromResponse is a helper function that extracts an error from the API response.
func errorFromResponse(response []byte) error {
	efield := errorField{}
	if err := json.Unmarshal(response, &efield); err != nil || (efield == errorField{}) {
		return errors.New("unexpected format of the returned error: please, report to developers")
	}
	return efield.Error
}
