package idology_test

import (
	"flag"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"gitlab.com/lambospeed/kyc/common"

	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "gitlab.com/lambospeed/kyc/integrations/idology"
)

// Use this to setup proxy when you run tests which interact with the API
// and you're not in front of a whitelisted host.
// Warning! This has sense only when using ssh-socks
// because the proxy must be running on a whitelisted host anyway.
var proxyURL = "socks5://localhost:8000"

var runLive bool

var _ = Describe("Service", func() {
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

		var validCustomer = &common.UserData{
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

		var config = Config{
			Host:     expectid.APIendpoint,
			Username: "modulus.dev2",
			Password: "}$tRPfT1sZQmU@uh8@",
		}

		BeforeEach(func() {
			proxy, _ := url.Parse(proxyURL)
			t := &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
			http.DefaultClient.Transport = t
		})

		It("should success", func() {
			if !runLive {
				Skip(runliveMsg)
			}

			service := New(config)
			result, details, err := service.ExpectID.CheckCustomer(validCustomer)

			Expect(err).NotTo(HaveOccurred())
			Expect(details).To(BeNil())
			Expect(result).NotTo(BeNil())
		})
	})
})

func init() {
	flag.BoolVar(&runLive, "runlive", false, "Run live tests against IDology API.")
}
