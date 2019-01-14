package coinfirm

import (
	"errors"
	"fmt"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/http"
	"modulus/kyc/integrations/coinfirm/model"
)

// Coinfirm represents the Coinfirm API client.
type Coinfirm struct {
	host     string
	email    string
	password string
	company  string
	token    string
	headers  http.Headers
}

// New constructs a new Coinfirm API client instance.
// It accepts Config object as the config param.
func New(c Config) *Coinfirm {
	return &Coinfirm{
		host:     c.Host,
		email:    c.Email,
		password: c.Password,
		company:  c.Company,
		headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// CheckCustomer implements CustomerChecker interface for the Coinfirm.
func (c *Coinfirm) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("customer is absent or no data received")
		return
	}

	details, docfile := prepareCustomerData(customer)

	token, code, err := c.newAuthToken()
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending auth request: %s", err)
		return
	}

	c.headers["Authorization"] = "Bearer " + token.Token

	newParticipant := model.NewParticipant{
		Email: customer.Email,
	}

	participant, code, err := c.newParticipant(newParticipant)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during registering customer: %s", err)
		return
	}

	code, err = c.sendParticipantDetails(participant.UUID, details)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending customer details: %s", err)
		return
	}

	if docfile != nil {
		code, err = c.sendDocFile(participant.UUID, docfile)
		if err != nil {
			if code != nil {
				res.ErrorCode = strconv.Itoa(*code)
			}
			err = fmt.Errorf("during sending customer document: %s", err)
			return
		}
	}

	status, code, err := c.getParticipantCurrentStatus(participant.UUID)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during checking customer status: %s", err)
		return
	}

	res, err = toResult(participant.UUID, status)

	return
}

// CheckStatus implements StatusChecker interface for the Coinfirm.
func (c *Coinfirm) CheckStatus(pID string) (res common.KYCResult, err error) {
	token, code, err := c.newAuthToken()
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending auth request: %s", err)
		return
	}

	c.headers["Authorization"] = "Bearer " + token.Token

	status, code, err := c.getParticipantCurrentStatus(pID)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during checking customer status: %s", err)
		return
	}

	res, err = toResult(pID, status)

	return
}
