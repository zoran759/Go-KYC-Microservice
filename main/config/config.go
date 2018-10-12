package config

import (
	"encoding/json"
	"fmt"
	"os"

	"modulus/kyc/common"
)

// KYC holds the current config for KYC providers.
// Beware that it isn't concurrent writes safe.
var KYC Config

// Options represents the configuration options for the KYC provider.
type Options map[string]string

// Config represents the configuration for KYC providers.
type Config map[common.KYCProvider]Options

// ErrMissingOption defines the missing option error.
type ErrMissingOption struct {
	provider common.KYCProvider
	option   string
}

// Error implements error interface for ErrMissingOption.
func (e ErrMissingOption) Error() string {
	return fmt.Sprintf("%s configuration error: missing or empty option %q", e.provider, e.option)
}

// FromFile loads the configuration from the file.
func FromFile(filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return
	}
	if info.Size() == 0 {
		err = fmt.Errorf("empty %s", filename)
		return
	}

	KYC := Config{}
	dec := json.NewDecoder(file)
	if err = dec.Decode(&KYC); err != nil {
		return
	}

	err = validate(KYC)

	return
}

func validate(config Config) (err error) {
	for provider, options := range config {
		switch provider {
		case common.IDology:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["Username"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Username"}
			}
			if len(options["Password"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Password"}
			}
			if len(options["UseSummaryResult"]) == 0 {
				return ErrMissingOption{provider: provider, option: "UseSummaryResult"}
			}
		case common.ShuftiPro:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["SecretKey"]) == 0 {
				return ErrMissingOption{provider: provider, option: "SecretKey"}
			}
			if len(options["ClientID"]) == 0 {
				return ErrMissingOption{provider: provider, option: "ClientID"}
			}
			if len(options["RedirectURL"]) == 0 {
				return ErrMissingOption{provider: provider, option: "RedirectURL"}
			}
		case common.SumSub:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["APIKey"]) == 0 {
				return ErrMissingOption{provider: provider, option: "APIKey"}
			}
		case common.Trulioo:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["NAPILogin"]) == 0 {
				return ErrMissingOption{provider: provider, option: "NAPILogin"}
			}
			if len(options["NAPIPassword"]) == 0 {
				return ErrMissingOption{provider: provider, option: "NAPIPassword"}
			}
		}
	}

	return
}
