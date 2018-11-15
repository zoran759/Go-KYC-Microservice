package verification

import "modulus/kyc/integrations/trulioo/configuration"

// Config represents the configuration for the service.
type Config struct {
	Host  string
	Token string
}

// Verification defines the interface for the verification services.
type Verification interface {
	Verify(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*Response, error)
}

// Mock represents the mock of the service for tests.
type Mock struct {
	VerifyFn func(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*Response, error)
}

// Verify implements Verification interface for Mock.
func (mock Mock) Verify(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*Response, error) {
	return mock.VerifyFn(countryAlpha2, consents, fields)
}
