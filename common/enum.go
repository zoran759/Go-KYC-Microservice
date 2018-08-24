package common

// A type that defines the result of the KYC check.
type KYCResult int

const Error KYCResult = -1
const Approved KYCResult = 1
const Denied KYCResult = 2
const Unclear KYCResult = 3

// A type that defines the finality of the result of the KYC check.
type KYCFinality int

const Final KYCFinality = 1
const NonFinal KYCFinality = 2
const Unknown KYCFinality = 3

// A type that defines user's gender.
type Gender int

const Male Gender = 1
const Female Gender = 2

// A type that defines user's document type.
type DocumentType string

// Different document types.
const (
	// An id card.
	IDCard DocumentType = "ID_CARD"
	// A passport.
	Passport DocumentType = "PASSPORT"
	// A driving license.
	Drivers DocumentType = "DRIVERS"
	// A bank card, like Visa or Maestro.
	BankCard DocumentType = "BANK_CARD"
	// An utility bill.
	UtilityBill DocumentType = "UTILITY_BILL"
	// A Russian individual insurance account number (SNILS).
	SNILS DocumentType = "SNILS"
	// A selfie image. No additional metadata should be sent.
	Selfie DocumentType = "SELFIE"
	// A profile image, aka avatar. No additional metadata should be sent.
	ProfileImage DocumentType = "PROFILE_IMAGE"
	// Photo from some identification document (like a photo from a passport). No additional metadata should be sent.
	IDDocPhoto DocumentType = "ID_DOC_PHOTO"
	// Agreement of some sort, e.g. for processing personal info.
	Agreement DocumentType = "AGREEMENT"
	// Some sort of contract.
	Contract DocumentType = "CONTRACT"
	// Residence permit or registration document in the foreign city/country.
	ResidencePermit DocumentType = "RESIDENCE_PERMIT"
	// A document from an employer, e.g. proof that a user works there.
	EmploymentCertificate DocumentType = "EMPLOYMENT_CERTIFICATE"
	// Translation of the driving license required in the target country.
	DriversTranslation DocumentType = "DRIVERS_TRANSLATION"
	// Should be used only when nothing else applies.
	Other DocumentType = "OTHER"
)
