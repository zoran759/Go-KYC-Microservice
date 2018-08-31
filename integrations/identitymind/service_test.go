package identitymind

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/lambospeed/kyc/integrations/identitymind/consumer"
)

var _ = Describe("The IdentityMind service", func() {
	Specify("should be properly created", func() {
		config := Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		}

		service := &Service{
			ConsumerKYC: consumer.NewClient(consumer.Config(config)),
		}

		testservice := New(config)

		Expect(testservice).NotTo(BeNil())
		Expect(reflect.TypeOf(testservice)).To(Equal(reflect.TypeOf((*Service)(nil))))

		testconsumer := testservice.ConsumerKYC
		Expect(testconsumer).ToNot(BeNil())

		Expect(testservice).To(Equal(service))
	})
})
