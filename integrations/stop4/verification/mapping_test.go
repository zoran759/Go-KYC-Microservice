package verification

import (
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCustomerToVerificationRequest(t *testing.T) {
	testTime := common.Time(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

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
		UtilityBill: &common.UtilityBill{
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
		Other: &common.Other{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			Number:        "Number",
		},
	}

	verificationRequest := MapCustomerToVerificationRequest(customer)

	assert.Equal(t, customer.CurrentAddress.String(), verificationRequest.CustomerInformation.Address1)
	assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), verificationRequest.CustomerInformation.Dob)
	assert.Equal(t, customer.Email, verificationRequest.CustomerInformation.Email)
	assert.Equal(t, customer.FirstName, verificationRequest.CustomerInformation.FirstName)
	assert.Equal(t, customer.LastName, verificationRequest.CustomerInformation.LastName)
	assert.Equal(t, customer.MiddleName, verificationRequest.CustomerInformation.MiddleName)
	assert.Equal(t, customer.Phone, verificationRequest.CustomerInformation.Phone1)
	assert.Equal(t, customer.MobilePhone, verificationRequest.CustomerInformation.Phone2)

}
