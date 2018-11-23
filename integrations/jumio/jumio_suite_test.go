package jumio_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJumio(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jumio Suite")
}
