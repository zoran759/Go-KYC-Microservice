package providers

import "modulus/kyc/common"

// ProviderList represents the list of implemented providers.
type ProviderList []common.KYCProvider

// Those methods implement sort.Interface for ProviderList.
func (p ProviderList) Len() int           { return len(p) }
func (p ProviderList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ProviderList) Less(i, j int) bool { return p[i] < p[j] }
