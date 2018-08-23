package expectid

import (
	"errors"

	"gitlab.com/lambospeed/kyc/common"
)

// Client defines the client for IDology ExpectID® API.
// It shouln't be instantiated directly.
// Use NewClient() constructor instead.
type Client struct {
	config Config
}

// NewClient constructs new client object.
func NewClient(config Config) *Client {
	return &Client{
		config: config,
	}
}

// CheckCustomer implements customer verification using IDology API.
func (c *Client) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
	//  TODO: implement this.
	if customer == nil {
		return common.Error, nil, errors.New("No customer supplied")
	}

	requestBody := c.makeRequestBody(customer)

	response, err := c.verify(requestBody)

	//  TODO: process response.
	_ = response

	if err != nil {
		return common.Error, nil, err
	}

	return common.Approved, nil, nil
}
