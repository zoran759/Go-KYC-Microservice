package config_test

import (
	"testing"

	"modulus/kyc/main/config"

	"github.com/stretchr/testify/assert"
)

func TestServicePort(t *testing.T) {
	assert := assert.New(t)

	port := config.ServicePort()

	assert.Equal(config.DefaultPort, port)

	updated, errs := config.Update(config.Config{
		config.ServiceSection: config.Options{
			"Port": "8999",
		},
	})

	assert.True(updated)
	assert.Empty(errs)

	port = config.ServicePort()

	assert.Equal("8999", port)
}
