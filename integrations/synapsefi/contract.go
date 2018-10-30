package synapsefi

import "modulus/kyc/integrations/synapsefi/verification"

type Connection verification.Config

type Config struct {
	Connection Connection
	TimeoutThreshold int64
	KYCFlow	string
}

const (
	DocStatusMissingOrInvalid = "MISSING|INVALID"
	DocStatusValid = "SUBMITTED|VALID"
	DocStatusInvalid = "SUBMITTED|INVALID"
	DocStatusPending = "SUBMITTED"
	DocStatusReviewing = "SUBMITTED|REVIEWING"
)

