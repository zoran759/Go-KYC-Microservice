package configuration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Error(t *testing.T) {
	assert.Equal(t, "test", Error{Message: "test"}.Error())
}
