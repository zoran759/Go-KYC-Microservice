package verification

// Config represents the configuration for the service.
type Config struct {
	Host        string
	SecretKey   string
	ClientID    string
	RedirectURL string
	CallbackURL string
}

// Verification defines the interface for the verification services.
type Verification interface {
	Verify(request OldRequest) (*Response, error)
}

// Mock represents the mock of the service for tests.
type Mock struct {
	VerifyFn func(request OldRequest) (*Response, error)
}

// Verify implements Verification interface for Mock.
func (mock Mock) Verify(request OldRequest) (*Response, error) {
	return mock.VerifyFn(request)
}
