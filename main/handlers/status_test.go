package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
)

var cfg = config.Config{
	string(common.IdentityMind): {
		"Host":     "https://sandbox.identitymind.com/im",
		"Username": "fakeuser",
		"Password": "fakepassword",
	},
	string(common.IDology): {
		"Host":             "https://web.idologylive.com/api/idiq.svc",
		"Username":         "fakeuser",
		"Password":         "fakepassword",
		"UseSummaryResult": "false",
	},
	string(common.ShuftiPro): {
		"Host":        "https://api.shuftipro.com",
		"ClientID":    "fakeID",
		"SecretKey":   "fakeKey",
		"RedirectURL": "https://api.shuftipro.com",
	},
	string(common.SumSub): {
		"Host":   "https://test-api.sumsub.com",
		"APIKey": "fakeKey",
	},
	string(common.Trulioo): {
		"Host":         "https://api.globaldatacompany.com",
		"NAPILogin":    "fakelogin",
		"NAPIPassword": "fakepassword",
	},
}

var response = []byte(`
{
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
    "applicantId": "testID"
}`)

var errorResponse = []byte(`
{
	"code": 401,
	"description": "Access denied"
}`)

func init() {
	if config.Cfg == nil {
		config.Cfg = cfg
	}
}

type FailedReader struct{}

func (r FailedReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Read failed")
}

func TestCheckStatus(t *testing.T) {
	assert := assert.New(t)

	cfg := config.Cfg[string(common.SumSub)]

	assert.NotNil(cfg)

	referenceID := "testID"

	request, err := json.Marshal(&common.CheckStatusRequest{
		Provider:    common.SumSub,
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)
	assert.NotEmpty(response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/resources/applicants/%s/status?key=%s", cfg["Host"], referenceID, cfg["APIKey"]),
		httpmock.NewBytesResponder(http.StatusOK, response),
	)

	// Testing valid request.
	req := httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w := httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp := common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.NotEmpty(resp.Result)
	assert.Empty(resp.Error)
	assert.Equal(common.KYCStatus2Status[common.Denied], resp.Result.Status)
	assert.NotEmpty(resp.Result.Details)
	assert.Equal(common.KYCFinality2Finality[common.NonFinal], resp.Result.Details.Finality)
	assert.Len(resp.Result.Details.Reasons, 1)
	assert.Equal("ID_INVALID", resp.Result.Details.Reasons[0])
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)

	// Testing reading request body failure.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", &FailedReader{})
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("Read failed", resp.Error)

	// Testing empty request.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", nil)
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("empty request", resp.Error)

	// Testing malformed request.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader([]byte("malformed request")))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal(`invalid character 'm' looking for beginning of value`, resp.Error)

	// Testing missing Provider field in the request.
	request, err = json.Marshal(&common.CheckStatusRequest{
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("missing KYC provider id in the request", resp.Error)

	// Testing missing CustomerID field in the request.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider: common.SumSub,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("missing verification id in the request", resp.Error)

	// Testing nonexistent KYC provider.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    "Nonexistent Provider",
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusNotFound, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("unknown KYC provider in the request: Nonexistent Provider", resp.Error)

	// Testing KYC provider without config.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    "Fake Provider",
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	if !common.KYCProviders["Fake Provider"] {
		common.KYCProviders["Fake Provider"] = true
	}

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("missing config for Fake Provider", resp.Error)

	// Testing KYC provider that doesn't support status polling.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    common.IDology,
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("IDology doesn't support status polling", resp.Error)

	// Testing KYC provider not implemented yet.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    "Fake Provider",
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)

	config.Cfg["Fake Provider"] = map[string]string{"test": "test"}

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("KYC provider not implemented yet: Fake Provider", resp.Error)

	// Testing error response from the KYC provider.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    common.SumSub,
		ReferenceID: referenceID,
	})

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/resources/applicants/%s/status?key=%s", cfg["Host"], referenceID, cfg["APIKey"]),
		httpmock.NewBytesResponder(http.StatusForbidden, errorResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.NotNil(resp.Result)
	assert.Equal(common.KYCStatus2Status[common.Error], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.NotEmpty(resp.Result.ErrorCode)
	assert.Equal("401", resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
	assert.NotEmpty(resp.Error)
	assert.Equal("Access denied", resp.Error)

	// Testing IdentityMind.
	cfg = config.Cfg[string(common.IdentityMind)]

	assert.NotNil(cfg)

	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:    common.IdentityMind,
		ReferenceID: referenceID,
	})

	assert.Nil(err)
	assert.NotEmpty(request)
	assert.NotEmpty(response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/account/consumer/v2/%s", cfg["Host"], referenceID),
		httpmock.NewBytesResponder(http.StatusOK, identitymindResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(err)
	assert.NotEmpty(resp.Result)
	assert.Empty(resp.Error)
	assert.Equal(common.KYCStatus2Status[common.Approved], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
}
