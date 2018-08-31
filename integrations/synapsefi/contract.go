package synapsefi

import "gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"

type Config struct {
	verification.Config
	TimeoutThreshold int64
}

const MissingOrInvalid = "MISSING|INVALID"
const Verified = "SUBMITTED|VERIFIED"
const Unverified = "SUBMITTED|UNVERIFIED"
