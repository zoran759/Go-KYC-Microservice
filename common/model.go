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

// HouseStreetApartmentAddress returns street address string in the form required for some providers.
// It includes house number, street name and apartment number.
func (a Address) HouseStreetApartmentAddress() string {
	writeSpace := func(b *strings.Builder) {
		if b.Len() == 0 {
			return
		}
		b.WriteString(" ")
	}

	b := strings.Builder{}
	b.WriteString(a.BuildingNumber)
	if len(a.Street) > 0 {
		writeSpace(&b)
		b.WriteString(a.Street)
	}
	if len(a.FlatNumber) > 0 {
		writeSpace(&b)
		b.WriteString(a.FlatNumber)
	}

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
