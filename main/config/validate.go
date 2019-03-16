package config

import (
	"fmt"

	"modulus/kyc/common"
)

// knownOptions keeps the collection of valid option names for each config section.
var knownOptions = map[common.KYCProvider]validOptions{
	common.Coinfirm: validOptions{
		"Host":     true,
		"Email":    true,
		"Password": true,
		"Company":  true,
	},
	common.ComplyAdvantage: validOptions{
		"Host":      true,
		"APIkey":    true,
		"Fuzziness": true,
	},
	common.IdentityMind: validOptions{
		"Host":     true,
		"Username": true,
		"Password": true,
	},
	common.IDology: validOptions{
		"Host":             true,
		"Username":         true,
		"Password":         true,
		"UseSummaryResult": true,
	},
	common.Jumio: validOptions{
		"BaseURL": true,
		"Token":   true,
		"Secret":  true,
	},
	common.ShuftiPro: validOptions{
		"Host":        true,
		"ClientID":    true,
		"SecretKey":   true,
		"CallbackURL": true,
	},
	common.SumSub: validOptions{
		"Host":   true,
		"APIKey": true,
	},
	common.SynapseFI: validOptions{
		"Host":         true,
		"ClientID":     true,
		"ClientSecret": true,
	},
	common.ThomsonReuters: validOptions{
		"Host":      true,
		"APIkey":    true,
		"APIsecret": true,
	},
	common.Trulioo: validOptions{
		"Host":         true,
		"NAPILogin":    true,
		"NAPIPassword": true,
	},
	// These aren't KYC providers.
	common.CipherTrace: validOptions{
		"URL":      true,
		"Key":      true,
		"Username": true,
	},
	// These are the options of the KYC service itself.
	common.KYCProvider(ServiceSection): validOptions{
		"Port":    true,
		"License": true,
	},
}

// validOptions represents a list of valid option names.
type validOptions map[string]bool

// filterOptions filters the provided options from unknown ones.
func filterOptions(provider common.KYCProvider, opts Options) (errs []string) {
	kopts := knownOptions[provider]
	if kopts == nil {
		return []string{string(provider) + " is missing configuration validation"}

	}
	for name := range opts {
		if !kopts[name] {
			delete(opts, name)
			errs = append(errs, fmt.Sprintf("%s: unknown option '%s' was filtered out", provider, name))
		}
	}
	return
}

// validateProvider checks the config correctness for the KYC provider.
func validateProvider(provider common.KYCProvider, opts Options) error {
	kopts := knownOptions[provider]
	if kopts == nil {
		return fmt.Errorf("%s is missing configuration validation", provider)
	}
	for name := range kopts {
		if len(opts[name]) == 0 {
			return MissingOptionError{provider: provider, option: name}
		}

	}
	return nil
}
