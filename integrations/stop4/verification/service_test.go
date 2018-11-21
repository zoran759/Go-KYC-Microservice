package verification

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"fmt"
	"errors"
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

func TestServiceVerify(t *testing.T) {
	config := Config{
		Host:       "https://private-f1649f-coreservices2.apiary-mock.com",
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
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				fmt.Sprintf(`{"status": 0, "id": "223", "score": 98, "rec":"Approve"}`),
			), nil
		},
	)

	response, err := testService.Verify(RegistrationRequest{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, 0, response.Status)
		assert.Equal(t, 98, response.Score)
		assert.Equal(t, "Approve", response.Rec)
	}
}

func TestServiceVerifyError(t *testing.T) {
	config := Config{
		Host:       "https://private-f1649f-coreservices2.apiary-mock.com",
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
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`,
			), nil
		},
	)

	response, err := testService.Verify(RegistrationRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = testService.Verify(RegistrationRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)


	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusBadRequest,
				fmt.Sprintf(`{"status": -13, "id": "223", "score": 98, "rec":"Decline"}`),
			), nil
		},
	)

	response, err = testService.Verify(RegistrationRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				fmt.Sprintf(`{"status": -13, "id": "223", "score": 98, "rec":"Decline"}`),
			), nil
		},
	)

	response, err = testService.Verify(RegistrationRequest{})
	assert.Nil(t, err)
	assert.Equal(t, "Wrong customer email format", response.Details)
	assert.Equal(t, -13, response.Status)
}

func TestServiceVerifyCroppedJson(t *testing.T) {
	config := Config{
		Host:       "https://private-f1649f-coreservices2.apiary-mock.com",
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
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
   "id":"233",
   "status":0,
   "description":"Success",
   "score":98,
   "rec":"Approve",
   "rules_triggered":[
      {
         "name":"Multi-Accounting : IP shared to Chargeback reason",
         "score":"100.00",
         "display_to_merchant":1
      }
   ],
   "scrubber_results":{
      "geo_check":"",
      "address_verification":"",
      "phone_verify":"",
      "idv_usa":"",
      "idv_global":"",
      "gav":"",
      "idv_br":"",
      "bav_usa":"",
      "bav_advanced":"",
      "cb_aml":"",
      "cb_bvs":"",
      "email_age":"",
      "compliance_watchlist":"",
      "iovation":"",
      "idv_advance":""
    },
    "facts":[{  
      "type":"1",
      "text":"Which one of the following addresses is associated with you?",
      "answers":[{
        "correct":"false",
        "text":"509 BIRDIE RD"
      },{  
        "correct":"false",
        "text":"667 ASHWOOD NORT CT"
      },{
        "correct":"true",
        "text":"291 LYNCH RD"
      }},{
      "type":"2",
      "text":"Which one of the following area codes is associated with you?",
      "answers":[{  
        "correct":"false",
        "text":"901"
      },{
        "correct":"true",
        "text":"407/321"
      },{
        "correct":"false",
        "text":"352"
      }}],
    "confidence_level":91.50
}`,
			), nil
		},
	)

	response, err := testService.Verify(RegistrationRequest{})
	if assert.NoError(t, err) && assert.NotNil(t, response) {
		assert.Equal(t, 0, response.Status)
		assert.Equal(t, 98, response.Score)
		assert.Equal(t, "Approve", response.Rec)
	}


	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://private-f1649f-coreservices2.apiary-mock.com/customerregistration",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
   "id":"233",
   "status":0,
   "description":"Success",
   "score":98,
   "rec":"Approve",
   "rules_triggered":[
      {
         "name":"Multi-Accounting : IP shared to Chargeback reason",
         "score":"100.00",
         "display_to_merchant":1
      }}
   ],
    "facts":[{  
      "type":"1",
      "text":"Which one of the following addresses is associated with you?",
      "answers":[{
        "correct":"false",
        "text":"509 BIRDIE RD"
      },{  
        "correct":"false",
        "text":"667 ASHWOOD NORT CT"
      },{
        "correct":"true",
        "text":"291 LYNCH RD"
      }},{
      "type":"2",
      "text":"Which one of the following area codes is associated with you?",
      "answers":[{  
        "correct":"false",
        "text":"901"
      },{
        "correct":"true",
        "text":"407/321"
      },{
        "correct":"false",
        "text":"352"
      }}],
    "confidence_level":91.50
}`,
			), nil
		},
	)

	response, err = testService.Verify(RegistrationRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestServiceVerifyAuthParams(t *testing.T) {
	config := Config{
		Host:       "https://private-f1649f-coreservices2.apiary-mock.com",
		MerchantID: "INT-123456",
		Password:   "123456",
	}

	testService := service{
		config: config,
	}

	request := RegistrationRequest{}
	testService.AttachAuthPartToRequest(&request)

	assert.Equal(t, "INT-123456", request.MerchantID)
	assert.Equal(t, "123456", request.Password)
}

func TestServiceConvertRequestToForm(t *testing.T) {

	request := RegistrationRequest{}
	request.CustomerInformation.FirstName =
		CustomerInformationField{
			FieldName:  "customer_information[first_name]",
			FieldVal: "Linus",
		}
	request.CustomerInformation.MiddleName =
		CustomerInformationField{
			FieldName:  "customer_information[middle_name]",
			FieldVal: "Benedict",
		}
	request.CustomerInformation.LastName =
		CustomerInformationField{
			FieldName:  "customer_information[last_name]",
			FieldVal: "Torvalds",
		}

	fields := make(map[string]string)
	convertModelToForm(request.CustomerInformation, fields)

	assert.Equal(t, "Linus", fields["customer_information[first_name]"])
	assert.Empty(t, fields["customer_information[email]"])
}

func TestServiceGenerateUserNumber(t *testing.T) {
	uNumber, err := generateUserNumber()
	uNumberControl, _ := generateUserNumber()

	assert.NotNil(t, uNumber)
	assert.Nil(t, err)
	assert.NotEqual(t, uNumber, uNumberControl)
}