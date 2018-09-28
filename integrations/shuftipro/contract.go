package shuftipro

import "modulus/kyc/integrations/shuftipro/verification"

// Config represents the service config.
type Config verification.Config

// Verification result codes from the API.
const (
	NotVerified = "SP0"
	Verified    = "SP1"
	Success     = "SP2"
)
