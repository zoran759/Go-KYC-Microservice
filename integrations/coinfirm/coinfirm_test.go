package coinfirm

import (
	"net/http"
	"testing"

	"modulus/kyc/common"
	mhttp "modulus/kyc/http"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := Coinfirm{
		host:     "host",
		email:    "email",
		password: "password",
		company:  "company",
		headers: mhttp.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	tc := New(Config{
		Host:     "host",
		Email:    "email",
		Password: "password",
		Company:  "company",
	})

	assert.Equal(t, c, *tc)
}

func TestCheckCustomerNil(t *testing.T) {
	assert := assert.New(t)

	res, err := c.CheckCustomer(nil)

	assert.Error(err)
	assert.Equal("customer is absent or no data received", err.Error())
	assert.Equal(common.Error, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCheckCustomerSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, tokenResp))
	httpmock.RegisterResponder(http.MethodPut, c.host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusOK, newParticipantResp))
	httpmock.RegisterResponder(http.MethodPut, c.host+"/kyc/forms/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewBytesResponder(http.StatusCreated, nil))
	httpmock.RegisterResponder(http.MethodGet, c.host+"/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusCreated, statusLowResp))

	res, err := c.CheckCustomer(&common.UserData{})

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCheckStatusSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, tokenResp))
	httpmock.RegisterResponder(http.MethodGet, c.host+"/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusCreated, statusLowResp))

	res, err := c.CheckStatus("33611d6d-2826-4c3e-a777-3f0397e283fc")

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
