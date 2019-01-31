package expectid

import (
	"net/http"
	"time"

	"modulus/kyc/common"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	approvedResponse = []byte(`<?xml version="1.0"?><response><id-number>2073386242</id-number><summary-result><key>id.success</key><message>PASS</message></summary-result><results><key>result.match</key><message>ID Located</message></results><idnotescore>0</idnotescore></response>`)

	rejectedResponse = []byte(`<?xml version="1.0"?><response><id-number>2073457900</id-number><summary-result><key>id.failure</key><message>FAIL</message></summary-result><results><key>result.no.match</key><message>ID Not Located</message></results><idliveq-error><key>id.not.eligible.for.questions</key><message>Not Eligible For Questions</message></idliveq-error></response>`)

	rejectedResponseWithNotes = []byte(`<?xml version="1.0"?><response><id-number>2073386264</id-number><summary-result><key>id.failure</key><message>FAIL</message></summary-result><results><key>result.match.restricted</key><message>result.match.restricted</message></results><qualifiers><qualifier><key>resultcode.coppa.alert</key><message>COPPA Alert</message></qualifier></qualifiers><idliveq-error><key>id.not.eligible.for.questions</key><message>Not Eligible For Questions</message></idliveq-error></response>`)

	errorResponse = []byte(`<response><error>Invalid username and password</error></response>`)
)

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("should construct proper client instance", func() {
			config := Config{
				Host:     "fake_host",
				Username: "fake_name",
				Password: "fake_password",
			}

			testclient := Client{
				config: config,
			}

			client := NewClient(config)

			Expect(client).To(Equal(testclient))
		})
	})

	Describe("CheckCustomer", func() {
		var client = NewClient(Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		})

		var newCustomer = func() *common.UserData {
			return &common.UserData{
				FirstName:   "John",
				LastName:    "Smith",
				DateOfBirth: common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
				CurrentAddress: common.Address{
					CountryAlpha2:     "US",
					State:             "Georgia",
					Town:              "Atlanta",
					Street:            "PeachTree Place",
					BuildingNumber:    "222333",
					PostCode:          "30318",
					StateProvinceCode: "GA",
				},
				IDCard: &common.IDCard{
					CountryAlpha2: "US",
					Number:        "112223333",
				},
			}
		}

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should return approved result", func() {
			httpmock.RegisterResponder(
				http.MethodPost,
				client.config.Host,
				httpmock.NewBytesResponder(http.StatusOK, approvedResponse),
			)

			customer := newCustomer()
			result, err := client.CheckCustomer(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Details).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))

		})

		It("should return rejected result", func() {
			httpmock.RegisterResponder(
				http.MethodPost,
				client.config.Host,
				httpmock.NewBytesResponder(http.StatusOK, rejectedResponse),
			)

			customer := newCustomer()
			customer.LastName = "Doe"
			result, err := client.CheckCustomer(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Details).To(BeNil())
			Expect(result.Status).To(Equal(common.Denied))
		})

		It("should return rejected result with ID Notes", func() {
			httpmock.RegisterResponder(
				http.MethodPost,
				client.config.Host,
				httpmock.NewBytesResponder(http.StatusOK, rejectedResponseWithNotes),
			)

			customer := newCustomer()
			customer.DateOfBirth = common.Time(time.Date(2009, time.February, 28, 0, 0, 0, 0, time.UTC))
			result, err := client.CheckCustomer(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(1))
			Expect(result.Details.Reasons[0]).To(Equal("COPPA Alert"))
			Expect(result.Status).To(Equal(common.Denied))
		})

		It("should return error", func() {
			httpmock.RegisterResponder(
				http.MethodPost,
				client.config.Host,
				httpmock.NewBytesResponder(http.StatusOK, errorResponse),
			)

			customer := newCustomer()
			result, err := client.CheckCustomer(customer)

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during verification: Invalid username and password"))
		})
	})
})
