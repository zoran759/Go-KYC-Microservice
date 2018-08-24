package idology_test

import (
	"reflect"

	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitlab.com/lambospeed/kyc/integrations/idology"
)

var _ = Describe("Service", func() {
	Specify("should be properly created", func() {
		config := Config{
			Host:     "fake_host",
			Username: "fake_username",
			Password: "fake_password",
		}

		service := &Service{
			ExpectID: expectid.NewClient(expectid.Config(config)),
		}

		testservice := New(config)

		Expect(testservice).NotTo(BeNil())
		Expect(reflect.TypeOf(testservice)).To(Equal(reflect.TypeOf((*Service)(nil))))

		expectID := testservice.ExpectID
		Expect(expectID).ToNot(BeNil())

		Expect(testservice).To(Equal(service))
	})
})
