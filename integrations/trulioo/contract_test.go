package trulioo

import (
	"testing"

	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/verification"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
)

func TestConfig_createToken(t *testing.T) {
	config := Config{
		NAPILogin:    "test",
		NAPIPassword: "test_test",
	}

	assert.Equal(t, base64.StdEncoding.EncodeToString(
		[]byte(config.NAPILogin+":"+config.NAPIPassword),
	), config.createToken())
}

func TestConfig_ToConfigurationConfig(t *testing.T) {
	config := Config{
		Host:         "http://host.com",
		NAPILogin:    "test",
		NAPIPassword: "test_test",
	}

	assert.Equal(t, configuration.Config{
		Host: "http://host.com/configuration/v1",
		Token: base64.StdEncoding.EncodeToString(
			[]byte(config.NAPILogin + ":" + config.NAPIPassword),
		),
	}, config.ToConfigurationConfig())
}

func TestConfig_ToVerificationConfig(t *testing.T) {
	config := Config{
		Host:         "http://host.com",
		NAPILogin:    "test",
		NAPIPassword: "test_test",
	}

	assert.Equal(t, verification.Config{
		Host: "http://host.com/verifications/v1",
		Token: base64.StdEncoding.EncodeToString(
			[]byte(config.NAPILogin + ":" + config.NAPIPassword),
		),
	}, config.ToVerificationConfig())
}
