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
		Host:        "host.com",
		SecretKey:   "SecretKey",
		ClientID:    "ClientID",
		RedirectURL: "redirect.com",
	}

	testService := service{
		config: config,
	}

	service := NewService(config)

	assert.Equal(t, testService, service)
}

func Test_service_Verify(t *testing.T) {
	config := Config{
		Host:        "https://api.shuftipro.com/",
		SecretKey:   "SecretKey",
		ClientID:    "ClientID",
		RedirectURL: "redirect.com",
	}

	testService := service{
		config: config,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	reference := ""
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.shuftipro.com/",
		func(request *http.Request) (*http.Response, error) {
			reference = request.Form.Get("reference")
			return httpmock.NewStringResponse(
				http.StatusOK,
				fmt.Sprintf(
					`{"status_code": "SP1", "message": "Verified", "reference": "%s", "signature": "sig"}`,
					request.Form.Get("reference")),
			), nil
		},
	)

	response, err := testService.Verify(OldRequest{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, reference, response.Reference)
		assert.Equal(t, "SP1", response.StatusCode)
		assert.Equal(t, "Verified", response.Message)
		assert.Equal(t, "sig", response.Signature)
	}
}

func Test_service_Verify_Error(t *testing.T) {
	config := Config{
		Host:        "https://api.shuftipro.com/",
		SecretKey:   "SecretKey",
		ClientID:    "ClientID",
		RedirectURL: "redirect.com",
	}

	testService := service{
		config: config,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.shuftipro.com/",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`,
			), nil
		},
	)

	response, err := testService.Verify(OldRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.shuftipro.com/",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = testService.Verify(OldRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)
}
