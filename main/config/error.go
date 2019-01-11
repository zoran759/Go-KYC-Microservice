package config

import (
	"fmt"
	"modulus/kyc/common"
)

// ErrMissingOption defines an error of the missing config option.
type ErrMissingOption struct {
	provider common.KYCProvider
	option   string
}

// Error implements error interface for ErrMissingOption.
func (e ErrMissingOption) Error() string {
	return fmt.Sprintf("%s configuration error: missing or empty option '%s'", e.provider, e.option)
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
