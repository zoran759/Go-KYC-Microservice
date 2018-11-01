package identitymind

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/identitymind/consumer"
)

// Assert that Service implements the CustomerChecker interface.
var _ common.CustomerChecker = (*Service)(nil)

// Assert that Service implements the StatusChecker interface.
var _ common.StatusChecker = (*Service)(nil)

// Service defines the model for the IdentityMind services.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type Service struct {
	consumer *consumer.Client
}

// New constructs new service object to use with IdentityMind services.
func New(config Config) *Service {
	return &Service{
		consumer: consumer.NewClient(consumer.Config(config)),
	}
}

// CheckCustomer implements CustomerChecker interface for the service.
func (s Service) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	res, err = s.consumer.CheckCustomer(customer)

	return
}

// CheckStatus implements StatusChecker interface for the service.
func (s Service) CheckStatus(customerID string) (res common.KYCResult, err error) {
	res, err = s.consumer.CheckStatus(customerID)

	return
}
