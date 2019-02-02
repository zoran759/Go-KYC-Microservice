package config

import "sync"

const (
	// ServiceSection is the hardcoded value of the KYC service config section name.
	ServiceSection = "Service"

	// DefaultPort is default port of the KYC service.
	DefaultPort = "8080"

	// DefaultFile is the default config file if no provided.
	DefaultFile = "kyc.cfg"

	// DefaultDevFile is the default config file for the development environment if no provided.
	DefaultDevFile = "kyc_dev.cfg"
)

// cfg holds the current config for the KYC service.
var cfg config

// config represents the configuration for the service.
type config struct {
	sync.RWMutex
	filename string
	config   privconfig
}

type privconfig map[string]Options

// Options represents the configuration options for the KYC provider.
type Options map[string]string

// Option tries to retrieve requested option from the specified config section.
// If the option doesn't exist the empty value is returned.
func Option(section, option string) (opt string) {
	cfg.Lock()
	s, ok := cfg.config[section]
	if ok {
		opt = s[option]
	}
	cfg.Unlock()
	return
}

// ServicePort returns the KYC service port.
func ServicePort() (port string) {
	if port = Option(ServiceSection, "Port"); len(port) == 0 {
		port = DefaultPort
	}
	return
}

// Update updates the provided section in the config with the provided options.
// If the section doesn't exist it will be created.
func Update(section string, opts Options) {
	cfg.Lock()
	copts := cfg.config[section]
	switch copts == nil {
	case true:
		cfg.config[section] = opts
	case false:
		for k, v := range opts {
			copts[k] = v
		}
		cfg.config[section] = copts
	}
	cfg.Unlock()
}

func init() {
	cfg.filename = DefaultDevFile
	cfg.config = privconfig{}
}
