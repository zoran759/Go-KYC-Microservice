package providers

import (
	"testing"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	assert := assert.New(t)

	p0 := common.Trulioo
	p1 := common.Coinfirm
	l := List{p0, p1}

	n := l.Len()
	assert.Equal(2, n)

	l.Swap(0, 1)
	assert.Equal(p1, l[0])
	assert.Equal(p0, l[1])

	less := l.Less(0, 1)
	assert.True(less)

	less = l.Less(1, 0)
	assert.False(less)
}
