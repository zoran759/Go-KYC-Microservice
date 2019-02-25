package config

import "modulus/kyc/common"

// validateProvider checks the config correctness for the KYC provider.
func validateProvider(provider common.KYCProvider, opts Options) error {
	switch provider {
	case common.Coinfirm:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["Email"]) == 0 {
			return MissingOptionError{provider: provider, option: "Email"}
		}
		if len(opts["Password"]) == 0 {
			return MissingOptionError{provider: provider, option: "Password"}
		}
		if len(opts["Company"]) == 0 {
			return MissingOptionError{provider: provider, option: "Company"}
		}
	case common.ComplyAdvantage:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["APIkey"]) == 0 {
			return MissingOptionError{provider: provider, option: "APIkey"}
		}
		if len(opts["Fuzziness"]) == 0 {
			return MissingOptionError{provider: provider, option: "Fuzziness"}
		}
	case common.IdentityMind:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["Username"]) == 0 {
			return MissingOptionError{provider: provider, option: "Username"}
		}
		if len(opts["Password"]) == 0 {
			return MissingOptionError{provider: provider, option: "Password"}
		}
	case common.IDology:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["Username"]) == 0 {
			return MissingOptionError{provider: provider, option: "Username"}
		}
		if len(opts["Password"]) == 0 {
			return MissingOptionError{provider: provider, option: "Password"}
		}
		if len(opts["UseSummaryResult"]) == 0 {
			return MissingOptionError{provider: provider, option: "UseSummaryResult"}
		}
	case common.Jumio:
		if len(opts["BaseURL"]) == 0 {
			return MissingOptionError{provider: provider, option: "BaseURL"}
		}
		if len(opts["Token"]) == 0 {
			return MissingOptionError{provider: provider, option: "Token"}
		}
		if len(opts["Secret"]) == 0 {
			return MissingOptionError{provider: provider, option: "Secret"}
		}
	case common.ShuftiPro:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["SecretKey"]) == 0 {
			return MissingOptionError{provider: provider, option: "SecretKey"}
		}
		if len(opts["ClientID"]) == 0 {
			return MissingOptionError{provider: provider, option: "ClientID"}
		}
		if len(opts["RedirectURL"]) == 0 {
			return MissingOptionError{provider: provider, option: "RedirectURL"}
		}
	case common.SumSub:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["APIKey"]) == 0 {
			return MissingOptionError{provider: provider, option: "APIKey"}
		}
	case common.SynapseFI:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["ClientID"]) == 0 {
			return MissingOptionError{provider: provider, option: "ClientID"}
		}
		if len(opts["ClientSecret"]) == 0 {
			return MissingOptionError{provider: provider, option: "ClientSecret"}
		}
	case common.ThomsonReuters:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["APIkey"]) == 0 {
			return MissingOptionError{provider: provider, option: "APIkey"}
		}
		if len(opts["APIsecret"]) == 0 {
			return MissingOptionError{provider: provider, option: "APIsecret"}
		}
	case common.Trulioo:
		if len(opts["Host"]) == 0 {
			return MissingOptionError{provider: provider, option: "Host"}
		}
		if len(opts["NAPILogin"]) == 0 {
			return MissingOptionError{provider: provider, option: "NAPILogin"}
		}
		if len(opts["NAPIPassword"]) == 0 {
			return MissingOptionError{provider: provider, option: "NAPIPassword"}
		}
	}
	return nil
}
