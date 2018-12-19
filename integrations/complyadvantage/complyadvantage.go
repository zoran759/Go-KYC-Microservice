package complyadvantage

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

// service represents the service.
type service struct {
	host string
	key  string
}

// New returns a new verification service object.
func New(c Config) common.CustomerChecker {
	return service{
		host: c.Host,
		key:  c.APIkey,
	}
}

// CheckCustomer implements CustomerChecker interface for the ComplyAdvantage.
func (s service) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	r := newRequest(customer)
	resp, status, err := s.performSearch(r)
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
func (s service) performSearch(r Request) (response *Response, status *int, err error) {
	body, err := json.Marshal(r)
	if err != nil {
		return
	}

	headers := http.Headers{
		"Content-Type":  "application/json; charset=utf-8",
		"Authorization": "Token " + s.key,
	}

	code, resp, err := http.Post(s.host+"/searches", headers, body)
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

	response = &Response{}

	err = json.Unmarshal(resp, response)

	return
}
