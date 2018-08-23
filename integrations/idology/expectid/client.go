package expectid

import (
	"errors"

	"gitlab.com/modulusglobal/kyc/common"
	"gitlab.com/modulusglobal/kyc/integrations/idology"
)

// Client defines the client for IDology ExpectID® API.
type client struct {
	idology.Config
}

// NewClient constructs new client object.
func NewClient(config idology.Config) idology.CustomerChecker {
	return &client{
		Config: config,
	}
}

// CheckCustomer implements customer verification using IDology API.
func (c *client) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
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
