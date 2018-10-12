package config

import (
	"encoding/json"
	"fmt"
	"modulus/kyc/common"
	"os"
)

// KYCConfig holds the current config for KYC providers.
// Beware that it isn't concurrent writes safe.
var KYCConfig *Config

// Options represents the configuration options for the KYC provider.
type Options map[string]string

// Config represents the configuration for KYC providers.
type Config map[common.KYCProvider]Options

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

	KYCConfig := &Config{}

	dec := json.NewDecoder(file)
	err = dec.Decode(KYCConfig)

	return
}
