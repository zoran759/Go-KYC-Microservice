package coinfirm

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"testing"

	"modulus/kyc/integrations/coinfirm/model"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var c = New(Config{
	Host:     "https://api.coinfirm.io/v2",
	Email:    "info@fuzioncapital.com",
	Password: "CAc8@6e12e823c602bcb85224e822609",
	Company:  "Fuzion",
})

var hdrs = headers()

var (
	tokenResp            = `{"token":"yFaReURiYkAECZsPt8dR1bzHpa2Y5kXpqsp4KunyH870OAoY577vI8mhABCj4vkK"}`
	malformedResp        = `the fake response to test unexpected answer`
	error400Resp         = `{"error":"Invalid email or password"}`
	newParticipantResp   = `{"uuid": "33611d6d-2826-4c3e-a777-3f0397e283fc"}`
	statusInprogressResp = `{"current_status": "inprogress","ico_owner_status": "accepted"}`
	statusLowResp        = `{"current_status": "low","ico_owner_status": "accepted"}`
)

/**************************************************************************************************
 * newAuthToken() tests                                                                           *
 **************************************************************************************************/
func TestNewAuthTokenSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, tokenResp))

	token, status, err := c.newAuthToken(hdrs)

	assert.NoError(err)
	assert.Nil(status)
	assert.Equal("yFaReURiYkAECZsPt8dR1bzHpa2Y5kXpqsp4KunyH870OAoY577vI8mhABCj4vkK", token.Token)
}

func TestNewAuthTokenUnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, malformedResp))

	token, status, err := c.newAuthToken(hdrs)

	assert.Error(err)
	assert.Equal("invalid character 'h' in literal true (expecting 'r')", err.Error())
	assert.Nil(status)
	assert.Empty(token.Token)
}

func TestNewAuthTokenSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	token, status, err := c.newAuthToken(hdrs)

	assert.Error(err)
	assert.Equal("Post https://api.coinfirm.io/v2/auth/login: no responder found", err.Error())
	assert.Nil(status)
	assert.Empty(token.Token)
}

func TestNewAuthTokenStatus400(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/auth/login", httpmock.NewStringResponder(http.StatusBadRequest, error400Resp))

	token, status, err := c.newAuthToken(hdrs)

	assert.Error(err)
	assert.Equal("Invalid email or password", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(token.Token)
}

func TestNewAuthTokenStatus400UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/auth/login", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	token, status, err := c.newAuthToken(hdrs)

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(token.Token)
}

/**************************************************************************************************
 * newParticipant() tests                                                                         *
 **************************************************************************************************/

func TestNewParticipantSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusOK, newParticipantResp))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(hdrs, nParticipant)

	assert.NoError(err)
	assert.Nil(status)
	assert.Equal("33611d6d-2826-4c3e-a777-3f0397e283fc", participant.UUID)
}

func TestNewParticipantSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	participant, status, err := c.newParticipant(hdrs, model.NewParticipant{})

	assert.Error(err)
	assert.Equal("Put https://api.coinfirm.io/v2/kyc/customers/Fuzion: no responder found", err.Error())
	assert.Nil(status)
	assert.Empty(participant)
}

func TestNewParticipantStatus400(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusBadRequest, `{"error":"Request body validation errors"}`))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(hdrs, nParticipant)

	assert.Error(err)
	assert.Equal("Request body validation errors", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(participant)
}

func TestNewParticipantStatus400UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(hdrs, nParticipant)

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(participant)
}

/**************************************************************************************************
 * sendParticipantDetails() tests                                                                 *
 **************************************************************************************************/

func TestSendParticipantDetailsSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/forms/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewBytesResponder(http.StatusCreated, nil))

	participant := model.ParticipantDetails{
		UserIP:        "192.168.0.117",
		Type:          model.Individual,
		FirstName:     "John",
		LastName:      "Doe",
		Email:         "john.doe@mail.com",
		Nationality:   "US",
		IDNumber:      "987654321",
		CountryAlpha3: "US",
		Postcode:      "15212",
		City:          "Pittsburgh",
		Street:        "Gifford St",
		BirthDate:     "1960-08-15",
	}

	status, err := c.sendParticipantDetails(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", participant)

	assert.NoError(err)
	assert.Nil(status)
}
func TestSendParticipantDetailsSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	status, err := c.sendParticipantDetails(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", model.ParticipantDetails{})

	assert.Error(err)
	assert.Equal("Put https://api.coinfirm.io/v2/kyc/forms/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc: no responder found", err.Error())
	assert.Nil(status)
}

