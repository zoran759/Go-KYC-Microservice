package consumer

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"gitlab.com/lambospeed/kyc/common"
)

var _ = Describe("Request", func() {
	Describe("setApplicantSSN", func() {
		It("should return empty result", func() {
			doc := &common.DocumentMetadata{
				Type:       common.Passport,
				Country:    "RU",
				DateIssued: common.Time(time.Date(2012, 12, 11, 0, 0, 0, 0, time.UTC)),
				Number:     "9214123456",
			}

			r := &KYCRequestData{}
			r.setApplicantSSN(doc)

			Expect(r.ApplicantSSN).To(BeEmpty())
		})

		It("should return proper result", func() {
			doc := &common.DocumentMetadata{
				Type:             common.IDCard,
				Country:          "US",
				DateIssued:       common.Time(time.Date(2012, 12, 11, 0, 0, 0, 0, time.UTC)),
				Number:           "123456789",
				CardFirst6Digits: "123456",
				CardLast4Digits:  "6789",
			}

			r := &KYCRequestData{}
			r.setApplicantSSN(doc)

			Expect(r.ApplicantSSN).To(Equal("US:123456789"))
		})
	})

	Describe("populateFields", func() {
		It("should fail with error message", func() {
			customer := &common.UserData{
				AccountName: "very long account name that is exceeding the limit of 60 symbols maximum",
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("account length 72 exceeded limit of 60 symbols"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				Email: "theoretically_this_should_never_happen@superduperpostaldomain.net",
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("email length 65 exceeded limit of 60 symbols"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				CurrentAddress: common.Address{
					BuildingNumber: "0123456789",
					Street:         "Very Loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong Street Name",
					FlatNumber:     "123",
				},
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("street address length 101 exceeded limit of 100 symbols"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				Documents: []common.Document{
					common.Document{
						Front: &common.DocumentFile{
							Filename:    "big_selfie",
							ContentType: "application/octet-stream",
							Data:        make([]byte, maxImageDataLength+1),
						},
						Metadata: common.DocumentMetadata{
							Type: common.Selfie,
						},
					},
				},
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during encoding selfi image data: too large image file"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				Documents: []common.Document{
					common.Document{
						Front: &common.DocumentFile{
							Filename:    "big_front_passport",
							ContentType: "application/octet-stream",
							Data:        make([]byte, maxImageDataLength+1),
						},
						Metadata: common.DocumentMetadata{
							Type: common.Passport,
						},
					},
				},
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during encoding passport front image: too large image file"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				Documents: []common.Document{
					common.Document{
						Front: &common.DocumentFile{},
						Back: &common.DocumentFile{
							Filename:    "big_back_passport",
							ContentType: "application/octet-stream",
							Data:        make([]byte, maxImageDataLength+1),
						},
						Metadata: common.DocumentMetadata{
							Type: common.Passport,
						},
					},
				},
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during encoding passport back image: too large image file"))
		})

		It("should properly populate request fields", func() {
			selfie := common.Document{
				Metadata: common.DocumentMetadata{
					Type: common.Selfie,
				},
				Front: &common.DocumentFile{
					Filename:    "selfie.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it's your selfie :)`),
				},
			}
			idcard := common.Document{
				Metadata: common.DocumentMetadata{
					Type:             common.IDCard,
					Country:          "US",
					DateIssued:       common.Time(time.Date(1960, 06, 23, 0, 0, 0, 0, time.UTC)),
					Number:           "159133253",
					CardFirst6Digits: "159133",
					CardLast4Digits:  "3253",
				},
				Front: &common.DocumentFile{
					Filename:    "ssn.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake SSN image data`),
				},
			}
			passport := common.Document{
				Metadata: common.DocumentMetadata{
					Type:             common.Passport,
					Country:          "US",
					DateIssued:       common.Time(time.Date(2015, 06, 15, 0, 0, 0, 0, time.UTC)),
					ValidUntil:       common.Time(time.Date(2025, 06, 14, 0, 0, 0, 0, time.UTC)),
					Number:           "987654321",
					CardFirst6Digits: "987654",
					CardLast4Digits:  "4321",
				},
				Front: &common.DocumentFile{
					Filename:    "passport_front.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake passport front image data`),
				},
				Back: &common.DocumentFile{
					Filename:    "passport_back.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake passport back image data`),
				},
			}
			drivers := common.Document{
				Metadata: common.DocumentMetadata{
					Type:             common.Drivers,
					Country:          "RU",
					DateIssued:       common.Time(time.Date(2010, 10, 7, 0, 0, 0, 0, time.UTC)),
					ValidUntil:       common.Time(time.Date(2020, 10, 6, 0, 0, 0, 0, time.UTC)),
					Number:           "210901975",
					CardFirst6Digits: "210901",
					CardLast4Digits:  "1975",
				},
				Front: &common.DocumentFile{
					Filename:    "drivers_front.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it is a fake drivers front image data`),
				},
				Back: &common.DocumentFile{
					Filename:    "drivers_back.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it is a fake drivers back image data`),
				},
			}
			residencePermit := common.Document{
				Metadata: common.DocumentMetadata{
					Type:             common.ResidencePermit,
					Country:          "US",
					DateIssued:       common.Time(time.Date(2017, 3, 27, 0, 0, 0, 0, time.UTC)),
					ValidUntil:       common.Time(time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC)),
					Number:           "197509210",
					CardFirst6Digits: "197509",
					CardLast4Digits:  "9210",
				},
				Front: &common.DocumentFile{
					Filename:    "permit_front.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake residence permit front image data`),
				},
				Back: &common.DocumentFile{
					Filename:    "permit_back.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake residence permit back image data`),
				},
			}

			customer := &common.UserData{
				FirstName:            "Jonh",
				PaternalLastName:     "Jones",
				LastName:             "Doe",
				MiddleName:           "Lee",
				LegalName:            "John Doe",
				LatinISO1Name:        "John Doe",
				AccountName:          "john_doe",
				Email:                "john.doe@email.doe",
				IPaddress:            "10.111.1.101",
				Gender:               common.Male,
				DateOfBirth:          common.Time(time.Date(1960, 06, 15, 0, 0, 0, 0, time.UTC)),
				PlaceOfBirth:         "Pittsburgh",
				CountryOfBirthAlpha2: "US",
				StateOfBirth:         "Pennsylvania",
				CountryAlpha2:        "US",
				Nationality:          "American",
				Phone:                "412-325-3686",
				MobilePhone:          "+110-325-368-6237",
				Location: &common.Location{
					Latitude:  "40.465769",
					Longitude: "-80.029928",
				},
				CurrentAddress: common.Address{
					CountryAlpha2:     "US",
					County:            "Neverhood",
					State:             "Pennsylvania",
					Town:              "Pittsburgh",
					Suburb:            "Cubby",
					Street:            "Gifford St",
					StreetType:        "Street",
					SubStreet:         "None",
					BuildingName:      "Home",
					BuildingNumber:    "1324",
					FlatNumber:        "1",
					PostOfficeBox:     "",
					PostCode:          "15212",
					StateProvinceCode: "PA",
					StartDate:         common.Time(time.Date(1960, 06, 15, 0, 0, 0, 0, time.UTC)),
				},
				SupplementalAddresses: []common.Address{
					common.Address{
						CountryAlpha2:     "US",
						County:            "Shire",
						State:             "Texas",
						Town:              "Austin",
						Suburb:            "Cubby",
						Street:            "Crescendo Ln",
						StreetType:        "Lane",
						SubStreet:         "None",
						BuildingName:      "Bungalo",
						BuildingNumber:    "10501",
						FlatNumber:        "1",
						PostOfficeBox:     "",
						PostCode:          "78747",
						StateProvinceCode: "TX",
						StartDate:         common.Time(time.Date(2000, 04, 22, 0, 0, 0, 0, time.UTC)),
					},
				},
				Documents: []common.Document{
					selfie,
					common.Document{
						Metadata: common.DocumentMetadata{
							Type: common.Selfie,
						},
						Front: &common.DocumentFile{
							Filename:    "selfie.jpg",
							ContentType: "image/jpeg",
							Data:        []byte{},
						},
					},
					common.Document{
						Metadata: common.DocumentMetadata{
							Type:             common.UtilityBill,
							Country:          "US",
							DateIssued:       common.Time(time.Date(2018, 9, 3, 0, 0, 0, 0, time.UTC)),
							Number:           "0123456001",
							CardFirst6Digits: "012345",
							CardLast4Digits:  "6001",
						},
						Front: &common.DocumentFile{
							Filename:    "ub_front.jpg",
							ContentType: "image/jpeg",
							Data:        []byte(`Fake utility bill permit front image data`),
						},
					},
				},
				Business: common.Business{
					Name:                      "Foobar",
					RegistrationNumber:        "0123456789",
					IncorporationDate:         common.Time(time.Date(2000, 03, 10, 0, 0, 0, 0, time.UTC)),
					IncorporationJurisdiction: "TX",
				},
			}

			r := &KYCRequestData{}

			err := r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.Documents = append(customer.Documents, residencePermit)
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.Documents = append(customer.Documents, idcard)
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.Documents = append(customer.Documents, drivers)
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.Documents = append(customer.Documents, passport)
			err = r.populateFields(customer)

			Expect(err).ToNot(HaveOccurred())
			Expect(r.AccountName).To(Equal(customer.AccountName))
			Expect(r.Email).To(Equal(customer.Email))
			Expect(r.IP).To(Equal(customer.IPaddress))
			Expect(r.BillingFirstName).To(Equal(customer.FirstName))
			Expect(r.BillingMiddleName).To(Equal(customer.MiddleName))
			Expect(r.BillingLastName).To(Equal(customer.LastName))
			Expect(r.BillingStreet).To(Equal(customer.CurrentAddress.BuildingNumber + " " + customer.CurrentAddress.Street + " " + customer.CurrentAddress.FlatNumber))
			Expect(r.BillingCountryAlpha2).To(Equal(customer.CountryAlpha2))
			Expect(r.BillingPostalCode).To(Equal(customer.CurrentAddress.PostCode))
			Expect(r.BillingCity).To(Equal(customer.CurrentAddress.Town))
			Expect(r.BillingState).To(Equal(customer.CurrentAddress.State))
			Expect(r.BillingGender).To(Equal(genderMap[customer.Gender]))
			Expect(r.CustomerLongitude).To(Equal(customer.Location.Longitude))
			Expect(r.CustomerLatitude).To(Equal(customer.Location.Latitude))
			Expect(r.CustomerPrimaryPhone).To(Equal(customer.Phone))
			Expect(r.CustomerMobilePhone).To(Equal(customer.MobilePhone))
			scanData, _ := toBase64(passport.Front.Data)
			Expect(r.ScanData).To(Equal(scanData))
			backside, _ := toBase64(passport.Back.Data)
			Expect(r.BacksideImageData).To(Equal(backside))
			face, _ := toBase64(selfie.Front.Data)
			Expect(r.FaceImages).To(HaveLen(10))
			Expect(r.FaceImages[0]).To(Equal(face))
			Expect(r.DocumentCountry).To(Equal(passport.Metadata.Country))
			Expect(r.DocumentType).To(Equal(Passport))
			Expect(r.DateOfBirth).To(Equal(customer.DateOfBirth.Format("2006-01-02")))
			Expect(r.ApplicantSSN).To(Equal(idcard.Metadata.Country + ":" + idcard.Metadata.Number))
			Expect(r.ApplicantSSNLast4).To(Equal(idcard.Metadata.CardLast4Digits))
		})
	})
})
