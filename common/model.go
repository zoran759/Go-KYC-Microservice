package common

import (
	"strings"
	"time"
)

// UserData defines the model for user data provided to KYC provider in order to check an individual.
type UserData struct {
	FirstName             string
	PaternalLastName      string
	LastName              string
	MiddleName            string
	LegalName             string
	LatinISO1Name         string
	Email                 string
	Gender                Gender
	DateOfBirth           Time
	PlaceOfBirth          string
	CountryOfBirthAlpha2  string
	StateOfBirth          string
	CountryAlpha2         string
	Nationality           string
	Phone                 string
	MobilePhone           string
	AddressString         string
	CurrentAddress        Address
	SupplementalAddresses []Address
	Documents             []Document
	Business              Business
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
	b := strings.Builder{}
	b.WriteString(a.BuildingNumber)
	if b.Len() > 0 {
		b.WriteString(" ")
	}
	b.WriteString(a.Street)

	return b.String()
}

// Time defines the model for time values.
type Time time.Time

// Format returns time string formatted according to the provided layout.
func (t Time) Format(layout string) string {
	if !time.Time(t).IsZero() {
		return time.Time(t).Format(layout)
	}
	return ""
}

// Business defines the model for a business.
type Business struct {
	Name                      string
	RegistrationNumber        string
	IncorporationDate         Time
	IncorporationJurisdiction string
}

// Document defines user's document.
type Document struct {
	Metadata DocumentMetadata
	Front    *DocumentFile
	Back     *DocumentFile
}

// DocumentMetadata defines a part of the Document model.
type DocumentMetadata struct {
	Type             DocumentType
	Country          string
	DateIssued       Time
	ValidUntil       Time
	Number           string
	CardFirst6Digits string
	CardLast4Digits  string
}

// DocumentFile defines document's file containing its original or an image.
type DocumentFile struct {
	Filename    string
	ContentType string
	Data        []byte
}

// DetailedKYCResult defines additional details about the verification process result.
type DetailedKYCResult struct {
	Finality KYCFinality
	Reasons  []string
}

// Request for the CheckCustomer handler
type CheckCustomerRequest struct {

	// KYC provider
	Provider KYCProvider

	// User Data
	UserData UserData
}

// Response for the CheckCustomer handler
type CheckCustomerResponse struct {

	// KYC result
	KYCResult KYCResult

	// KYC detailed result
	DetailedKYCResult *DetailedKYCResult
}
