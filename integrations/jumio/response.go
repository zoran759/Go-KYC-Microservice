package jumio

import (
	"errors"
	"strings"

	"modulus/kyc/common"
)

// Response defines the model for the performNetverify API response.
type Response struct {
	// Timestamp (UTC) of the response format: YYYY-MM-DDThh:mm:ss.SSSZ.
	Timestamp string `json:"timestamp"`
	// Jumio's reference number for each scan. Max. length 36.
	JumioIDScanReference string `json:"jumioIdScanReference"`
}

// StatusResponse defines the model for the response on retrieving scan status.
type StatusResponse struct {
	// Timestamp of the response in the format YYYY-MM-DDThh:mm:ss.SSSZ.
	Timestamp string `json:"timestamp"`
	// Jumio’s reference number for each scan. Max. lenght 36.
	ScanReference string     `json:"scanReference"`
	Status        ScanStatus `json:"status"`
}

// DocumentDetails represents the part of DetailsResponse.
type DocumentDetails struct {
	Status         DocumentStatus `json:"status"`
	Type           IDType         `json:"type"`
	IDSubtype      string         `json:"idSubtype"`
	IssuingCountry string         `json:"issuingCountry"`
	FirstName      string         `json:"firstName"`
	LastName       string         `json:"lastName"`
	DOB            string         `json:"dob"`
	Expiry         string         `json:"expiry"`
	IssuingDate    string         `json:"issuingDate"`
	Number         string         `json:"number"`
	USState        string         `json:"usState"`
	PersonalNumber string         `json:"personalNumber"`
}

// TransactionDetails represents the part of DetailsResponse.
type TransactionDetails struct {
	Status                    ScanStatus `json:"status"`
	Source                    string     `json:"source"`
	Date                      string     `json:"date"`
	ClientIP                  string     `json:"clientIp"`
	CustomerID                string     `json:"customerId"`
	MerchantScanReference     string     `json:"merchantScanReference"`
	MerchantReportingCriteria string     `json:"merchantReportingCriteria"`
}

// RejectReasonDetails represents the reject reason details.
type RejectReasonDetails struct {
	Code        string `json:"detailsCode"`
	Description string `json:"detailsDescription"`
}

// String returns the string representation of the RejectReasonDetails.
func (r RejectReasonDetails) String() string {
	return r.Code + " " + r.Description
}

// RejectReason represents the reject reason.
type RejectReason struct {
	Code        string                `json:"rejectReasonCode"`
	Description string                `json:"rejectReasonDescription"`
	Details     []RejectReasonDetails `json:"rejectReasonDetails"`
}

// IdentityVerification represents the identity verification.
type IdentityVerification struct {
	Similarity             string `json:"similarity"`
	Validity               string `json:"validity"`
	Reason                 string `json:"reason"`
	HandwrittenNoteMatches string `json:"handwrittenNoteMatches"`
}

// String returns the string representation of the IdentityVerification.
func (i IdentityVerification) String() string {
	b := strings.Builder{}
	b.WriteString("Identity Verification: similarity = ")
	b.WriteString(i.Similarity)
	b.WriteString(" | validity = ")
	b.WriteString(i.Validity)
	if len(i.Reason) > 0 {
		b.WriteString(" | reason = ")
		b.WriteString(i.Reason)
	}

	return b.String()
}

// VerificationDetails represents the part of DetailsResponse.
type VerificationDetails struct {
	MrzCheck             string                `json:"mrzCheck"`
	RejectReason         *RejectReason         `json:"rejectReason"`
	IdentityVerification *IdentityVerification `json:"identityVerification"`
}

// DetailsResponse defines the model for the response on retrieving scan details.
type DetailsResponse struct {
	Timestamp     string               `json:"timestamp"`
	ScanReference string               `json:"scanReference"`
	Document      *DocumentDetails     `json:"document"`
	Transaction   *TransactionDetails  `json:"transaction"`
	Verification  *VerificationDetails `json:"verification"`
}

// toResult processes the response and generates the verification result.
func (r *DetailsResponse) toResult() (result common.KYCResult, err error) {
	switch r.Document.Status {
	case ApprovedVerified:
		result.Status = common.Approved
		return
	case DeniedFraud, DeniedUnsupportedIDType, DeniedUnsupportedIDCountry:
		result.Status = common.Denied
	case ErrorNotReadableID, NoIDUploaded:
		result.Details = &common.KYCDetails{
			Reasons: []string{"Document status: " + string(r.Document.Status)},
		}
	}

	if r.Verification != nil {
		reasons := []string{}
		if r.Verification.RejectReason != nil {
			result.ErrorCode = r.Verification.RejectReason.Code
			reasons = append(reasons, r.Verification.RejectReason.Description)
			for _, details := range r.Verification.RejectReason.Details {
				reasons = append(reasons, details.String())
			}
		}
		if r.Verification.IdentityVerification != nil {
			reasons = append(reasons, r.Verification.IdentityVerification.String())
		}
		if len(reasons) != 0 {
			if result.Details == nil {
				result.Details = &common.KYCDetails{}
			}
			result.Details.Reasons = append(result.Details.Reasons, reasons...)
		}
	}

	if r.Transaction.Status == FailedStatus {
		err = errors.New("for some reason Jumio returned the 'FAILED' status for the verification transaction")
	}

	return
}
