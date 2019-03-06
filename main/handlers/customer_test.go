package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
)

var identitymindResponse = []byte(`
{
    "ednaScoreCard": {
        "er": {
            "reportedRule": {
                "description": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
                "details": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
                "name": "Unknown Fallthrough",
                "resultCode": "ACCEPT",
                "ruleId": 1002,
                "testResults": []
            }
        },
        "sc": []
    },
    "erd": "Unknown User",
    "frd": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
    "frn": "Unknown Fallthrough",
    "frp": "ACCEPT",
    "mtid": "26860023",
    "rcd": "1002,101,202,111,131,50005,150",
    "res": "ACCEPT",
    "state": "A",
    "tid": "26860023",
    "upr": "UNKNOWN",
    "user": "UNKNOWN"
}`)

var idologyResponse = []byte(`
<?xml version="1.0"?>
<response>
	<id-number>2073386264</id-number>
	<summary-result>
		<key>id.failure</key>
		<message>FAIL</message>
	</summary-result>
	<results>
		<key>result.match.restricted</key>
		<message>result.match.restricted</message>
	</results>
	<qualifiers>
		<qualifier>
			<key>resultcode.coppa.alert</key>
			<message>COPPA Alert</message>
		</qualifier>
	</qualifiers>
	<idliveq-error>
		<key>id.not.eligible.for.questions</key>
		<message>Not Eligible For Questions</message>
	</idliveq-error>
</response>`)

var idologyErrorResponse = []byte(`
<response>
	<error>Invalid username and password</error>
</response>`)

var shuftiproResponse = []byte(`{"status_code": "SP1", "message": "Verified", "reference": "tester", "signature": "sig"}`)

var sumsubResponse = []byte(`
{
	"id" : "596eb3c93a0eb985b8ade34d",
	"createdAt" : "2017-07-19 03:20:09",
	"inspectionId" : "596eb3c83a0eb985b8ade349",
	"jobId" : "a8f77946-14ff-4398-aa23-a1027e16f627",
	"info" : {
	  "firstName" : "Serge",
	  "middleName" : "Sergeevich",
	  "lastName" : "Sergeew",
	  "dob" : "2000-03-04",
	  "placeOfBirth" : "Saint-Petersburg",
	  "country" : "RUS",
	  "phone" : "+7-911-2081223"
	},
	"email" : "ivanov@gmail.com"
}`)

var truliooConsentsResponse = []byte(`
[
    "Australia Driver Licence",
    "Australia Passport",
    "Birth Registry",
    "Visa Verification",
    "DVS Driver License Search",
    "DVS Medicare Search",
    "DVS Passport Search",
    "DVS Visa Search",
    "DVS ImmiCard Search",
    "DVS Citizenship Certificate Search",
    "DVS Certificate of Registration by Descent Search",
    "Credit Agency"
]`)

var truliooResponse = []byte(`
{
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
}`)

func init() {
	if config.Cfg == nil {
		config.Cfg = cfg
	}
}

