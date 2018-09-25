package jumio

import "gitlab.com/lambospeed/kyc/common"

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

// RejectReason represents the reject reason.
type RejectReason struct {
	Code        string              `json:"rejectReasonCode"`
	Description string              `json:"rejectReasonDescription"`
	Details     RejectReasonDetails `json:"rejectReasonDetails"`
}

// IdentityVerification represents the identity verification.
type IdentityVerification struct {
	Similarity             string `json:"similarity"`
	Validity               string `json:"validity"`
	Reason                 string `json:"reason"`
	HandwrittenNoteMatches string `json:"handwrittenNoteMatches"`
}

// VerificationDetails represents the part of DetailsResponse.
type VerificationDetails struct {
	MrzCheck             string                `json:"mrzCheck"`
	FaceMatch            string                `json:"faceMatch"`
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
func (r *DetailsResponse) toResult() (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	// TODO: implement this.

	return
}
