package config

import (
	"modulus/kyc/common"
)

// KYC holds the current config for KYC providers.
// Beware that it isn't concurrent writes safe.
var KYC Config

// Options represents the configuration options for the KYC provider.
type Options map[string]string

// Config represents the configuration for KYC providers.
type Config map[common.KYCProvider]Options
