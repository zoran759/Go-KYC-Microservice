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
	AccountName           string
	Email                 string
	IPaddress             string
	Gender                Gender
	DateOfBirth           Time
	PlaceOfBirth          string
	CountryOfBirthAlpha2  string
	StateOfBirth          string
	CountryAlpha2         string
	Nationality           string
	Phone                 string
	MobilePhone           string
	Location              *Location
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

// Location defines the model for the geopositional data.
type Location struct {
	Latitude  string
	Longitude string
}
