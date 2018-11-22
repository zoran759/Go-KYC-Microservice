package verification

import (
	"encoding/base64"

	"modulus/kyc/common"
)

// MapCustomerToVerificationRequest maps the values of common data to the service specific values.
func MapCustomerToVerificationRequest(customer common.UserData) Request {

	request := Request{
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
		Country:     customer.CountryAlpha2,
		VerificationServices: Services{
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			MiddleName:  customer.MiddleName,
			DateOfBirth: customer.DateOfBirth.Format("2006-01-02"),
		},
	}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		request.VerificationData.FaceImage = base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data)
	}
	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		request.VerificationServices.Address = customer.CurrentAddress.String()
		request.VerificationData.UtilityBill = base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data)
	}
	if customer.Passport != nil && customer.Passport.Image != nil {
		request.VerificationServices.DocumentType = "passport"
		request.VerificationServices.DocumentExpiryDate = customer.Passport.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.Passport.Number
		request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.Passport.Image.Data)

		return request
	}
	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		request.VerificationServices.DocumentType = "driving_license"
		request.VerificationServices.DocumentExpiryDate = customer.DriverLicense.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.DriverLicense.Number
		request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data)
		if customer.DriverLicense.BackImage != nil {
			request.VerificationData.BackImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.BackImage.Data)
		}

		return request
	}
	if customer.IDCard != nil && customer.IDCard.Image != nil {
		request.VerificationServices.DocumentType = "id_card"
		request.VerificationServices.DocumentIDNumber = customer.IDCard.Number
		request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data)

		return request
	}
	if customer.SNILS != nil && customer.SNILS.Image != nil {
		request.VerificationServices.DocumentType = "id_card"
		request.VerificationServices.DocumentIDNumber = customer.SNILS.Number
		request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data)

		return request
	}
	if customer.CreditCard != nil && customer.CreditCard.Image != nil {
		request.VerificationServices.DocumentType = "credit_card"
		request.VerificationServices.DocumentExpiryDate = customer.CreditCard.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.CreditCard.Number
		if len(customer.CreditCard.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.CreditCard.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.CreditCard.Number[len(customer.CreditCard.Number)-4:]
		}
		request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.CreditCard.Image.Data)
	}

	return request
}
