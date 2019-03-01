package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"

	"modulus/kyc/common"
)

// These are config keywords.
const (
	comment   = '#'
	namestart = '['
	namestop  = ']'
	sep       = '='
)

// List of kind values.
const (
	iscomment kind = iota
	isname
	isopt
	iserror
)

// kind represents a kind of read string.
type kind int

// parseConfig reads string by string from the input and parses it
// into valid Config or returns an error if occured.
func parseConfig(r io.Reader) (Config, error) {
	if r == nil {
		return nil, errors.New("the config source is nil")
	}

	scanner := bufio.NewScanner(r)

	cfg := Config{}
	opts := Options{}
	name := ""
	s := ""
	count := 0
	for scanner.Scan() {
		count++
		s = strings.TrimSpace(scanner.Text())
		if len(s) == 0 {
			continue
		}
		switch kindOf(s) {
		case iscomment:
			continue
		case isname:
			// We will omit empty sections because it doesn't make sense to keep them.
			if len(name) > 0 && len(opts) > 0 {
				cfg[name] = opts
				opts = Options{}
			}
			name = s[1 : len(s)-1]
			if err := validateName(name); err != nil {
				err := ParseError{
					strnum:  count,
					content: scanner.Text(),
					err:     err.Error(),
				}
				return nil, err
			}
		case isopt:
			if len(name) == 0 {
				err := ParseError{
					strnum:  count,
					content: scanner.Text(),
					err:     "standalone option string",
				}
				return nil, err
			}
			i := bytes.IndexByte([]byte(s), sep)
			key := s[:i]
			val := ""
			if i < len(s)-1 {
				val = s[i+1:]
			}
			opts[key] = val
		case iserror:
			err := ParseError{
				strnum:  count,
				content: scanner.Text(),
				err:     "not proper config string",
			}
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		err := ParseError{
			strnum:  count,
			content: scanner.Text(),
			err:     err.Error(),
		}
		return nil, err
	}
	if len(name) == 0 {
		err := ParseError{
			strnum:  count,
			content: scanner.Text(),
			err:     "config is empty",
		}
		return nil, err
	}
	cfg[name] = opts

	return cfg, nil
}

func kindOf(s string) kind {
	if s[0] == comment {
		return iscomment
	}
	if s[0] == namestart {
		if s[len(s)-1] != namestop {
			return iserror
		}
		return isname
	}
	if bytes.IndexByte([]byte(s), sep) > 0 {
		return isopt
	}

	return iserror
}

// validateName validates KYC provider name from a config.
func validateName(name string) (err error) {
	if len(name) == 0 {
		err = errors.New("empty section name")
		return err
	}
	if name != ServiceSection && !common.KYCProviders[common.KYCProvider(name)] {
		err = errors.New("unknown KYC provider name in the config")
		return err
	}

	return
}
