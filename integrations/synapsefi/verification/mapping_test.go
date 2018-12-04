package verification

import (
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCustomerToCreateUserRequest(t *testing.T) {
	testTime := common.Time(time.Date(1967, 1, 2, 0, 0, 0, 0, time.UTC))

	customer := common.UserData{
		FirstName:            "FirstName",
		PaternalLastName:     "PaternalLastName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		LatinISO1Name:        "LATIN",
		Email:                "Email",
		IPaddress:            "127.0.0.1",
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

	userRequest := MapCustomerToCreateUserRequest(customer, true)

	assert.Len(t, userRequest.PhoneNumbers, 2)
	assert.Equal(t, "Phone", userRequest.PhoneNumbers[0])
	assert.Equal(t, "MobilePhone", userRequest.PhoneNumbers[1])

	assert.Len(t, userRequest.Logins, 1)
	assert.Equal(t, Login{
		Email: "Email",
		Scope: "READ_AND_WRITE",
	}, userRequest.Logins[0])

	assert.Len(t, userRequest.LegalNames, 1)
	assert.Equal(t, "FirstName MiddleName LastName", userRequest.LegalNames[0])

	assert.Len(t, userRequest.Documents, 1)

	document := userRequest.Documents[0]

	assert.Equal(t, "FirstName MiddleName LastName", document.OwnerName)
	assert.Equal(t, "Email", document.Email)
	assert.Equal(t, "Phone", document.PhoneNumber)
	assert.Equal(t, "127.0.0.1", document.IPAddress)
	assert.Equal(t, "M", document.EntityType)
	assert.Equal(t, "Not Known", document.EntityScope)
	assert.Equal(t, 2, document.DayOfBirth)
	assert.Equal(t, 1, document.MonthOfBirth)
	assert.Equal(t, 1967, document.YearOfBirth)
	assert.Equal(t, "BuildingNumber1 Street1 StreetType1", document.AddressStreet)
	assert.Equal(t, "Town1", document.AddressCity)
	assert.Equal(t, "SPC1", document.AddressSubdivision)
	assert.Equal(t, "PostCode1", document.AddressPostalCode)
	assert.Equal(t, "Country1", document.AddressCountryCode)

	if assert.Len(t, document.PhysicalDocs, 2) {
		assert.Equal(t,
			"GOVT_ID_INT",
			document.PhysicalDocs[0].DocumentType,
		)
		assert.Equal(t,
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			document.PhysicalDocs[0].DocumentValue,
		)

		assert.Equal(t,
			"SELFIE",
			document.PhysicalDocs[1].DocumentType,
		)
		assert.Equal(t,
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			document.PhysicalDocs[1].DocumentValue,
		)
	}
}

func TestMapDocumentsToCreateUserRequest(t *testing.T) {
	testTime := common.Time(time.Date(1967, 1, 2, 0, 0, 0, 0, time.UTC))

	customer := common.UserData{
		FirstName:            "FirstName",
		PaternalLastName:     "PaternalLastName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		LatinISO1Name:        "LATIN",
		IPaddress:            "127.0.0.1",
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

	docsRequest := MapDocumentsToCreateUserRequest(customer)

	assert.Equal(t, "FirstName MiddleName LastName", docsRequest.Documents.OwnerName)
	assert.Equal(t, "Email", docsRequest.Documents.Email)
	assert.Equal(t, "Phone", docsRequest.Documents.PhoneNumber)
	assert.Equal(t, "127.0.0.1", docsRequest.Documents.IPAddress)
	assert.Equal(t, "M", docsRequest.Documents.EntityType)
	assert.Equal(t, "Not Known", docsRequest.Documents.EntityScope)
	assert.Equal(t, 2, docsRequest.Documents.DayOfBirth)
	assert.Equal(t, 1, docsRequest.Documents.MonthOfBirth)
	assert.Equal(t, 1967, docsRequest.Documents.YearOfBirth)
	assert.Equal(t, "BuildingNumber1 Street1 StreetType1", docsRequest.Documents.AddressStreet)
	assert.Equal(t, "Town1", docsRequest.Documents.AddressCity)
	assert.Equal(t, "SPC1", docsRequest.Documents.AddressSubdivision)
	assert.Equal(t, "PostCode1", docsRequest.Documents.AddressPostalCode)
	assert.Equal(t, "CountryAlpha2", docsRequest.Documents.AddressCountryCode)

	if assert.Len(t, docsRequest.Documents.PhysicalDocs, 2) {
		assert.Equal(t,
			"GOVT_ID_INT",
			docsRequest.Documents.PhysicalDocs[0].DocumentType,
		)
		assert.Equal(t,
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			docsRequest.Documents.PhysicalDocs[0].DocumentValue,
		)

		assert.Equal(t,
			"SELFIE",
			docsRequest.Documents.PhysicalDocs[1].DocumentType,
		)
		assert.Equal(t,
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			docsRequest.Documents.PhysicalDocs[1].DocumentValue,
		)
	}
}

func Test_mapCustomerGender(t *testing.T) {
	assert.Equal(t, "M", mapCustomerGender(common.Male))
	assert.Equal(t, "F", mapCustomerGender(common.Female))
	assert.Equal(t, "O", mapCustomerGender(common.Gender(0)))

}

func Test_mapDocumentType(t *testing.T) {

	assert.Equal(t, "GOVT_ID_INT", mapDocumentType("IDCard"))
	assert.Equal(t, "GOVT_ID_INT", mapDocumentType("DriverLicense"))
	assert.Equal(t, "GOVT_ID_INT", mapDocumentType("Passport"))

	assert.Equal(t, "SELFIE", mapDocumentType("Selfie"))

	assert.Equal(t, "PROOF_OF_ADDRESS", mapDocumentType("UtilityBill"))
	assert.Equal(t, "OTHER", mapDocumentType("ResidencePermit"))

	assert.Equal(t, "LEGAL_AGREEMENT", mapDocumentType("Agreement"))
	assert.Equal(t, "LEGAL_AGREEMENT", mapDocumentType("Contract"))

	assert.Equal(t, "OTHER", mapDocumentType("DriversTranslation"))
}

func TestMapUserToOauth(t *testing.T) {

	refreshToken := "sometoken"

	oauthRequest := MapUserToOauth(refreshToken)

	assert.Equal(t, "sometoken", oauthRequest.RefreshToken)
}

func TestMapResponseError(t *testing.T) {
	err := []byte(`{"error":{"en":"Invalid/expired oauth_key."},"error_code":"110","http_code":"401","success":false}`)

	response, _ := MapResponseError(err)

	assert.Equal(t, errors.New("invalid/expired oauth_key"), response)
}