func TestCheckCustomer(t *testing.T) {
	assert := assert.New(t)

	request, err := json.Marshal(&common.CheckCustomerRequest{
		Provider: common.IDology,
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)
	assert.NotEmpty(response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://web.idologylive.com/api/idiq.svc",
		httpmock.NewBytesResponder(http.StatusOK, idologyResponse),
	)

	// Testing valid request.
	req := httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w := httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp := common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotEmpty(resp.Result)
	assert.Empty(resp.Error)
	assert.Equal(common.KYCStatus2Status[common.Denied], resp.Result.Status)
	assert.NotEmpty(resp.Result.Details)
	assert.Equal(common.KYCFinality2Finality[common.Unknown], resp.Result.Details.Finality)
	assert.Len(resp.Result.Details.Reasons, 1)
	assert.Equal("COPPA Alert", resp.Result.Details.Reasons[0])
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)

	// Testing reading request body failure.
	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", &FailedReader{})
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("Read failed", resp.Error)

	// Testing empty request.
	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", nil)
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("empty request", resp.Error)

	// Testing malformed request.
	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader([]byte("malformed request")))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal(`invalid character 'm' looking for beginning of value`, resp.Error)

	// Testing missing Provider field in the request.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusBadRequest, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("missing KYC provider id in the request", resp.Error)

	// Testing nonexistent KYC provider.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: "Nonexistent Provider",
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusNotFound, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("unknown KYC provider in the request: Nonexistent Provider", resp.Error)

	// Testing KYC provider without config.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: "Provider Without Config",
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	if !common.KYCProviders["Provider Without Config"] {
		common.KYCProviders["Provider Without Config"] = true
	}

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("missing config for Provider Without Config", resp.Error)

	// Testing KYC provider not implemented yet.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: "Not Implemented Provider",
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	if !common.KYCProviders["Not Implemented Provider"] {
		common.KYCProviders["Not Implemented Provider"] = true
	}
	config.Cfg["Not Implemented Provider"] = map[string]string{"test": "test"}

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusUnprocessableEntity, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal("KYC provider not implemented yet: Not Implemented Provider", resp.Error)

	// Testing error response from the KYC provider.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.IDology,
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://web.idologylive.com/api/idiq.svc",
		httpmock.NewBytesResponder(http.StatusForbidden, idologyErrorResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotNil(resp.Result)
	assert.Equal(common.KYCStatus2Status[common.Error], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
	assert.NotEmpty(resp.Error)
	assert.Equal("during verification: Invalid username and password", resp.Error)

	// Testing IdentityMind.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.IdentityMind,
		UserData: &common.UserData{
			AccountName: "tester",
		},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://sandbox.identitymind.com/im/account/consumer",
		httpmock.NewBytesResponder(http.StatusOK, identitymindResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotNil(resp.Result)
	assert.Equal(common.KYCStatus2Status[common.Approved], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
	assert.Empty(resp.Error)

	// Testing ShuftiPro.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.ShuftiPro,
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://api.shuftipro.com",
		httpmock.NewBytesResponder(http.StatusOK, shuftiproResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotNil(resp.Result)
	assert.Equal(common.KYCStatus2Status[common.Approved], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
	assert.Empty(resp.Error)

	// Testing Sum&Substance.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.SumSub,
		UserData: &common.UserData{
			IDCard: &common.IDCard{
				Number: "xyz",
			},
		},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	sumsubCfg := cfg[string(common.SumSub)]

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("%s/resources/applicants?key=%s", sumsubCfg["Host"], sumsubCfg["APIKey"]),
		httpmock.NewBytesResponder(http.StatusOK, sumsubResponse),
	)

	httpmock.RegisterResponder(
		http.MethodPost,
		fmt.Sprintf("%s/resources/applicants/596eb3c93a0eb985b8ade34d/info/idDoc?key=%s", sumsubCfg["Host"], sumsubCfg["APIKey"]),
		httpmock.NewStringResponder(http.StatusOK, `{"ok":1}`),
	)

	expectedQuery := map[string]string{
		"reason": "docs_sent",
		"key":    sumsubCfg["APIKey"],
	}
	httpmock.RegisterResponderWithQuery(
		http.MethodPost,
		fmt.Sprintf("%s/resources/applicants/596eb3c93a0eb985b8ade34d/status/pending", sumsubCfg["Host"]),
		expectedQuery,
		httpmock.NewStringResponder(http.StatusOK, `{"ok":1}`),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotEmpty(resp.Result)
	assert.Empty(resp.Error)
	assert.Equal(common.KYCStatus2Status[common.Unclear], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.NotNil(resp.Result.StatusCheck)
	assert.Equal(common.SumSub, resp.Result.StatusCheck.Provider)
	assert.Equal("596eb3c93a0eb985b8ade34d", resp.Result.StatusCheck.ReferenceID)
	assert.NotZero(t, time.Time(resp.Result.StatusCheck.LastCheck))

	// Testing Trulioo.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.Trulioo,
		UserData: &common.UserData{
			CountryAlpha2: "AU",
		},
	})

	assert.NoError(err)
	assert.NotEmpty(request)

	truliooCfg := cfg[string(common.Trulioo)]

	httpmock.RegisterResponder(
		http.MethodGet,
		truliooCfg["Host"]+"/configuration/v1/consents/Identity%20Verification/AU",
		httpmock.NewBytesResponder(http.StatusOK, truliooConsentsResponse),
	)

	httpmock.RegisterResponder(
		http.MethodPost,
		truliooCfg["Host"]+"/verifications/v1/verify",
		httpmock.NewBytesResponder(http.StatusOK, truliooResponse),
	)

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusOK, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.NotNil(resp.Result)
	assert.Equal(common.KYCStatus2Status[common.Error], resp.Result.Status)
	assert.Nil(resp.Result.Details)
	assert.Empty(resp.Result.ErrorCode)
	assert.Nil(resp.Result.StatusCheck)
	assert.NotEmpty(resp.Error)
	assert.Equal("Test error;", resp.Error)

	// Testing IDology config error.
	request, err = json.Marshal(&common.CheckCustomerRequest{
		Provider: common.IDology,
		UserData: &common.UserData{},
	})

	assert.NoError(err)
	assert.NotEmpty(request)
	assert.NotEmpty(response)

	config.Cfg[string(common.IDology)] = map[string]string{
		"Host":     "https://web.idologylive.com/api/idiq.svc",
		"Username": "fakeuser",
		"Password": "fakepassword",
	}

	req = httptest.NewRequest(http.MethodPost, "/CheckCustomer", bytes.NewReader(request))
	w = httptest.NewRecorder()

	handlers.CheckCustomer(w, req)

	assert.Equal(http.StatusInternalServerError, w.Code)
	assert.Equal("application/json; charset=utf-8", w.Header().Get("Content-Type"))

	resp = common.KYCResponse{}

	err = json.Unmarshal(w.Body.Bytes(), &resp)

	assert.NoError(err)
	assert.Nil(resp.Result)
	assert.NotEmpty(resp.Error)
	assert.Equal(`IDology config error: strconv.ParseBool: parsing "": invalid syntax`, resp.Error)
}
