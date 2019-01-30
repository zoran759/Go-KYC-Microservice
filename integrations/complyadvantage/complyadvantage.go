package complyadvantage

import (
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

var _ common.KYCPlatform = ComplyAdvantage{}

// ComplyAdvantage represents the ComplyAdvantage KYC service.
type ComplyAdvantage struct {
	config Config
}

// New returns a new verification service object.
func New(config Config) ComplyAdvantage {
	return ComplyAdvantage{
		config: config,
	}
}

// CheckCustomer implements KYCPlatform interface for the ComplyAdvantage.
func (c ComplyAdvantage) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	r := c.newRequest(customer)
	resp, status, err := c.performSearch(r)
	if err != nil {
		if status != nil {
			result.ErrorCode = fmt.Sprintf("%d", *status)
		}
		return
	}

	result, err = resp.toResult()

	return
}

// performSearch performs a search request to the ComplyAdvantage API.
func (c ComplyAdvantage) performSearch(r Request) (response Response, status *int, err error) {
	body, err := json.Marshal(r)
	if err != nil {
		return
	}

	headers := http.Headers{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": "Token " + c.config.APIkey,
	}

	code, resp, err := http.Post(c.config.Host+"/searches", headers, body)
	if err != nil {
		return
	}

	if code != stdhttp.StatusOK {
		status = &code

		eresp := &ErrorResponse{}
		err = json.Unmarshal(resp, eresp)
		if err != nil {
			return
		}

		err = eresp
		return
	}

	response = Response{}

	err = json.Unmarshal(resp, &response)

	return
}

// CheckStatus implements KYCPlatform interface for the ComplyAdvantage.
func (c ComplyAdvantage) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	err = errors.New("ComplyAdvantage doesn't support a verification status check")
	return
}
