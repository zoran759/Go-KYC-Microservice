package applicants

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNewService(t *testing.T) {
	testService := service{
		host:   "test_host",
		apiKey: "api_key",
	}

	service := NewService(Config{
		Host:   "test_host",
		APIKey: "api_key",
	})

	assert.Equal(t, testService, service)
}

func Test_service_CreateApplicantSuccess(t *testing.T) {
	applicantsService := service{
		host:   "https://test_host.subsub.com",
		apiKey: "api_key",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
  "id" : "596eb3c93a0eb985b8ade34d",
  "createdAt" : "2017-07-19 03:20:09",
  "inspectionId" : "596eb3c83a0eb985b8ade349",
  "jobId" : "a8f77946-14ff-4398-aa23-a1027e16f627",
  "info" : {
    "firstName" : "Nikita",
    "middleName" : "Nikitaive",
    "lastName" : "Roman",
    "dob" : "2000-03-04",
    "placeOfBirth" : "Dnipro",
    "country" : "UA",
    "phone" : "+380-68-0055416"
  },
  "email" : "nikita@gmail.com"
}`,
			), nil
		},
	)

	response, err := applicantsService.CreateApplicant("", ApplicantInfo{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, "596eb3c93a0eb985b8ade34d", response.ID)
		assert.Equal(t, "2017-07-19 03:20:09", response.CreatedAt)
		assert.Equal(t, "596eb3c83a0eb985b8ade349", response.InspectionID)
		assert.Equal(t, "a8f77946-14ff-4398-aa23-a1027e16f627", response.JobID)
		assert.Equal(t, "nikita@gmail.com", response.Email)

		info := response.Info
		assert.Equal(t, "Serge", info.FirstName)
		assert.Equal(t, "Sergeew", info.LastName)
		assert.Equal(t, "Sergeevich", info.MiddleName)
		assert.Equal(t, "2000-03-04", info.DateOfBirth)
		assert.Equal(t, "Saint-Petersburg", info.PlaceOfBirth)
		assert.Equal(t, "UA", info.Country)
		assert.Equal(t, "+380-68-0055416", info.Phone)
	}
}

func Test_service_CreateApplicantError(t *testing.T) {
	applicantsService := service{
		host:   "https://test_host.subsub.com",
		apiKey: "api_key",
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{
    "description": "Null applicant was provided",
    "code": 400,
    "correlationId": "req-40fd1459-9674-4ad0-b7bf-cbb6447a71f8"
}`), nil
		},
	)

	response, err := applicantsService.CreateApplicant("", ApplicantInfo{})
	if assert.Error(t, err) && assert.Nil(t, response) {
		assert.Equal(t, "Null applicant was provided", err.Error())
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "{"), nil
		},
	)

	response, err = applicantsService.CreateApplicant("", ApplicantInfo{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test_host.subsub.com/resources/applicants?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = applicantsService.CreateApplicant("", ApplicantInfo{})
	if assert.Error(t, err) && assert.Nil(t, response) {
		assert.Equal(t, "Post https://test_host.subsub.com/resources/applicants?key=api_key: test_error", err.Error())
	}
}
