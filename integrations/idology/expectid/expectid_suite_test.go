package expectid

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExpectID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExpectID Suite")
}
