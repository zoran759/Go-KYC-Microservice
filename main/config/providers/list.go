package providers

import "modulus/kyc/common"

// List represents the list of implemented providers.
type List []common.KYCProvider

// Those methods implement sort.Interface for ProviderList.
func (l List) Len() int           { return len(l) }
func (l List) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l List) Less(i, j int) bool { return l[i] < l[j] }
