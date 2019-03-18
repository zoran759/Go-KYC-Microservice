package license

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	type testCase struct {
		name string
		lic  string
		err  error
	}

	testCases := []testCase{
		testCase{
			name: "Invalid license",
			lic:  "fake license",
			err:  errors.New("during license validation: LicenseKey is not registered"),
		},
		testCase{
			name: "Valid license",
			lic:  "dev-testing",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Update(tc.lic)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestValid(t *testing.T) {
	valid := Valid()
	assert.False(t, valid)

	SetDevMode()
	valid = Valid()
	assert.True(t, valid)
}
