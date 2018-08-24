package idology

import (
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Config holds configuration settings for the verifiers.
type Config struct {
	Host             string
	Username         string
	Password         string
	UseSummaryResult bool
	UseAlertList     bool
}

var _ *Config = (*Config)(&expectid.Config{})
