package verification

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

func Test_service_StartVerification(t *testing.T) {
	service := NewService(Config{
		Host:   "https://test-api.sumsub.com",
		APIKey: "api_key",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status/pending?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"ok": 1}`,
			), nil
		},
	)

	response, err := service.StartVerification("test_applicant_id")
	if assert.NoError(t, err) {
		assert.True(t, response)
	}
}

func Test_service_StartVerificationError(t *testing.T) {
	service := NewService(Config{
		Host:   "https://test-api.sumsub.com",
		APIKey: "api_key",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status/pending?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"description": "Invalid id '5b7298530975a1df03bdd13'", "code": 400}`,
			), nil
		},
	)

	response, err := service.StartVerification("test_applicant_id")
	if assert.Error(t, err) {
		assert.False(t, response)
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status/pending?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"`,
			), nil
		},
	)

	response, err = service.StartVerification("test_applicant_id")
	if assert.Error(t, err) {
		assert.False(t, response)
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status/pending?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = service.StartVerification("test_applicant_id")
	if assert.Error(t, err) {
		assert.False(t, response)
	}
}

func Test_service_CheckApplicantStatus(t *testing.T) {
	service := NewService(Config{
		Host:   "https://test-api.sumsub.com",
		APIKey: "api_key",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
    "id": "5b7298530a975a1df03bdd17",
    "inspectionId": "5b7298530a975a1df03bdd14",
    "jobId": "81c0b38e-904b-4d55-bd7f-870952feb6d2",
    "createDate": "2018-08-14 10:21:33+0000",
    "reviewDate": "2018-08-15 05:23:47+0000",
    "reviewResult": {
        "reviewAnswer": "RED",
        "label": "OTHER",
        "rejectLabels": [
            "ID_INVALID"
        ],
        "reviewRejectType": "RETRY"
    },
    "reviewStatus": "completed",
    "notificationFailureCnt": 0,
    "applicantId": "5b7298530a975a1df03bdd13"
}`,
			), nil
		},
	)

	status, result, err := service.CheckApplicantStatus("test_applicant_id")
	if assert.NoError(t, err) && assert.Equal(t, "completed", status) && assert.NotNil(t, result) {
		assert.Equal(t, "RED", result.ReviewAnswer)
		assert.Equal(t, "OTHER", result.Label)
		assert.Equal(t, []string{"ID_INVALID"}, result.RejectLabels)
		assert.Equal(t, "RETRY", result.ReviewRejectType)
	}
}

func Test_service_CheckApplicantStatusError(t *testing.T) {
	service := NewService(Config{
		Host:   "https://test-api.sumsub.com",
		APIKey: "api_key",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"description": "Invalid id '5b7298530975a1df03bdd13'", "code": 400}`,
			), nil
		},
	)

	_, response, err := service.CheckApplicantStatus("test_applicant_id")
	if assert.Error(t, err) {
		assert.Nil(t, response)
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{"`,
			), nil
		},
	)

	_, response, err = service.CheckApplicantStatus("test_applicant_id")
	if assert.Error(t, err) {
		assert.Nil(t, response)
	}

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://test-api.sumsub.com/resources/applicants/test_applicant_id/status?key=api_key",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	_, response, err = service.CheckApplicantStatus("test_applicant_id")
	if assert.Error(t, err) {
		assert.Nil(t, response)
	}
}
