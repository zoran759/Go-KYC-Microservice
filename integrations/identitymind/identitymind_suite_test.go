package identitymind_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIdentityMind(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IdentityMind Suite")
}
