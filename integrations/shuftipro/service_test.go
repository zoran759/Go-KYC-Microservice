package shuftipro

import (
	"errors"
	"flag"
	"io/ioutil"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/shuftipro/verification"

	"github.com/stretchr/testify/assert"
)

var testImageUpload = flag.Bool("use-images", false, "test document images uploading")

func TestNew(t *testing.T) {

}

func TestShuftiPro_CheckCustomer(t *testing.T) {
	service := ShuftiPro{
		verification: verification.Mock{
			VerifyFn: func(request verification.OldRequest) (*verification.Response, error) {
				return &verification.Response{
					StatusCode: "SP0",
					Message:    "Not verified",
				}, nil
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Nil(t, result.Details)
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.OldRequest) (*verification.Response, error) {
			return &verification.Response{
				StatusCode: "SP1",
				Message:    "Verified",
			}, nil
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, result.Status)
		assert.Nil(t, result.Details)
	}
}

func TestShuftiPro_CheckCustomer_Error(t *testing.T) {
	service := ShuftiPro{
		verification: verification.Mock{
			VerifyFn: func(request verification.OldRequest) (*verification.Response, error) {
				return &verification.Response{
					StatusCode: "SP22",
					Message:    "Invalid checksum value.",
				}, nil
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result.Status)
		assert.Nil(t, result.Details)
		assert.Equal(t, "Invalid checksum value.", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.OldRequest) (*verification.Response, error) {
			return &verification.Response{
				StatusCode: "SP2",
			}, nil
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result.Status)
		assert.Nil(t, result.Details)
		assert.Equal(t, "There are no documents provided or they are invalid", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.OldRequest) (*verification.Response, error) {
			return nil, errors.New("test_error")
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result.Status)
		assert.Nil(t, result.Details)
		assert.Equal(t, "test_error", err.Error())
	}

	result, err = service.CheckCustomer(nil)
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result.Status)
		assert.Nil(t, result.Details)
		assert.Equal(t, "No customer supplied", err.Error())
	}
}

func TestShuftiProImageUpload(t *testing.T) {
	if !*testImageUpload {
		t.Skip("use '-use-images' flag to activate images uploading test")
	}

	testIDCard, _ := ioutil.ReadFile("../../test_data/realId.jpg")
	testSelfie, _ := ioutil.ReadFile("../../test_data/realFace.jpg")

	assert := assert.New(t)

	if !assert.NotEmpty(testIDCard, "testIDCard must contain the content of the image data file 'realId.jpg'") {
		return
	}
	if !assert.NotEmpty(testSelfie, "testSelfie must contain the content of the image data file 'realFace.jpg'") {
		return
	}

	customer := &common.UserData{
		FirstName:     "John",
		LastName:      "Doe",
		Email:         "jd@email.com",
		DateOfBirth:   common.Time(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "GB",
		Phone:         "+440000000000",
		CurrentAddress: common.Address{
			CountryAlpha2:  "GB",
			County:         "Westminster",
			Town:           "London",
			Street:         "Downing St.",
			BuildingNumber: "10",
			PostCode:       "SW1A 2AA",
		},
		IDCard: &common.IDCard{
			Number: "6980XYZ4521XYZ",
			Image: &common.DocumentFile{
				Filename:    "realId.jpg",
				ContentType: "image/jpeg",
				Data:        testIDCard,
			},
		},
		UtilityBill: &common.UtilityBill{
			Image: &common.DocumentFile{
				Filename:    "util_bill.jpg",
				ContentType: "image/jpeg",
				Data:        testIDCard,
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "realFace.png",
				ContentType: "image/jpeg",
				Data:        testSelfie,
			},
		},
	}

	config := Config{
		Host:        "https://api.shuftipro.com",
		ClientID:    "ac93f3a0fee5afa2d9399d5d0f257dc92bbde89b1e48452e1bfac3c5c1dc99db",
		SecretKey:   "lq34eOTxDe1e6G8a1P7Igqo5YK3ABCDF",
		RedirectURL: "https://api.shuftipro.com",
	}

	service := New(config)

	result, err := service.CheckCustomer(customer)

	if assert.Nil(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}
