package common

import (
	"strings"
	"time"
)

// UserData defines the model for user data provided to KYC provider in order to check an individual.
type UserData struct {
	FirstName                string
	LastName                 string
	MaternalLastName         string
	MiddleName               string
	FullName                 string
	LegalName                string
	LatinISO1Name            string
	AccountName              string
	Email                    string
	IPaddress                string
	Gender                   Gender
	DateOfBirth              Time
	PlaceOfBirth             string
	CountryOfBirthAlpha2     string
	StateOfBirth             string
	CountryAlpha2            string
	Nationality              string
	Phone                    string
	MobilePhone              string
	BankAccountNumber        string
	VehicleRegistrationPlate string
	CurrentAddress           Address
	SupplementalAddresses    []Address
	Location                 *Location
	Business                 *Business
	Passport                 *Passport
	IDCard                   *IDCard
	SNILS                    *SNILS
	HealthID                 *HealthID
	SocialServiceID          *SocialServiceID
	TaxID                    *TaxID
	DriverLicense            *DriverLicense
	DriverLicenseTranslation *DriverLicenseTranslation
	CreditCard               *CreditCard
	DebitCard                *DebitCard
	UtilityBill              *UtilityBill
	ResidencePermit          *ResidencePermit
	Agreement                *Agreement
	EmploymentCertificate    *EmploymentCertificate
	Contract                 *Contract
	DocumentPhoto            *DocumentPhoto
	Selfie                   *Selfie
	Avatar                   *Avatar
	Other                    *Other
	VideoAuth                *VideoAuth
}

// Location defines the model for the geopositional data.
type Location struct {
	Latitude  string
	Longitude string
}

// Business defines the model for a business.
type Business struct {
	Name                      string
	RegistrationNumber        string
	IncorporationDate         Time
	IncorporationJurisdiction string
}

// DocumentFile defines document's file containing its original or an image.
type DocumentFile struct {
	Filename    string
	ContentType string
	Data        []byte
}

// KYCDetails defines additional details about the verification result.
type KYCDetails struct {
	Finality KYCFinality
	Reasons  []string
}

// KYCResult represents the verification result.
type KYCResult struct {
	Status      KYCStatus
	Details     *KYCDetails
	ErrorCode   string
	StatusCheck *KYCStatusCheck
}

// KYCStatusCheck contains data required to do status check requests if needed.
type KYCStatusCheck struct {
	Provider    KYCProvider
	ReferenceID string
	LastCheck   time.Time
}

/*******************************************************************/
/* Below are the models representing different types of documents. */
/* Please, add new models for documents after this note.           */
/*******************************************************************/

// Passport represents the passport.
type Passport struct {
	Number        string
	Mrz1          string
	Mrz2          string
	CountryAlpha2 string
	State         string
	IssuedDate    Time
	ValidUntil    Time
	Image         *DocumentFile
}

// IDCard represents the id card.
type IDCard struct {
	Number        string
	CountryAlpha2 string
	IssuedDate    Time
	Image         *DocumentFile
}

// SNILS represents the Russian individual insurance account number.
type SNILS struct {
	Number     string
	IssuedDate Time
	Image      *DocumentFile
}

// DriverLicense represents the driver/driving license.
type DriverLicense struct {
	Number        string
	Version       string
	CountryAlpha2 string
	State         string
	IssuedDate    Time
	ValidUntil    Time
	FrontImage    *DocumentFile
	BackImage     *DocumentFile
}

// DriverLicenseTranslation represents the translated driver/driving license.
type DriverLicenseTranslation struct {
	Number        string
	CountryAlpha2 string
	State         string
	IssuedDate    Time
	ValidUntil    Time
	FrontImage    *DocumentFile
	BackImage     *DocumentFile
}

// CreditCard represents the banking credit card.
type CreditCard struct {
	Number     string
	ValidUntil Time
	Image      *DocumentFile
}

// DebitCard represents the banking debit card.
type DebitCard struct {
	Number     string
	ValidUntil Time
	Image      *DocumentFile
}

// UtilityBill represents the utility bill.
type UtilityBill struct {
	Number        string
	CountryAlpha2 string
	Image         *DocumentFile
}

// ResidencePermit represents the residence permit.
type ResidencePermit struct {
	CountryAlpha2 string
	IssuedDate    Time
	ValidUntil    Time
	Image         *DocumentFile
}

// Agreement represents an agreement of some sort, e.g. for processing personal info.
type Agreement struct {
	Image *DocumentFile
}

// Contract represents a contract of some sort.
type Contract struct {
	Image *DocumentFile
}

// EmploymentCertificate represents a document from an employer, e.g. proof that a user works there.
type EmploymentCertificate struct {
	IssuedDate Time
	Image      *DocumentFile
}

// Selfie represents the selfie.
type Selfie struct {
	Image *DocumentFile
}

// Avatar represents the profile image aka avatar.
type Avatar struct {
	Image *DocumentFile
}

// DocumentPhoto represents a photo from some document (like a photo from a passport).
type DocumentPhoto struct {
	Image *DocumentFile
}

// Other represents the model for other documents.
type Other struct {
	Number        string
	CountryAlpha2 string
	State         string
	IssuedDate    Time
	ValidUntil    Time
	Image         *DocumentFile
}

// HealthID represents National Health Service Identification Information.
type HealthID struct {
	Number string
	Image  *DocumentFile
}

// SocialServiceID represents National Social Service Identification Information
// (Social Security Number, Social Insurance Number, National Insurance Number).
type SocialServiceID struct {
	Number     string
	IssuedDate Time
	Image      *DocumentFile
}

// TaxID represents National Taxpayer Identification Information.
type TaxID struct {
	Number string
	Image  *DocumentFile
}

// VideoAuth represents authorization video of the customer.
type VideoAuth DocumentFile

// Fullname builds and returns full name of the customer.
func (u *UserData) Fullname() string {
	if len(u.FullName) != 0 {
		return u.FullName
	}

	insertWhitespace := func(b *strings.Builder) {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
	}

	b := &strings.Builder{}
	b.WriteString(u.FirstName)
	if len(u.MiddleName) > 0 {
		insertWhitespace(b)
		b.WriteString(u.MiddleName)
	}
	if len(u.LastName) > 0 {
		insertWhitespace(b)
		b.WriteString(u.LastName)
	}
	if len(u.MaternalLastName) > 0 {
		insertWhitespace(b)
		b.WriteString(u.MaternalLastName)
	}

	return b.String()
}
