package documents

import (
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCommonDocumentToDocument(t *testing.T) {
	testTime := common.Time(time.Now())
	customer := common.UserData{
		FirstName:            "FirstName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
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
		Passport: &common.Passport{
			Number:        "Number",
			CountryAlpha2: "CountryAlpha2",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
	}

	documents := MapCommonCustomerDocuments(customer)

	if assert.Equal(t, 1, len(documents)) {
		commonDocument := customer.Passport
		file := documents[0].File
		metadata := documents[0].Metadata

		assert.Equal(t, commonDocument.Image.Data, file.Data)
		assert.Equal(t, commonDocument.Image.ContentType, file.ContentType)
		assert.Equal(t, commonDocument.Image.Filename, file.Filename)

		assert.Equal(t, commonDocument.Number, metadata.Number)
		assert.Equal(t, commonDocument.ValidUntil.Format("2006-01-02"), metadata.ValidUntil)
		assert.Equal(t, commonDocument.IssuedDate.Format("2006-01-02"), metadata.DateIssued)
		assert.Equal(t, commonDocument.CountryAlpha2, metadata.Country)
		assert.Empty(t, metadata.DocumentSubType)
		assert.Equal(t, "PASSPORT", metadata.DocumentType)

		assert.Equal(t, customer.FirstName, metadata.FirstName)
		assert.Equal(t, customer.MiddleName, metadata.MiddleName)
		assert.Equal(t, customer.LastName, metadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, metadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), metadata.DateOfBirth)
	}
}

func TestMapCommonDoubleSideDocument(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := common.UserData{
		FirstName:            "FirstName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
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
		DriverLicense: &common.DriverLicense{
			Number:        "Number",
			CountryAlpha2: "CountryAlpha2",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			FrontImage: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
			BackImage: &common.DocumentFile{
				Filename:    "Filename2",
				ContentType: "ContentType2",
				Data:        []byte{7, 6, 5, 4, 3, 2, 1},
			},
		},
	}

	documents := MapCommonCustomerDocuments(customer)

	if assert.Equal(t, 2, len(documents)) {
		commonDocument := customer.DriverLicense
		frontFile := documents[1].File
		frontMetadata := documents[1].Metadata

		assert.Equal(t, commonDocument.FrontImage.Data, frontFile.Data)
		assert.Equal(t, commonDocument.FrontImage.ContentType, frontFile.ContentType)
		assert.Equal(t, commonDocument.FrontImage.Filename, frontFile.Filename)

		assert.Equal(t, commonDocument.Number, frontMetadata.Number)
		assert.Equal(t, commonDocument.ValidUntil.Format("2006-01-02"), frontMetadata.ValidUntil)
		assert.Equal(t, commonDocument.IssuedDate.Format("2006-01-02"), frontMetadata.DateIssued)
		assert.Equal(t, commonDocument.CountryAlpha2, frontMetadata.Country)
		assert.Equal(t, FrontSide, frontMetadata.DocumentSubType)
		assert.Equal(t, "DRIVERS", frontMetadata.DocumentType)

		assert.Equal(t, customer.FirstName, frontMetadata.FirstName)
		assert.Equal(t, customer.MiddleName, frontMetadata.MiddleName)
		assert.Equal(t, customer.LastName, frontMetadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, frontMetadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), frontMetadata.DateOfBirth)

		backFile := documents[0].File
		backMetadata := documents[0].Metadata

		assert.Equal(t, commonDocument.BackImage.Data, backFile.Data)
		assert.Equal(t, commonDocument.BackImage.ContentType, backFile.ContentType)
		assert.Equal(t, commonDocument.BackImage.Filename, backFile.Filename)

		assert.Equal(t, commonDocument.Number, backMetadata.Number)
		assert.Equal(t, commonDocument.ValidUntil.Format("2006-01-02"), backMetadata.ValidUntil)
		assert.Equal(t, commonDocument.IssuedDate.Format("2006-01-02"), backMetadata.DateIssued)
		assert.Equal(t, commonDocument.CountryAlpha2, backMetadata.Country)
		assert.Equal(t, BackSide, backMetadata.DocumentSubType)
		assert.Equal(t, "DRIVERS", backMetadata.DocumentType)

		assert.Equal(t, customer.FirstName, backMetadata.FirstName)
		assert.Equal(t, customer.MiddleName, backMetadata.MiddleName)
		assert.Equal(t, customer.LastName, backMetadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, backMetadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), backMetadata.DateOfBirth)
	}
}

func TestMapNilCommonDocuments(t *testing.T) {
	customer := common.UserData{}

	assert.Nil(t, MapCommonCustomerDocuments(customer))
}
