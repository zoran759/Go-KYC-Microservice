package strings

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointerize(t *testing.T) {
	testString := "test"
	pointerString := Pointerize(testString)

	if assert.NotNil(t, pointerString) {
		assert.Equal(t, testString, *pointerString)

	}

	assert.Nil(t, Pointerize(""))
}
