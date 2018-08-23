package expectid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Client", func() {
	Describe("new client should get correct config", func() {
		It("should success", func() {
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
