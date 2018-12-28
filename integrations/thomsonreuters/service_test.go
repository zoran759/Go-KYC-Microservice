package thomsonreuters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// Test URL parsing error.
	svc := New(Config{
		Host:      "::",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s := svc.(service)

	assert.Empty(s)

	// Test malformed Host.
	svc = New(Config{
		Host:      "host",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s = svc.(service)

	assert.Empty(s)

	// Test valid config.
	svc = New(Config{
		Host:      "https://rms-world-check-one-api-pilot.thomsonreuters.com/v1",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s = svc.(service)

	assert.NotEmpty(s)
	assert.Equal("https", s.scheme)
	assert.Equal("rms-world-check-one-api-pilot.thomsonreuters.com", s.host)
	assert.Equal("/v1/", s.path)
	assert.Equal("key", s.key)
	assert.Equal("secret", s.secret)
}
