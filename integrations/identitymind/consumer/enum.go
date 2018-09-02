package consumer

// DeviceFingerprintType defines the device fingerprint type.
type DeviceFingerprintType string

// Possible values of DeviceFingerprintType.
const (
	Augur        DeviceFingerprintType = "AU"
	Iovation     DeviceFingerprintType = "IO"
	ThreatMetrix DeviceFingerprintType = "CB"
	InAuth       DeviceFingerprintType = "IA"
	BlueCava     DeviceFingerprintType = "BC"
)

// CardType defines the type of the card.
type CardType string

// Possible values of CardType.
const (
	Credit  CardType = "CREDIT"
	Debit   CardType = "DEBIT"
	Prepaid CardType = "PREPAID"
	Unknown CardType = "UNKNOWN"
)

// DocumentType defines the Type of the document for usage in the document verification.
type DocumentType string

// Possible values of DocumentType.
const (
	Passport               DocumentType = "PP"
	DriverLicence          DocumentType = "DL"
	GovernmentIssuedIDCard DocumentType = "ID"
	ResidencePermit        DocumentType = "RP"
	UtilityBill            DocumentType = "UB"
)

// EDNAPolicyResult defines the result of the user reputation evaluation.
type EDNAPolicyResult string

// Possible values of EDNAPolicyResult.
const (
	EPTrusted       EDNAPolicyResult = "TRUSTED"
	EPWeaklyTrusted EDNAPolicyResult = "WEAKLY_TRUSTED"
	EPUnknown       EDNAPolicyResult = "UNKNOWN"
	EPSuspicious    EDNAPolicyResult = "SUSPICIOUS"
	EPBad           EDNAPolicyResult = "BAD"
)

// FraudPolicyResult defines the result of the fraud evaluation.
type FraudPolicyResult string

// Possible values of FraudPolicyResult.
const (
	Accept       FraudPolicyResult = "ACCEPT"
	ManualReview FraudPolicyResult = "MANUAL_REVIEW"
	Deny         FraudPolicyResult = "DENY"
)

// AutomatedReviewPolicyResult defines the result of the automated review evaluation.
type AutomatedReviewPolicyResult string

// Possible values of AutomatedReviewPolicyResult.
const (
	ARPError         AutomatedReviewPolicyResult = "ERROR"
	ARPNoPolicy      AutomatedReviewPolicyResult = "NO_POLICY"
	ARPDisabled      AutomatedReviewPolicyResult = "DISABLED"
	ARPFiltered      AutomatedReviewPolicyResult = "FILTERED"
	ARPPending       AutomatedReviewPolicyResult = "PENDING"
	ARPFail          AutomatedReviewPolicyResult = "FAIL"
	ARPIndeterminate AutomatedReviewPolicyResult = "INDETERMINATE"
	ARPSuccess       AutomatedReviewPolicyResult = "SUCCESS"
)

// KYCState defines the current state of the KYC.
type KYCState string

// Possible values of KYCState.
const (
	Accepted    KYCState = "A"
	UnderReview KYCState = "R"
	Rejected    KYCState = "D"
)
