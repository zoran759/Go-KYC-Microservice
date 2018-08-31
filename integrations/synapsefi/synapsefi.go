package synapsefi

import (
	"github.com/pkg/errors"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"
	"math"
	"time"
)

type SynapseFI struct {
	verification     verification.Verification
	timeoutThreshold int64
}

func New(config Config) SynapseFI {
	return SynapseFI{
		verification:     verification.NewService(config.Config),
		timeoutThreshold: config.TimeoutThreshold,
	}
}

func (service SynapseFI) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
	if customer == nil {
		return common.Error, nil, errors.New("No customer supplied")
	}
	createUserRequest := verification.MapCustomerToCreateUserRequest(*customer)

	response, err := service.verification.CreateUser(createUserRequest)
	if err != nil {
		return common.Error, nil, err
	}

	if response.DocumentStatus.PhysicalDoc == "SUBMITTED|REVIEWING" || response.DocumentStatus.PhysicalDoc == "SUBMITTED" {
		// Begin polling SynapseFI for validation results
		startingPower := 3
		startingDate := time.Now()
		for {
			time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
			startingPower += 1

			if response.DocumentStatus.PhysicalDoc == "SUBMITTED|REVIEWING" || response.DocumentStatus.PhysicalDoc == "SUBMITTED" {
				continue
			}

			return mapResponseToResult(*response)

			if time.Now().Unix()-startingDate.Unix() >= service.timeoutThreshold {
				return common.Error, nil, errors.New("request timed out")
			}
		}
	}

	return mapResponseToResult(*response)
}
