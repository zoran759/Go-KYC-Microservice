package common

import "strings"

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
