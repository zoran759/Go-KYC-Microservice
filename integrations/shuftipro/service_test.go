package shuftipro

import (
	"testing"

	"errors"
	"modulus/kyc/common"
	"modulus/kyc/integrations/shuftipro/verification"

	"github.com/stretchr/testify/assert"
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

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Nil(t, result.Details)
	}

	service.verification = verification.Mock{
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
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
			VerifyFn: func(request verification.Request) (*verification.Response, error) {
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
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
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
		VerifyFn: func(request verification.Request) (*verification.Response, error) {
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
