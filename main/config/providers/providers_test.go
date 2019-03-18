package providers

import (
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/sumsub"

	"github.com/stretchr/testify/assert"
)

func TestStorePlatform(t *testing.T) {
	s := sumsub.New(sumsub.Config{
		Host:   "host",
		APIKey: "key",
	})

	StorePlatform(common.SumSub, s)

	p, err := GetPlatform(common.SumSub)
	s2 := p.(sumsub.SumSub)

	assert.Equal(t, s, s2)
	assert.NoError(t, err)
}

func TestGetPlatform(t *testing.T) {
	assert := assert.New(t)

	s := jumio.New(jumio.Config{
		BaseURL: "url",
		Token:   "token",
		Secret:  "secret",
	})

	StorePlatform(common.Jumio, s)

	p, err := GetPlatform(common.Jumio)
	s2 := p.(jumio.Jumio)

	assert.Equal(s, s2)
	assert.NoError(err)

	p, err = GetPlatform(common.Coinfirm)

	assert.Nil(p)
	assert.Error(err)
	assert.Equal("the provider 'Coinfirm' is unknown or not configured in the service", err.Error())
}
