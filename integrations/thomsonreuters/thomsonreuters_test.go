package thomsonreuters

import (
	"net/http"
	"testing"

	"modulus/kyc/common"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// Test URL parsing error.
	tr := New(Config{
		Host:      "::",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.Empty(tr)

	// Test malformed Host.
	tr = New(Config{
		Host:      "host",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.Empty(tr)

	// Test valid config.
	tr = New(Config{
		Host:      "https://rms-world-check-one-api-pilot.thomsonreuters.com/v1",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.NotEmpty(tr)
	assert.Equal("https", tr.scheme)
	assert.Equal("rms-world-check-one-api-pilot.thomsonreuters.com", tr.host)
	assert.Equal("/v1/", tr.path)
	assert.Equal("key", tr.key)
	assert.Equal("secret", tr.secret)
}

func TestCheckCustomerApproved(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))
	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups/0a3687d0-65b4-1cc3-9975-f20b0000066f/caseTemplate", httpmock.NewStringResponder(http.StatusOK, caseTemplateResponse))
	httpmock.RegisterResponder(http.MethodPost, tr.scheme+"://"+tr.host+tr.path+"cases/screeningRequest", httpmock.NewStringResponder(http.StatusOK, syncScreeningResponseApproved))

	customer := &common.UserData{}

	res, err := tr.CheckCustomer(customer)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCheckCustomerDenied(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))
	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups/0a3687d0-65b4-1cc3-9975-f20b0000066f/caseTemplate", httpmock.NewStringResponder(http.StatusOK, caseTemplateResponse))
	httpmock.RegisterResponder(http.MethodPost, tr.scheme+"://"+tr.host+tr.path+"cases/screeningRequest", httpmock.NewStringResponder(http.StatusOK, syncScreeningResponseDenied))

	customer := &common.UserData{}

	res, err := tr.CheckCustomer(customer)

	assert.NoError(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 3)
	assert.Equal("Case ID: 24da33ec-9ad9-463c-9ef7-9e0dce1bfcbb", res.Details.Reasons[0])
	assert.Equal("Matched Term: Сергей Владимирович Железняк", res.Details.Reasons[1])
	assert.Equal("Category: POLITICAL INDIVIDUAL", res.Details.Reasons[2])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
