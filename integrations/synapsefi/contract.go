package synapsefi

import "modulus/kyc/integrations/synapsefi/verification"

// Connection represents the config of verification.
type Connection verification.Config

// Config represents the service config.
type Config struct {
	Connection
	TimeoutThreshold int64
	KYCFlow          string
}

// Possible document status values.
const (
	DocStatusMissingOrInvalid = "MISSING|INVALID"
	DocStatusValid            = "SUBMITTED|VALID"
	DocStatusInvalid          = "SUBMITTED|INVALID"
	DocStatusPending          = "SUBMITTED"
	DocStatusReviewing        = "SUBMITTED|REVIEWING"
)
