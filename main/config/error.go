package config

import (
	"fmt"

	"modulus/kyc/common"
)

// MissingOptionError defines an error of the missing config option.
type MissingOptionError struct {
	provider common.KYCProvider
	option   string
}

// Error implements error interface for MissingOptionError.
func (e MissingOptionError) Error() string {
	return fmt.Sprintf("missing or empty option '%s' for the %s provider", e.option, e.provider)
}

// OptionError represents an option error.
type OptionError struct {
	provider common.KYCProvider
	option   string
	text     string
}

// Error implements error interface for OptionError.
func (e OptionError) Error() string {
	return fmt.Sprintf("%s '%s' option error: %s", e.provider, e.option, e.text)
}

// ParseError represents a config parser error.
type ParseError struct {
	strnum  int
	content string
	err     string
}

// Error implements error interface for ParseError.
func (p ParseError) Error() string {
	return fmt.Sprintf("parsing failed at line %d '%s': %s", p.strnum, p.content, p.err)
}
