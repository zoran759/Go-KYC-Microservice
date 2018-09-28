package shuftipro

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/shuftipro/verification"

	"github.com/pkg/errors"
)

// ShuftiPro represents the verification service.
type ShuftiPro struct {
	verification verification.Verification
}

// New constructs a new verification service object.
func New(config Config) ShuftiPro {
	return ShuftiPro{
		verification: verification.NewService(verification.Config(config)),
	}
}

// CheckCustomer implements CustomerChecker interface for ShuftiPro.
func (service ShuftiPro) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	verificationRequest := verification.MapCustomerToVerificationRequest(*customer)

	response, err := service.verification.Verify(verificationRequest)
	if err != nil {
		return
	}

	switch response.StatusCode {
	case Verified:
		result.Status = common.Approved
		return
	case NotVerified:
		result.Status = common.Denied
		return
	// This status means that online verification is being performed instead of offline verification(which we need).
	// It happens when documents are not provided or they are invalid.
	case Success:
		err = errors.New("There are no documents provided or they are invalid")
		return
	default:
		err = errors.New(response.Message)
		return
	}
}
