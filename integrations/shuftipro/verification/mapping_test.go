package verification

import (
	"testing"

	"gitlab.com/lambospeed/kyc/common"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestMapCustomerToVerificationRequest(t *testing.T) {
	testTime := common.Time(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

	customer := common.UserData{
		FirstName:        "FirstName",
		PaternalLastName: "PaternalLastName",
		LastName:         "LastName",
		MiddleName:       "MiddleName",
		LegalName:        "LegalName",
		LatinISO1Name:    "LATIN",
		Email:            "Email",
		Gender:           common.Male,
		DateOfBirth:      testTime,
		PlaceOfBirth:     "PlaceOfBirth",
		CountryOfBirth:   "CountryOfBirth",
		StateOfBirth:     "StateOfBirth",
		CountryAlpha2:    "CountryAlpha2",
		CountryAlpha3:    "CountryAlpha3",
		CountryName:      "CountryName",
		Nationality:      "Nationality",
		Phone:            "Phone",
		MobilePhone:      "MobilePhone",
		AddressString:    "AddressString",
		CurrentAddress: common.Address{
			Country:           "Country1",
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
					Type:       "Type",
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

	verificationRequest := MapCustomerToVerificationRequest(customer)

	assert.Equal(t, customer.CountryAlpha2, verificationRequest.Country)
	assert.Equal(t, customer.Phone, verificationRequest.PhoneNumber)
	assert.Equal(t, customer.Email, verificationRequest.Email)

	assert.Equal(t, customer.FirstName, verificationRequest.VerificationServices.FirstName)
	assert.Equal(t, customer.MiddleName, verificationRequest.VerificationServices.MiddleName)
	assert.Equal(t, customer.LastName, verificationRequest.VerificationServices.LastName)
	assert.Equal(t, customer.AddressString, verificationRequest.VerificationServices.Address)
	assert.Equal(t, "2000-01-01", verificationRequest.VerificationServices.DateOfBirth)

}
