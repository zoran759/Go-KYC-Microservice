package expectid

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/modulusglobal/kyc/integrations/idology"
)

var _ = Describe("Client", func() {
	Describe("new client should get correct config", func() {
		It("should success", func() {
			config := idology.Config{
				Host:     "fake_host",
				Username: "fake_name",
				Password: "fake_password",
			}

			testclient := &client{
				Config: config,
			}

			newclient := NewClient(config)
			Expect(newclient).NotTo(BeNil())

			client := newclient.(*client)
			Expect(client).To(Equal(testclient))
		})
	})
})
