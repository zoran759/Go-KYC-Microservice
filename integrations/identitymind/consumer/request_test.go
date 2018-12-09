package consumer

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"modulus/kyc/common"
)

var _ = Describe("Request", func() {
	Describe("setApplicantSSN", func() {
		It("should return empty result", func() {
			ssn := &common.IDCard{}

			r := &KYCRequestData{}
			r.setApplicantSSN(ssn)

			Expect(r.ApplicantSSN).To(BeEmpty())
		})

		It("should return proper result", func() {
			ssn := &common.IDCard{
				Number:        "123456789",
				CountryAlpha2: "US",
				IssuedDate:    common.Time(time.Date(2012, 12, 11, 0, 0, 0, 0, time.UTC)),
			}

			r := &KYCRequestData{}
			r.setApplicantSSN(ssn)

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
				AccountName: "john_doe",
				Email:       "theoretically_this_should_never_happen@superduperpostaldomain.net",
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("email length 65 exceeded limit of 60 symbols"))
		})

		It("should fail with error message", func() {
			customer := &common.UserData{
				AccountName: "john_doe",
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
				AccountName: "john_doe",
				Selfie: &common.Selfie{
					Image: &common.DocumentFile{
						Filename:    "big_selfie",
						ContentType: "application/octet-stream",
						Data:        make([]byte, maxImageDataLength+1),
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
				AccountName: "john_doe",
				Passport: &common.Passport{
					Image: &common.DocumentFile{
						Filename:    "big_front_passport",
						ContentType: "application/octet-stream",
						Data:        make([]byte, maxImageDataLength+1),
					},
				},
			}

			r := &KYCRequestData{}
			err := r.populateFields(customer)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during encoding passport image: too large image file"))
		})

		It("should properly populate request fields", func() {
			selfie := &common.Selfie{
				Image: &common.DocumentFile{
					Filename:    "selfie.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it's your selfie :)`),
				},
			}
			idcard := &common.IDCard{
				Number:        "159133253",
				CountryAlpha2: "US",
				IssuedDate:    common.Time(time.Date(1960, 06, 23, 0, 0, 0, 0, time.UTC)),
				Image: &common.DocumentFile{
					Filename:    "ssn.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake SSN image data`),
				},
			}
			passport := &common.Passport{
				Number:        "987654321",
				CountryAlpha2: "US",
				IssuedDate:    common.Time(time.Date(2015, 06, 15, 0, 0, 0, 0, time.UTC)),
				ValidUntil:    common.Time(time.Date(2025, 06, 14, 0, 0, 0, 0, time.UTC)),
				Image: &common.DocumentFile{
					Filename:    "passport.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake passport image data`),
				},
			}
			drivers := &common.DriverLicense{
				Number:        "210901975",
				CountryAlpha2: "RU",
				IssuedDate:    common.Time(time.Date(2010, 10, 7, 0, 0, 0, 0, time.UTC)),
				ValidUntil:    common.Time(time.Date(2020, 10, 6, 0, 0, 0, 0, time.UTC)),
				FrontImage: &common.DocumentFile{
					Filename:    "drivers_front.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it is a fake drivers front image data`),
				},
				BackImage: &common.DocumentFile{
					Filename:    "drivers_back.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Smile, - it is a fake drivers back image data`),
				},
			}
			residencePermit := &common.ResidencePermit{
				CountryAlpha2: "US",
				IssuedDate:    common.Time(time.Date(2017, 3, 27, 0, 0, 0, 0, time.UTC)),
				ValidUntil:    common.Time(time.Date(2020, 3, 26, 0, 0, 0, 0, time.UTC)),
				Image: &common.DocumentFile{
					Filename:    "permit.jpg",
					ContentType: "image/jpeg",
					Data:        []byte(`Fake residence permit image data`),
				},
			}

			customer := &common.UserData{
				FirstName:            "Jonh",
				MaternalLastName:     "Jones",
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
				Selfie: selfie,
				UtilityBill: &common.UtilityBill{
					CountryAlpha2: "US",
					Image: &common.DocumentFile{
						Filename:    "ub.jpg",
						ContentType: "image/jpeg",
						Data:        []byte(`Fake utility bill permit image data`),
					},
				},
				Business: &common.Business{
					Name:                      "Foobar",
					RegistrationNumber:        "0123456789",
					IncorporationDate:         common.Time(time.Date(2000, 03, 10, 0, 0, 0, 0, time.UTC)),
					IncorporationJurisdiction: "TX",
				},
			}

			r := &KYCRequestData{}

			err := r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.ResidencePermit = residencePermit
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.IDCard = idcard
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			customer.DriverLicense = drivers
			err = r.populateFields(customer)
			Expect(err).ToNot(HaveOccurred())

			r = &KYCRequestData{}

			customer.Passport = passport
			err = r.populateFields(customer)

			Expect(err).ToNot(HaveOccurred())
			Expect(r.AccountName).To(Equal(customer.AccountName))
			Expect(r.Email).To(Equal(customer.Email))
			Expect(r.IP).To(Equal(customer.IPaddress))
			Expect(r.BillingFirstName).To(Equal(customer.FirstName))
			Expect(r.BillingMiddleName).To(Equal(customer.MiddleName))
			Expect(r.BillingLastName).To(Equal(customer.LastName))
			Expect(r.BillingStreet).To(Equal(customer.CurrentAddress.HouseStreetApartment()))
			Expect(r.BillingCountryAlpha2).To(Equal(customer.CountryAlpha2))
			Expect(r.BillingPostalCode).To(Equal(customer.CurrentAddress.PostCode))
			Expect(r.BillingCity).To(Equal(customer.CurrentAddress.Town))
			Expect(r.BillingState).To(Equal(customer.CurrentAddress.State))
			Expect(r.BillingGender).To(Equal(gender2API[customer.Gender]))
			Expect(r.CustomerLongitude).To(Equal(customer.Location.Longitude))
			Expect(r.CustomerLatitude).To(Equal(customer.Location.Latitude))
			Expect(r.CustomerPrimaryPhone).To(Equal(customer.Phone))
			Expect(r.CustomerMobilePhone).To(Equal(customer.MobilePhone))
			scanData, _ := toBase64(passport.Image.Data)
			Expect(r.ScanData).To(Equal(scanData))
			Expect(r.BacksideImageData).To(BeEmpty())
			face, _ := toBase64(selfie.Image.Data)
			Expect(r.FaceImages).To(HaveLen(1))
			Expect(r.FaceImages[0]).To(Equal(face))
			Expect(r.DocumentCountry).To(Equal(passport.CountryAlpha2))
			Expect(r.DocumentType).To(Equal(Passport))
			Expect(r.DateOfBirth).To(Equal(customer.DateOfBirth.Format("2006-01-02")))
			Expect(r.ApplicantSSN).To(Equal(idcard.CountryAlpha2 + ":" + idcard.Number))
			Expect(r.ApplicantSSNLast4).To(Equal(idcard.Number[len(idcard.Number)-4:]))
		})
	})
})
