package verification

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
)

func TestMapCustomerToDataFields(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := common.UserData{
		FirstName:            "FirstName",
		PaternalLastName:     "PaternalLastName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		LatinISO1Name:        "LATIN",
		Email:                "Email",
		Gender:               common.Male,
		DateOfBirth:          testTime,
		PlaceOfBirth:         "PlaceOfBirth",
		CountryOfBirthAlpha2: "CountryOfBirth",
		StateOfBirth:         "StateOfBirth",
		CountryAlpha2:        "CountryAlpha2",
		Nationality:          "Nationality",
		Phone:                "Phone",
		MobilePhone:          "MobilePhone",
		CurrentAddress: common.Address{
			CountryAlpha2:     "Country1",
			County:            "County1",
			State:             "State1",
			Town:              "Town1",
			Suburb:            "Suburb1",
			Street:            "Street1",
			StreetType:        "StreetType1",
			SubStreet:         "SubStreet1",
			BuildingName:      "BuildingName1",
			BuildingNumber:    "BuildingNumber1",
			FlatNumber:        "FlatNumber1",
			PostCode:          "PostCode1",
			StateProvinceCode: "SPC1",
			PostOfficeBox:     "POB1",
			StartDate:         testTime,
			EndDate:           testTime,
		},
		Business: common.Business{
			Name:                      "BusinessName",
			RegistrationNumber:        "RegNumber",
			IncorporationDate:         testTime,
			IncorporationJurisdiction: "IncorporationJurisdiction",
		},
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:       common.IDCard,
					Country:    "Country",
					DateIssued: testTime,
					ValidUntil: testTime,
					Number:     "Number",
				},
				Front: &common.DocumentFile{
					Filename:    "Filename",
					ContentType: "ContentType",
					Data:        []byte{1, 2, 3, 4, 5, 6, 7},
				},
				Back: &common.DocumentFile{
					Filename:    "Filename2",
					ContentType: "ContentType2",
					Data:        []byte{7, 6, 5, 4, 3, 2, 1},
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:       common.Selfie,
					Country:    "Country",
					DateIssued: testTime,
					ValidUntil: testTime,
					Number:     "Number",
				},
				Front: &common.DocumentFile{
					Filename:    "Filename",
					ContentType: "ContentType",
					Data:        []byte{1, 2, 3, 4, 5, 6, 7},
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:       common.Passport,
					Country:    "Country",
					DateIssued: testTime,
					ValidUntil: testTime,
					Number:     "Number",
				},
				Front: &common.DocumentFile{
					Filename:    "Filename",
					ContentType: "ContentType",
					Data:        []byte{1, 2, 3, 4, 5, 6, 7},
				},
				Back: &common.DocumentFile{
					Filename:    "Filename2",
					ContentType: "ContentType2",
					Data:        []byte{7, 6, 5, 4, 3, 2, 1},
				},
			},
		},
	}

	dataFields := MapCustomerToDataFields(customer)

	if assert.NotNil(t, dataFields.PersonInfo) {
		personInfo := dataFields.PersonInfo

		assert.Equal(t, customer.FirstName, personInfo.FirstGivenName)
		assert.Equal(t, customer.PaternalLastName, personInfo.FirstSurName)
		assert.Equal(t, customer.LastName, personInfo.SecondSurname)
		assert.Equal(t, customer.MiddleName, personInfo.MiddleName)
		assert.Equal(t, customer.LatinISO1Name, personInfo.ISOLatin1Name)
		assert.Equal(t, "M", personInfo.Gender)
		assert.Equal(t, time.Time(customer.DateOfBirth).Day(), personInfo.DayOfBirth)
		assert.Equal(t, int(time.Time(customer.DateOfBirth).Month()), personInfo.MonthOfBirth)
		assert.Equal(t, time.Time(customer.DateOfBirth).Year(), personInfo.YearOfBirth)
	}

	if assert.NotNil(t, dataFields.Communication) {
		communication := dataFields.Communication

		assert.Equal(t, customer.Phone, communication.Telephone)
		assert.Equal(t, customer.MobilePhone, communication.MobileNumber)
		assert.Equal(t, customer.Email, communication.EmailAddress)
	}

	if assert.NotNil(t, dataFields.Business) {
		business := dataFields.Business

		assert.Equal(t, customer.Business.Name, business.BusinessName)
		assert.Equal(t, customer.Business.RegistrationNumber, business.BusinessRegistrationNumber)
		assert.Equal(t, customer.Business.IncorporationJurisdiction, business.JurisdictionOfIncorporation)
		assert.Equal(t, time.Time(customer.Business.IncorporationDate).Day(), business.DayOfIncorporation)
		assert.Equal(t, int(time.Time(customer.Business.IncorporationDate).Month()), business.MonthOfIncorporation)
		assert.Equal(t, time.Time(customer.Business.IncorporationDate).Year(), business.YearOfIncorporation)
	}

	if assert.NotNil(t, dataFields.Location) {
		location := dataFields.Location

		assert.Equal(t, customer.CurrentAddress.CountryAlpha2, location.Country)
		assert.Equal(t, customer.CurrentAddress.County, location.County)
		assert.Equal(t, customer.CurrentAddress.Town, location.City)
		assert.Equal(t, customer.CurrentAddress.Suburb, location.Suburb)
		assert.Equal(t, customer.CurrentAddress.Street, location.StreetName)
		assert.Equal(t, customer.CurrentAddress.StreetType, location.StreetType)
		assert.Equal(t, customer.CurrentAddress.BuildingName, location.BuildingName)
		assert.Equal(t, customer.CurrentAddress.BuildingNumber, location.BuildingNumber)
		assert.Equal(t, customer.CurrentAddress.FlatNumber, location.UnitNumber)
		assert.Equal(t, customer.CurrentAddress.PostCode, location.PostalCode)
		assert.Equal(t, customer.CurrentAddress.StateProvinceCode, location.StateProvinceCode)
		assert.Equal(t, customer.CurrentAddress.PostOfficeBox, location.POBox)
	}

	if assert.NotNil(t, dataFields.Document) {
		document := dataFields.Document
		commonDocument := customer.Documents[0]
		selfie := customer.Documents[1]

		assert.Equal(t, "IdentityCard", document.DocumentType)
		assert.Equal(t, commonDocument.Front.Data, document.DocumentFrontImage)
		assert.Equal(t, commonDocument.Back.Data, document.DocumentBackImage)
		assert.Equal(t, selfie.Front.Data, document.LivePhoto)
	}

	dataFields = MapCustomerToDataFields(common.UserData{})

	assert.Nil(t, dataFields.Communication)
	assert.Nil(t, dataFields.Business)
	assert.Nil(t, dataFields.Location)
	assert.Nil(t, dataFields.Document)
}

