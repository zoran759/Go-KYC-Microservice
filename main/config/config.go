package config

import (
	"sync"

	"modulus/kyc/common"
	"modulus/kyc/main/config/providers"
	"modulus/kyc/main/license"
)

const (
	// ServiceSection is the hardcoded value of the KYC service config section name.
	ServiceSection = "Config"

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
	cfg.Unlock()

	if len(port) == 0 {
		port = DefaultPort
	}
	return
}

// Update updates the config with the options provided.
func Update(c Config) (updated bool, errs []string) {
	if len(c) == 0 {
		errs = append(errs, "no config update data provided")
		return
	}

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
		errs1 := filterOptions(common.KYCProvider(sect), opts)
		if errs1 != nil {
			errs = append(errs, errs1...)
		}
		if len(opts) == 0 {
			continue
		}
		if sect == ServiceSection {
			lic, ok := opts["License"]
			if ok {
				if err := license.Update(lic); err != nil {
					errs = append(errs, err.Error())
				}
			}
		}
		if !NotProviders[sect] {
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
		if !updated {
			updated = true
		}
	}

	if len(updlist) > 0 {
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
	return !providers.Providers[common.KYCProvider(sect)] && !NotProviders[sect]
}

// IsKnownName returns the result of check whether the given name is known to the service.
func IsKnownName(name string) bool {
	return providers.Providers[common.KYCProvider(name)] || (NotProviders[name] && name != ServiceSection)
}

// GetOptions returns the Options of the configuration section specified or nil if the section isn't exist.
func GetOptions(section string) Options {
	cfg.Lock()
	opts := cfg.config[section]
	cfg.Unlock()
	return opts
}

// GetConfig returns current configuration.
func GetConfig() Config {
	cfg.Lock()
	c := cfg.config
	cfg.Unlock()
	return c
}

func init() {
	cfg.filename = DefaultDevFile
	cfg.config = Config{}
}
