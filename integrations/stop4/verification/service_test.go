package verification

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNewService(t *testing.T) {
	config := Config{
		Host:       "https://api.verifyglobalrisk.com",
		MerchantID: "INT-123456",
		Password:   "123456",
	}

	testService := service{
		config: config,
	}

	service := NewService(config)

	assert.Equal(t, testService, service)
}

func Test_service_Verify(t *testing.T) {
	config := Config{
		Host:       "https://api.verifyglobalrisk.com/",
		MerchantID: "INT-123456",
		Password:   "123456",
	}

	testService := service{
		config: config,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.verifyglobalrisk.com/",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				fmt.Sprintf(`{"status": 0, "id": "223", "score": 98, "rec":"Approve"}`),
			), nil
		},
	)

	response, err := testService.Verify(Request{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, 0, response.Status)
		assert.Equal(t, 98, response.Score)
		assert.Equal(t, "Approve", response.Rec)
	}
}

func Test_service_Verify_Error(t *testing.T) {
	config := Config{
		Host:       "https://api.verifyglobalrisk.com/",
		MerchantID: "INT-123456",
		Password:   "123456",
	}

	testService := service{
		config: config,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.verifyglobalrisk.com/",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`,
			), nil
		},
	)

	response, err := testService.Verify(Request{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.verifyglobalrisk.com/",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = testService.Verify(Request{})
	assert.Error(t, err)
	assert.Nil(t, response)
}
