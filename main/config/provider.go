package config

import (
	"fmt"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm"
	"modulus/kyc/integrations/complyadvantage"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/idology"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/synapsefi"
	"modulus/kyc/integrations/thomsonreuters"
	"modulus/kyc/integrations/trulioo"
)

// createPlatform constructs the new KYC verification platform.
func createPlatform(provider common.KYCProvider) (platform common.KYCPlatform, err error) {
	opts := cfg.config[string(provider)]
	if err = validateProvider(provider, opts); err != nil {
		return
	}

	switch provider {
	case common.Coinfirm:
		platform = coinfirm.New(coinfirm.Config{
			Host:     opts["Host"],
			Email:    opts["Email"],
			Password: opts["Password"],
			Company:  opts["Company"],
		})
	case common.ComplyAdvantage:
		fuzziness, err1 := strconv.ParseFloat(opts["Fuzziness"], 32)
		if err1 != nil {
			err = OptionError{
				provider: provider,
				option:   "Fuzziness",
				text:     err1.Error(),
			}
			return
		}
		platform = complyadvantage.New(complyadvantage.Config{
			Host:      opts["Host"],
			APIkey:    opts["APIkey"],
			Fuzziness: float32(fuzziness),
		})
	case common.IdentityMind:
		platform = identitymind.New(identitymind.Config{
			Host:     opts["Host"],
			Username: opts["Username"],
			Password: opts["Password"],
		})
	case common.IDology:
		useSummaryResult, err1 := strconv.ParseBool(opts["UseSummaryResult"])
		if err1 != nil {
			err = OptionError{
				provider: provider,
				option:   "UseSummaryResult",
				text:     err1.Error(),
			}
			return
		}
		platform = idology.New(idology.Config{
			Host:             opts["Host"],
			Username:         opts["Username"],
			Password:         opts["Password"],
			UseSummaryResult: useSummaryResult,
		})
	case common.Jumio:
		platform = jumio.New(jumio.Config{
			BaseURL: opts["BaseURL"],
			Token:   opts["Token"],
			Secret:  opts["Secret"],
		})
	case common.ShuftiPro:
		platform = shuftipro.New(shuftipro.Config{
			Host:        opts["Host"],
			SecretKey:   opts["SecretKey"],
			ClientID:    opts["ClientID"],
			RedirectURL: opts["RedirectURL"],
		})
	case common.SumSub:
		platform = sumsub.New(sumsub.Config{
			Host:   opts["Host"],
			APIKey: opts["APIKey"],
		})
	case common.SynapseFI:
		platform = synapsefi.New(synapsefi.Config{
			Host:         opts["Host"],
			ClientID:     opts["ClientID"],
			ClientSecret: opts["ClientSecret"],
		})
	case common.ThomsonReuters:
		platform = thomsonreuters.New(thomsonreuters.Config{
			Host:      opts["Host"],
			APIkey:    opts["APIkey"],
			APIsecret: opts["APIsecret"],
		})
	case common.Trulioo:
		platform = trulioo.New(trulioo.Config{
			Host:         opts["Host"],
			NAPILogin:    opts["NAPILogin"],
			NAPIPassword: opts["NAPIPassword"],
		})
	}

	if platform == nil {
		err = fmt.Errorf("unexpected error during creating verification platform for %s provider", provider)
	}

	return
}
