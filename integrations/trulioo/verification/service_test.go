package verification

import (
	"testing"

	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
)

func TestNewService(t *testing.T) {
	config := Config{
		Host:  "host",
		Token: "token",
	}

	testService := service{
		config: config,
	}

	service := NewService(config)

	assert.Equal(t, testService, service)
}

func Test_service_Verify(t *testing.T) {
	service := service{
		config: Config{
			Host: "https://api.globaldatacompany.com/verification/v1",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.globaldatacompany.com/verification/v1/verify",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
							"TransactionID":"2b780267-9c5f-474f-9442-9449fd2d2eaa",
							"UploadedDt":"2018-08-20T07:11:27",
							"CountryCode":"US",
							"ProductName":"Identity Verification",
							"Record":{
								"TransactionRecordID":"02b39dac-55f2-4019-8cac-5de931669191",
								"RecordStatus":"match",
								"DatasourceResults":[{
									"DatasourceName":"Citizen",
									"DatasourceFields":[{"FieldName":"YearOfBirth","Status":"match"},{"FieldName":"MiddleName","Status":"missing"}],
									"Errors": [{"Code": "400", "Message": "Test error"}]
								},
								{
									"DatasourceName":"Credit Agency 2",
									"DatasourceFields":[{"FieldName":"FirstInitial","Status":"match"},{"FieldName":"socialservice","Status":"missing"}]
								}]
							},
							"Errors": [{"Code": "400", "Message": "Test error"}]
						}`), nil
		},
	)

	response, err := service.Verify("US", configuration.Consents{}, DataFields{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, "US", response.CountryCode)
		assert.Equal(t, "02b39dac-55f2-4019-8cac-5de931669191", response.Record.TransactionRecordID)
		assert.Equal(t, "match", response.Record.RecordStatus)
		assert.Len(t, response.Record.DatasourceResults, 2)
		assert.Equal(t, "Citizen", response.Record.DatasourceResults[0].DatasourceName)
		assert.Equal(t, []DatasourceField{
			{
				FieldName: "YearOfBirth",
				Status:    "match",
			},
			{
				FieldName: "MiddleName",
				Status:    "missing",
			},
		}, response.Record.DatasourceResults[0].DatasourceFields)
		assert.Equal(t, Errors{
			{
				Code:    "400",
				Message: "Test error",
			},
		}, response.Record.DatasourceResults[0].Errors)
		assert.Equal(t, "Credit Agency 2", response.Record.DatasourceResults[1].DatasourceName)
		assert.Equal(t, []DatasourceField{
			{
				FieldName: "FirstInitial",
				Status:    "match",
			},
			{
				FieldName: "socialservice",
				Status:    "missing",
			},
		}, response.Record.DatasourceResults[1].DatasourceFields)

		assert.Len(t, response.Errors, 1)
		assert.Equal(t, Error{
			Code:    "400",
			Message: "Test error",
		}, response.Errors[0])
	}
}

func Test_service_Verify_Error(t *testing.T) {
	service := service{
		config: Config{
			Host: "https://api.globaldatacompany.com/verification/v1",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.globaldatacompany.com/verification/v1/verify",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`), nil
		},
	)

	response, err := service.Verify("", configuration.Consents{}, DataFields{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.globaldatacompany.com/verification/v1/verify",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = service.Verify("", configuration.Consents{}, DataFields{})
	assert.Error(t, err)
	assert.Nil(t, response)
}
