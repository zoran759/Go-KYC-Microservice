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
			Address:     customer.CurrentAddress.String(),
			DateOfBirth: customer.DateOfBirth.Format("2006-01-02"),
		},
	}

	if customer.Selfie != nil {
		request.VerificationData.FaceImage = base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data)
	}
	if customer.UtilityBill != nil {
		request.VerificationData.UtilityBill = base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data)
	}
	if customer.Passport != nil {
		request.VerificationServices.DocumentType = "passport"
		request.VerificationServices.DocumentExpiryDate = customer.Passport.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.Passport.Number
		if len(customer.Passport.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.Passport.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.Passport.Number[len(customer.Passport.Number)-4:]
		}
		if customer.Passport.Image != nil {
			request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.Passport.Image.Data)
		}

		return request
	}
	if customer.DriverLicense != nil {
		request.VerificationServices.DocumentType = "driving_license"
		request.VerificationServices.DocumentExpiryDate = customer.DriverLicense.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.DriverLicense.Number
		if len(customer.DriverLicense.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.DriverLicense.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.DriverLicense.Number[len(customer.DriverLicense.Number)-4:]
		}
		if customer.DriverLicense.FrontImage != nil {
			request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data)
		}
		if customer.DriverLicense.BackImage != nil {
			request.VerificationData.BackImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.BackImage.Data)
		}

		return request
	}
	if customer.IDCard != nil {
		request.VerificationServices.DocumentType = "id_card"
		request.VerificationServices.DocumentIDNumber = customer.IDCard.Number
		if len(customer.IDCard.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.IDCard.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.IDCard.Number[len(customer.IDCard.Number)-4:]
		}
		if customer.IDCard.Image != nil {
			request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data)
		}

		return request
	}
	if customer.SNILS != nil {
		request.VerificationServices.DocumentType = "id_card"
		request.VerificationServices.DocumentIDNumber = customer.SNILS.Number
		if len(customer.SNILS.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.SNILS.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.SNILS.Number[len(customer.SNILS.Number)-4:]
		}
		if customer.SNILS.Image != nil {
			request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data)
		}

		return request
	}
	if customer.CreditCard != nil {
		request.VerificationServices.DocumentType = "credit_card"
		request.VerificationServices.DocumentExpiryDate = customer.CreditCard.ValidUntil.Format("2006-01-02")
		request.VerificationServices.DocumentIDNumber = customer.CreditCard.Number
		if len(customer.CreditCard.Number) > 5 {
			request.VerificationServices.CardFirst6Digits = customer.CreditCard.Number[:6]
			request.VerificationServices.CardLast4Digits = customer.CreditCard.Number[len(customer.CreditCard.Number)-4:]
		}
		if customer.CreditCard.Image != nil {
			request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(customer.CreditCard.Image.Data)
		}
	}

	return request
}
