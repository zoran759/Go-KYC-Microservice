package config

import "modulus/kyc/common"

// validate ensures the config correctness for all KYC providers containing in the given config.
func validate(config Config) (err error) {
	for provider, options := range config {
		switch provider {
		case common.ComplyAdvantage:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["APIkey"]) == 0 {
				return ErrMissingOption{provider: provider, option: "APIkey"}
			}
			if len(options["Fuzziness"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Fuzziness"}
			}
		case common.IdentityMind:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["Username"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Username"}
			}
			if len(options["Password"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Password"}
			}
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
		case common.Jumio:
			if len(options["BaseURL"]) == 0 {
				return ErrMissingOption{provider: provider, option: "BaseURL"}
			}
			if len(options["Token"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Token"}
			}
			if len(options["Secret"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Secret"}
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
		case common.SynapseFI:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["ClientID"]) == 0 {
				return ErrMissingOption{provider: provider, option: "ClientID"}
			}
			if len(options["ClientSecret"]) == 0 {
				return ErrMissingOption{provider: provider, option: "ClientSecret"}
			}
		case common.ThomsonReuters:
			if len(options["Host"]) == 0 {
				return ErrMissingOption{provider: provider, option: "Host"}
			}
			if len(options["APIkey"]) == 0 {
				return ErrMissingOption{provider: provider, option: "APIkey"}
			}
			if len(options["APIsecret"]) == 0 {
				return ErrMissingOption{provider: provider, option: "APIsecret"}
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
