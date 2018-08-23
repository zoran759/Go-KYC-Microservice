package idology

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Config holds configuration settings for the verifiers.
type Config struct {
	Host,
	Username,
	Password string
}

// Verifier defines the verifier of IDology services.
// It shouln't be instantiated directly.
// Use New() constructor instead for convenience.
type Verifier struct {
	ExpectID CustomerChecker
	// FIXME: AlertList has to be implemented yet.
	AlertList CustomerChecker
}

// New return new verifier to use with IDology services.
func New(config Config) *Verifier {
	return &Verifier{
		ExpectID: expectid.NewClient(expectid.Config(config)),
	}
}

// CustomerChecker defines the customer verification interface for IDology.
type CustomerChecker interface {
	CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error)
}

var _ *Config = (*Config)(&expectid.Config{})
var _ CustomerChecker = (*expectid.Client)(nil)
