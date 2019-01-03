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
	tomson := New(Config{
		Host:      "::",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.Empty(tomson)

	// Test malformed Host.
	tomson = New(Config{
		Host:      "host",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.Empty(tomson)

	// Test valid config.
	tomson = New(Config{
		Host:      "https://rms-world-check-one-api-pilot.thomsonreuters.com/v1",
		APIkey:    "key",
		APIsecret: "secret",
	})

	assert.NotEmpty(tomson)
	assert.Equal("https", tomson.scheme)
	assert.Equal("rms-world-check-one-api-pilot.thomsonreuters.com", tomson.host)
	assert.Equal("/v1/", tomson.path)
	assert.Equal("key", tomson.key)
	assert.Equal("secret", tomson.secret)
}

func TestCheckCustomer(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	customer := &common.UserData{}

	res, err := tomson.CheckCustomer(customer)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
