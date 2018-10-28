package jumio

import (
	"encoding/base64"
	"fmt"
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

		It("should fail using the nil customer", func() {
			Expect(service).NotTo(BeNil())

			_, err := service.CheckCustomer(nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no customer supplied"))
		})

		It("should fail using the empty customer", func() {
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

			It("should fail with an error", func() {
				Expect(service).NotTo(BeNil())

				_, err := service.CheckCustomer(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during sending request: Post https://netverify.com/api/netverify/v2/performNetverify: no responder found"))
			})

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

	Describe("CheckStatus", func() {
		var service = New(Config{
			BaseURL: USbaseURL,
			Token:   "test_token",
			Secret:  "test_secret",
		})

		var (
			referenceID            = "jumioID"
			retrieveScanStatusURL  = USbaseURL + scanStatusEndpoint + referenceID
			retrieveScanDetailsURL = fmt.Sprintf(USbaseURL+scanDetailsEndpoint, referenceID)
		)

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail using empty reference id", func() {
			Expect(service).NotTo(BeNil())

			_, err := service.CheckStatus("")

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("empty Jumio’s reference number of an existing scan"))
		})

		It("should fail with an error", func() {
			Expect(service).NotTo(BeNil())

			_, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: Get https://netverify.com/api/netverify/v2/scans/jumioID: no responder found"))
		})

		It("should fail with an http error", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusNotFound, nil))

			result, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: http error"))
			Expect(result.ErrorCode).To(Equal("404"))
		})

		It("should fail due to malformed response", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":7}`)))

			result, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: json: cannot unmarshal number into Go struct field StatusResponse.status of type jumio.ScanStatus"))
			Expect(result.ErrorCode).To(BeEmpty())
		})

		It("should success with pending status", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":"PENDING"}`)))

			result, err := service.CheckStatus(referenceID)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal(common.Unclear))
			Expect(result.Details).To(BeNil())
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).NotTo(BeNil())
			Expect(result.StatusCheck.Provider).To(Equal(common.Jumio))
			Expect(result.StatusCheck.ReferenceID).To(Equal(referenceID))
			Expect(time.Time(result.StatusCheck.LastCheck)).NotTo(BeZero())
		})

		It("should fail to retrieve details with an http error", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":"DONE"}`)))
			httpmock.RegisterResponder(http.MethodGet, retrieveScanDetailsURL, httpmock.NewBytesResponder(http.StatusInternalServerError, nil))

			result, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: http error"))
			Expect(result.ErrorCode).To(Equal("500"))
		})

		It("should success to retrieve details with approved status", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":"DONE"}`)))
			httpmock.RegisterResponder(http.MethodGet, retrieveScanDetailsURL, httpmock.NewBytesResponder(http.StatusOK, approvedResponse))

			result, err := service.CheckStatus(referenceID)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).To(BeNil())
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("should fail with an unknown status", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":"OTHER"}`)))

			_, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown status of the verification: OTHER"))
		})

		It("should fail to retrieve details with an error", func() {
			Expect(service).NotTo(BeNil())

			httpmock.RegisterResponder(http.MethodGet, retrieveScanStatusURL, httpmock.NewBytesResponder(http.StatusOK, []byte(`{"status":"DONE"}`)))

			_, err := service.CheckStatus(referenceID)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: Get https://netverify.com/api/netverify/v2/scans/jumioID/data: no responder found"))
		})
	})
})
