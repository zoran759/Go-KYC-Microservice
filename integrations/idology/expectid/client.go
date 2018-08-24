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
		err = errors.New("no customer supplied")
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

	result, details, err = c.responseToResult(response)

	return
}

// responseToResult processes the response and generates the verification result.
func (c *Client) responseToResult(r *Response) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	detailsCreateIfNil := func(details *common.DetailedKYCResult) {
		if details == nil {
			details = &common.DetailedKYCResult{
				Finality: common.Unknown,
			}
		}
	}

	switch c.config.UseSummaryResult {
	case true:
		switch r.SummaryResult.Key {
		case Success:
			result = common.Approved
		case Failure:
			result = common.Denied
		case Partial:
			result = common.Unclear
		}
	case false:
		switch r.Results.Key {
		case Match:
			result = common.Approved
		case NoMatch, MatchRestricted:
			result = common.Denied
		}
	}

	if r.Restriction != nil {
		detailsCreateIfNil(details)
		details.Reasons = []string{
			r.Restriction.Message,
			r.Restriction.PatriotAct.List,
			fmt.Sprintf("Patriot Act score: %d", r.Restriction.PatriotAct.Score),
		}
	}

	if len(r.Qualifiers) > 0 {
		detailsCreateIfNil(details)
		for _, q := range r.Qualifiers {
			details.Reasons = append(details.Reasons, q.Message)
		}
	}

	return
}
