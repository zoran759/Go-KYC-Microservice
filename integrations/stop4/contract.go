package stop4

import "modulus/kyc/integrations/stop4/verification"

// Config represents the service config.
type Config verification.Config

// Verification result codes from the API.
const (
	Success = "0"
)
