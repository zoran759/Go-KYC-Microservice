package consumer

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/lambospeed/kyc/common"
)

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("should construct proper client instance", func() {
			config := Config{
				Host:     "fake_host",
				Username: "fake_name",
				Password: "fake_password",
			}

			testclient := &Client{
				config: config,
			}

			client := NewClient(config)

			Expect(client).NotTo(BeNil())
			Expect(client).To(Equal(testclient))
		})
	})

	Describe("CheckCustomer", func() {
		var client = NewClient(Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		})

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail with error message", func() {
			Expect(client).ToNot(BeNil())

			result, details, err := client.CheckCustomer(nil)

			Expect(result).To(Equal(common.Error))
			Expect(details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no customer supplied"))
		})
	})
})
