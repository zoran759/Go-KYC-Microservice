package config_test

import (
	"testing"

	"modulus/kyc/main/config"

	"github.com/stretchr/testify/assert"
)

var c = config.Config{
	"Foo": config.Options{
		"Bar": "bar option",
		"Baz": "baz option",
	},
}

func TestOption(t *testing.T) {
	assert := assert.New(t)

	opt := c.Option("NonExistent", "option")
	assert.Empty(opt)

	opt = c.Option("Foo", "option")
	assert.Empty(opt)

	opt = c.Option("Foo", "Bar")
	assert.Equal("bar option", opt)

	opt = c.Option("Foo", "Baz")
	assert.Equal("baz option", opt)
}

func TestServicePort(t *testing.T) {
	assert := assert.New(t)

	cfg := c

	port := cfg.ServicePort()
	assert.Equal(config.DefaultPort, port)

	cfg[config.ServiceSection] = config.Options{
		"Port": "8999",
	}

	port = cfg.ServicePort()
	assert.Equal("8999", port)
}
