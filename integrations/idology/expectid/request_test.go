package expectid

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/lambospeed/kyc/common"
)

var _ = Describe("Request", func() {
	Describe("makeRequestBody", func() {
		It("should return proper request body", func() {
			customer := &common.UserData{
				FirstName:   "John",
				LastName:    "Smith",
				DateOfBirth: common.Time(time.Date(1975, time.February, 28, 0, 0, 0, 0, time.UTC)),
				CurrentAddress: common.Address{
					CountryAlpha2:     "US",
					State:             "Georgia",
					Town:              "Atlanta",
					Street:            "PeachTree Place",
					BuildingNumber:    "222333",
					PostCode:          "30318",
					StateProvinceCode: "GA",
				},
				Documents: []common.Document{
					common.Document{
						Metadata: common.DocumentMetadata{
							Type:            common.IDCard,
							Country:         "USA",
							Number:          "112223333",
							CardLast4Digits: "3333",
						},
					},
				},
			}

			client := NewClient(Config{
				Host:     "host",
				Username: "test",
				Password: "test",
			})

			body := client.makeRequestBody(customer)

			Expect(body).NotTo(BeEmpty())
			Expect(body).To(Equal("address=222333+PeachTree+Place&altAddress=&altCity=&altState=&altZip=&amount=&c_custom_field_1=&c_custom_field_2=&c_custom_field_3=&captureQueryId=&city=Atlanta&dobMonth=+2&dobYear=1975&email=&firstName=John&idIssuer=&idNumber=&idType=&invoice=&ipAddress=&lastName=Smith&password=test&paymentMethod=&purchaseDate=&score=&shipping=&sku=&ssn=112223333&ssnLast4=3333&state=GA&tax=&telephone=&total=&uid=&username=test&zip=30318"))
		})
	})
})
