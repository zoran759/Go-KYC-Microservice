package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingOptionError(t *testing.T) {
	err := MissingOptionError{
		provider: "Foobar",
		option:   "Barbaz",
	}

	text := "missing or empty option 'Barbaz' for the Foobar"

	assert.Equal(t, text, err.Error())
}

func TestParseError(t *testing.T) {
	err := ParseError{
		strnum:  17,
		content: "[Foobar",
		err:     "not proper config string",
	}

	text := "parsing failed at line 17 '[Foobar': not proper config string"

	assert.Equal(t, text, err.Error())
}
