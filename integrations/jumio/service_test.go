package jumio

import (
	"encoding/base64"
	"net/http"
	"time"

	"modulus/kyc/common"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var scanrefResponse = []byte(`
{
	"timestamp": "2018-08-16T10:37:44.623Z",
	"jumioIdScanReference": "sample-jumio-scan-reference"
}`)

var _ = Describe("Service", func() {
	Describe("New", func() {
		Specify("should properly create service object", func() {
			config := Config{
				BaseURL: "fake_baseURL",
				Token:   "fake_token",
				Secret:  "fake_secret",
			}

			s := &service{
				baseURL:     config.BaseURL,
				credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
			}

			cc := New(config)
			ts := cc.(*service)

			Expect(s).To(Equal(ts))
		})
	})

	Describe("CheckCustomer", func() {
		var service = New(Config{
			BaseURL: USbaseURL,
			Token:   "test_token",
			Secret:  "test_secret",
		})

		It("should fail with the nil customer", func() {
			Expect(service).NotTo(BeNil())

			_, err := service.CheckCustomer(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no customer supplied"))
		})

		It("should fail with the empty customer", func() {
			Expect(service).NotTo(BeNil())

			_, err := service.CheckCustomer(&common.UserData{})

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid customer data: missing acceptable document for the verification (anyone of passport, driving license or id card)"))
		})

		Context("with http-requests using httpmock", func() {
			BeforeEach(func() {
				httpmock.Activate()
			})

			AfterEach(func() {
				httpmock.DeactivateAndReset()
			})

			var customer = &common.UserData{
				FirstName:   "Bruce",
				LastName:    "Wayne",
				DateOfBirth: common.Time(time.Date(1950, 03, 17, 0, 0, 0, 0, time.UTC)),
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "batman.png",
						ContentType: "image/png",
						Data:        []byte{},
					},
				},
				Passport: &common.Passport{
					Number:        "1234567890",
					CountryAlpha2: "US",
					State:         "WA",
					IssuedDate:    common.Time(time.Date(2010, 05, 15, 0, 0, 0, 0, time.UTC)),
					ValidUntil:    common.Time(time.Date(2020, 05, 14, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "passport.png",
						ContentType: "image/png",
						Data:        []byte{},
					},
				},
			}

			It("should fail with http error", func() {
				Expect(service).NotTo(BeNil())

				httpmock.RegisterResponder(http.MethodPost, USbaseURL+performNetverifyEndpoint, httpmock.NewBytesResponder(http.StatusBadRequest, nil))

				result, err := service.CheckCustomer(customer)

				Expect(result.ErrorCode).To(Equal("400"))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during sending request: http error"))
			})

			It("should success with status check data", func() {
				Expect(service).NotTo(BeNil())

				httpmock.RegisterResponder(http.MethodPost, USbaseURL+performNetverifyEndpoint, httpmock.NewBytesResponder(http.StatusOK, scanrefResponse))
				httpmock.RegisterResponder(http.MethodPost, USbaseURL+scanStatusEndpoint, httpmock.NewBytesResponder(http.StatusOK, scanrefResponse))

				result, err := service.CheckCustomer(customer)

				Expect(err).NotTo(HaveOccurred())
				Expect(result.Status).To(Equal(common.Unclear))
				Expect(result.Details).To(BeNil())
				Expect(result.ErrorCode).To(BeEmpty())
				Expect(result.StatusCheck).NotTo(BeNil())
				Expect(result.StatusCheck.Provider).To(Equal(common.Jumio))
				Expect(result.StatusCheck.ReferenceID).To(Equal("sample-jumio-scan-reference"))
				Expect(result.StatusCheck.LastCheck).NotTo(BeZero())
			})
		})
	})
})
