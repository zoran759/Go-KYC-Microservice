package main

import (
	"gitlab.com/modulusglobal/kyc/common"
	"gitlab.com/modulusglobal/kyc/integrations/trulioo"
	"log"
	"time"
)

func main() {

	customer := &common.UserData{
		FirstName:        "Smith",
		PaternalLastName: "James",
		LastName:         "James",
		MiddleName:       "M",
		CountryAlpha2:    "US",
		Phone:            "221-214-4456",
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
					Data: []byte{0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff,
						0x92, 0x92, 0x92, 0x92, 0x92,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00},
				},
				Back: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data: []byte{0xff, 0xff, 0xff, 0xff, 0xff,
						0xff, 0xff, 0xff, 0xff, 0xff,
						0x92, 0x92, 0x92, 0x92, 0x92,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00,
						0x00, 0x00, 0x00, 0x00, 0x00},
				},
			},
		},
	}

	/*

		// Example integration for SumSub

		//
		//sumsubService := sumsub.New(sumsub.Config{
		//	Host:             "https://test-api.sumsub.com",
		//	APIKey:           "GKTBNXNEPJHCXY",
		//	TimeoutThreshold: int64(time.Hour.Seconds()),
		//})
		//
		//log.Println(sumsubService.CheckCustomer(customer))
	*/

	// Example Trulioo integration
	service := trulioo.New(trulioo.Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "Modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	log.Println(service.CheckCustomer(customer))
}
