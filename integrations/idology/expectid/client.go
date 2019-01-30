package expectid

import (
	"errors"
	"fmt"

	"modulus/kyc/common"
)

// Client defines the client for IDology ExpectID® API.
// It shouldn't be instantiated directly.
// Use NewClient() constructor instead.
type Client struct {
	config Config
}

// NewClient constructs new client object.
func NewClient(config Config) Client {
	return Client{
		config: config,
	}
}

// CheckCustomer implements customer verification using IDology API.
func (c Client) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	requestBody := c.makeRequestBody(customer)

	response, err := c.sendRequest(requestBody)
	if err != nil {
		return
	}

	if response.Error != nil {
		err = fmt.Errorf("during verification: %s", *response.Error)
		return
	}

	result, err = response.toResult(c.config.UseSummaryResult)

	return
}
