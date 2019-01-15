package coinfirm

import (
	"encoding/json"
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

var (
	tokenResp          = `{"token":"yFaReURiYkAECZsPt8dR1bzHpa2Y5kXpqsp4KunyH870OAoY577vI8mhABCj4vkK"}`
	malformedResp      = `the fake response to test unexpected answer`
	error400Resp       = `{"error":"Invalid email or password"}`
	newParticipantResp = `{"uuid": "33611d6d-2826-4c3e-a777-3f0397e283fc"}`
)

/**************************************************************************************************
 * newAuthToken() tests                                                                           *
 **************************************************************************************************/
func TestNewAuthTokenSuccess(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, tokenResp))

	token, status, err := c.newAuthToken()

	assert.NoError(err)
	assert.Nil(status)
	assert.Equal("yFaReURiYkAECZsPt8dR1bzHpa2Y5kXpqsp4KunyH870OAoY577vI8mhABCj4vkK", token.Token)
}

func TestNewAuthTokenUnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusOK, malformedResp))

	token, status, err := c.newAuthToken()

	assert.Error(err)
	assert.Equal("invalid character 'h' in literal true (expecting 'r')", err.Error())
	assert.Nil(status)
	assert.Empty(token.Token)
}

func TestNewAuthTokenSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	token, status, err := c.newAuthToken()

	assert.Error(err)
	assert.Equal("Post https://api.coinfirm.io/v2/auth/login: no responder found", err.Error())
	assert.Nil(status)
	assert.Empty(token.Token)
}

func TestNewAuthTokenStatus400(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusBadRequest, error400Resp))

	token, status, err := c.newAuthToken()

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

	httpmock.RegisterResponder(http.MethodPost, c.host+"/auth/login", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	token, status, err := c.newAuthToken()

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(token.Token)
}

/**************************************************************************************************
 * newParticipant() tests                                                                         *
 **************************************************************************************************/

func TestNewParticipantCuccess(t *testing.T) {
	assert := assert.New(t)

	c := New(Config{
		Host:     "https://api.coinfirm.io/v2",
		Email:    "info@fuzioncapital.com",
		Password: "CAc8@6e12e823c602bcb85224e822609",
		Company:  "Fuzion",
	})

	// token, status, err := c.newAuthToken()
	// assert.NoError(err)
	// assert.Nil(status)
	// assert.NotEmpty(token)

	token := model.AuthResponse{}
	_ = json.Unmarshal([]byte(tokenResp), &token)
	c.headers["Authorization"] = "Bearer " + token.Token

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusOK, newParticipantResp))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(nParticipant)

	assert.NoError(err)
	assert.Nil(status)
	assert.Equal("33611d6d-2826-4c3e-a777-3f0397e283fc", participant.UUID)
}

func TestNewParticipantSendError(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	participant, status, err := c.newParticipant(model.NewParticipant{})

	assert.Error(err)
	assert.Equal("Put https://api.coinfirm.io/v2/kyc/customers/Fuzion: no responder found", err.Error())
	assert.Nil(status)
	assert.Empty(participant)
}

func TestNewParticipantStatus400(t *testing.T) {
	assert := assert.New(t)

	c := New(Config{
		Host:     "https://api.coinfirm.io/v2",
		Email:    "info@fuzioncapital.com",
		Password: "CAc8@6e12e823c602bcb85224e822609",
		Company:  "Fuzion",
	})

	token := model.AuthResponse{}
	_ = json.Unmarshal([]byte(tokenResp), &token)
	c.headers["Authorization"] = "Bearer " + token.Token

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusBadRequest, `{"error":"Request body validation errors"}`))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(nParticipant)

	assert.Error(err)
	assert.Equal("Request body validation errors", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(participant)
}

func TestNewParticipantStatus400UnexpectedResp(t *testing.T) {
	assert := assert.New(t)

	c := New(Config{
		Host:     "https://api.coinfirm.io/v2",
		Email:    "info@fuzioncapital.com",
		Password: "CAc8@6e12e823c602bcb85224e822609",
		Company:  "Fuzion",
	})

	token := model.AuthResponse{}
	_ = json.Unmarshal([]byte(tokenResp), &token)
	c.headers["Authorization"] = "Bearer " + token.Token

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodPut, c.host+"/kyc/customers/Fuzion", httpmock.NewStringResponder(http.StatusBadRequest, malformedResp))

	nParticipant := model.NewParticipant{
		Email: "sarbash.s@ya.ru",
	}

	participant, status, err := c.newParticipant(nParticipant)

	assert.Error(err)
	assert.Equal("http error", err.Error())
	assert.NotNil(status)
	assert.Equal(400, *status)
	assert.Empty(participant)
}

/**************************************************************************************************
 * sendParticipantDetails() tests                                                                 *
 **************************************************************************************************/

/**************************************************************************************************
 * sendDocFile() tests                                                                            *
 **************************************************************************************************/

/**************************************************************************************************
 * getParticipantCurrentStatus() tests                                                            *
 **************************************************************************************************/
