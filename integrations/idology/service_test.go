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
// If you run tests on a whitelisted host leave this empty.
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

	// Below are the tests that should be run either on a whitelisted host
	// or using some method to forward requests through a whitelisted host.
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
			if len(proxyURL) > 0 {
				proxy, _ := url.Parse(proxyURL)
				t := &nethttp.Transport{
					Proxy: nethttp.ProxyURL(proxy),
				}
				http.Client.Transport = t
			}
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
			It("should deny and return COPPA Alert", func() {
				skipFunc()

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.DateOfBirth = common.Time(
					time.Date(2009, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "2240 Magnolia, Atlanta, GA 30318"
				noteCustomer.CurrentAddress.Street = "Magnolia"
				noteCustomer.CurrentAddress.BuildingNumber = "2240"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "222333 Magnolia, Atlanta, GA 30318"
				noteCustomer.CurrentAddress.Street = "Magnolia"
				noteCustomer.CurrentAddress.BuildingNumber = "222333"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "2240 PeachTree Place, Atlanta, GA 30318"
				noteCustomer.CurrentAddress.BuildingNumber = "2240"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "PO Box 123, Atlanta, GA 30318"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "222333 PeachTree Place, Atlanta, GA 30316"
				noteCustomer.CurrentAddress.PostCode = "30316"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.DateOfBirth = common.Time(
					time.Date(1970, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.DateOfBirth = common.Time(
					time.Date(1976, time.February, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.DateOfBirth = common.Time(
					time.Date(1975, time.May, 28, 0, 0, 0, 0, time.UTC),
				)

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.Documents[0].Metadata.Number = "112223345"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.Documents[0].Metadata.Number = "112223334"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.AddressString = "222333 PeachTree Place, Atlanta, AL 30318"
				noteCustomer.CurrentAddress.State = "Alabama"
				noteCustomer.CurrentAddress.StateProvinceCode = "AL"
				noteCustomer.Documents = nil

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.FirstName = "Jane"
				noteCustomer.AddressString = "5432 Any Place, La Crescenta, CA 91214"
				noteCustomer.CurrentAddress.State = "California"
				noteCustomer.CurrentAddress.Town = "La Crescenta"
				noteCustomer.CurrentAddress.Street = "Any Place"
				noteCustomer.CurrentAddress.BuildingNumber = "5432"
				noteCustomer.CurrentAddress.PostCode = "91214"
				noteCustomer.CurrentAddress.StateProvinceCode = "CA"
				noteCustomer.Documents[0].Metadata.Number = "112221111"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.LastName = "Black"
				noteCustomer.AddressString = "345 Some Avenu, Atlanta, GA 30303"
				noteCustomer.CurrentAddress.Street = "Some Avenu"
				noteCustomer.CurrentAddress.BuildingNumber = "345"
				noteCustomer.CurrentAddress.PostCode = "30303"
				noteCustomer.Documents = nil

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.FirstName = "Jane"
				noteCustomer.LastName = "Brown"
				noteCustomer.AddressString = "9000 Any Street, La Crescenta, CA 91224"
				noteCustomer.CurrentAddress.State = "California"
				noteCustomer.CurrentAddress.Town = "La Crescenta"
				noteCustomer.CurrentAddress.Street = "Any Street"
				noteCustomer.CurrentAddress.BuildingNumber = "9000"
				noteCustomer.CurrentAddress.PostCode = "91224"
				noteCustomer.CurrentAddress.StateProvinceCode = "CA"
				noteCustomer.Documents[0].Metadata.Number = "112221010"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{}
				*noteCustomer = *customer
				noteCustomer.FirstName = "Jane"
				noteCustomer.LastName = "Black"
				noteCustomer.AddressString = "12345 Magnolia Way, Atlanta, GA 30303"
				noteCustomer.CurrentAddress.Street = "Magnolia Way"
				noteCustomer.CurrentAddress.BuildingNumber = "12345"
				noteCustomer.CurrentAddress.PostCode = "30303"

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{
					FirstName:     "Jane",
					LastName:      "Williams",
					DateOfBirth:   common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
					AddressString: "8888 Any Street, Dallas, GA 30132",
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

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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

				noteCustomer := &common.UserData{
					FirstName:     "John",
					LastName:      "Bredenkamp",
					DateOfBirth:   common.Time(time.Date(1940, time.August, 1, 0, 0, 0, 0, time.UTC)),
					AddressString: "147 Brentwood Drive, Nashville, TN 37214",
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
								Type:    common.IDCard,
								Country: "USA",
								Number:  "555667777",
							},
						},
					},
				}

				result, details, err := service.ExpectID.CheckCustomer(noteCustomer)

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
}
