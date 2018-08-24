package idology

import (
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Config holds configuration settings for the verifiers.
type Config struct {
	Host,
	Username,
	Password string
}

var _ *Config = (*Config)(&expectid.Config{})
