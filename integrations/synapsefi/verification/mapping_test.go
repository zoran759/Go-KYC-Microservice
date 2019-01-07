package verification

import (
	"encoding/base64"
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCustomerToUser(t *testing.T) {
	testTime := common.Time(time.Date(1967, 1, 2, 0, 0, 0, 0, time.UTC))

	customer := &common.UserData{
		FirstName:            "FirstName",
		MaternalLastName:     "MaternalLastName",
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

	user := MapCustomerToUser(customer)

	assert := assert.New(t)

	assert.Len(user.PhoneNumbers, 2)
	assert.Equal("Phone", user.PhoneNumbers[0])
	assert.Equal("MobilePhone", user.PhoneNumbers[1])

	assert.Len(user.Logins, 1)
	assert.Equal(Login{
		Email: "Email",
		Scope: "READ_AND_WRITE",
	}, user.Logins[0])

	assert.Len(user.LegalNames, 1)
	assert.Equal("LegalName", user.LegalNames[0])

	assert.Len(user.Documents, 1)

	document := user.Documents[0]

	assert.Equal("LegalName", document.OwnerName)
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

	if assert.Len(document.VirtualDocs, 2) {
		assert.Equal("PERSONAL_IDENTIFICATION", document.VirtualDocs[0].Type)
		assert.Equal("Number", document.VirtualDocs[0].Value)

		assert.Equal("PASSPORT", document.VirtualDocs[1].Type)
		assert.Equal("Number", document.VirtualDocs[1].Value)
	}
}

func TestMapCustomerToPhysicalDocs(t *testing.T) {
	testTime := common.Time(time.Date(1967, 1, 2, 0, 0, 0, 0, time.UTC))

	customer := &common.UserData{
		FirstName:            "FirstName",
		MaternalLastName:     "MaternalLastName",
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

	docs := MapCustomerToPhysicalDocs(customer)

	assert := assert.New(t)

	if assert.Len(docs, 5) {
		assert.Equal("GOVT_ID_INT", docs[0].Type)
		assert.Equal(
			"data:ContentType;base64,"+base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			docs[0].Value,
		)

		assert.Equal("GOVT_ID_INT", docs[1].Type)
		assert.Equal(
			"data:ContentType;base64,"+base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
			docs[1].Value,
		)

		assert.Equal("PROOF_OF_ADDRESS", docs[2].Type)
		assert.Equal(
			"data:ContentType;base64,"+base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
			docs[2].Value,
		)

		assert.Equal("SELFIE", docs[3].Type)
		assert.Equal(
			"data:ContentType;base64,"+base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			docs[3].Value,
		)

		assert.Equal("VIDEO_AUTHORIZATION", docs[4].Type)
		assert.Equal("data:video/mp4;base64,"+base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
			docs[4].Value)
	}
}

func Test_mapCustomerGender(t *testing.T) {
	assert.Equal(t, "M", mapCustomerGender(common.Male))
	assert.Equal(t, "F", mapCustomerGender(common.Female))
	assert.Equal(t, "O", mapCustomerGender(common.Gender(0)))

}

func Test_mapDocType(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("GOVT_ID_INT", mapDocType("IDCard"))
	assert.Equal("GOVT_ID_INT", mapDocType("DriverLicense"))
	assert.Equal("GOVT_ID_INT", mapDocType("Passport"))

	assert.Equal("SELFIE", mapDocType("Selfie"))

	assert.Equal("PROOF_OF_ADDRESS", mapDocType("UtilityBill"))
	assert.Equal("OTHER", mapDocType("ResidencePermit"))

	assert.Equal("LEGAL_AGREEMENT", mapDocType("Agreement"))
	assert.Equal("LEGAL_AGREEMENT", mapDocType("Contract"))

	assert.Equal("OTHER", mapDocType("DriversTranslation"))
}

func TestMapResponseError(t *testing.T) {
	resp := []byte(`{"error":{"en":"Invalid/expired oauth_key."},"error_code":"110","http_code":"401","success":false}`)

	code, err := MapErrorResponse(resp)

	assert.Equal(t, "110", *code)
	assert.Equal(t, "http status: 401 | error code: 110 | error: Invalid/expired oauth_key.", err.Error())

	resp = []byte("wrong response")

	code, err = MapErrorResponse(resp)

	assert.Nil(t, code)
	assert.Equal(t, "invalid character 'w' looking for beginning of value", err.Error())
}
