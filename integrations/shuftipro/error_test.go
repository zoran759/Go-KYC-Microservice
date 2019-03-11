package shuftipro

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	type testCase struct {
		name   string
		err    Error
		result string
	}

	testCases := []testCase{
		testCase{
			name:   "Empty error",
			err:    Error{},
			result: "",
		},
		testCase{
			name: "Full error",
			err: Error{
				Service: "Foo",
				Key:     "Bar",
				Message: "Foobar test",
			},
			result: "service: 'Foo' | key: 'Bar' | message: 'Foobar test'",
		},
		testCase{
			name: "Only service field",
			err: Error{
				Service: "Foo",
				Key:     "",
				Message: "",
			},
			result: "service: 'Foo'",
		},
		testCase{
			name: "Only key field",
			err: Error{
				Service: "",
				Key:     "Bar",
				Message: "",
			},
			result: "key: 'Bar'",
		},
		testCase{
			name: "Only message field",
			err: Error{
				Service: "",
				Key:     "",
				Message: "Foobar test",
			},
			result: "message: 'Foobar test'",
		},
		testCase{
			name: "Partial error",
			err: Error{
				Service: "",
				Key:     "Bar",
				Message: "Foobar test",
			},
			result: "key: 'Bar' | message: 'Foobar test'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.result, tc.err.Error())
		})
	}
}
