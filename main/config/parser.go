package config

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"strings"
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
	issecterror
)

// kind represents a kind of read string.
type kind int

var parserLogPrefix = "[config parser]"

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
			// Comments are simply skipped.
			continue
		case isname:
			// We will omit empty sections because it doesn't make sense to keep them.
			// If the section is wrong it will be skipped with all belonged options.
			addSectionIfValid(cfg, name, &opts)
			name = s[1 : len(s)-1]
		case isopt:
			// Skip standalone options.
			if len(name) == 0 {
				log.Printf("%s missing section name for the option at line %d: '%s'\n", parserLogPrefix, count, s)
				continue
			}
			i := bytes.IndexByte([]byte(s), sep)
			key := s[:i]
			val := ""
			if i < len(s)-1 {
				val = s[i+1:]
			}
			opts[key] = val
		case issecterror:
			// If line seems like malformed section name,
			// save current section to prevent options mixing.
			addSectionIfValid(cfg, name, &opts)
			name = ""
		case iserror:
			// Skip errors, make the parser fault tolerant.
			log.Printf("%s error at line %d: '%s'\n", parserLogPrefix, count, s)
			continue
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

	addSectionIfValid(cfg, name, &opts)

	if len(cfg) == 0 {
		return nil, errors.New("config doesn't contains a proper configuration")
	}

	return cfg, nil
}

// kindOf determines what kind is the string.
func kindOf(s string) kind {
	if s[0] == comment {
		return iscomment
	}
	if s[0] == namestart {
		if s[len(s)-1] != namestop {
			return issecterror
		}
		return isname
	}
	if bytes.IndexByte([]byte(s), sep) > 0 {
		return isopt
	}

	return iserror
}

// validName checks if the KYC provider name is valid.
func validName(name string) bool {
	return !unknownSection(name)
}

// addSectionIfValid adds new section into the configuration if all provided data is acceptable.
func addSectionIfValid(conf Config, name string, opts *Options) {
	if len(*opts) == 0 {
		return
	}
	if validName(name) {
		conf[name] = *opts
	} else {
		log.Printf("%s unknown section name '%s'\n", parserLogPrefix, name)
	}
	*opts = Options{}
}
