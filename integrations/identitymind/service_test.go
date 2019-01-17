package identitymind

import (
	"flag"
	"io/ioutil"
	"reflect"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/identitymind/consumer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	useSandbox = flag.Bool("sandbox", false, "activate sandbox testing")
	testImages = flag.Bool("test-images", false, "test document images upload")
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

	Describe("CheckCustomer Sandbox Testing", func() {
		var skipTest = func() {
			if !*useSandbox {
				Skip("use '-sandbox' command line flag to activate sandbox testing")
			}
		}

		var service = New(Config{
			Host:     SandboxBaseURL,
			Username: "cointrading",
			Password: "0c67f85b7d882326ff00ca77b5b98071c1609d2b",
		})

		It("Should return bad reputation for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "bad_brad",
				FirstName:   "Brad",
				LastName:    "Pit",
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: BAD"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: User previously failed validation"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("Should return suspicious reputation for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "suspicious_sue",
				FirstName:   "Sue",
				LastName:    "Rushton",
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: SUSPICIOUS"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: User previously failed validation"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("Should return trusted reputation for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "trusted_tom",
				FirstName:   "Tom",
				LastName:    "Pennington",
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: TRUSTED"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Unvalidated, but long-lived good User"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("Should return unknown reputation for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "unknown_boriss",
				FirstName:   "Boriss",
				LastName:    "Godunoff",
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: UNKNOWN"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Unknown User"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("Should return denied policy result for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "denied_sergey",
				FirstName:   "Sergey",
				LastName:    "Sarbash",
				CurrentAddress: common.Address{
					CountryAlpha2: "US",
					Town:          "Detroit",
				},
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Denied))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: UNKNOWN"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: DENY"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Unknown User"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: DENY"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("Should return review policy result for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "reviewed_sergey",
				FirstName:   "Sergey",
				LastName:    "Sarbash",
				CurrentAddress: common.Address{
					CountryAlpha2: "US",
					Town:          "Monte Rio",
				},
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Unclear))
			Expect(result.Details).To(BeNil())
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).NotTo(BeNil())
			Expect(result.StatusCheck.Provider).To(Equal(common.IdentityMind))
			Expect(result.StatusCheck.ReferenceID).NotTo(BeEmpty())
			Expect(time.Time(result.StatusCheck.LastCheck)).NotTo(BeZero())
		})

		It("Should return accepted policy result for the customer", func() {
			skipTest()

			Expect(service).NotTo(BeNil())

			customer := &common.UserData{
				AccountName: "accepted_sergey",
				FirstName:   "Sergey",
				LastName:    "Sarbash",
				CurrentAddress: common.Address{
					CountryAlpha2: "RU",
					Town:          "Naberezhnyye Chelny",
				},
			}

			result, err := service.CheckCustomer(customer)

			Expect(err).To(BeNil())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: UNKNOWN"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Unknown User"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})
	})

	Describe("testing document image upload", func() {
		It("should upload document image and successfully retrieve it from the API", func() {
			if !*testImages {
				Skip("use '-test-images' command line flag to activate document images upload testing")
			}

			service := New(Config{
				Host:     SandboxBaseURL,
				Username: "cointrading",
				Password: "0c67f85b7d882326ff00ca77b5b98071c1609d2b",
			})

			Expect(service).NotTo(BeNil())

			passport, err := ioutil.ReadFile("../../test_data/passport.jpg")

			Expect(err).NotTo(HaveOccurred())
			Expect(passport).NotTo(BeEmpty())

			customer := &common.UserData{
				AccountName: "trusted_tom",
				FirstName:   "Tom",
				LastName:    "Pennington",
				Passport: &common.Passport{
					Number:        "0123456789",
					CountryAlpha2: "US",
					State:         "WA",
					IssuedDate:    common.Time(time.Date(2015, 05, 25, 0, 0, 0, 0, time.UTC)),
					ValidUntil:    common.Time(time.Date(2025, 05, 24, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "passport.jpg",
						ContentType: "image/jpeg",
						Data:        passport,
					},
				},
			}

			_, err = service.CheckCustomer(customer)

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
