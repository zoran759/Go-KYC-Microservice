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

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var cfg = config.Config{
	common.IdentityMind: {
		"Host":     "https://sandbox.identitymind.com/im",
		"Username": "fakeuser",
		"Password": "fakepassword",
	},
	common.IDology: {
		"Host":             "https://web.idologylive.com/api/idiq.svc",
		"Username":         "fakeuser",
		"Password":         "fakepassword",
		"UseSummaryResult": "false",
	},
	common.ShuftiPro: {
		"Host":        "https://api.shuftipro.com",
		"ClientID":    "fakeID",
		"SecretKey":   "fakeKey",
		"RedirectURL": "https://api.shuftipro.com",
	},
	common.SumSub: {
		"Host":   "https://test-api.sumsub.com",
		"APIKey": "fakeKey",
	},
	common.Trulioo: {
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
	if config.KYC == nil {
		config.KYC = cfg
	}
}

type FailedReader struct{}

func (r FailedReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("Read failed")
}

func TestCheckStatus(t *testing.T) {
	cfg := config.KYC[common.SumSub]

	assert.NotNil(t, cfg)

	customerID := "testID"

	request, err := json.Marshal(&common.CheckStatusRequest{
		Provider:   common.SumSub,
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)
	assert.NotEmpty(t, response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/resources/applicants/%s/status?key=%s", cfg["Host"], customerID, cfg["APIKey"]),
		httpmock.NewBytesResponder(http.StatusOK, response),
	)

	// Testing valid request.
	req := httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w := httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp := common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Result)
	assert.Empty(t, resp.Error)
	assert.Equal(t, common.Denied, resp.Result.Status)
	assert.NotEmpty(t, resp.Result.Details)
	assert.Equal(t, common.NonFinal, resp.Result.Details.Finality)
	assert.Len(t, resp.Result.Details.Reasons, 1)
	assert.Equal(t, "ID_INVALID", resp.Result.Details.Reasons[0])
	assert.Empty(t, resp.Result.ErrorCode)
	assert.Nil(t, resp.Result.StatusPolling)

	// Testing reading request body failure.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", &FailedReader{})
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "Read failed", resp.Error)

	// Testing empty request.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", nil)
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "empty request", resp.Error)

	// Testing malformed request.
	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader([]byte("malformed request")))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, `invalid character 'm' looking for beginning of value`, resp.Error)

	// Testing missing Provider field in the request.
	request, err = json.Marshal(&common.CheckStatusRequest{
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "missing KYC provider id in the request", resp.Error)

	// Testing missing CustomerID field in the request.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider: common.SumSub,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "missing verification id in the request", resp.Error)

	// Testing nonexistent KYC provider.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   "Nonexistent Provider",
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "unknown KYC provider in the request: Nonexistent Provider", resp.Error)

	// Testing KYC provider without config.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   "Fake Provider",
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	if !common.KYCProviders["Fake Provider"] {
		common.KYCProviders["Fake Provider"] = true
	}

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "missing config for Fake Provider", resp.Error)

	// Testing KYC provider that doesn't support status polling.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   common.IDology,
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "IDology doesn't support status polling", resp.Error)

	// Testing KYC provider not implemented yet.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   "Fake Provider",
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)

	config.KYC["Fake Provider"] = map[string]string{"test": "test"}

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.Nil(t, resp.Result)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "KYC provider not implemented yet: Fake Provider", resp.Error)

	// Testing error response from the KYC provider.
	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   common.SumSub,
		CustomerID: customerID,
	})

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/resources/applicants/%s/status?key=%s", cfg["Host"], customerID, cfg["APIKey"]),
		httpmock.NewBytesResponder(http.StatusForbidden, errorResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.NotNil(t, resp.Result)
	assert.Equal(t, common.Error, resp.Result.Status)
	assert.Nil(t, resp.Result.Details)
	assert.NotEmpty(t, resp.Result.ErrorCode)
	assert.Equal(t, "401", resp.Result.ErrorCode)
	assert.Nil(t, resp.Result.StatusPolling)
	assert.NotEmpty(t, resp.Error)
	assert.Equal(t, "Access denied", resp.Error)

	// Testing IdentityMind.
	cfg = config.KYC[common.IdentityMind]

	assert.NotNil(t, cfg)

	request, err = json.Marshal(&common.CheckStatusRequest{
		Provider:   common.IdentityMind,
		CustomerID: customerID,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, request)
	assert.NotEmpty(t, response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		fmt.Sprintf("%s/account/consumer/v2/%s", cfg["Host"], customerID),
		httpmock.NewBytesResponder(http.StatusOK, identitymindResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckStatus", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckStatus(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.Nil(t, err)
	assert.NotEmpty(t, resp.Result)
	assert.Empty(t, resp.Error)
	assert.Equal(t, common.Approved, resp.Result.Status)
	assert.NotEmpty(t, resp.Result.Details)
	assert.Equal(t, common.Unknown, resp.Result.Details.Finality)
	assert.NotEmpty(t, resp.Result.Details.Reasons)
	assert.Empty(t, resp.Result.ErrorCode)
	assert.Nil(t, resp.Result.StatusPolling)
}
