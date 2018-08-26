package expectid

import (
	"flag"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// It is here for compatibility with recursive running of tests for IDology.
var runLive = flag.Bool("runlive", false, "Don't use this flag in this context.")

func TestExpectID(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ExpectID Suite")
}
