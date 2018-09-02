package synapsefi

import "gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"

type Config struct {
	verification.Config
	TimeoutThreshold int64
}

const MissingOrInvalid = "MISSING|INVALID"
const Valid = "SUBMITTED|VALID"
const Invalid = "SUBMITTED|INVALID"
