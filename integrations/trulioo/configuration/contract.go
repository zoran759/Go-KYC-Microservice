package configuration

// Config represents the configuration for the configuration provider.
type Config struct {
	Host  string
	Token string
}

// Configuration represents the configuration interface.
type Configuration interface {
	Consents(countryAlpha2 string) (Consents, *int, error)
}

// Mock represents the mock for the configuration provider.
type Mock struct {
	ConsentsFn func(countryAlpha2 string) (Consents, *int, error)
}

// Consents implements the Configuration interface for the Mock.
func (mock Mock) Consents(countryAlpha2 string) (Consents, *int, error) {
	return mock.ConsentsFn(countryAlpha2)
}
