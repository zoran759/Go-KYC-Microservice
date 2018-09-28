package shuftipro

import (
	"github.com/pkg/errors"
	"modulus/kyc/common"
	"modulus/kyc/integrations/shuftipro/verification"
)

type ShuftiPro struct {
	verification verification.Verification
}

func New(config Config) ShuftiPro {
	return ShuftiPro{
		verification: verification.NewService(verification.Config(config)),
	}
}

func (service ShuftiPro) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
	if customer == nil {
		return common.Error, nil, errors.New("No customer supplied")
	}

	verificationRequest := verification.MapCustomerToVerificationRequest(*customer)

	response, err := service.verification.Verify(verificationRequest)
	if err != nil {
		return common.Error, nil, err
	}

	switch response.StatusCode {
	case Verified:
		return common.Approved, nil, nil
	case NotVerified:
		return common.Denied, nil, nil
	// This status means that online verification is being performed instead of offline verification(which we need).
	// It happens when documents are not provided or they are invalid.
	case Success:
		return common.Error, nil, errors.New("There are no documents provided or they are invalid")
	default:
		return common.Error, nil, errors.New(response.Message)
	}
}
