package model_test

import (
	"testing"

	"modulus/kyc/integrations/coinfirm/model"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	e := "test error response"

	err := model.ErrorResponse{
		Err: e,
	}

	assert.Equal(t, e, err.Error())
}
