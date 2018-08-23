package expectid

import (
	"errors"
	"fmt"

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
func (c *Client) CheckCustomer(customer *common.UserData) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	if customer == nil {
		result = common.Error
		err = errors.New("No customer supplied")
		return
	}

	requestBody := c.makeRequestBody(customer)

	response, err := c.verify(requestBody)
	if err != nil {
		result = common.Error
		return
	}

	if response.Error != nil {
		result = common.Error
		err = fmt.Errorf("during verification: %q", *response.Error)
		return
	}

	// FIXME: The <summary-result> and <results> tags are not the same.
	// At this moment I don't count on <summary-result>.
	// I need to clarify that.

	switch response.Results.Key {
	case Match, MatchRestricted:
		result = common.Approved
		details = &common.DetailedKYCResult{
			Finality: common.Unknown,
			Reasons:  []string{response.Results.Message},
		}
	case NoMatch:
		result = common.Denied
	}

	return
}
