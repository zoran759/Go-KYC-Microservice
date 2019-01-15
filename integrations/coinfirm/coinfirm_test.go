package coinfirm

import (
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/http"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := Coinfirm{
		host:     "host",
		email:    "email",
		password: "password",
		company:  "company",
		headers: http.Headers{
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
