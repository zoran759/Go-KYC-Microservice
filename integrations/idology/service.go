package idology

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Service defines the verification services of IDology API.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type Service struct {
	ExpectID common.CustomerChecker
	// FIXME: AlertList has to be implemented yet.
}

// New returns new verifier to use with IDology services.
func New(config Config) *Service {
	return &Service{
		ExpectID: expectid.NewClient(expectid.Config(config)),
	}
}

// Ensure implementation conformance to the interface.
var _ common.CustomerChecker = (*expectid.Client)(nil)
