package main

import (
	"io/ioutil"
	"log"
	"modulus/kyc/common"
	"modulus/kyc/integrations/shuftipro"
	"testing"
	"time"
)

func Test_Shufti(*testing.T) {

	id, err := ioutil.ReadFile("../test_data/realId.jpg")
	if err != nil {
		panic(err)
	}

	face, err := ioutil.ReadFile("../test_data/realFace.jpg")
	if err != nil {
		panic(err)
	}

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
		IDCard: &common.IDCard{
			CountryAlpha2: "RU",
			Image: &common.DocumentFile{
				Filename:    "passport.png",
				ContentType: "image/png",
				Data:        id,
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "passport.png",
				ContentType: "image/png",
				Data:        face,
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
