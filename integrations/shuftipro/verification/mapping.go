package verification

import (
	"gitlab.com/lambospeed/kyc/common"
	"encoding/base64"
	"strings"
)

func MapCustomerToVerificationRequest(customer common.UserData) Request {

	request := Request{
		Email:       customer.Email,
		PhoneNumber: customer.Phone,
		Country:     strings.ToLower(customer.CountryAlpha2),
		VerificationServices: Services{
			FirstName:   customer.FirstName,
			LastName:    customer.LastName,
			MiddleName:  customer.MiddleName,
			Address:     customer.AddressString,
			DateOfBirth: customer.DateOfBirth.Format("2006-01-02"),
		},
	}

	if customer.Documents != nil && len(customer.Documents) > 0 {
		for _, document := range customer.Documents {
			if request.VerificationServices.DocumentType != "" &&
				request.VerificationData.FaceImage != "" &&
				request.VerificationData.UtilityBill != "" {
				break
			} else if request.VerificationData.FaceImage == "" && document.Metadata.Type == common.Selfie {
				if document.Front != nil {
					request.VerificationData.FaceImage = base64.StdEncoding.EncodeToString(document.Front.Data)
				}
			} else if request.VerificationData.UtilityBill == "" && document.Metadata.Type == common.UtilityBill {
				if document.Front != nil {
					request.VerificationData.FaceImage = base64.StdEncoding.EncodeToString(document.Front.Data)
				}
			} else if request.VerificationServices.DocumentType == "" {
				if mappedType := mapDocumentType(document.Metadata.Type); mappedType != "" {
					request.VerificationServices.DocumentType = mappedType
					request.VerificationServices.DocumentExpiryDate = document.Metadata.ValidUntil.Format("2006-01-02")
					request.VerificationServices.DocumentIDNumber = document.Metadata.Number
					request.VerificationServices.CardFirst6Digits = document.Metadata.CardFirst6Digits
					request.VerificationServices.CardLast4Digits = document.Metadata.CardLast4Digits

					if document.Front != nil && document.Front.Data != nil {
						request.VerificationData.FrontImage = base64.StdEncoding.EncodeToString(document.Front.Data)
					}

					if document.Back != nil && document.Back.Data != nil {
						request.VerificationData.BackImage = base64.StdEncoding.EncodeToString(document.Back.Data)
					}
				}
			}
		}

	}

	return request
}

func mapDocumentType(documentType common.DocumentType) string {
	switch documentType {
	case common.Passport:
		return "passport"
	case common.Drivers:
		return "driving_license"
	case common.IDCard:
		return "id_card"
	case common.BankCard:
		return "credit_card"
	default:
		return ""
	}
}