func TestSendParticipantDetailsStatus400(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/forms/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusBadRequest, `{"error":"Request body validation errors"}`))

	participant := model.ParticipantDetails{
		UserIP:        "192.168.0.117",
		Type:          model.Individual,
		FirstName:     "John",
		LastName:      "Doe",
		CountryAlpha3: "US",
		Postcode:      "15212",
		City:          "Pittsburgh",
		Street:        "Gifford St",
	}

	status, err := c.sendParticipantDetails(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", participant)

	assert.Error(err)
	assert.Equal("Request body validation errors", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
}

func TestSendParticipantDetailsStatus400UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.config.Host+"/kyc/forms/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	status, err := c.sendParticipantDetails(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", model.ParticipantDetails{})

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
}

/**************************************************************************************************
 * sendDocFile() tests                                                                            *
 **************************************************************************************************/

func TestSendDocFileSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/kyc/files/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewBytesResponder(http.StatusCreated, nil))

	data, _ := ioutil.ReadFile("../../test_data/realId.jpg")

	docfile := &model.File{
		Type:       model.FileID,
		Extension:  "jpg",
		DataBase64: base64.StdEncoding.EncodeToString(data),
	}

	status, err := c.sendDocFile(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", docfile)

	assert.NoError(err)
	assert.Nil(status)
}

func TestSendDocFileSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	status, err := c.sendDocFile(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", &model.File{})

	assert.Error(err)
	assert.Equal("Post https://api.coinfirm.io/v2/kyc/files/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc: no responder found", err.Error())
	assert.Nil(status)
}

func TestSendDocFileStatus400(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/kyc/files/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusBadRequest, `{"error":"Request body validation errors"}`))

	status, err := c.sendDocFile(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", &model.File{})

	assert.Error(err)
	assert.Equal("Request body validation errors", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
}

func TestSendDocFileStatus400UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.config.Host+"/kyc/files/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	status, err := c.sendDocFile(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc", &model.File{})

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
}

/**************************************************************************************************
 * getParticipantCurrentStatus() tests                                                            *
 **************************************************************************************************/

func TestGetParticipantCurrentStatusSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, c.config.Host+"/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusOK, statusInprogressResp))

	status, code, err := c.getParticipantCurrentStatus(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc")

	assert.NoError(err)
	assert.Nil(code)
	assert.Equal(model.InProgress, status.CurrentStatus)
	assert.Equal("accepted", status.ICOOwnerStatus)
}

func TestGetParticipantCurrentStatusSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	status, code, err := c.getParticipantCurrentStatus(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc")

	assert.Error(err)
	assert.Equal("Get https://api.coinfirm.io/v2/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc: no responder found", err.Error())
	assert.Nil(code)
	assert.Empty(status)
}

func TestGetParticipantCurrentStatusStatus404(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, c.config.Host+"/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusNotFound, `{"error":"Resource not found"}`))

	status, code, err := c.getParticipantCurrentStatus(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc")

	assert.Error(err)
	assert.Equal("Resource not found", err.Error())
	assert.NotNil(code)
	assert.Equal(404, *code)
	assert.Empty(status)
}

func TestGetParticipantCurrentStatusStatus404UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, c.config.Host+"/kyc/status/Fuzion/33611d6d-2826-4c3e-a777-3f0397e283fc", httpmock.NewStringResponder(http.StatusNotFound, malformedResp))

	status, code, err := c.getParticipantCurrentStatus(hdrs, "33611d6d-2826-4c3e-a777-3f0397e283fc")

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(code)
	assert.Equal(404, *code)
	assert.Empty(status)
}
