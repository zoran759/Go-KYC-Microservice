package numbers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointerizeInt(t *testing.T) {
	assert.Equal(t, 100, *PointerizeInt(100))
	assert.Nil(t, PointerizeInt(0))
}
