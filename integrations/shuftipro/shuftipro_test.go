package shuftipro

import (
	"errors"
	stdhttp "net/http"
	"testing"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNew(t *testing.T) {
	config := Config{
		Host:        "host",
		ClientID:    "client_id",
		SecretKey:   "secret_key",
		CallbackURL: "callback_url",
	}

	sh1 := ShuftiPro{
		client: NewClient(config),
	}

	sh2 := New(config)

	assert.Equal(t, sh1, sh2)
}

func TestShuftiProCheckCustomer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	s := New(Config{
		Host:        "https://shuftipro.com/api",
		ClientID:    "client_id",
		SecretKey:   "secret_key",
		CallbackURL: "callback_url",
	})

	type testCase struct {
		name      string
		customer  *common.UserData
		responder httpmock.Responder
		result    common.KYCResult
		err       error
	}

	testCases := []testCase{
		testCase{
			name:     "Nil customer",
			customer: nil,
			result:   common.KYCResult{},
			err:      errors.New("No customer supplied"),
		},
		testCase{
			name: "Approved result",
			customer: &common.UserData{
				FirstName: "John",
				LastName:  "Doe",
			},
			responder: httpmock.NewStringResponder(stdhttp.StatusOK, reqAcceptedResponse),
			result: common.KYCResult{
				Status: common.Approved,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.RegisterResponder(stdhttp.MethodPost, s.client.host, tc.responder)
			res, err := s.CheckCustomer(tc.customer)
			assert := assert.New(t)
			assert.Equal(tc.result, res)
			if tc.err != nil {
				assert.Equal(tc.err.Error(), err.Error())
			} else {
				assert.Equal(tc.err, err)
			}
		})
	}
}

func TestShuftiProCheckStatus(t *testing.T) {
	result := common.KYCResult{}
	res, err := ShuftiPro{}.CheckStatus("")
	assert.Equal(t, result, res)
	assert.Equal(t, "Shufti Pro doesn't support a verification status check", err.Error())
}
