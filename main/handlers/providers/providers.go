package providers

import (
	"fmt"
	"strconv"
	"sync"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm"
	"modulus/kyc/integrations/complyadvantage"
	"modulus/kyc/integrations/example"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/idology"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/synapsefi"
	"modulus/kyc/integrations/thomsonreuters"
	"modulus/kyc/integrations/trulioo"
	"modulus/kyc/main/config"
)

var mutex sync.RWMutex

// providers keeps keyed list of active providers.
var providers Providers

// Providers represents a keyed list of KYC service providers.
type Providers map[common.KYCProvider]common.KYCPlatform

// ProviderList represents the list of implemented providers.
type ProviderList []common.KYCProvider

// Those methods implement sort.Interface for ProviderList.
func (p ProviderList) Len() int           { return len(p) }
func (p ProviderList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ProviderList) Less(i, j int) bool { return p[i] < p[j] }

// Create makes new keyed list of active providers.
func Create(cfg config.Config) (err error) {
	mutex.Lock()
	defer mutex.Unlock()

	providers = Providers{
		common.Example: &example.Service{},
	}

	for name, opts := range cfg {
		switch common.KYCProvider(name) {
		case common.Coinfirm:
			providers[common.Coinfirm] = coinfirm.New(coinfirm.Config{
				Host:     opts["Host"],
				Email:    opts["Email"],
				Password: opts["Password"],
				Company:  opts["Company"],
			})
		case common.ComplyAdvantage:
			fuzziness, err1 := strconv.ParseFloat(opts["Fuzziness"], 32)
			if err1 != nil {
				err = fmt.Errorf("%s config error: %s", name, err1)
				return
			}
			providers[common.ComplyAdvantage] = complyadvantage.New(complyadvantage.Config{
				Host:      opts["Host"],
				APIkey:    opts["APIkey"],
				Fuzziness: float32(fuzziness),
			})
		case common.IdentityMind:
			providers[common.IdentityMind] = identitymind.New(identitymind.Config{
				Host:     opts["Host"],
				Username: opts["Username"],
				Password: opts["Password"],
			})
		case common.IDology:
			useSummaryResult, err1 := strconv.ParseBool(opts["UseSummaryResult"])
			if err1 != nil {
				err = fmt.Errorf("%s config error: %s", name, err1)
				return
			}
			providers[common.IDology] = idology.New(idology.Config{
				Host:             opts["Host"],
				Username:         opts["Username"],
				Password:         opts["Password"],
				UseSummaryResult: useSummaryResult,
			})
		case common.Jumio:
			providers[common.Jumio] = jumio.New(jumio.Config{
				BaseURL: opts["BaseURL"],
				Token:   opts["Token"],
				Secret:  opts["Secret"],
			})
		case common.ShuftiPro:
			providers[common.ShuftiPro] = shuftipro.New(shuftipro.Config{
				Host:        opts["Host"],
				SecretKey:   opts["SecretKey"],
				ClientID:    opts["ClientID"],
				RedirectURL: opts["RedirectURL"],
			})
		case common.SumSub:
			providers[common.SumSub] = sumsub.New(sumsub.Config{
				Host:   opts["Host"],
				APIKey: opts["APIKey"],
			})
		case common.SynapseFI:
			providers[common.SynapseFI] = synapsefi.New(synapsefi.Config{
				Host:         opts["Host"],
				ClientID:     opts["ClientID"],
				ClientSecret: opts["ClientSecret"],
			})
		case common.ThomsonReuters:
			providers[common.ThomsonReuters] = thomsonreuters.New(thomsonreuters.Config{
				Host:      opts["Host"],
				APIkey:    opts["APIkey"],
				APIsecret: opts["APIsecret"],
			})
		case common.Trulioo:
			providers[common.Trulioo] = trulioo.New(trulioo.Config{
				Host:         opts["Host"],
				NAPILogin:    opts["NAPILogin"],
				NAPIPassword: opts["NAPIPassword"],
			})
		}
	}

	return
}

// Provider returns a provider by the name.
func Provider(provider common.KYCProvider) common.KYCPlatform {
	mutex.RLock()
	defer mutex.RUnlock()

	return providers[provider]
}
