package identitymind

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/identitymind/consumer"
)

// Assert that IdentityMind implements the KYCPlatform interface.
var _ common.KYCPlatform = IdentityMind{}

// IdentityMind defines the model for the IdentityMind services.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type IdentityMind struct {
	consumer consumer.Client
}

// New constructs new service object to use with IdentityMind services.
func New(config Config) IdentityMind {
	return IdentityMind{
		consumer: consumer.NewClient(consumer.Config(config)),
	}
}

// CheckCustomer implements CustomerChecker interface for the service.
func (i IdentityMind) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	res, err = i.consumer.CheckCustomer(customer)
	return
}

// CheckStatus implements StatusChecker interface for the service.
func (i IdentityMind) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	res, err = i.consumer.CheckStatus(referenceID)
	return
}
