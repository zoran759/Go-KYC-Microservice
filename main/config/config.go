package config

const (
	// ServiceSection is the hardcoded value of the KYC service config section name.
	ServiceSection = "Config"

	// DefaultPort is default port of the KYC service.
	DefaultPort = "8080"
)

// Cfg holds the current config for the KYC service.
// Beware that it isn't concurrent writes safe.
var Cfg Config

// Options represents the configuration options for the KYC provider.
type Options map[string]string

// Config represents the configuration for the service.
type Config map[string]Options

// Option tries to retrieve requested option from the specified config section.
// If the option doesn't exist the empty value is returned.
func (c Config) Option(section, option string) (opt string) {
	s, ok := c[section]
	if !ok {
		return
	}

	opt = s[option]

	return
}

// ServicePort returns the KYC service port.
func (c Config) ServicePort() (port string) {
	if port = c.Option(ServiceSection, "Port"); len(port) == 0 {
		port = DefaultPort
	}
	return
}
