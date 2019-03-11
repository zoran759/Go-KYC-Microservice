package shuftipro

import (
	"encoding/base64"
	"errors"
	stdhttp "net/http"
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/http"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

var reqInvalidResponse = `{
    "reference": "17374217",
    "event": "request.invalid",
    "error": {
        "service": "document",
        "key": "dob",
        "message": "The dob does not match the format Y-m-d."
    },
    "verification_url": ""
}`

var reqUnauthorizedResponse = `{
	"reference": "",
	"event": "request.unauthorized",
	"error": {
	  "service": "",
	  "key": "",
	  "message": "Authorization keys are missing/invalid."
	},
	"token": "",
	"verification_url": ""
}`

var reqAcceptedResponse = `{
	"reference": "17374217",
	"event": "verification.accepted",
	"error": "",
	"verification_url": "",
	"verification_result": {
	  "document": {
		"name": 1,
		"dob": 1,
		"expiry_date": 1,
		"issue_date": 1,
		"document_number": 1,
		"document": 1
	  },
	  "address": {
		"name": 1,
		"full_address": 1
	  }
	},
	"verification_data": {
	  "document": {
		"name": {
		  "first_name": "John",
		  "middle_name": "Carter",
		  "last_name": "Doe"
		},
		"dob": "1978-03-13",
		"issue_date": "2015-10-10",
		"expiry_date": "2025-12-31",
		"document_number": "1456-0989-5567-0909",
		"selected_type": [
		  "id_card"
		],
		"supported_types": [
		  "id_card",
		  "driving_license",
		  "passport"
		]
	  },
	  "address": {
		"name": {
		  "first_name": "John",
		  "middle_name": "Carter",
		  "last_name": "Doe"
		},
		"full_address": "3339 Maryland Avenue, Largo, Florida",
		"selected_type": [
		  "id_card"
		],
		"supported_types": [
		  "id_card",
		  "bank_statement"
		]
	  }
	}
}`

var reqDeclinedResponse = `{
	"reference": "95156124",
	"event": "verification.declined",
	"error": "",
	"verification_url": "",
	"verification_result": {
	  "document": {
		"name": null,
		"dob": null,
		"expiry_date": null,
		"issue_date": null,
		"document_number": null,
		"document": null
	  },
	  "address": {
		"name": null,
		"full_address": null
	  },
	  "face": 0,
	  "background_checks": null
	},
	"verification_data": {
	  "document": {
		"name": {
		  "first_name": "John",
		  "middle_name": "Carter",
		  "last_name": "Doe"
		},
		"dob": "1978-03-13",
		"issue_date": "2015-10-10",
		"expiry_date": "2025-12-31",
		"document_number": "1456-0989-5567-0909",
		"selected_type": [
		  "id_card"
		],
		"supported_types": [
		  "id_card",
		  "driving_license",
		  "passport"
		]
	  },
	  "address": {
		"name": {
		  "first_name": "John",
		  "middle_name": "Carter",
		  "last_name": "Doe"
		},
		"full_address": "3339 Maryland Avenue, Largo, Florida",
		"selected_type": [
		  "id_card"
		],
		"supported_types": [
		  "id_card",
		  "bank_statement"
		]
	  }
	},
	"declined_reason": "Face is not verified."
}`

var changedResponse = `{
	"reference": "17374217",
	"event": 5,
	"error": "",
	"verification_url": ""
  }`

var changedErrorFormatResponse = `{
    "reference": "17374217",
    "event": "request.invalid",
    "error": "The dob does not match the format Y-m-d.",
    "verification_url": ""
}`

var unexpectedErrorFormatResponse = `{
    "reference": "17374217",
    "event": "request.invalid",
    "error": {
        "errors":[
			"The dob does not match the format Y-m-d."
		]
    },
    "verification_url": ""
}`

func TestNewClient(t *testing.T) {
	config := Config{
		Host:        "host",
		ClientID:    "ClientID",
		SecretKey:   "SecretKey",
		CallbackURL: "CallbackURL",
	}

	client1 := Client{
		host: config.Host,
		headers: http.Headers{
			"Content-Type":  "application/json",
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.ClientID+":"+config.SecretKey)),
		},
		callbackURL: config.CallbackURL,
	}

	client2 := NewClient(config)

	assert.Equal(t, client1, client2)
}

func TestCheckCustomer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	client := NewClient(Config{
		Host:        "https://shuftipro.com/api",
		ClientID:    "ClientID",
		SecretKey:   "SecretKey",
		CallbackURL: "callback_url",
	})
	customer := &common.UserData{}

	type testCase struct {
		name      string
		responder httpmock.Responder
		result    common.KYCResult
		err       error
	}

	testCases := []testCase{
		testCase{
			name:      "Invalid request 400",
			responder: httpmock.NewStringResponder(stdhttp.StatusBadRequest, reqInvalidResponse),
			result: common.KYCResult{
				ErrorCode: "400",
			},
			err: Error{
				Service: "document",
				Key:     "dob",
				Message: "The dob does not match the format Y-m-d.",
			},
		},
		testCase{
			name:      "Unauthorized request 401",
			responder: httpmock.NewStringResponder(stdhttp.StatusUnauthorized, reqUnauthorizedResponse),
			result: common.KYCResult{
				ErrorCode: "401",
			},
			err: Error{
				Service: "",
				Key:     "",
				Message: "Authorization keys are missing/invalid.",
			},
		},
		testCase{
			name:      "Accepted verification 200",
			responder: httpmock.NewStringResponder(stdhttp.StatusOK, reqAcceptedResponse),
			result: common.KYCResult{
				Status: common.Approved,
			},
		},
		testCase{
			name:      "Declined verification 200",
			responder: httpmock.NewStringResponder(stdhttp.StatusOK, reqDeclinedResponse),
			result: common.KYCResult{
				Status: common.Denied,
				Details: &common.KYCDetails{
					Reasons: []string{"Face is not verified."},
				},
			},
		},
		testCase{
			name:   "Error no responder",
			result: common.KYCResult{},
			err:    errors.New("Post https://shuftipro.com/api: no responder found"),
		},
		testCase{
			name:      "Test changed response format",
			responder: httpmock.NewStringResponder(stdhttp.StatusOK, changedResponse),
			result:    common.KYCResult{},
			err:       errors.New("json: cannot unmarshal number into Go struct field Response.event of type shuftipro.Event"),
		},
		testCase{
			name:      "Test changed error field format",
			responder: httpmock.NewStringResponder(stdhttp.StatusBadRequest, changedErrorFormatResponse),
			result: common.KYCResult{
				ErrorCode: "400",
			},
			err: errors.New("request parameters provided in the request are invalid; check the error code in the result"),
		},
		testCase{
			name:      "Test unexpected error field format",
			responder: httpmock.NewStringResponder(stdhttp.StatusBadRequest, unexpectedErrorFormatResponse),
			result: common.KYCResult{
				ErrorCode: "400",
			},
			err: errors.New("unexpected format of the returned error: please, report to developers"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			httpmock.RegisterResponder(stdhttp.MethodPost, client.host, tc.responder)
			res, err := client.CheckCustomer(customer)
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
