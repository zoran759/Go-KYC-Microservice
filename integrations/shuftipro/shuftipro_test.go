package shuftipro

import (
	"errors"
	"flag"
	"io/ioutil"
	stdhttp "net/http"
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/jarcoal/httpmock.v1"
)

var demo = flag.Bool("demo", false, "Test integration with Shufti Pro API using the demo")

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
		Host:        "https://shuftipro.com/api/",
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
			name: "Nil customer",
			err:  errors.New("No customer supplied"),
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

func TestShuftiProIntegration(t *testing.T) {
	if !*demo {
		t.Skip("Use -demo command line flag to test the demo")

	}

	require := require.New(t)

	realFace, err := ioutil.ReadFile("../../test_data/real-face.jpg")
	require.NoError(err)
	fakeFace, err := ioutil.ReadFile("../../test_data/fake-face.jpg")
	require.NoError(err)

	require.NotEmpty(realFace)
	require.NotEmpty(fakeFace)

	s := New(Config{
		Host:        "https://shuftipro.com/api/",
		ClientID:    "d76612f86d26846065f5c37dfef7a7dd04eaa724f923773fe02f9d8b0bec0877",
		SecretKey:   "5gCvERlJy4w6Lf1Bcn6ztQSKP0Lrqxhp",
		CallbackURL: "https://shuftipro.com/api",
	})

	type testCase struct {
		name     string
		customer *common.UserData
		result   common.KYCResult
		err      error
	}

	testCases := []testCase{
		testCase{
			name: "Approved result",
			customer: &common.UserData{
				FirstName:     "John",
				LastName:      "Livone",
				DateOfBirth:   common.Time(time.Date(1989, 9, 6, 0, 0, 0, 0, time.UTC)),
				CountryAlpha2: "GB",
				Email:         "john.livone@example.com",
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "real-face.jpg",
						ContentType: "image/jpeg",
						Data:        realFace,
					},
				},
			},
			result: common.KYCResult{
				Status: common.Approved,
			},
		},
		testCase{
			name: "Denied result",
			customer: &common.UserData{
				FirstName:     "John",
				LastName:      "Doe",
				DateOfBirth:   common.Time(time.Date(1989, 9, 6, 0, 0, 0, 0, time.UTC)),
				CountryAlpha2: "GB",
				Email:         "john.doe@example.com",
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "fake-face.jpg",
						ContentType: "image/jpeg",
						Data:        fakeFace,
					},
				},
			},
			result: common.KYCResult{
				Status: common.Denied,
				Details: &common.KYCDetails{
					Reasons: []string{"Face is not verified."},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
