package verification

// Config represents the configuration for the service.
type Config struct {
	Host       string
	MerchantID string
	Password   string
}

// Verification defines the interface for the verification services.
type Verification interface {
	Verify(request RegistrationRequest) (*Response, error)
}

// Mock represents the mock of the service for tests.
type Mock struct {
	VerifyFn func(request RegistrationRequest) (*Response, error)
}

// Verify implements Verification interface for Mock.
func (mock Mock) Verify(request RegistrationRequest) (*Response, error) {
	return mock.VerifyFn(request)
}
