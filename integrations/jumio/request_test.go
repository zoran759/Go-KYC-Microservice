package jumio

import (
	"modulus/kyc/common"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Request", func() {
	Describe("populateFields", func() {
		var r = &Request{}

		It("should fail with the empty user data", func() {
			customer := &common.UserData{}

			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("missing acceptable document for the verification (anyone of passport, driving license or id card)"))
		})

		Context("using erroneous customer.Selfie", func() {
			It("should fail with wrong image format", func() {
				customer := &common.UserData{
					Selfie: &common.Selfie{
						Image: &common.DocumentFile{
							Filename:    "selfie.bmp",
							ContentType: "image/bmp",
							Data:        []byte("Smile, it's your selfie :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable selfie image format: image/bmp"))
			})

			It("should fail with too long image data", func() {
				customer := &common.UserData{
					Selfie: &common.Selfie{
						Image: &common.DocumentFile{
							Filename:    "selfie.png",
							ContentType: "image/png",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding selfi image data: too large image file"))
			})
		})

		Context("using erroneous customer.Passport", func() {
			It("should fail with wrong image format", func() {
				customer := &common.UserData{
					Passport: &common.Passport{
						Image: &common.DocumentFile{
							Filename:    "passport.bmp",
							ContentType: "image/bmp",
							Data:        []byte("Smile, it's your passport :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable passport image format: image/bmp"))
			})

			It("should fail with too long image data", func() {
				customer := &common.UserData{
					Passport: &common.Passport{
						Image: &common.DocumentFile{
							Filename:    "passport.png",
							ContentType: "image/png",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding passport image: too large image file"))
			})
		})

		Context("using erroneous customer.DriverLicense", func() {
			It("should fail with wrong front image format", func() {
				customer := &common.UserData{
					DriverLicense: &common.DriverLicense{
						FrontImage: &common.DocumentFile{
							Filename:    "driver license.bmp",
							ContentType: "image/bmp",
							Data:        []byte("Smile, it's your driver license :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable driver license front image format: image/bmp"))
			})

			It("should fail with too long front image data", func() {
				customer := &common.UserData{
					DriverLicense: &common.DriverLicense{
						FrontImage: &common.DocumentFile{
							Filename:    "driver license.png",
							ContentType: "image/png",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding driver license front image: too large image file"))
			})

			It("should fail with wrong back image format", func() {
				customer := &common.UserData{
					DriverLicense: &common.DriverLicense{
						FrontImage: &common.DocumentFile{
							Filename:    "driver license front.jpg",
							ContentType: "image/jpeg",
							Data:        []byte("Smile, it's your driver license :)"),
						},
						BackImage: &common.DocumentFile{
							Filename:    "driver license back.gif",
							ContentType: "image/gif",
							Data:        []byte("Smile, it's your driver license :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable driver license back image format: image/gif"))
			})

			It("should fail with too long back image data", func() {
				customer := &common.UserData{
					DriverLicense: &common.DriverLicense{
						FrontImage: &common.DocumentFile{
							Filename:    "driver license front.jpg",
							ContentType: "image/jpeg",
							Data:        []byte("Smile, it's your driver license :)"),
						},
						BackImage: &common.DocumentFile{
							Filename:    "driver license back.jpg",
							ContentType: "image/jpeg",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding driver license back image: too large image file"))
			})
		})

		Context("using erroneous customer.IDCard", func() {
			It("should fail with wrong image format", func() {
				customer := &common.UserData{
					IDCard: &common.IDCard{
						Image: &common.DocumentFile{
							Filename:    "id card.bmp",
							ContentType: "image/bmp",
							Data:        []byte("Smile, it's your id card :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable id card image format: image/bmp"))
			})

			It("should fail with too long image data", func() {
				customer := &common.UserData{
					IDCard: &common.IDCard{
						Image: &common.DocumentFile{
							Filename:    "id card.png",
							ContentType: "image/png",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding id card image: too large image file"))
			})
		})

		Context("using erroneous customer.SNILS", func() {
			It("should fail with wrong image format", func() {
				customer := &common.UserData{
					SNILS: &common.SNILS{
						Image: &common.DocumentFile{
							Filename:    "SNILS.bmp",
							ContentType: "image/bmp",
							Data:        []byte("Smile, it's your id card :)"),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("unacceptable SNILS image format: image/bmp"))
			})

			It("should fail with too long image data", func() {
				customer := &common.UserData{
					SNILS: &common.SNILS{
						Image: &common.DocumentFile{
							Filename:    "SNILS.png",
							ContentType: "image/png",
							Data:        make([]byte, maxImageDataLength+1),
						},
					},
				}

				err := r.populateFields(customer)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("during encoding SNILS image: too large image file"))
			})
		})

		It("should fail with the valid customer.Selfie but without any document", func() {
			customer := &common.UserData{
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "selfie.png",
						ContentType: "image/png",
						Data:        []byte{},
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("missing acceptable document for the verification (anyone of passport, driving license or id card)"))
		})

		It("should success with the valid customer.Passport", func() {
			customer := &common.UserData{
				CountryAlpha2: "US",
				Passport: &common.Passport{
					Number:     "1234567890",
					State:      "WA",
					IssuedDate: common.Time(time.Date(2015, 05, 15, 0, 0, 0, 0, time.UTC)),
					ValidUntil: common.Time(time.Date(2025, 05, 14, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "passport.png",
						ContentType: "image/png",
						Data:        []byte("Smile, it's your passport :)"),
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.MerchantIDScanReference).NotTo(BeEmpty())
			Expect(r.FrontsideImageMimeType).To(Equal("image/png"))
			Expect(r.Country).To(Equal("USA"))
			Expect(r.IDType).To(Equal(Passport))
			Expect(r.Expiry).To(Equal(common.Time(time.Date(2025, 05, 14, 0, 0, 0, 0, time.UTC)).Format("2006-01-02")))
			Expect(r.Number).To(Equal("1234567890"))
			Expect(r.USState).To(Equal("WA"))
		})

		It("should success with the valid customer.DriverLicense", func() {
			customer := &common.UserData{
				CountryAlpha2: "US",
				DriverLicense: &common.DriverLicense{
					Number:     "123456789",
					State:      "WA",
					IssuedDate: common.Time(time.Date(2012, 12, 22, 0, 0, 0, 0, time.UTC)),
					ValidUntil: common.Time(time.Date(2022, 12, 21, 0, 0, 0, 0, time.UTC)),
					FrontImage: &common.DocumentFile{
						Filename:    "driver license front.jpg",
						ContentType: "image/jpeg",
						Data:        []byte("Smile, it's your driver license :)"),
					},
					BackImage: &common.DocumentFile{
						Filename:    "driver license back.jpg",
						ContentType: "image/jpeg",
						Data:        []byte("Smile, it's your driver license :)"),
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.MerchantIDScanReference).NotTo(BeEmpty())
			Expect(r.FrontsideImageMimeType).To(Equal("image/jpeg"))
			Expect(r.BacksideImageMimeType).To(Equal("image/jpeg"))
			Expect(r.Country).To(Equal("USA"))
			Expect(r.IDType).To(Equal(DrivingLicense))
			Expect(r.Expiry).To(Equal(common.Time(time.Date(2022, 12, 21, 0, 0, 0, 0, time.UTC)).Format("2006-01-02")))
			Expect(r.Number).To(Equal("123456789"))
			Expect(r.USState).To(Equal("WA"))
		})

		It("should success with the valid customer.IDCard", func() {
			customer := &common.UserData{
				CountryAlpha2: "CA",
				IDCard: &common.IDCard{
					Number:     "0123456789",
					IssuedDate: common.Time(time.Date(2018, 07, 13, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "id card.jpg",
						ContentType: "image/jpeg",
						Data:        []byte("Smile, it's your id card :)"),
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.MerchantIDScanReference).NotTo(BeEmpty())
			Expect(r.FrontsideImageMimeType).To(Equal("image/jpeg"))
			Expect(r.Country).To(Equal("CAN"))
			Expect(r.IDType).To(Equal(IDCard))
			Expect(r.Number).To(Equal("0123456789"))
		})

		It("should success with the valid customer.SNILS", func() {
			customer := &common.UserData{
				CountryAlpha2: "RU",
				SNILS: &common.SNILS{
					Number:     "11112223333",
					IssuedDate: common.Time(time.Date(2018, 03, 14, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "SNILS.png",
						ContentType: "image/png",
						Data:        []byte("Smile, it's your id card :)"),
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.MerchantIDScanReference).NotTo(BeEmpty())
			Expect(r.FrontsideImageMimeType).To(Equal("image/png"))
			Expect(r.Country).To(Equal("RUS"))
			Expect(r.IDType).To(Equal(IDCard))
			Expect(r.Number).To(Equal("11112223333"))
		})

		It("should success with the valid user data", func() {
			customer := &common.UserData{
				FirstName:     "Bruce",
				LastName:      "Wayne",
				DateOfBirth:   common.Time(time.Date(1950, 03, 17, 0, 0, 0, 0, time.UTC)),
				CountryAlpha2: "US",
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "batman.png",
						ContentType: "image/png",
						Data:        []byte{},
					},
				},
				Passport: &common.Passport{
					Number:     "1234567890",
					State:      "WA",
					IssuedDate: common.Time(time.Date(2010, 05, 15, 0, 0, 0, 0, time.UTC)),
					ValidUntil: common.Time(time.Date(2020, 05, 14, 0, 0, 0, 0, time.UTC)),
					Image: &common.DocumentFile{
						Filename:    "passport.png",
						ContentType: "image/png",
						Data:        []byte{},
					},
				},
			}

			err := r.populateFields(customer)

			Expect(err).NotTo(HaveOccurred())
			Expect(r.MerchantIDScanReference).NotTo(BeEmpty())
			Expect(r.FirstName).To(Equal("Bruce"))
			Expect(r.LastName).To(Equal("Wayne"))
			Expect(r.DOB).To(Equal(common.Time(time.Date(1950, 03, 17, 0, 0, 0, 0, time.UTC)).Format("2006-01-02")))
			Expect(r.FrontsideImageMimeType).To(Equal("image/png"))
			Expect(r.Country).To(Equal("USA"))
			Expect(r.IDType).To(Equal(Passport))
			Expect(r.Expiry).To(Equal(common.Time(time.Date(2020, 05, 14, 0, 0, 0, 0, time.UTC)).Format("2006-01-02")))
			Expect(r.Number).To(Equal("1234567890"))
			Expect(r.USState).To(Equal("WA"))
		})
	})
})
