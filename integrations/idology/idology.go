package idology

import (
	"errors"
	"modulus/kyc/common"
	"modulus/kyc/integrations/idology/expectid"
)

// Assert that Service implements the CustomerChecker interface.
var _ common.KYCPlatform = IDology{}

// IDology represents the IDology API services.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type IDology struct {
	expectID expectid.Client
}

// New constructs new service object to use with IDology services.
func New(config Config) IDology {
	return IDology{
		expectID: expectid.NewClient(expectid.Config(config)),
	}
}

// CheckCustomer implements KYCPlatform interface for the IDology.
func (i IDology) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	res, err = i.expectID.CheckCustomer(customer)
	return
}

// CheckStatus implements KYCPlatform interface for the IDology.
func (i IDology) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	err = errors.New("IDology doesn't support a verification status check")
	return
}
