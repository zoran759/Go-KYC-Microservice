package model_test

import (
	"testing"

	"modulus/kyc/integrations/coinfirm/model"

	"github.com/stretchr/testify/assert"
)

func TestNewCryptoAddress(t *testing.T) {
	addr := "fakeCryptoAddress"
	ca := model.NewCryptoAddress(addr)

	assert.Equal(t, ca, model.CryptoAddress("address: `"+addr+"`"))
}
