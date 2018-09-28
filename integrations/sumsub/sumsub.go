package sumsub

import (
	"log"
	"math"
	"time"

	"github.com/pkg/errors"
	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub/applicants"
	"modulus/kyc/integrations/sumsub/documents"
	"modulus/kyc/integrations/sumsub/verification"
)

// SumSub defines the verification service.
type SumSub struct {
	applicants       applicants.Applicants
	documents        documents.Documents
	verification     verification.Verification
	timeoutThreshold int64
}

// New constructs new verification service object.
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

// CheckCustomer implements CustomerChecker interface for Sum&Substance KYC provider.
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
		return common.Error, nil, errors.New("verification hasn't started for unknown reason. Please try again")
	}
	// Begin polling sumsub for validation results
	startingPower := 3
	startingDate := time.Now()
	for {
		time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
		startingPower++

		status, result, err := service.verification.CheckApplicantStatus(applicantResponse.ID)
		if err != nil {
			log.Printf("Sumsub polling error: %s for applicant with id: %s", err, applicantResponse.ID)
			continue
		}

		if status == CompleteStatus {
			var detailedResult *common.DetailedKYCResult

			if result.ReviewAnswer != GreenScore && result.RejectLabels != nil && len(result.RejectLabels) > 0 {
				detailedResult = &common.DetailedKYCResult{
					Reasons: result.RejectLabels,
				}

				switch result.ReviewRejectType {
				case FinalRejectType:
					detailedResult.Finality = common.Final
				case RetryRejectType:
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
