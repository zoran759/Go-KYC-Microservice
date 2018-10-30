package synapsefi

import (
	"github.com/pkg/errors"
	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"
	"log"
	"time"
)

type SynapseFI struct {
	verification     verification.Verification
	timeoutThreshold int64
	kycFlow 		 string
}

func New(config Config) SynapseFI {
	kycFlow := "simple"
	if len(config.KYCFlow) > 0 {
		kycFlow = config.KYCFlow
	}

	timeoutThreshold := int64(time.Hour.Seconds())
	if config.TimeoutThreshold > 0 {
		timeoutThreshold = config.TimeoutThreshold
	}

	return SynapseFI{
		verification:     verification.NewService(verification.Config(config.Connection)),
		timeoutThreshold: timeoutThreshold,
		kycFlow: kycFlow,
	}


}

func (service SynapseFI) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	createUserRequest := verification.MapCustomerToCreateUserRequest(*customer, true)

	response, err := service.verification.CreateUser(createUserRequest)
	if err != nil {
		return result, err
	}

	if service.kycFlow != "" && service.kycFlow != "simple" {
		log.Printf("Alternative flow, auth user");

		uID := response.ID

		createOauthRequest := verification.MapUserToOauth(response.RefreshToken)
		responseAuth, err := service.verification.GetOauthKey(uID, createOauthRequest)
		if err != nil {
			return result, err
		}
		log.Printf("OAuth response: %+v", responseAuth);

		createDocumentRequest := verification.MapDocumentsToCreateUserRequest(*customer)
		response, err = service.verification.AddDocument(uID, responseAuth.OAuthKey, createDocumentRequest)
		if err != nil {
			return result, err
		}
	}

	log.Printf("Initial status: %+v", response)

	switch response.DocumentStatus.PhysicalDoc {
	case DocStatusInvalid:
		fallthrough
	case DocStatusValid:
		result, err = mapResponseToResult(response)
		return result, err

	case DocStatusMissingOrInvalid:
		fallthrough
	case DocStatusPending:
		fallthrough
	case DocStatusReviewing:
		startingPower := 60
		startingDate := time.Now()
		for {
			if time.Now().Unix()-startingDate.Unix() >= service.timeoutThreshold {
				log.Printf("Timeout exceeded")
				continue
			}

			time.Sleep(time.Duration(startingPower) * time.Second)

			getUserResponse, err := service.verification.GetUser(response.ID)
			if err != nil {
				log.Printf("SynapseFI polling error: %s for user with id: %s", err, response.ID)
				continue
			}
			log.Printf("Response: %+v", getUserResponse)

			if getUserResponse.DocumentStatus.PhysicalDoc == DocStatusPending || getUserResponse.DocumentStatus.PhysicalDoc == DocStatusReviewing {
				continue
			}

			result, err = mapResponseToResult(getUserResponse)
			return result, err
		}

	default:
		result.Status = common.Denied
		result.Details = &common.KYCDetails{
			Finality: common.Unknown,
		}
		result.Details.Reasons = append(result.Details.Reasons, DocStatusMissingOrInvalid)

		return result, err
	}
}