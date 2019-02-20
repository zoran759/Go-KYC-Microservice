package coinfirm

import (
	"encoding/base64"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm/model"

	"github.com/stretchr/testify/assert"
)

func TestPrepareCustomerData(t *testing.T) {
	assert := assert.New(t)

	customer := &common.UserData{
		IPaddress:     "192.168.0.137",
		FirstName:     "John",
		LastName:      "Doe",
		MiddleName:    "James",
		Email:         "john.doe@mail.com",
		Nationality:   "US",
		Phone:         "+1 502 224 6786",
		MobilePhone:   "+912111222333",
		CountryAlpha2: "US",
		DateOfBirth:   common.Time(time.Date(1975, 9, 21, 0, 0, 0, 0, time.UTC)),
		CurrentAddress: common.Address{
			CountryAlpha2:     "US",
			County:            "Neverhood",
			State:             "Pennsylvania",
			Town:              "Pittsburgh",
			Suburb:            "Cubby",
			Street:            "Gifford St",
			StreetType:        "Street",
			SubStreet:         "None",
			BuildingName:      "Home",
			BuildingNumber:    "1324",
			FlatNumber:        "1",
			PostCode:          "15212",
			StateProvinceCode: "PA",
		},
		Passport: &common.Passport{
			Number:        "987654321",
			CountryAlpha2: "US",
			IssuedDate:    common.Time(time.Date(2015, 06, 15, 0, 0, 0, 0, time.UTC)),
			ValidUntil:    common.Time(time.Date(2025, 06, 14, 0, 0, 0, 0, time.UTC)),
			Image: &common.DocumentFile{
				Filename:    "passport.jpg",
				ContentType: "image/jpeg",
				Data:        []byte(`Fake passport image data`),
			},
		},
	}

	details, docfile := prepareCustomerData(customer)

	assert.Equal(customer.IPaddress, details.UserIP)
	assert.Equal(model.Individual, details.Type)
	assert.Equal(customer.FirstName, details.FirstName)
	assert.Equal(customer.LastName, details.LastName)
	assert.Equal(customer.MiddleName, details.MiddleName)
	assert.Equal(customer.Email, details.Email)
	assert.Equal(common.CountryAlpha2ToAlpha3[customer.Nationality], details.Nationality)
	assert.Equal(customer.Phone, details.Phone)
	assert.Equal(common.CountryAlpha2ToAlpha3[customer.CountryAlpha2], details.CountryAlpha3)
	assert.Equal(customer.CurrentAddress.PostCode, details.Postcode)
	assert.Equal(customer.CurrentAddress.Town, details.City)
	assert.Equal(customer.CurrentAddress.Street, details.Street)
	assert.Equal("1975-09-21", details.BirthDate)
	assert.Equal(customer.Passport.Number, details.IDNumber)
	assert.NotNil(docfile)

	customer.Phone = ""

	details, _ = prepareCustomerData(customer)

	assert.Equal(customer.MobilePhone, details.Phone)
}

