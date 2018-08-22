package verification

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrors_Error(t *testing.T) {
	errors := Errors{
		{
			Message: "message",
		},
		{
			Message: "message1",
		},
		{
			Message: "",
		},
	}

	assert.Equal(t, "message;message1;", errors.Error())
}
