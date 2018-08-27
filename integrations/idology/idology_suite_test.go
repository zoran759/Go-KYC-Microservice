package idology_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestIDology(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "IDology Suite")
}
