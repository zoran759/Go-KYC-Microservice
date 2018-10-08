package trulioo

import (
	"encoding/base64"
	"modulus/kyc/integrations/trulioo/configuration"
	"modulus/kyc/integrations/trulioo/verification"
)

// Config represents the service config.
type Config struct {
	Host         string
	NAPILogin    string
	NAPIPassword string
}

func (config Config) createToken() string {
	return base64.StdEncoding.EncodeToString(
		[]byte(config.NAPILogin + ":" + config.NAPIPassword),
	)
}

// ToConfigurationConfig converts the service config to the specific config required to use for certain requests.
func (config Config) ToConfigurationConfig() configuration.Config {
	return configuration.Config{
		Host:  config.Host + "/configuration/v1",
		Token: config.createToken(),
	}
}

// ToVerificationConfig converts the service config to the specific config required to use for certain requests.
func (config Config) ToVerificationConfig() verification.Config {
	return verification.Config{
		Host:  config.Host + "/verifications/v1",
		Token: config.createToken(),
	}
}

// Result constants.
const (
	Match   = "match"
	NoMatch = "nomatch"
)
