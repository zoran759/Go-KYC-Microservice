package consumer

import (
	"errors"

	"gitlab.com/lambospeed/kyc/common"
)

const (
	consumerEndpoint    = "/account/consumer"
	filesUploadEndpoint = "/account/consumer/%s/files"
)

// Client defines the client for IdentityMind API.
// It shouldn't be instantiated directly.
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

// CheckCustomer implements customer verification using IdentityMind API.
func (c *Client) CheckCustomer(customer *common.UserData) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	// TODO: implement this.
	if customer == nil {
		result = common.Error
		err = errors.New("no customer supplied")
		return
	}

	return
}

// Ensure implementation conformance to the interface.
var _ common.CustomerChecker = (*Client)(nil)
