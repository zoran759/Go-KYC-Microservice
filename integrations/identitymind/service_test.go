package identitymind

import (
	"reflect"

	"modulus/kyc/integrations/identitymind/consumer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("The IdentityMind service", func() {
	Specify("should be properly created", func() {
		config := Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		}

		service := &Service{
			consumer: consumer.NewClient(consumer.Config(config)),
		}

		Expect(service.consumer).NotTo(BeNil())

		testservice := New(config)

		Expect(testservice).NotTo(BeNil())
		Expect(testservice.consumer).ToNot(BeNil())
		Expect(reflect.TypeOf(testservice)).To(Equal(reflect.TypeOf((*Service)(nil))))

		Expect(*testservice).To(Equal(*service))
		Expect(*testservice.consumer).To(Equal(*service.consumer))
	})
})
