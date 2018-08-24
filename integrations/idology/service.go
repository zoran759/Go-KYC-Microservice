package idology

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Service defines the verification services of IDology API.
// It shouln't be instantiated directly.
// Use New() constructor instead for convenience.
type Service struct {
	ExpectID common.CustomerChecker
	// FIXME: AlertList has to be implemented yet.
}

// New return new verifier to use with IDology services.
func New(config Config) *Service {
	return &Service{
		ExpectID: expectid.NewClient(expectid.Config(config)),
	}
}

// Ensure implementation conformance to the interface.
var _ common.CustomerChecker = (*expectid.Client)(nil)
