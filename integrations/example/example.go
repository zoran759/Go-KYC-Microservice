package example

import (
	"errors"
	"modulus/kyc/common"
)

// Service represents the example KYC provider.
type Service struct{}

// CheckCustomer implements CustomerChecker interface for the example KYC provider.
func (s *Service) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}
	if len(customer.FirstName) == 0 {
		err = errors.New("missing the required field: FirstName")
		return
	}

	// KYCResult simulation.
	switch customer.FirstName {
	case "Abby":
		res.Status = common.Approved
	case "Delilah":
		res = deniedResult()
	case "Destiny":
		res = deniedResultWithFinality()
	case "Urbi":
		res = unclearResult()
	case "Erika":
		res.ErrorCode = "429"
		err = errors.New("during sending request: http error")
	default:
		res = errorResult()
	}

	return
}

// CheckStatus implements StatusChecker interface for the example KYC provider.
func (s *Service) CheckStatus(customerID string) (res common.KYCResult, err error) {

	return
}

func errorResult() (res common.KYCResult) {
	res.Details = &common.KYCDetails{
		Reasons: []string{
			"Not readable document",
			"This is the example error reason",
		},
	}
	res.ErrorCode = "42"

	return
}

func deniedResult() (res common.KYCResult) {
	res.Status = common.Denied
	res.Details = &common.KYCDetails{
		Reasons: []string{
			"This is the example reason of denial",
			"Delilah wants to trick Samson",
		},
	}

	return
}

func deniedResultWithFinality() (res common.KYCResult) {
	res.Status = common.Denied
	res.Details = &common.KYCDetails{
		Finality: common.Final,
		Reasons: []string{
			"Date of birth does not match",
			"Selfie is manipulated",
			"Destiny wish you good luck",
		},
	}

	return
}

func unclearResult() (res common.KYCResult) {
	res.Status = common.Unclear
	res.StatusPolling = &common.StatusPolling{
		Provider:   common.Example,
		CustomerID: "lily_was_here",
	}

	return
}
