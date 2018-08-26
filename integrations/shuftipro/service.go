package shuftipro

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/shuftipro/verification"
	"github.com/pkg/errors"
	"log"
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
	case Success:
		log.Println(service.verification.CheckStatus(response.Reference))
		fallthrough
	default:
		return common.Error, nil, errors.New(response.Message)
	}
}