func TestMapCustomerToDataFieldsNoSupportedDocuments(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := common.UserData{
		FirstName:            "FirstName",
		PaternalLastName:     "PaternalLastName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		LatinISO1Name:        "LATIN",
		Email:                "Email",
		Gender:               common.Male,
		DateOfBirth:          testTime,
		PlaceOfBirth:         "PlaceOfBirth",
		CountryOfBirthAlpha2: "CountryOfBirth",
		StateOfBirth:         "StateOfBirth",
		CountryAlpha2:        "CountryAlpha2",
		Nationality:          "Nationality",
		Phone:                "Phone",
		MobilePhone:          "MobilePhone",
		CurrentAddress: common.Address{
			CountryAlpha2:     "Country1",
			County:            "County1",
			State:             "State1",
			Town:              "Town1",
			Suburb:            "Suburb1",
			Street:            "Street1",
			StreetType:        "StreetType1",
			SubStreet:         "SubStreet1",
			BuildingName:      "BuildingName1",
			BuildingNumber:    "BuildingNumber1",
			FlatNumber:        "FlatNumber1",
			PostCode:          "PostCode1",
			StateProvinceCode: "SPC1",
			PostOfficeBox:     "POB1",
			StartDate:         testTime,
			EndDate:           testTime,
		},
		Business: common.Business{
			Name:                      "BusinessName",
			RegistrationNumber:        "RegNumber",
			IncorporationDate:         testTime,
			IncorporationJurisdiction: "IncorporationJurisdiction",
		},
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:       common.SNILS,
					Country:    "Country",
					DateIssued: testTime,
					ValidUntil: testTime,
					Number:     "Number",
				},
				Front: &common.DocumentFile{
					Filename:    "Filename",
					ContentType: "ContentType",
					Data:        []byte{1, 2, 3, 4, 5, 6, 7},
				},
				Back: &common.DocumentFile{
					Filename:    "Filename2",
					ContentType: "ContentType2",
					Data:        []byte{7, 6, 5, 4, 3, 2, 1},
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:       common.Selfie,
					Country:    "Country",
					DateIssued: testTime,
					ValidUntil: testTime,
					Number:     "Number",
				},
				Front: &common.DocumentFile{
					Filename:    "Filename",
					ContentType: "ContentType",
					Data:        []byte{1, 2, 3, 4, 5, 6, 7},
				},
			},
		},
	}

	dataFields := MapCustomerToDataFields(customer)

	assert.Nil(t, dataFields.Document)
}

func Test_mapDocumentType(t *testing.T) {
	assert.Equal(t, "DrivingLicence", mapCustomerDocumentType(common.Drivers))
	assert.Equal(t, "IdentityCard", mapCustomerDocumentType(common.IDCard))
	assert.Equal(t, "Passport", mapCustomerDocumentType(common.Passport))
	assert.Equal(t, "ResidencePermit", mapCustomerDocumentType(common.ResidencePermit))
}

func Test_mapGender(t *testing.T) {
	assert.Equal(t, "M", mapGender(common.Male))
	assert.Equal(t, "F", mapGender(common.Female))
	assert.Empty(t, mapGender(common.Gender(10)))
}
