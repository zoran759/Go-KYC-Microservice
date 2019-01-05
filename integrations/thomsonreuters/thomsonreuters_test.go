package thomsonreuters

import (
	"testing"

	"modulus/kyc/common"

	"github.com/jarcoal/httpmock"
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

func TestCheckCustomer(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	customer := &common.UserData{}

	res, err := tr.CheckCustomer(customer)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
