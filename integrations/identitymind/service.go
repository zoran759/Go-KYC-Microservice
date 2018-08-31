package identitymind

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/identitymind/consumer"
)

// Service defines the model for the IdentityMind services.
// It shouldn't be instantiated directly.
// Use New() constructor instead.
type Service struct {
	ConsumerKYC common.CustomerChecker
}

// New constructs new service object to use with IdentityMind services.
func New(config Config) *Service {
	return &Service{
		ConsumerKYC: consumer.NewClient(consumer.Config(config)),
	}
}
