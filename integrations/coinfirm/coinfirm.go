package coinfirm

import (
	"errors"
	"fmt"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/http"
	"modulus/kyc/integrations/coinfirm/model"
)

var _ common.KYCPlatform = Coinfirm{}

// Coinfirm represents the Coinfirm API client.
type Coinfirm struct {
	config Config
}

// New constructs a new Coinfirm API client instance.
// It accepts Config object as the config param.
func New(config Config) Coinfirm {
	return Coinfirm{
		config: config,
	}
}

// CheckCustomer implements CustomerChecker interface for the Coinfirm.
func (c Coinfirm) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("customer is absent or no data received")
		return
	}

	details, docfile := prepareCustomerData(customer)

	headers := headers()

	token, code, err := c.newAuthToken(headers)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending auth request: %s", err)
		return
	}

	headers["Authorization"] = "Bearer " + token.Token

	newParticipant := model.NewParticipant{
		Email: customer.Email,
	}

	participant, code, err := c.newParticipant(headers, newParticipant)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during registering customer: %s", err)
		return
	}

	code, err = c.sendParticipantDetails(headers, participant.UUID, details)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending customer details: %s", err)
		return
	}

	if docfile != nil {
		code, err = c.sendDocFile(headers, participant.UUID, docfile)
		if err != nil {
			if code != nil {
				res.ErrorCode = strconv.Itoa(*code)
			}
			err = fmt.Errorf("during sending customer document: %s", err)
			return
		}
	}

	status, code, err := c.getParticipantCurrentStatus(headers, participant.UUID)
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
func (c Coinfirm) CheckStatus(pID string) (res common.KYCResult, err error) {
	headers := headers()

	token, code, err := c.newAuthToken(headers)
	if err != nil {
		if code != nil {
			res.ErrorCode = strconv.Itoa(*code)
		}
		err = fmt.Errorf("during sending auth request: %s", err)
		return
	}

	headers["Authorization"] = "Bearer " + token.Token

	status, code, err := c.getParticipantCurrentStatus(headers, pID)
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

// headers is the helper returning mandatory headers but they're requiring to complement with authorization.
func headers() http.Headers {
	return http.Headers{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
}
