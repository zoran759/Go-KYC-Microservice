package model_test

import (
	"fmt"
	"testing"

	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	errs := model.Errors{
		model.ErrorEntity{
			Error:    "error1",
			Cause:    "cause1",
			ObjectID: "objId1",
		},
		model.ErrorEntity{
			Error:    "error2",
			Cause:    "cause2",
			ObjectID: "objId2",
		},
		model.ErrorEntity{
			Error:    "error3",
			Cause:    "cause3",
			ObjectID: "objId3",
		},
	}

	s := fmt.Sprintf("%v", errs)

	assert.Equal(t, "error1 (cause1) objId1 | error2 (cause2) objId2 | error3 (cause3) objId3", s)
}
