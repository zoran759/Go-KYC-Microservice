package coinfirm

import (
	"testing"

	"modulus/kyc/http"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := Coinfirm{
		host:     "host",
		email:    "email",
		password: "password",
		company:  "company",
		headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}

	tc := New(Config{
		Host:     "host",
		Email:    "email",
		Password: "password",
		Company:  "company",
	})

	assert.Equal(t, c, *tc)
}
