package shuftipro

import (
	"testing"

	"errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/shuftipro/verification"
)

func TestNew(t *testing.T) {

}

func TestShuftiPro_CheckCustomer(t *testing.T) {
	service := ShuftiPro{
		verification: verification.Mock{
			VerifyFn: func(request verification.Request) (*verification.Response, error) {
				return &verification.Response{
					StatusCode: "SP0",
					Message:    "Not verified",
				}, nil
			},
		},
	}

	result, detailedResult, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result)
		assert.Nil(t, detailedResult)
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
			return &verification.Response{
				StatusCode: "SP1",
				Message:    "Verified",
			}, nil
		},
	}

	result, detailedResult, err = service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, result)
		assert.Nil(t, detailedResult)
	}
}

func TestShuftiPro_CheckCustomer_Error(t *testing.T) {
	service := ShuftiPro{
		verification: verification.Mock{
			VerifyFn: func(request verification.Request) (*verification.Response, error) {
				return &verification.Response{
					StatusCode: "SP22",
					Message:    "Invalid checksum value.",
				}, nil
			},
		},
	}

	result, detailedResult, err := service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result)
		assert.Nil(t, detailedResult)
		assert.Equal(t, "Invalid checksum value.", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
			return &verification.Response{
				StatusCode: "SP2",
			}, nil
		},
	}

	result, detailedResult, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result)
		assert.Nil(t, detailedResult)
		assert.Equal(t, "There are no documents provided or they are invalid", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
			return nil, errors.New("test_error")
		},
	}

	result, detailedResult, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result)
		assert.Nil(t, detailedResult)
		assert.Equal(t, "test_error", err.Error())
	}

	result, detailedResult, err = service.CheckCustomer(nil)
	if assert.Error(t, err) {
		assert.Equal(t, common.Error, result)
		assert.Nil(t, detailedResult)
		assert.Equal(t, "No customer supplied", err.Error())
	}
}
