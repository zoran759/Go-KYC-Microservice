package sumsub

import (
	"fmt"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub/applicants"
	"modulus/kyc/integrations/sumsub/documents"
	"modulus/kyc/integrations/sumsub/verification"

	"github.com/pkg/errors"
)

// SumSub defines the verification service.
type SumSub struct {
	applicants   applicants.Applicants
	documents    documents.Documents
	verification verification.Verification
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
	}
}

// CheckCustomer implements CustomerChecker interface for Sum&Substance KYC provider.
func (service SumSub) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	// Process customer documents. At least one document has to be provided for verification.
	mappedDocuments := documents.MapCommonCustomerDocuments(*customer)
	if len(mappedDocuments) == 0 {
		err = errors.New("at least one document has to be provided for verification")
		return
	}

	// Create an applicant.
	applicantResponse, err := service.applicants.CreateApplicant(
		customer.Email,
		applicants.MapCommonCustomerToApplicant(*customer),
	)
	if err != nil {
		if applicantResponse != nil && applicantResponse.Code != nil {
			res.ErrorCode = fmt.Sprintf("%d", *applicantResponse.Code)
		}
		return
	}

	if len(applicantResponse.ID) == 0 {
		err = errors.New("missing applicant id in the API response")
		return
	}

	// Upload applicant's documents.
	for _, document := range mappedDocuments {
		_, errorCode, err1 := service.documents.UploadDocument(applicantResponse.ID, document)
		if err1 != nil {
			if errorCode != nil {
				res.ErrorCode = fmt.Sprintf("%d", *errorCode)
			}
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

	// Request applicant check.
	if err = service.verification.RequestApplicantCheck(applicantResponse.ID); err != nil {
		err = fmt.Errorf("during requesting applicant check: %s", err)
		return
	}

	// Save status polling data.
	res.Status = common.Unclear
	res.StatusCheck = &common.KYCStatusCheck{
		Provider:    common.SumSub,
		ReferenceID: applicantResponse.ID,
		LastCheck:   time.Now(),
	}

	return
}

// CheckStatus implements StatusChecker interface for Sum&Substance KYC provider.
func (service SumSub) CheckStatus(refID string) (res common.KYCResult, err error) {
	status, result, err := service.verification.CheckApplicantStatus(refID)
	if err != nil {
		if result != nil && result.ErrorCode != 0 {
			res.ErrorCode = fmt.Sprintf("%d", result.ErrorCode)
		}
		return
	}

	switch status {
	case "completed", "completedSent", "completedSentFailure":
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
		case YellowScore:
			res.Status = common.Unclear
			res.Details = detailedResult
		case GreenScore:
			res.Status = common.Approved
		case ErrorScore:
			res.Details = detailedResult
		case IgnoredScore:
			res.Details = detailedResult
		}
	case "init":
		err = errors.New("documents upload failed. Please, try to upload a document for this applicant")
	case "pending", "queued":
		res.Status = common.Unclear
		res.StatusCheck = &common.KYCStatusCheck{
			Provider:    common.SumSub,
			ReferenceID: refID,
			LastCheck:   time.Now(),
		}
	case "awaitingUser":
		err = errors.New("waiting some additional documents from the applicant (e.g. a selfie or a better passport image)")
	}

	return
}
