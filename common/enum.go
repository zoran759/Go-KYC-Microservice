package common

// KYCResult defines the result of the KYC check.
type KYCResult int

// Possible KYCResult values.
const (
	Error    KYCResult = -1
	Approved KYCResult = 1
	Denied   KYCResult = 2
	Unclear  KYCResult = 3
)

// KYCFinality defines the finality of the result of the KYC check.
type KYCFinality int

// Possible KYCFinality values.
const (
	Final    KYCFinality = 1
	NonFinal KYCFinality = 2
	Unknown  KYCFinality = 3
)

// Gender defines user's gender.
type Gender int

// Gender values.
const (
	Male   Gender = 1
	Female Gender = 2
)

// DocumentType defines user's document type.
type DocumentType string

// Different document types.
const (
	// An id card. (It'll be used as SSN in IDology.)
	IDCard DocumentType = "ID_CARD"
	// An id card that's written in English.
	IDCardEng DocumentType = "ID_CARD_ENG"
	// A passport.
	Passport DocumentType = "PASSPORT"
	// A passport that's written in English.
	PassportEng DocumentType = "PASSPORT_ENG"
	// A driving license.
	Drivers DocumentType = "DRIVERS"
	// A driving license that's written in English.
	DriversEng DocumentType = "DRIVERS_ENG"
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
	// User's proof of income (i.e. Pay Stub)
	ProofOfIncome DocumentType = "PROOF_OF_INCOME"
	// User's proof of account ownership (i.e. Bank statement)
	ProofOfAccount DocumentType = "PROOF_OF_ACCOUNT"
	// ACH authorization signed by the user
	ACHAuthorization DocumentType = "AUTHORIZATION"
	// Background check of the user
	BackgroundCheck DocumentType = "BK_CHECK"
	// SSN Card of the user
	SSN DocumentType = "SSN"
	// Form 147C issued to the business
	EINDocument DocumentType = "EIN_DOC"
	// W-9 with EIN Number
	W9Document DocumentType = "W9_DOC"
	// W-8 Document
	W8Document DocumentType = "W8_DOC"
	// W-2 Document
	W2Document DocumentType = "W2_DOC"
	// Voided Check of the Individual/business
	VoidedCheck DocumentType = "VOIDED_CHECK"
	// Articles of Incorporation
	ArticlesOfIncorporation DocumentType = "AOI"
	// Bylaw document
	BylawsDocument DocumentType = "BYLAWS_DOC"
	// Letter of Engagement
	LetterOfEngagement DocumentType = "LOE"
	// CIP & Business description document
	CIPDoc DocumentType = "CIP_DOC"
	// Subscription agreement
	SubscriptionAgreement DocumentType = "SUBSCRIPTION_AGREEMENT"
	// Promissory Note
	PromissoryNote DocumentType = "PROMISSORY_NOTE"
	// Reg GG Form
	RegGG DocumentType = "REG_GG"
	// DBA or Fictitious Name Documentation
	DBADoc DocumentType = "DBA_DOC"
	// Deposit Agreement
	DepositAgreement DocumentType = "DEPOSIT_AGREEMENT"
	// Should be used only when nothing else applies.
	Other DocumentType = "OTHER"
)
