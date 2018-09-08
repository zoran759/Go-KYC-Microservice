package idology

import (
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

const (
	// KYCendpoint holds IDology ExpectID® API Endpoint.
	KYCendpoint = "https://web.idologylive.com/api/idiq.svc"
)

// Config holds configuration settings for the verifiers.
type Config struct {
	Host             string
	Username         string
	Password         string
	UseSummaryResult bool
}

var _ *Config = (*Config)(&expectid.Config{})
