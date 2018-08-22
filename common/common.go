package common

import (
	"time"
)

// User data provided to KYC in order to check an individual.
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
	CountryOfBirth        string
	StateOfBirth          string
	CountryAlpha2         string
	CountryAlpha3         string
	CountryName           string
	Nationality           string
	Phone                 string
	MobilePhone           string
	AddressString         string
	CurrentAddress        Address
	SupplementalAddresses []Address
	Documents             []Document
	Business              Business
}

type Address struct {
	Country           string
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

type Time time.Time

func (t Time) Format(layout string) string {
	if !time.Time(t).IsZero() {
		return time.Time(t).Format(layout)
	}
	return ""
}

type Business struct {
	Name                      string
	RegistrationNumber        string
	IncorporationDate         Time
	IncorporationJurisdiction string
}

type Document struct {
	Metadata DocumentMetadata
	Front    *DocumentFile
	Back     *DocumentFile
}

type DocumentMetadata struct {
	//TODO: Make document type an enum
	Type       string
	SubType    string
	Country    string
	DateIssued Time
	ValidUntil Time
	Number     string
}

type DocumentFile struct {
	Filename    string
	ContentType string
	Data        []byte
}

type DetailedKYCResult struct {
	Finality KYCFinality
	Reasons  []string
}
