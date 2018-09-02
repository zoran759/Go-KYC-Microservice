package synapsefi

import (
	"github.com/pkg/errors"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"
	"log"
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
			if time.Now().Unix()-startingDate.Unix() >= service.timeoutThreshold {
				return common.Error, nil, errors.New("request timed out")
			}

			time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
			startingPower += 1

			getUserResponse, err := service.verification.GetUser(response.ID)
			if err != nil {
				log.Printf("SynapseFI polling error: %s for user with id: %s", err, response.ID)
				continue
			}

			if getUserResponse.DocumentStatus.PhysicalDoc == "SUBMITTED|REVIEWING" || getUserResponse.DocumentStatus.PhysicalDoc == "SUBMITTED" {
				continue
			}

			return mapResponseToResult(*getUserResponse)
		}
	}

	return mapResponseToResult(*response)
}
