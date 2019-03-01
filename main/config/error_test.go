package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrMissingOption(t *testing.T) {
	err := ErrMissingOption{
		provider: "Foobar",
		option:   "Barbaz",
	}

	text := "Foobar configuration error: missing or empty option 'Barbaz'"

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
