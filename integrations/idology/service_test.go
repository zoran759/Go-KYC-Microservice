package idology_test

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.com/lambospeed/kyc/integrations/idology"

	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/idology/expectid"
)

// This is "-proxy" command-line flag to set the proxy for requests to the API.
// Use this to setup proxy when you run tests which interact with the API
// and you're not in front of a whitelisted host.
// Warning! Anyway, the proxy must be running on a whitelisted host otherwise it'll not help.
var proxyURL string

// This is "-runlive" command-line flag to activate the sandbox testing (see "ExpectID Sandbox Guide.pdf").
var runLive bool

var _ = BeforeSuite(func() {
	if runLive && len(proxyURL) > 0 {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			log.Fatalln(err)
		}
		t := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		http.DefaultClient.Transport = t
	}
})

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

	// Below are the tests that should be run either on a whitelisted host
	// or using some method to forward requests through a whitelisted host.
	Context("using sandbox testing of IDology API", func() {
		var runliveMsg = "use '-runlive' command-line flag to activate this test"

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
				Documents: []common.Document{
					common.Document{
						Metadata: common.DocumentMetadata{
							Type:            common.IDCard,
							Country:         "USA",
							Number:          "112223333",
							CardLast4Digits: "3333",
						},
					},
				},
			}
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

		Context("when using wrong credentials in config", func() {
			It("should error", func() {
				skipFunc()

				failedService := New(Config{
					Host:     expectid.APIendpoint,
					Username: "modulus.dev2",
					Password: "wrong_password",
				})

				customer := newCustomer()
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

				customer := newCustomer()
				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(details).To(BeNil())
				Expect(result).To(Equal(common.Approved))
			})
		})

		Context("when the test data for triggering ID Notes is provided", func() {
			It("should deny and return COPPA Alert", func() {
				skipFunc()

				customer := newCustomer()
				customer.DateOfBirth = common.Time(
					time.Date(2009, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Denied))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("COPPA Alert"))
			})

			// "Address Does Not Match" test actually returns more qualifiers.
			It("should approve and return Address Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.Street = "Magnolia"
				customer.CurrentAddress.BuildingNumber = "2240"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(3))
				Expect(details.Reasons[0]).To(Equal("Address Does Not Match"))
				Expect(details.Reasons[1]).To(Equal("Street Number Does Not Match"))
				Expect(details.Reasons[2]).To(Equal("Street Name Does Not Match"))
			})

			// "Street Name Does Not Match" test actually returns more qualifiers.
			It("should approve and return Street Name Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.Street = "Magnolia"
				customer.CurrentAddress.BuildingNumber = "222333"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(2))
				Expect(details.Reasons[0]).To(Equal("Address Does Not Match"))
				Expect(details.Reasons[1]).To(Equal("Street Name Does Not Match"))
			})

			// "Street Number Does Not Match" test actually returns more qualifiers.
			It("should approve and return Street Number Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.Street = "PeachTree Place"
				customer.CurrentAddress.BuildingNumber = "2240"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(2))
				Expect(details.Reasons[0]).To(Equal("Address Does Not Match"))
				Expect(details.Reasons[1]).To(Equal("Street Number Does Not Match"))
			})

			// "Input Address is a PO Box" test actually returns more qualifiers.
			It("should approve and return Input Address is a PO Box", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.Street = "PO Box 123"
				customer.CurrentAddress.BuildingNumber = ""

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(2))
				Expect(details.Reasons[0]).To(Equal("Address Does Not Match"))
				Expect(details.Reasons[1]).To(Equal("Input Address is a PO Box"))
			})

			It("should approve and return ZIP Code Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.PostCode = "30316"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("ZIP Code Does Not Match"))
			})

			It("should approve and return YOB Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.DateOfBirth = common.Time(
					time.Date(1970, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("YOB Does Not Match"))
			})

			It("should approve and return YOB Does Not Match, Within 1 Year Tolerance", func() {
				skipFunc()

				customer := newCustomer()
				customer.DateOfBirth = common.Time(
					time.Date(1976, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("YOB Does Not Match, Within 1 Year Tolerance"))
			})

			It("should approve and return MOB Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.DateOfBirth = common.Time(
					time.Date(1975, time.May, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("MOB Does Not Match"))
			})

			// "Newer Record Found" test doesn't return what is expected. Skipped.

			It("should approve and return SSN Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.Documents[0].Metadata.Number = ""
				customer.Documents[0].Metadata.CardLast4Digits = "3345"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("SSN Does Not Match"))
			})

			It("should approve and return SSN Does Not Match, Within Tolerance", func() {
				skipFunc()

				customer := newCustomer()
				customer.Documents[0].Metadata.Number = ""
				customer.Documents[0].Metadata.CardLast4Digits = "3334"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("SSN Does Not Match, Within Tolerance"))
			})

			It("should approve and return State Does Not Match", func() {
				skipFunc()

				customer := newCustomer()
				customer.CurrentAddress.State = "Alabama"
				customer.CurrentAddress.StateProvinceCode = "AL"
				customer.Documents = nil

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("State Does Not Match"))
			})

			It("should approve and return Single Address in File", func() {
				skipFunc()

				customer := &common.UserData{
					FirstName:   "Jane",
					LastName:    "Smith",
					DateOfBirth: common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
					CurrentAddress: common.Address{
						CountryAlpha2:     "US",
						State:             "California",
						Town:              "La Crescenta",
						Street:            "Any Place",
						BuildingNumber:    "5432",
						PostCode:          "91214",
						StateProvinceCode: "CA",
					},
					Documents: []common.Document{
						common.Document{
							Metadata: common.DocumentMetadata{
								Type:            common.IDCard,
								Country:         "USA",
								CardLast4Digits: "1111",
							},
						},
					},
				}

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("Single Address in File"))
			})

			// "Thin File" test actually returns more qualifiers.
			It("should approve and return Thin File", func() {
				skipFunc()

				customer := newCustomer()
				customer.LastName = "Black"
				customer.CurrentAddress.Street = "Some Avenu"
				customer.CurrentAddress.BuildingNumber = "345"
				customer.CurrentAddress.PostCode = "30303"
				customer.Documents = nil

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(4))
				Expect(details.Reasons[0]).To(Equal("No DOB Available"))
				Expect(details.Reasons[1]).To(Equal("SSN Not Found"))
				Expect(details.Reasons[2]).To(Equal("Thin File"))
				Expect(details.Reasons[3]).To(Equal("Data Strength Alert"))
			})

			// "DOB Not Available" test actually returns slightly different result.
			It("should approve and return DOB Not Available", func() {
				skipFunc()

				customer := &common.UserData{
					FirstName: "Jane",
					LastName:  "Brown",
					CurrentAddress: common.Address{
						CountryAlpha2:     "US",
						State:             "California",
						Town:              "La Crescenta",
						Street:            "Any Street",
						BuildingNumber:    "9000",
						PostCode:          "91224",
						StateProvinceCode: "CA",
					},
					Documents: []common.Document{
						common.Document{
							Metadata: common.DocumentMetadata{
								Type:            common.IDCard,
								Country:         "USA",
								CardLast4Digits: "1010",
							},
						},
					},
				}

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(2))
				Expect(details.Reasons[0]).To(Equal("No DOB Available"))
				Expect(details.Reasons[1]).To(Equal("Data Strength Alert"))
			})

			// "SSN Not Available" test actually returns slightly different result.
			It("should approve and return SSN Not Available", func() {
				skipFunc()

				customer := newCustomer()
				customer.FirstName = "Jane"
				customer.LastName = "Black"
				customer.CurrentAddress.Street = "Magnolia Way"
				customer.CurrentAddress.BuildingNumber = "12345"
				customer.CurrentAddress.PostCode = "30303"

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(1))
				Expect(details.Reasons[0]).To(Equal("SSN Not Found"))
			})

			// "Subject Deceased" test doesn't return what is expected.

			// "SSN Issue Prior to DOB" test doesn't return what is expected.
			// "SSN Invalid" test doesn't return what is expected.
			// Are they kidding me??? These two are identical in the table!

			// "Warm Address" test actually returns slightly different result.
			It("should approve and return Warm Address", func() {
				skipFunc()

				customer := &common.UserData{
					FirstName:   "Jane",
					LastName:    "Williams",
					DateOfBirth: common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
					CurrentAddress: common.Address{
						CountryAlpha2:     "US",
						State:             "Georgia",
						Town:              "Dallas",
						Street:            "Any Street",
						BuildingNumber:    "8888",
						PostCode:          "30132",
						StateProvinceCode: "GA",
					},
				}

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Approved))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(2))
				Expect(details.Reasons[0]).To(Equal("Warm Address Alert (hotel)"))
				Expect(details.Reasons[1]).To(Equal("Data Strength Alert"))
			})
		})

		Context("when the test data for triggering Patriot Act Alert is provided", func() {
			It("should deny and return Patriot Act Alert", func() {
				skipFunc()

				customer := &common.UserData{
					FirstName:   "John",
					LastName:    "Bredenkamp",
					DateOfBirth: common.Time(time.Date(1940, time.August, 1, 0, 0, 0, 0, time.UTC)),
					CurrentAddress: common.Address{
						CountryAlpha2:     "US",
						State:             "Tennessee",
						Town:              "Nashville",
						Street:            "Brentwood Drive",
						BuildingNumber:    "147",
						PostCode:          "37214",
						StateProvinceCode: "TN",
					},
					Documents: []common.Document{
						common.Document{
							Metadata: common.DocumentMetadata{
								Type:            common.IDCard,
								Country:         "USA",
								Number:          "555667777",
								CardLast4Digits: "7777",
							},
						},
					},
				}

				result, details, err := service.ExpectID.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(Equal(common.Denied))
				Expect(details).NotTo(BeNil())
				Expect(details.Finality).To(Equal(common.Unknown))
				Expect(details.Reasons).NotTo(BeNil())
				Expect(details.Reasons).To(HaveLen(4))
				Expect(details.Reasons[0]).To(Equal("Patriot Act Alert"))
				Expect(details.Reasons[1]).To(Equal("Office of Foreign Asset Control"))
				Expect(details.Reasons[2]).To(Equal("Patriot Act score: 100"))
				Expect(details.Reasons[3]).To(Equal("PA DOB Match"))
			})
		})
	})
})

func init() {
	flag.BoolVar(&runLive, "runlive", false, "Run live tests against IDology API.")
	flag.StringVar(&proxyURL, "proxy", "", "Set a proxy when you're not in front of a whitelisted host.")
}
