package documents

import (
	"testing"

	"gitlab.com/modulusglobal/kyc/common"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestMapCommonDocumentToDocument(t *testing.T) {
	testTime := common.Time(time.Now())
	customer := common.UserData{
		FirstName:      "FirstName",
		LastName:       "LastName",
		MiddleName:     "MiddleName",
		LegalName:      "LegalName",
		Email:          "Email",
		Gender:         "Gender",
		DateOfBirth:    testTime,
		PlaceOfBirth:   "PlaceOfBirth",
		CountryOfBirth: "CountryOfBirth",
		StateOfBirth:   "StateOfBirth",
		CountryAlpha2:  "CountryAlpha2",
		Nationality:    "Nationality",
		Phone:          "Phone",
		MobilePhone:    "MobilePhone",
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:       "Type",
					Country:    "CountryAlpha2",
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

	documents := MapCommonCustomerDocuments(customer)

	if assert.Equal(t, 1, len(documents)) {
		commonDocument := customer.Documents[0]
		file := documents[0].File
		metadata := documents[0].Metadata

		assert.Equal(t, commonDocument.Front.Data, file.Data)
		assert.Equal(t, commonDocument.Front.ContentType, file.ContentType)
		assert.Equal(t, commonDocument.Front.Filename, file.Filename)

		assert.Equal(t, commonDocument.Metadata.Number, *metadata.Number)
		assert.Equal(t, commonDocument.Metadata.ValidUntil.Format("2006-01-02"), *metadata.ValidUntil)
		assert.Equal(t, commonDocument.Metadata.DateIssued.Format("2006-01-02"), *metadata.DateIssued)
		assert.Equal(t, commonDocument.Metadata.Country, metadata.Country)
		assert.Nil(t, metadata.DocumentSubType)
		assert.Equal(t, commonDocument.Metadata.Type, metadata.DocumentType)

		assert.Equal(t, customer.FirstName, *metadata.FirstName)
		assert.Equal(t, customer.MiddleName, *metadata.MiddleName)
		assert.Equal(t, customer.LastName, *metadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, *metadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), *metadata.DateOfBirth)
	}
}

func TestMapCommonDoubleSideDocument(t *testing.T) {
	testTime := common.Time(time.Now())

	customer := common.UserData{
		FirstName:      "FirstName",
		LastName:       "LastName",
		MiddleName:     "MiddleName",
		LegalName:      "LegalName",
		Email:          "Email",
		Gender:         "Gender",
		DateOfBirth:    testTime,
		PlaceOfBirth:   "PlaceOfBirth",
		CountryOfBirth: "CountryOfBirth",
		StateOfBirth:   "StateOfBirth",
		CountryAlpha2:  "CountryAlpha2",
		Nationality:    "Nationality",
		Phone:          "Phone",
		MobilePhone:    "MobilePhone",
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:       "Type",
					Country:    "CountryAlpha2",
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

	documents := MapCommonCustomerDocuments(customer)

	if assert.Equal(t, 2, len(documents)) {
		commonDocument := customer.Documents[0]
		frontFile := documents[1].File
		frontMetadata := documents[1].Metadata

		assert.Equal(t, commonDocument.Front.Data, frontFile.Data)
		assert.Equal(t, commonDocument.Front.ContentType, frontFile.ContentType)
		assert.Equal(t, commonDocument.Front.Filename, frontFile.Filename)

		assert.Equal(t, commonDocument.Metadata.Number, *frontMetadata.Number)
		assert.Equal(t, commonDocument.Metadata.ValidUntil.Format("2006-01-02"), *frontMetadata.ValidUntil)
		assert.Equal(t, commonDocument.Metadata.DateIssued.Format("2006-01-02"), *frontMetadata.DateIssued)
		assert.Equal(t, commonDocument.Metadata.Country, frontMetadata.Country)
		assert.Equal(t, FrontSide, *frontMetadata.DocumentSubType)
		assert.Equal(t, commonDocument.Metadata.Type, frontMetadata.DocumentType)

		assert.Equal(t, customer.FirstName, *frontMetadata.FirstName)
		assert.Equal(t, customer.MiddleName, *frontMetadata.MiddleName)
		assert.Equal(t, customer.LastName, *frontMetadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, *frontMetadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), *frontMetadata.DateOfBirth)

		backFile := documents[0].File
		backMetadata := documents[0].Metadata

		assert.Equal(t, commonDocument.Back.Data, backFile.Data)
		assert.Equal(t, commonDocument.Back.ContentType, backFile.ContentType)
		assert.Equal(t, commonDocument.Back.Filename, backFile.Filename)

		assert.Equal(t, commonDocument.Metadata.Number, *backMetadata.Number)
		assert.Equal(t, commonDocument.Metadata.ValidUntil.Format("2006-01-02"), *backMetadata.ValidUntil)
		assert.Equal(t, commonDocument.Metadata.DateIssued.Format("2006-01-02"), *backMetadata.DateIssued)
		assert.Equal(t, commonDocument.Metadata.Country, backMetadata.Country)
		assert.Equal(t, BackSide, *backMetadata.DocumentSubType)
		assert.Equal(t, commonDocument.Metadata.Type, backMetadata.DocumentType)

		assert.Equal(t, customer.FirstName, *backMetadata.FirstName)
		assert.Equal(t, customer.MiddleName, *backMetadata.MiddleName)
		assert.Equal(t, customer.LastName, *backMetadata.LastName)
		assert.Equal(t, customer.PlaceOfBirth, *backMetadata.PlaceOfBirth)
		assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), *backMetadata.DateOfBirth)
	}
}

func TestMapNilCommonDocuments(t *testing.T) {
	customer := common.UserData{}

	assert.Nil(t, MapCommonCustomerDocuments(customer))
}
