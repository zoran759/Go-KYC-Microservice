package stop4

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/stop4/verification"
	"strconv"

	"github.com/pkg/errors"
)

// Stop4 represents the verification service.
type Stop4 struct {
	verification verification.Verification
}

// New constructs a new verification service object.
func New(config Config) Stop4 {
	return Stop4{
		verification: verification.NewService(verification.Config(config)),
	}
}

// CheckCustomer implements CustomerChecker interface for stop4.
func (service Stop4) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	verificationRequest := verification.MapCustomerToVerificationRequest(*customer)

	response, err := service.verification.Verify(verificationRequest)
	if err != nil {
		return
	}

	if response.Status >= 0 {
		result.Status = common.Approved
	} else {
		result.Status = common.Denied
		result.Details = &common.KYCDetails{
			Finality: common.Unknown,
		}
		result.Details.Reasons = append(result.Details.Reasons, response.Details)

		result.ErrorCode = strconv.FormatInt(int64(response.Status), 10)
	}
	return
}
