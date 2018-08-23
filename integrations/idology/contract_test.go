package idology_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitlab.com/lambospeed/kyc/integrations/idology"
)

var _ = Describe("Contract", func() {
	Describe("Verifier", func() {
		It("should successfully create Verifier instance", func() {
			config := Config{
				Host:     "fake_host",
				Username: "fake_username",
				Password: "fake_password",
			}

			verifier := New(config)

			Expect(verifier).NotTo(BeNil())
			Expect(reflect.TypeOf(verifier)).To(Equal(reflect.TypeOf((*Verifier)(nil))))

			expectID := verifier.ExpectID
			Expect(expectID).ToNot(BeNil())
		})
	})
})
