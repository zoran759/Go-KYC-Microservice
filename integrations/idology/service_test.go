package idology_test

import (
	"flag"
	nethttp "net/http"
	"net/url"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.com/lambospeed/kyc/integrations/idology"

	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/http"
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// Use this to setup proxy when you run tests which interact with the API
// and you're not in front of a whitelisted host.
// Warning! This has sense only when using ssh local “dynamic” application-level port forwarding.
// So you have to have ssh-access to a whitelisted host as well.
var proxyURL = "socks5://localhost:8000"

var runLive bool

var _ = Describe("The IDology KYC service", func() {
	Specify("should be properly created", func() {
		config := Config{
			Host:             "fake_host",
			Username:         "fake_username",
			Password:         "fake_password",
			UseSummaryResult: true,
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

	Context("when sending requests to IDology API", func() {
		var runliveMsg = "use '-runlive' command-line flag to activate this test"

		var customer = &common.UserData{
			FirstName:     "John",
			LastName:      "Smith",
			DateOfBirth:   common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
			AddressString: "222333 PeachTree Place, Atlanta, GA 30318",
			CurrentAddress: common.Address{
				CountryAlpha2:     "US",
				State:             "Georgia",
				Town:              "Atlanta",
				Street:            "PeachTree Place",
				BuildingNumber:    "222333",
				PostCode:          "30318",
				StateProvinceCode: "GA",
			},
			Documents: []common.Document{
				common.Document{
					Metadata: common.DocumentMetadata{
						Type:    common.IDCard,
						Country: "USA",
						Number:  "112223333",
					},
				},
			},
		}

		var service = New(Config{
			Host:     expectid.APIendpoint,
			Username: "modulus.dev2",
			Password: "}$tRPfT1sZQmU@uh8@",
		})

		var skipFunc = func() {
			if !runLive {
				Skip(runliveMsg)
			}
		}

		BeforeEach(func() {
			proxy, _ := url.Parse(proxyURL)
			t := &nethttp.Transport{
				Proxy: nethttp.ProxyURL(proxy),
			}
			http.Client.Transport = t
		})

		Context("when using wrong credentials in config", func() {
			It("should error", func() {
				skipFunc()

				failedService := New(Config{
					Host:     expectid.APIendpoint,
					Username: "modulus.dev2",
					Password: "wrong_password",
				})

				result, details, err := failedService.ExpectID.CheckCustomer(customer)

				Expect(result).To(Equal(common.Error))
				Expect(details).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during verification: Invalid username and password"))
			})
		})

		Context("when the test data for the successful response is provided", func() {
			It("should return clean result", func() {
				skipFunc()

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
				Expect(result).To(Equal(common.Approved))
			})
		})

		Context("when the test data for triggering ID Notes is provided", func() {
			It("should return COPPA Alert", func() {
				skipFunc()

				coppaCustomer := &common.UserData{}
				*coppaCustomer = *customer
				coppaCustomer.DateOfBirth = common.Time(
					time.Date(2009, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(coppaCustomer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Denied))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("COPPA Alert"))
			})
			// TODO: implement this.
		})
	})
})

func init() {
	flag.BoolVar(&runLive, "runlive", false, "Run live tests against IDology API.")
}
