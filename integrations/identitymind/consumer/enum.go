package consumer

// DeviceFingerprintType defines Device fingerprint type.
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

// DocumentType defines the Type of the Document for usage in Document Verification.
type DocumentType string

// Possible values of DocumentType.
const (
	Passport               DocumentType = "PP"
	DriverLicence          DocumentType = "DL"
	GovernmentIssuedIDCard DocumentType = "ID"
	ResidencePermit        DocumentType = "RP"
	UtilityBill            DocumentType = "UB"
)
