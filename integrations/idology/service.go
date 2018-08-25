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

// New constructs new service object to use with IDology services.
func New(config Config) *Service {
	return &Service{
		ExpectID: expectid.NewClient(expectid.Config(config)),
	}
}
