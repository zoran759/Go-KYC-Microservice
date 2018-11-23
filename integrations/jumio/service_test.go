package jumio

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"modulus/kyc/common"
	mhttp "modulus/kyc/http"

	"github.com/jarcoal/httpmock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testImageUpload = flag.Bool("use-images", false, "test document images uploading")

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

	Describe("TestImageUpload", func() {
		It("should success", func() {
			if !*testImageUpload {
				Skip("use '-use-images' flag to activate images uploading test")
			}

			testPassport, _ := ioutil.ReadFile("../../test_data/passport.jpg")
			testSelfie, _ := ioutil.ReadFile("../../test_data/selfie.jpg")

			Expect(testPassport).NotTo(BeEmpty(), "testPassport must contain the content of the image data file 'passport.jpg'")
			Expect(testSelfie).NotTo(BeEmpty(), "testSelfie must contain the content of the image data file 'selfie.png'")

			customer := &common.UserData{
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: common.Time(time.Date(1975, 06, 15, 0, 0, 0, 0, time.UTC)),
				Passport: &common.Passport{
					Number:        "0123456789",
					CountryAlpha2: "US",
					State:         "GA",
					IssuedDate:    common.Time(time.Date(2015, 06, 20, 0, 0, 0, 0, time.UTC)),
					ValidUntil:    common.Time(time.Date(2025, 06, 19, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "passport.jpg",
						ContentType: "image/jpeg",
						Data:        testPassport,
					},
				},
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "selfie.jpg",
						ContentType: "image/jpeg",
						Data:        testSelfie,
					},
				},
			}

			config := Config{
				BaseURL: "https://netverify.com/api/netverify/v2",
				Token:   "0afc675f-400b-421f-a3d5-19520ba2d8e7",
				Secret:  "nZtnsNyc7mO1rtqbjT59XC5GnG6IOQS6",
			}

			service := New(config)

			result, err := service.CheckCustomer(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(result.StatusCheck).NotTo(BeNil(), "status polling data has to be provided")

			refID := result.StatusCheck.ReferenceID

			// According to workflow firstly check the status.
			_, _ = service.CheckStatus(refID)

			// Get back the downloaded document images.

			// Get list of available images of the scan.
			type imgref struct {
				Classifier string `json:"classifier"`
				Href       string `json:"href"`
			}
			type imgresp struct {
				Timestamp     string   `json:"timestamp"`
				ScanReference string   `json:"scanReference"`
				Images        []imgref `json:"images"`
			}

			headers := mhttp.Headers{
				"Accept":        "application/json",
				"User-Agent":    "Modulus Exchange/1.0",
				"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
			}
			_, body, err := mhttp.Get(fmt.Sprintf("https://netverify.com/api/netverify/v2/scans/%s/images", refID), headers)

			Expect(err).NotTo(HaveOccurred())
			Expect(body).NotTo(BeEmpty())

			imglist := imgresp{}
			err = json.Unmarshal(body, &imglist)

			Expect(err).NotTo(HaveOccurred())

			// Get back the downloaded document images.
			delete(headers, "Accept")

			type img struct {
				Kind    string
				Content []byte
			}
			imgs := []img{}
			kind := ""
			for _, ref := range imglist.Images {
				_, body, err := mhttp.Get(ref.Href, headers)

				Expect(err).NotTo(HaveOccurred())
				Expect(body).NotTo(BeEmpty())

				if ref.Classifier == "front" {
					kind = "passport"
				} else {
					kind = "selfie"
				}
				imgs = append(imgs, img{
					Kind:    kind,
					Content: body,
				})
			}

			Expect(imgs).NotTo(BeEmpty())
			Expect(imgs).To(HaveLen(2))

			// The received images must be equal to the sent ones.
			for _, i := range imgs {
				switch i.Kind {
				case "passport":
					Expect(i.Content).To(Equal(testPassport))
				case "selfie":
					Expect(i.Content).To(Equal(testSelfie))
				default:
					Fail(fmt.Sprintf("unexpected kind of the document: %s", i.Kind))
				}
			}
		})
	})
})
