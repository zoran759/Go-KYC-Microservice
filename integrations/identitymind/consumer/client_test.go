package consumer

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})
