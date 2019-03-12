package shuftipro

import (
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	c := NewClient(Config{
		Host:        "host",
		ClientID:    "client_id",
		SecretKey:   "secret_key",
		CallbackURL: "callback_url",
	})

	selfie := &common.Selfie{
		Image: &common.DocumentFile{
			Filename:    "selfie.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	passport := &common.Passport{
		Number:     "123456",
		IssuedDate: common.Time(time.Date(2000, 9, 23, 0, 0, 0, 0, time.UTC)),
		ValidUntil: common.Time(time.Date(2025, 8, 28, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "passport.png",
			ContentType: "image/png",
			Data:        []byte("fake image data"),
		},
	}

	driverLicense := &common.DriverLicense{
		Number:     "ABC456798",
		IssuedDate: common.Time(time.Date(2015, 5, 15, 0, 0, 0, 0, time.UTC)),
		ValidUntil: common.Time(time.Date(2025, 5, 14, 0, 0, 0, 0, time.UTC)),
		FrontImage: &common.DocumentFile{
			Filename:    "driver_front.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
		BackImage: &common.DocumentFile{
			Filename:    "driver_back.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	idcard := &common.IDCard{
		Number:     "789 123 654",
		IssuedDate: common.Time(time.Date(1985, 9, 21, 0, 0, 0, 0, time.UTC)),
		ValidUntil: common.Time(time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "idcard.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	creditCard := &common.CreditCard{
		Number:     "4000 1234 5678 9010",
		ValidUntil: common.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "ccard.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	debitCard := &common.DebitCard{
		Number:     "4000 1234 5678 9010",
		ValidUntil: common.Time(time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)),
		Image: &common.DocumentFile{
			Filename:    "dcard.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	utilityBill := &common.UtilityBill{
		Image: &common.DocumentFile{
			Filename:    "utilbill.jpg",
			ContentType: "image/jpeg",
			Data:        []byte("fake image data"),
		},
	}

	address := common.Address{
		CountryAlpha2:  "GB",
		County:         "Westminster",
		Town:           "London",
		Street:         "Downing St",
		BuildingNumber: "10",
		PostCode:       "SW1A 2AA",
	}

	type testCase struct {
		name     string
		customer *common.UserData
		request  *Request
		err      error
	}

	testCases := []testCase{
		testCase{
			name: "Main customer info",
			customer: &common.UserData{
				FirstName:     "John",
				MiddleName:    "Middleman",
				LastName:      "Doe",
				CountryAlpha2: "US",
				Email:         "john.doe@example.com",
				DateOfBirth:   common.Time(time.Date(1980, 8, 28, 0, 0, 0, 0, time.UTC)),
				Selfie:        selfie,
			},
			request: &Request{
				CountryAlpha2: "US",
				Email:         "john.doe@example.com",
				CallbackURL:   c.callbackURL,
				BackgroundChecks: &BackgroundChecks{
					Name: Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
					DateOfBirth: "1980-08-28",
				},
				Face: &Face{
					Proof: toBase64(selfie.Image),
				},
			},
		},
		testCase{
			name: "Passport",
			customer: &common.UserData{
				FirstName:     "John",
				MiddleName:    "Middleman",
				LastName:      "Doe",
				CountryAlpha2: "US",
				Email:         "john.doe@example.com",
				DateOfBirth:   common.Time(time.Date(1980, 8, 28, 0, 0, 0, 0, time.UTC)),
				Selfie:        selfie,
				Passport:      passport,
			},
			request: &Request{
				CountryAlpha2: "US",
				Email:         "john.doe@example.com",
				CallbackURL:   c.callbackURL,
				BackgroundChecks: &BackgroundChecks{
					Name: Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
					DateOfBirth: "1980-08-28",
				},
				Face: &Face{
					Proof: toBase64(selfie.Image),
				},
				Document: &Document{
					Proof:          toBase64(passport.Image),
					SupportedTypes: []DocumentType{Passport},
					Name: &Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
					DateOfBirth: "1980-08-28",
					Number:      passport.Number,
					IssueDate:   "2000-09-23",
					ExpiryDate:  "2025-08-28",
				},
			},
		},
		testCase{
			name: "Driver license",
			customer: &common.UserData{
				FirstName:     "John",
				LastName:      "Doe",
				DriverLicense: driverLicense,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Document: &Document{
					Proof:           toBase64(driverLicense.FrontImage),
					AdditionalProof: toBase64(driverLicense.BackImage),
					SupportedTypes:  []DocumentType{DrivingLicense},
					Name: &Name{
						FirstName: "John",
						LastName:  "Doe",
					},
					Number:     driverLicense.Number,
					IssueDate:  "2015-05-15",
					ExpiryDate: "2025-05-14",
				},
			},
		},
		testCase{
			name: "ID card",
			customer: &common.UserData{
				FirstName: "John",
				LastName:  "Doe",
				IDCard:    idcard,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Document: &Document{
					Proof:          toBase64(idcard.Image),
					SupportedTypes: []DocumentType{IDcard},
					Name: &Name{
						FirstName: "John",
						LastName:  "Doe",
					},
					Number:     idcard.Number,
					IssueDate:  "1985-09-21",
					ExpiryDate: "2025-09-20",
				},
			},
		},
		testCase{
			name: "Credit card",
			customer: &common.UserData{
				FirstName:  "John",
				LastName:   "Doe",
				CreditCard: creditCard,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Document: &Document{
					Proof:          toBase64(creditCard.Image),
					SupportedTypes: []DocumentType{CreditOrDebitCard},
					Name: &Name{
						FirstName: "John",
						LastName:  "Doe",
					},
					Number:     creditCard.Number,
					ExpiryDate: "2023-10-01",
				},
			},
		},
		testCase{
			name: "Debit card",
			customer: &common.UserData{
				FirstName: "John",
				LastName:  "Doe",
				DebitCard: debitCard,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Document: &Document{
					Proof:          toBase64(debitCard.Image),
					SupportedTypes: []DocumentType{CreditOrDebitCard},
					Name: &Name{
						FirstName: "John",
						LastName:  "Doe",
					},
					Number:     debitCard.Number,
					ExpiryDate: "2023-10-01",
				},
			},
		},
		testCase{
			name: "Address with utility bill",
			customer: &common.UserData{
				FirstName:      "John",
				MiddleName:     "Middleman",
				LastName:       "Doe",
				CurrentAddress: address,
				UtilityBill:    utilityBill,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Address: &Address{
					Proof:          toBase64(utilityBill.Image),
					SupportedTypes: []DocumentType{UtilityBill},
					FullAddress:    address.String(),
					Name: &Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
				},
			},
		},
		testCase{
			name: "Address with id card",
			customer: &common.UserData{
				FirstName:      "John",
				MiddleName:     "Middleman",
				LastName:       "Doe",
				CurrentAddress: address,
				IDCard:         idcard,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
				Document: &Document{
					Proof:          toBase64(idcard.Image),
					SupportedTypes: []DocumentType{IDcard},
					Name: &Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
					Number:     idcard.Number,
					IssueDate:  "1985-09-21",
					ExpiryDate: "2025-09-20",
				},
				Address: &Address{
					Proof:          toBase64(idcard.Image),
					SupportedTypes: []DocumentType{IDcard},
					FullAddress:    address.String(),
					Name: &Name{
						FirstName:  "John",
						MiddleName: "Middleman",
						LastName:   "Doe",
					},
				},
			},
		},
		testCase{
			name: "Address without proof",
			customer: &common.UserData{
				FirstName:      "John",
				MiddleName:     "Middleman",
				LastName:       "Doe",
				CurrentAddress: address,
			},
			request: &Request{
				CallbackURL: c.callbackURL,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := c.NewRequest(tc.customer)
			// Reference is uniquely generated for each request.
			tc.request.Reference = req.Reference
			assert.Equal(t, tc.request, req)
			assert.Equal(t, tc.err, err)
		})
	}
}
