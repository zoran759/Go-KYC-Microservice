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
	UKNHSNumber              string
	UKNINumber               string
	CurrentAddress           Address
	SupplementalAddresses    []Address
	Location                 *Location
	Business                 *Business
	Passport                 *Passport
	IDCard                   *IDCard
	SNILS                    *SNILS
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
}

// Address defines user's address.
type Address struct {
	CountryAlpha2     string
	County            string
	State             string
	Town              string
	Suburb            string
	Street            string
	StreetType        string
	SubStreet         string
	BuildingName      string
	BuildingNumber    string
	FlatNumber        string
	PostOfficeBox     string
	PostCode          string
	StateProvinceCode string
	StartDate         Time
	EndDate           Time
}

// StreetAddress is a helper func that returns street part of the address.
func (a Address) StreetAddress() string {
	// ATM, USPS standard is used. Maybe, we need to take into count Country's specifics.
	b := strings.Builder{}
	b.WriteString(a.BuildingNumber)
	if len(a.Street) > 0 {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
		b.WriteString(a.Street)
	}

	return b.String()
}

// HouseStreetApartment returns street address string in the form required for some providers.
// It includes house number, street name and apartment number.
func (a Address) HouseStreetApartment() string {
	insertWhitespace := func(b *strings.Builder) {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
	}

	b := &strings.Builder{}
	b.WriteString(a.BuildingNumber)
	if len(a.Street) > 0 {
		insertWhitespace(b)
		b.WriteString(a.Street)
	}
	if len(a.FlatNumber) > 0 {
		insertWhitespace(b)
		b.WriteString(a.FlatNumber)
	}

	return b.String()
}

// String returns string representation of the address.
func (a Address) String() string {
	// ATM, USPS standard is used. Maybe, we need to take into count Country's specifics.
	insertWhitespace := func(b *strings.Builder) {
		if b.Len() > 0 {
			b.WriteString(" ")
		}
	}

	b := &strings.Builder{}
	if len(a.PostOfficeBox) > 0 {
		b.WriteString("PO BOX ")
		b.WriteString(a.PostOfficeBox)
	} else {
		b.WriteString(a.BuildingNumber)
		if len(a.Street) > 0 {
			insertWhitespace(b)
			b.WriteString(a.Street)
		}
		if len(a.FlatNumber) > 0 {
			insertWhitespace(b)
			b.WriteString(a.FlatNumber)
		}
	}
	if len(a.County) > 0 {
		insertWhitespace(b)
		b.WriteString(a.County)
	}
	if len(a.Town) > 0 {
		insertWhitespace(b)
		b.WriteString(a.Town)
	}
	if len(a.StateProvinceCode) > 0 {
		insertWhitespace(b)
		b.WriteString(a.StateProvinceCode)
	}
	if len(a.PostCode) > 0 {
		insertWhitespace(b)
		b.WriteString(a.PostCode)
	}
	if len(a.CountryAlpha2) > 0 {
		insertWhitespace(b)
		if a3, ok := CountryAlpha2ToAlpha3[a.CountryAlpha2]; ok {
			b.WriteString(a3)
		}
	}

	return b.String()
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

// FullName builds and returns full name of the customer.
func (u *UserData) FullName() string {
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
