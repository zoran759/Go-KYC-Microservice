package verification

import (
	"gitlab.com/lambospeed/kyc/common"
	"encoding/base64"
)

func MapCustomerToVerificationRequest(customer common.UserData) Request {

	request := Request{
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
		Country:     customer.CountryAlpha2,
		VerificationServices: Services{
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			MiddleName:  customer.MiddleName,
			Address:     customer.AddressString,
			DateOfBirth: customer.DateOfBirth.Format("2006-01-02"),
		},
	}

	if customer.Documents != nil && len(customer.Documents) > 0 {
		document := customer.Documents[0]
		request.VerificationServices.DocumentType = document.Metadata.Type
		request.VerificationServices.DocumentExpiryDate = document.Metadata.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = document.Metadata.Number

		request.VerificationData = Data{
			FrontImage: base64.StdEncoding.EncodeToString(document.Front.Data),
			BackImage:  base64.StdEncoding.EncodeToString(document.Back.Data),
		}
	}

	return request
}
