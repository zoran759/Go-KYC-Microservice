package providers

import (
	"fmt"
	"sync"

	"modulus/kyc/common"
	"modulus/kyc/integrations/example"
)

// pvs keeps the pool of active providers.
var providers activeProviders

// activeProviders represents the pool of active providers.
type activeProviders struct {
	sync.RWMutex
	pool pool
}

// pool represents a pool of KYC providers.
type pool map[common.KYCProvider]common.KYCPlatform

// Providers enumerates the implemented KYC providers.
var Providers = map[common.KYCProvider]bool{
	common.Coinfirm:        true,
	common.ComplyAdvantage: true,
	common.IdentityMind:    true,
	common.IDology:         true,
	common.Jumio:           true,
	common.ShuftiPro:       true,
	common.SumSub:          true,
	common.SynapseFI:       true,
	common.ThomsonReuters:  true,
	common.Trulioo:         true,
}

// StorePlatform adds or replaces a provider in the pool.
func StorePlatform(name common.KYCProvider, provider common.KYCPlatform) {
	providers.Lock()
	providers.pool[name] = provider
	providers.Unlock()
}

// GetPlatform returns an active provider from the pool or error if occured.
func GetPlatform(provider common.KYCProvider) (common.KYCPlatform, error) {
	providers.Lock()
	defer providers.Unlock()

	p, ok := providers.pool[provider]
	if !ok {
		return nil, fmt.Errorf("the provider '%s' is unknown or not configured in the service", provider)
	}
	return p, nil
}

func init() {
	providers.pool = pool{
		common.Example: example.Example{},
	}
}
