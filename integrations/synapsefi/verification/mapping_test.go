package verification

import (
	"testing"

	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
	"time"
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

	userRequest := MapCustomerToCreateUserRequest(customer)

	assert.Len(t, userRequest.PhoneNumbers, 2)
	assert.Equal(t, "Phone", userRequest.PhoneNumbers[0])
	assert.Equal(t, "MobilePhone", userRequest.PhoneNumbers[1])

	assert.Len(t, userRequest.Logins, 1)
	assert.Equal(t, Login{
		Email: "Email",
		Scope: "READ",
	}, userRequest.Logins[0])

	assert.Len(t, userRequest.LegalNames, 1)
	assert.Equal(t, "FirstName MiddleName LastName", userRequest.LegalNames[0])

	assert.Len(t, userRequest.Documents, 1)

	document := userRequest.Documents[0]

	assert.Equal(t, "FirstName MiddleName LastName", document.OwnerName)
	assert.Equal(t, "Email", document.Email)
	assert.Equal(t, "Phone", document.PhoneNumber)
	assert.Equal(t, "0.0.0.0", document.IPAddress)
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

	if assert.Len(t, document.PhysicalDocs, 3) {
		assert.Equal(t,
			"GOVT_ID_INT",
			document.PhysicalDocs[0].DocumentType,
		)
		assert.Equal(t,
			base64.StdEncoding.EncodeToString(customer.Documents[0].Front.Data),
			document.PhysicalDocs[0].DocumentValue,
		)

		assert.Equal(t,
			"SELFIE",
			document.PhysicalDocs[1].DocumentType,
		)
		assert.Equal(t,
			base64.StdEncoding.EncodeToString(customer.Documents[1].Front.Data),
			document.PhysicalDocs[1].DocumentValue,
		)

		assert.Equal(t,
			"GOVT_ID_INT",
			document.PhysicalDocs[2].DocumentType,
		)
		assert.Equal(t,
			base64.StdEncoding.EncodeToString(customer.Documents[2].Front.Data),
			document.PhysicalDocs[2].DocumentValue,
		)
	}
}

func Test_mapCustomerGender(t *testing.T) {
	assert.Equal(t, "M", mapCustomerGender(common.Male))
	assert.Equal(t, "F", mapCustomerGender(common.Female))
	assert.Equal(t, "O", mapCustomerGender(common.Gender(0)))

}

func Test_mapDocumentType(t *testing.T) {
	assert.Equal(t, "GOVT_ID", mapDocumentType(common.IDCardEng))
	assert.Equal(t, "GOVT_ID", mapDocumentType(common.DriversEng))
	assert.Equal(t, "GOVT_ID", mapDocumentType(common.PassportEng))

	assert.Equal(t, "GOVT_ID_INT", mapDocumentType(common.IDCard))
	assert.Equal(t, "GOVT_ID_INT", mapDocumentType(common.Drivers))
	assert.Equal(t, "GOVT_ID_INT", mapDocumentType(common.Passport))

	assert.Equal(t, "SELFIE", mapDocumentType(common.Selfie))

	assert.Equal(t, "PROOF_OF_ADDRESS", mapDocumentType(common.UtilityBill))
	assert.Equal(t, "PROOF_OF_ADDRESS", mapDocumentType(common.ResidencePermit))

	assert.Equal(t, "PROOF_OF_INCOME", mapDocumentType(common.ProofOfIncome))

	assert.Equal(t, "PROOF_OF_ACCOUNT", mapDocumentType(common.ProofOfAccount))

	assert.Equal(t, "AUTHORIZATION", mapDocumentType(common.ACHAuthorization))

	assert.Equal(t, "BK_CHECK", mapDocumentType(common.BackgroundCheck))

	assert.Equal(t, "SSN_CARD", mapDocumentType(common.SSN))

	assert.Equal(t, "EIN_DOC", mapDocumentType(common.EINDocument))

	assert.Equal(t, "W9_DOC", mapDocumentType(common.W9Document))

	assert.Equal(t, "W8_DOC", mapDocumentType(common.W8Document))

	assert.Equal(t, "W2_DOC", mapDocumentType(common.W2Document))

	assert.Equal(t, "VOIDED_CHECK", mapDocumentType(common.VoidedCheck))

	assert.Equal(t, "AOI", mapDocumentType(common.ArticlesOfIncorporation))

	assert.Equal(t, "BYLAWS_DOC", mapDocumentType(common.BylawsDocument))

	assert.Equal(t, "LOE", mapDocumentType(common.LetterOfEngagement))

	assert.Equal(t, "CIP_DOC", mapDocumentType(common.CIPDoc))

	assert.Equal(t, "SUBSCRIPTION_AGREEMENT", mapDocumentType(common.SubscriptionAgreement))

	assert.Equal(t, "PROMISSORY_NOTE", mapDocumentType(common.PromissoryNote))

	assert.Equal(t, "LEGAL_AGREEMENT", mapDocumentType(common.Agreement))
	assert.Equal(t, "LEGAL_AGREEMENT", mapDocumentType(common.Contract))

	assert.Equal(t, "REG_GG", mapDocumentType(common.RegGG))

	assert.Equal(t, "DBA_DOC", mapDocumentType(common.DBADoc))

	assert.Equal(t, "DEPOSIT_AGGREEMENT", mapDocumentType(common.DepositAgreement))

	assert.Equal(t, "OTHER", mapDocumentType(common.Other))
	assert.Equal(t, "OTHER", mapDocumentType(common.DriversTranslation))
}
