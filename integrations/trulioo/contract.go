package trulioo

import (
	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/verification"
	"encoding/base64"
)

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

func (config Config) ToConfigurationConfig() configuration.Config {
	return configuration.Config{
		Host:  config.Host + "/configuration/v1",
		Token: config.createToken(),
	}
}

func (config Config) ToVerificationConfig() verification.Config {
	return verification.Config{
		Host:  config.Host + "/verifications/v1",
		Token: config.createToken(),
	}
}

const Match = "match"
const NoMatch = "nomatch"
