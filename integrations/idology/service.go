package idology

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/idology/expectid"
)

// Assert that Service implements the CustomerChecker interface.
var _ common.CustomerChecker = (*Service)(nil)

// Service represents the IDology API services.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type Service struct {
	expectID *expectid.Client
}

// New constructs new service object to use with IDology services.
func New(config Config) *Service {
	return &Service{
		expectID: expectid.NewClient(expectid.Config(config)),
	}
}

// CheckCustomer implements CustomerChecker interface for the Service
func (s *Service) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	res, err = s.expectID.CheckCustomer(customer)

	return
}