func TestPrepareCustomerDocument(t *testing.T) {
	assert := assert.New(t)

	passport := &common.Passport{
		Number:        "987654321",
		CountryAlpha2: "US",
		IssuedDate:    common.Time(time.Date(2015, 06, 15, 0, 0, 0, 0, time.UTC)),
		ValidUntil:    common.Time(time.Date(2025, 06, 14, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "passport.jpg",
			ContentType: "image/jpeg",
			Data:        []byte(`Fake passport image data`),
		},
	}
	drivers := &common.DriverLicense{
		Number:        "210901975",
		CountryAlpha2: "RU",
		IssuedDate:    common.Time(time.Date(2010, 10, 7, 0, 0, 0, 0, time.UTC)),
		ValidUntil:    common.Time(time.Date(2020, 10, 6, 0, 0, 0, 0, time.UTC)),
		FrontImage: &common.DocumentFile{
			Filename:    "drivers_front.pdf",
			ContentType: "application/pdf",
			Data:        []byte(`Smile, - it is a fake drivers front image data`),
		},
		BackImage: &common.DocumentFile{
			Filename:    "drivers_back.jpg",
			ContentType: "image/jpeg",
			Data:        []byte(`Smile, - it is a fake drivers back image data`),
		},
	}
	drivertrans := &common.DriverLicenseTranslation{
		Number:        "210901975",
		CountryAlpha2: "RU",
		IssuedDate:    common.Time(time.Date(2010, 10, 7, 0, 0, 0, 0, time.UTC)),
		ValidUntil:    common.Time(time.Date(2020, 10, 6, 0, 0, 0, 0, time.UTC)),
		FrontImage: &common.DocumentFile{
			Filename:    "drivers_front.psd",
			ContentType: "image/x-photoshop",
			Data:        []byte(`Smile, - it is a fake drivers front image data`),
		},
		BackImage: &common.DocumentFile{
			Filename:    "drivers_back.jpg",
			ContentType: "image/jpeg",
			Data:        []byte(`Smile, - it is a fake drivers back image data`),
		},
	}
	idcard := &common.IDCard{
		Number:        "159133253",
		CountryAlpha2: "US",
		IssuedDate:    common.Time(time.Date(1960, 06, 23, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "ssn.jpg",
			ContentType: "image/jpeg",
			Data:        []byte(`Fake SSN image data`),
		},
	}
	snils := &common.SNILS{
		Number: "Number",
		Image: &common.DocumentFile{
			Filename:    "snils.bmp",
			ContentType: "image/x-ms-bmp",
			Data:        []byte{1, 2, 3, 4, 5, 6, 7},
		},
	}
	utilityBill := &common.UtilityBill{
		CountryAlpha2: "US",
		Image: &common.DocumentFile{
			Filename:    "ub.png",
			ContentType: "image/png",
			Data:        []byte(`Fake utility bill permit image data`),
		},
	}

	customer := &common.UserData{
		Passport: passport,
	}

	docnum, docfile := prepareCustomerDocument(customer)

	assert.Equal(customer.Passport.Number, docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileID, docfile.Type)
	assert.Equal("jpg", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.Passport.Image.Data), docfile.DataBase64)

	customer = &common.UserData{
		DriverLicense: drivers,
	}

	docnum, docfile = prepareCustomerDocument(customer)

	assert.Equal(customer.DriverLicense.Number, docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileID, docfile.Type)
	assert.Equal("pdf", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data), docfile.DataBase64)

	customer = &common.UserData{
		DriverLicenseTranslation: drivertrans,
	}

	docnum, docfile = prepareCustomerDocument(customer)

	assert.Equal(customer.DriverLicenseTranslation.Number, docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileID, docfile.Type)
	assert.Equal("psd", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.DriverLicenseTranslation.FrontImage.Data), docfile.DataBase64)

	customer = &common.UserData{
		IDCard: idcard,
	}

	docnum, docfile = prepareCustomerDocument(customer)

	assert.Equal(customer.IDCard.Number, docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileID, docfile.Type)
	assert.Equal("jpg", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data), docfile.DataBase64)

	customer = &common.UserData{
		SNILS: snils,
	}

	docnum, docfile = prepareCustomerDocument(customer)

	assert.Equal(customer.SNILS.Number, docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileID, docfile.Type)
	assert.Equal("bmp", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data), docfile.DataBase64)

	customer = &common.UserData{
		UtilityBill: utilityBill,
	}

	docnum, docfile = prepareCustomerDocument(customer)

	assert.Empty(docnum)
	assert.NotNil(docfile)
	assert.Equal(model.FileAddress, docfile.Type)
	assert.Equal("png", docfile.Extension)
	assert.Equal(base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data), docfile.DataBase64)
}

func TestToResult(t *testing.T) {
	assert := assert.New(t)

	pID := "test_ref_id"
	resp := model.StatusResponse{
		CurrentStatus: model.Empty,
	}

	res, err := toResult(pID, resp)

	assert.Error(err)
	assert.Equal("unexpected status value: empty", err.Error())
	assert.Equal(common.Error, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

	resp.CurrentStatus = model.New
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Unclear, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.NotNil(res.StatusCheck)
	assert.Equal(common.Coinfirm, res.StatusCheck.Provider)
	assert.Equal(pID, res.StatusCheck.ReferenceID)
	assert.NotZero(res.StatusCheck.LastCheck)

	resp.CurrentStatus = model.InProgress
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Unclear, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.NotNil(res.StatusCheck)
	assert.Equal(common.Coinfirm, res.StatusCheck.Provider)
	assert.Equal(pID, res.StatusCheck.ReferenceID)
	assert.NotZero(res.StatusCheck.LastCheck)

	resp.CurrentStatus = model.Incomplete
	res, err = toResult(pID, resp)

	assert.Error(err)
	assert.Equal("data provided by participant is incomplete or does not meet the requirements set in KYC form", err.Error())
	assert.Equal(common.Error, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

	resp.CurrentStatus = model.Low
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

	resp.CurrentStatus = model.Medium
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

	resp.CurrentStatus = model.High
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 1)
	assert.Equal("Coinfirm analysts evaluated the risk associated to participant as high", res.Details.Reasons[0])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)

	resp.CurrentStatus = model.Fail
	res, err = toResult(pID, resp)

	assert.NoError(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 1)
	assert.Equal("Coinfirm analysts evaluated the risk associated to participant as unacceptable", res.Details.Reasons[0])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCommonDocTypeToFileTypeUnspecified(t *testing.T) {
	filetype := commonDocTypeToFileType(unspecified)

	assert.Empty(t, filetype)
}

func TestExtFromContentType(t *testing.T) {
	assert := assert.New(t)

	ext := extFromContentType("")

	assert.Empty(ext)

	ext = extFromContentType("<error>")

	assert.Empty(ext)
}
