package stop4

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/stop4/verification"
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	service := New(Config{
		Host: "https://private-f1649f-coreservices2.apiary-mock.com",
		MerchantID: "testmerchant123",
		Password: "testpassword123",
	})

	assert.NotNil(t, service)
}

func TestStop4_CheckCustomer(t *testing.T) {
	service := Stop4{
		verification: verification.Mock{
			VerifyFn: func(request verification.RegistrationRequest) (*verification.Response, error) {
				return &verification.Response{
					Status: int(common.Approved),
					Score:  98,
					Rec:    "Approve",
				}, nil
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, result.Status)
		assert.Nil(t, result.Details)
	}
}

func TestStop4_CheckCustomer_Error(t *testing.T) {
	service := Stop4{
		verification: verification.Mock{
			VerifyFn: func(request verification.RegistrationRequest) (*verification.Response, error) {
				return &verification.Response{
					Status: -2,
					Score:  98,
					Details: "Registration Date is not present",
					Rec:    "Denied",
				}, nil
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Equal(t, "-2", result.ErrorCode)
		assert.Equal(t, "Registration Date is not present", result.Details.Reasons[0])
	}
}

func TestStop4_CheckCustomerError(t *testing.T) {
	service := Stop4{
		verification: verification.Mock{
			VerifyFn: func(request verification.RegistrationRequest) (*verification.Response, error) {
				return nil, errors.New("test_error")
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)

	result, err = service.CheckCustomer(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)
}

