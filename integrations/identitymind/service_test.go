package identitymind

import (
	"io/ioutil"
	"reflect"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/identitymind/consumer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testImage []byte

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
		var service = New(Config{
			Host:     SandboxBaseURL,
			Username: "modulusglobal",
			Password: "64117e699462ce859d970648461a625bc6a6f3cb",
		})

		It("Should return bad reputation for the customer", func() {
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
			Expect(result.StatusPolling).To(BeNil())
		})

		It("Should return suspicious reputation for the customer", func() {
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
			Expect(result.StatusPolling).To(BeNil())
		})

		It("Should return trusted reputation for the customer", func() {
			Expect(service).NotTo(BeNil())

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
						Data:        testImage,
					},
				},
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
			Expect(result.StatusPolling).To(BeNil())
		})

		It("Should return unknown reputation for the customer", func() {
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
			Expect(result.StatusPolling).To(BeNil())
		})

		It("Should return denied policy result for the customer", func() {
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
			Expect(result.StatusPolling).To(BeNil())
		})

		It("Should return review policy result for the customer", func() {
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
			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusPolling).NotTo(BeNil())
			Expect(result.StatusPolling.Provider).To(Equal(common.IdentityMind))
			Expect(result.StatusPolling.CustomerID).NotTo(BeEmpty())
		})

		It("Should return accepted policy result for the customer", func() {
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
			Expect(result.StatusPolling).To(BeNil())
		})
	})
})

func init() {
	testImage, _ = ioutil.ReadFile("../../test_data/passport.jpg")
}
