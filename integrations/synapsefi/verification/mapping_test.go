package verification

import (
	"encoding/base64"
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
		VideoAuth: &common.VideoAuth{
			Filename:    "my_video.mp4",
			ContentType: "video/mp4",
			Data: []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70,
				0x6D, 0x70, 0x34, 0x32, 0x00, 0x00, 0x00, 0x00,
				0x69, 0x73, 0x6F, 0x6D, 0x6D, 0x70, 0x34, 0x32,
				0x00, 0x14, 0xDF, 0x3B, 0x6D, 0x6F, 0x6F, 0x76},
		},
	}

	userRequest := MapCustomerToCreateUserRequest(customer, true)

	assert := assert.New(t)

	assert.Len(userRequest.PhoneNumbers, 2)
	assert.Equal("Phone", userRequest.PhoneNumbers[0])
	assert.Equal("MobilePhone", userRequest.PhoneNumbers[1])

	assert.Len(userRequest.Logins, 1)
	assert.Equal(Login{
		Email: "Email",
		Scope: "READ_AND_WRITE",
	}, userRequest.Logins[0])

	assert.Len(userRequest.LegalNames, 1)
	assert.Equal("FirstName MiddleName LastName", userRequest.LegalNames[0])

	assert.Len(userRequest.Documents, 1)

	document := userRequest.Documents[0]

	assert.Equal("FirstName MiddleName LastName", document.OwnerName)
	assert.Equal("Email", document.Email)
	assert.Equal("Phone", document.PhoneNumber)
	assert.Equal("127.0.0.1", document.IPAddress)
	assert.Equal("M", document.EntityType)
	assert.Equal("Not Known", document.EntityScope)
	assert.Equal(2, document.DayOfBirth)
	assert.Equal(1, document.MonthOfBirth)
	assert.Equal(1967, document.YearOfBirth)
	assert.Equal("BuildingNumber1 Street1", document.AddressStreet)
	assert.Equal("Town1", document.AddressCity)
	assert.Equal("SPC1", document.AddressSubdivision)
	assert.Equal("PostCode1", document.AddressPostalCode)
	assert.Equal("Country1", document.AddressCountryCode)

	if assert.Len(document.PhysicalDocs, 3) {
		assert.Equal(
			"GOVT_ID_INT",
			document.PhysicalDocs[0].DocumentType,
		)
		assert.Equal(
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			document.PhysicalDocs[0].DocumentValue,
		)

		assert.Equal(
			"SELFIE",
			document.PhysicalDocs[1].DocumentType,
		)
		assert.Equal(
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			document.PhysicalDocs[1].DocumentValue,
		)

		assert.Equal("VIDEO_AUTHORIZATION", document.PhysicalDocs[2].DocumentType)
		assert.Equal("data:video/mp4;base64,"+base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
			document.PhysicalDocs[2].DocumentValue)
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
		VideoAuth: &common.VideoAuth{
			Filename:    "my_video.mp4",
			ContentType: "video/mp4",
			Data: []byte{0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70,
				0x6D, 0x70, 0x34, 0x32, 0x00, 0x00, 0x00, 0x00,
				0x69, 0x73, 0x6F, 0x6D, 0x6D, 0x70, 0x34, 0x32,
				0x00, 0x14, 0xDF, 0x3B, 0x6D, 0x6F, 0x6F, 0x76},
		},
	}

	docsRequest := MapDocumentsToCreateUserRequest(customer)

	assert := assert.New(t)

	assert.Equal("FirstName MiddleName LastName", docsRequest.Documents.OwnerName)
	assert.Equal("Email", docsRequest.Documents.Email)
	assert.Equal("Phone", docsRequest.Documents.PhoneNumber)
	assert.Equal("127.0.0.1", docsRequest.Documents.IPAddress)
	assert.Equal("M", docsRequest.Documents.EntityType)
	assert.Equal("Not Known", docsRequest.Documents.EntityScope)
	assert.Equal(2, docsRequest.Documents.DayOfBirth)
	assert.Equal(1, docsRequest.Documents.MonthOfBirth)
	assert.Equal(1967, docsRequest.Documents.YearOfBirth)
	assert.Equal("BuildingNumber1 Street1", docsRequest.Documents.AddressStreet)
	assert.Equal("Town1", docsRequest.Documents.AddressCity)
	assert.Equal("SPC1", docsRequest.Documents.AddressSubdivision)
	assert.Equal("PostCode1", docsRequest.Documents.AddressPostalCode)
	assert.Equal("CountryAlpha2", docsRequest.Documents.AddressCountryCode)

	if assert.Len(docsRequest.Documents.PhysicalDocs, 3) {
		assert.Equal(
			"GOVT_ID_INT",
			docsRequest.Documents.PhysicalDocs[0].DocumentType,
		)
		assert.Equal(
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			docsRequest.Documents.PhysicalDocs[0].DocumentValue,
		)

		assert.Equal(
			"SELFIE",
			docsRequest.Documents.PhysicalDocs[1].DocumentType,
		)
		assert.Equal(
			"data:image/png;base64,"+base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			docsRequest.Documents.PhysicalDocs[1].DocumentValue,
		)

		assert.Equal("VIDEO_AUTHORIZATION", docsRequest.Documents.PhysicalDocs[2].DocumentType)
		assert.Equal("data:video/mp4;base64,"+base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
			docsRequest.Documents.PhysicalDocs[2].DocumentValue)
	}
}

func Test_mapCustomerGender(t *testing.T) {
	assert.Equal(t, "M", mapCustomerGender(common.Male))
	assert.Equal(t, "F", mapCustomerGender(common.Female))
	assert.Equal(t, "O", mapCustomerGender(common.Gender(0)))

}

func Test_mapDocumentType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("GOVT_ID_INT", mapDocumentType("IDCard"))
	assert.Equal("GOVT_ID_INT", mapDocumentType("DriverLicense"))
	assert.Equal("GOVT_ID_INT", mapDocumentType("Passport"))

	assert.Equal("SELFIE", mapDocumentType("Selfie"))

	assert.Equal("PROOF_OF_ADDRESS", mapDocumentType("UtilityBill"))
	assert.Equal("OTHER", mapDocumentType("ResidencePermit"))

	assert.Equal("LEGAL_AGREEMENT", mapDocumentType("Agreement"))
	assert.Equal("LEGAL_AGREEMENT", mapDocumentType("Contract"))

	assert.Equal("OTHER", mapDocumentType("DriversTranslation"))
}

func TestMapUserToOauth(t *testing.T) {

	refreshToken := "sometoken"

	oauthRequest := MapUserToOauth(refreshToken)

	assert.Equal(t, "sometoken", oauthRequest.RefreshToken)
}

func TestMapResponseError(t *testing.T) {
	err := []byte(`{"error":{"en":"Invalid/expired oauth_key."},"error_code":"110","http_code":"401","success":false}`)

	response, _ := MapResponseError(err)

	assert.Equal(t, "Invalid/expired oauth_key.", response.Error())

	err = []byte("wrong response")

	resp, err1 := MapResponseError(err)

	assert.Nil(t, resp)
	assert.Equal(t, "invalid character 'w' looking for beginning of value", err1.Error())
}
