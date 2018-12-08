package synapsefi

import "modulus/kyc/integrations/synapsefi/verification"

// Config represents the service config.
type Config verification.Config

// Possible document status values.
const (
	DocStatusMissingOrInvalid = "MISSING|INVALID"
	DocStatusValid            = "SUBMITTED|VALID"
	DocStatusInvalid          = "SUBMITTED|INVALID"
	DocStatusPending          = "SUBMITTED"
	DocStatusReviewing        = "SUBMITTED|REVIEWING"
)
