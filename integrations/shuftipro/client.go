package shuftipro

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

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
		e, ok := response.Error.(Error)
		if ok {
			err = e
			return
		}
		err = fmt.Errorf("%scheck the error code in the result", event2description[response.Event])
		return
	}

	res = response.ToKYCResult()

	return
}
