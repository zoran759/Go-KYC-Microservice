package sumsub

import (
	"github.com/pkg/errors"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/applicants"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/documents"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/verification"
	"log"
	"math"
	"time"
)

type SumSub struct {
	applicants       applicants.Applicants
	documents        documents.Documents
	verification     verification.Verification
	timeoutThreshold int64
}

func New(config Config) SumSub {
	return SumSub{
		applicants: applicants.NewService(applicants.Config{
			Host:   config.Host,
			APIKey: config.APIKey,
		}),
		documents: documents.NewService(documents.Config{
			Host:   config.Host,
			APIKey: config.APIKey,
		}),
		verification: verification.NewService(verification.Config{
			Host:   config.Host,
			APIKey: config.APIKey,
		}),
		timeoutThreshold: int64(time.Second.Seconds()) * config.TimeoutThreshold,
	}
}

func (service SumSub) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
	if customer == nil {
		return common.Error, nil, errors.New("No customer supplied")
	}

	applicantResponse, err := service.applicants.CreateApplicant(
		customer.Email,
		applicants.MapCommonCustomerToApplicant(*customer),
	)
	if err != nil {
		return common.Error, nil, err
	}

	mappedDocuments := documents.MapCommonCustomerDocuments(*customer)
	if mappedDocuments != nil {
		for _, document := range mappedDocuments {
			_, err := service.documents.UploadDocument(applicantResponse.ID, document)

			if err != nil {
				return common.Error, nil, errors.Wrapf(
					err,
					"Unable to upload document with filename: %s, type: %s, side: %s",
					document.File.Filename,
					document.Metadata.DocumentType,
					document.Metadata.DocumentSubType,
				)
			}
		}
	}

	started, err := service.verification.StartVerification(applicantResponse.ID)
	if err != nil {
		return common.Error, nil, err
	}
	if !started {
		return common.Error, nil, errors.New("Verification hasn't started for unknown reason. Please try again.")
	}
	// Begin polling sumsub for validation results
	startingPower := 3
	startingDate := time.Now()
	for {
		time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
		startingPower += 1

		status, result, err := service.verification.CheckApplicantStatus(applicantResponse.ID)
		if err != nil {
			log.Printf("Sumsub polling error: %s for applicant with id: %s", err, applicantResponse.ID)
			continue
		}

		if status == CompleteStatus {
			var detailedResult *common.DetailedKYCResult = nil

			if result.ReviewAnswer != GreenScore && result.RejectLabels != nil && len(result.RejectLabels) > 0 {
				detailedResult = &common.DetailedKYCResult{
					Reasons: result.RejectLabels,
				}

				switch result.ReviewRejectType {
				case FinalRejectType:
					detailedResult.Finality = common.Final
				case RetryRejectTYpe:
					detailedResult.Finality = common.NonFinal
				default:
					detailedResult.Finality = common.Unknown
				}
			}

			switch result.ReviewAnswer {
			case RedScore:
				return common.Denied, detailedResult, nil
			case YellowScore:
				return common.Unclear, detailedResult, nil
			case GreenScore:
				return common.Approved, nil, nil
			case ErrorScore:
				return common.Error, detailedResult, nil
			case IgnoredScore:
				return common.Error, detailedResult, nil
			}
		}

		if time.Now().Unix()-startingDate.Unix() >= service.timeoutThreshold {
			return common.Error, nil, errors.New("request timed out")
		}
	}
}
