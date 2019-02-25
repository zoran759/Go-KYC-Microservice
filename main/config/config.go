package config

import (
	"modulus/kyc/common"
	"modulus/kyc/main/config/providers"
	"sync"
)

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

// config represents the global configuration for the service.
type config struct {
	sync.RWMutex
	filename string
	config   Config
}

// Config represents a config.
type Config map[string]Options

// Options represents configuration options.
type Options map[string]string

// ServicePort returns the KYC service port.
func ServicePort() (port string) {
	cfg.Lock()
	opts := cfg.config[ServiceSection]
	port = opts["Port"]
	if len(port) == 0 {
		port = DefaultPort
	}
	cfg.Unlock()
	return
}

// Update updates the config with the options provided.
func Update(c Config) (updated bool, errs []string) {
	cfg.Lock()
	defer cfg.Unlock()

	updlist := []common.KYCProvider{}

	for sect, opts := range c {
		if unknownSection(sect) {
			errs = append(errs, "unknown config section: "+sect)
			continue
		}
		if len(opts) == 0 {
			errs = append(errs, "empty config section: "+sect)
			continue
		}
		if sect != ServiceSection {
			updlist = append(updlist, common.KYCProvider(sect))
		}
		oo := cfg.config[sect]
		if oo == nil {
			oo = Options{}
		}
		for o, v := range opts {
			oo[o] = v
		}
		cfg.config[sect] = oo
	}

	if len(updlist) > 0 {
		updated = true
		for _, p := range updlist {
			platform, err := createPlatform(p)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			providers.StorePlatform(p, platform)
		}
	}

	return
}

// unknownSection returns the result of check whether the given section name is unknown to the service.
func unknownSection(sect string) bool {
	return !providers.Providers[common.KYCProvider(sect)] && sect != ServiceSection
}

func init() {
	cfg.filename = DefaultDevFile
	cfg.config = Config{}
}
