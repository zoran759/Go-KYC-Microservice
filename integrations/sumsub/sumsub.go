package sumsub

import (
	"log"
	"math"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub/applicants"
	"modulus/kyc/integrations/sumsub/documents"
	"modulus/kyc/integrations/sumsub/verification"

	"github.com/pkg/errors"
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
func (service SumSub) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	applicantResponse, err := service.applicants.CreateApplicant(
		customer.Email,
		applicants.MapCommonCustomerToApplicant(*customer),
	)
	if err != nil {
		return
	}

	mappedDocuments := documents.MapCommonCustomerDocuments(*customer)
	if mappedDocuments != nil {
		for _, document := range mappedDocuments {
			_, err1 := service.documents.UploadDocument(applicantResponse.ID, document)

			if err1 != nil {
				err = errors.Wrapf(
					err1,
					"Unable to upload document with filename: %s, type: %s, side: %s",
					document.File.Filename,
					document.Metadata.DocumentType,
					document.Metadata.DocumentSubType,
				)
				return
			}
		}
	}

	started, err := service.verification.StartVerification(applicantResponse.ID)
	if err != nil {
		return
	}
	if !started {
		err = errors.New("verification hasn't started for unknown reason. Please try again")
		return
	}
	// Begin polling sumsub for validation results
	startingPower := 3
	startingDate := time.Now()
	for {
		time.Sleep(time.Duration(math.Exp(float64(startingPower))) * time.Second)
		startingPower++

		status, result, err1 := service.verification.CheckApplicantStatus(applicantResponse.ID)
		if err1 != nil {
			log.Printf("Sumsub polling error: %s for applicant with id: %s", err, applicantResponse.ID)
			continue
		}

		if status == CompleteStatus {
			var detailedResult *common.KYCDetails

			if result.ReviewAnswer != GreenScore && result.RejectLabels != nil && len(result.RejectLabels) > 0 {
				detailedResult = &common.KYCDetails{
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
				res.Status = common.Denied
				res.Details = detailedResult
				return
			case YellowScore:
				res.Status = common.Unclear
				res.Details = detailedResult
				return
			case GreenScore:
				res.Status = common.Approved
				return
			case ErrorScore:
				res.Details = detailedResult
				return
			case IgnoredScore:
				res.Details = detailedResult
				return
			}
		}

		if time.Now().Unix()-startingDate.Unix() >= service.timeoutThreshold {
			err = errors.New("request timed out")
			return
		}
	}
}
