package jumio

// ScanStatus represents the verification scan status.
type ScanStatus string

// Possible values of Status.
const (
	PendingStatus ScanStatus = "PENDING"
	DoneStatus    ScanStatus = "DONE"
	FailedStatus  ScanStatus = "FAILED"
)

// IDType represents document types returning and acceptable by the API.
type IDType string

// Possible values of IDType.
const (
	Passport       IDType = "PASSPORT"
	DrivingLicense IDType = "DRIVING_LICENSE"
	IDCard         IDType = "ID_CARD"
	Visa           IDType = "VISA"
	Unsupported    IDType = "UNSUPPORTED"
)

// DocumentStatus represents the scan result for the document.
type DocumentStatus string

// Possible values of DocumentStatus.
const (
	ApprovedVerified           DocumentStatus = "APPROVED_VERIFIED"
	DeniedFraud                DocumentStatus = "DENIED_FRAUD"
	DeniedUnsupportedIDType    DocumentStatus = "DENIED_UNSUPPORTED_ID_TYPE"
	DeniedUnsupportedIDCountry DocumentStatus = "DENIED_UNSUPPORTED_ID_COUNTRY"
	ErrorNotReadableID         DocumentStatus = "ERROR_NOT_READABLE_ID"
	NoIDUploaded               DocumentStatus = "NO_ID_UPLOADED"
)
