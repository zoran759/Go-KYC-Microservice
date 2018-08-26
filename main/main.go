package main

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/shuftipro"
	"io/ioutil"
	"log"
	"time"
)

func main() {

	file, _ := ioutil.ReadFile("../../testdata/snils.jpg")

	customer := &common.UserData{
		FirstName:        "Smith",
		PaternalLastName: "James",
		LastName:         "James",
		MiddleName:       "M",
		CountryAlpha2:    "US",
		Phone:            "+3221-214-4456",
		Email:            "jsmith@yahoo.com",
		DateOfBirth:      common.Time(time.Date(1982, 4, 3, 0, 0, 0, 0, time.UTC)),
		SupplementalAddresses: []common.Address{
			{},
		},
		CurrentAddress: common.Address{
			PostCode:          "90010",
			Town:              "Chicago",
			BuildingNumber:    "452",
			FlatNumber:        "2",
			State:             "MI",
			StateProvinceCode: "MI",
			StreetType:        "Avenue",
			Street:            "Michigan",
		},
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:    "PASSPORT",
					Country: "RUS",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        file,
				},
				Back: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        file,
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:    "SELFIE",
					Country: "RUS",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        file,
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:    common.UtilityBill,
					Country: "RUS",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        file,
				},
			},
		},
	}
	/*
		// Example integration for SumSub

		sumsubService := sumsub.New(sumsub.Config{
			Host:             "https://test-api.sumsub.com",
			APIKey:           "GKTBNXNEPJHCXY",
			TimeoutThreshold: int64(time.Hour.Seconds()),
		})

		log.Println(sumsubService.CheckCustomer(customer))
	*/
	/*
		// Example Trulioo integration
		service := trulioo.New(trulioo.Config{
			Host:         "https://api.globaldatacompany.com",
			NAPILogin:    "Modulus.dev",
			NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
		})

		log.Println(service.CheckCustomer(customer))
	*/

	//Example Shufti Pro integration
	service := shuftipro.New(shuftipro.Config{
		Host:        "https://api.shuftipro.com",
		ClientID:    "ac93f3a0fee5afa2d9399d5d0f257dc92bbde89b1e48452e1bfac3c5c1dc99db",
		SecretKey:   "lq34eOTxDe1e6G8a1P7Igqo5YK3ABCDF",
		RedirectURL: "https://api.shuftipro.com",
	})

	log.Println(service.CheckCustomer(customer))
}
