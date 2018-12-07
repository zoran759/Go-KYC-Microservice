package verification

import (
	"encoding/base64"
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCustomerToDataFields(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := &common.UserData{
		FirstName:            "FirstName",
		LastName:             "LastName",
		MaternalLastName:     "MaternalLastName",
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
		Business: &common.Business{
			Name:                      "BusinessName",
			RegistrationNumber:        "RegNumber",
			IncorporationDate:         testTime,
			IncorporationJurisdiction: "IncorporationJurisdiction",
		},

		IDCard: &common.IDCard{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			Number:        "Number",
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Passport: &common.Passport{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			Number:        "Number",
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
	}

	dataFields := MapCustomerToDataFields(customer)

	if assert.NotNil(t, dataFields.PersonInfo) {
		personInfo := dataFields.PersonInfo

		assert.Equal(t, customer.FirstName, personInfo.FirstGivenName)
		assert.Equal(t, customer.LastName, personInfo.FirstSurName)
		assert.Equal(t, customer.MaternalLastName, personInfo.SecondSurname)
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
		commonDocument := customer.Passport
		selfie := customer.Selfie

		assert.Equal(t, "Passport", document.DocumentType)
		assert.Equal(t, base64.StdEncoding.EncodeToString(commonDocument.Image.Data), document.DocumentFrontImage)
		assert.Equal(t, "", document.DocumentBackImage)
		assert.Equal(t, base64.StdEncoding.EncodeToString(selfie.Image.Data), document.LivePhoto)
	}

	dataFields = MapCustomerToDataFields(&common.UserData{})

	assert.Nil(t, dataFields.Communication)
	assert.Nil(t, dataFields.Business)
	assert.Nil(t, dataFields.Location)
	assert.Nil(t, dataFields.Document)
}

func TestMapCustomerToDataFieldsNoSupportedDocuments(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := common.UserData{
		FirstName:            "FirstName",
		MaternalLastName:     "MaternalLastName",
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
		Business: &common.Business{
			Name:                      "BusinessName",
			RegistrationNumber:        "RegNumber",
			IncorporationDate:         testTime,
			IncorporationJurisdiction: "IncorporationJurisdiction",
		},
		SNILS: &common.SNILS{
			IssuedDate: testTime,
			Number:     "Number",
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
	}

	dataFields := MapCustomerToDataFields(&customer)

	assert.Nil(t, dataFields.Document)
}

func Test_mapGender(t *testing.T) {
	assert.Equal(t, "M", mapGender(common.Male))
	assert.Equal(t, "F", mapGender(common.Female))
	assert.Empty(t, mapGender(common.Gender(10)))
}
