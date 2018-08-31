package verification

import (
	"encoding/base64"
	"fmt"
	"gitlab.com/lambospeed/kyc/common"
	"time"
)

func MapCustomerToCreateUserRequest(customer common.UserData) CreateUserRequest {
	phoneNumbers := make([]string, 0)

	if customer.Phone != "" {
		phoneNumbers = append(phoneNumbers, customer.Phone)
	}
	if customer.MobilePhone != "" {
		phoneNumbers = append(phoneNumbers, customer.MobilePhone)
	}

	return CreateUserRequest{
		Logins: []Login{
			{
				Email: customer.Email,
				Scope: "READ",
			},
		},
		PhoneNumbers: phoneNumbers,
		LegalNames: []string{
			fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
		},
		Documents: mapCustomerDocuments(customer),
		Extra: Extra{
			CIPTag:     1,
			IsBusiness: false,
		},
	}
}

func mapCustomerDocuments(customer common.UserData) []Document {
	document := Document{
		OwnerName:    fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
		Email:        customer.Email,
		PhoneNumber:  customer.Phone,
		IPAddress:    "0.0.0.0",
		EntityType:   mapCustomerGender(customer.Gender),
		EntityScope:  "Not Known",
		DayOfBirth:   time.Time(customer.DateOfBirth).Day(),
		MonthOfBirth: int(time.Time(customer.DateOfBirth).Month()),
		YearOfBirth:  time.Time(customer.DateOfBirth).Year(),
		AddressStreet: fmt.Sprintf(
			"%s %s %s",
			customer.CurrentAddress.BuildingNumber,
			customer.CurrentAddress.Street,
			customer.CurrentAddress.StreetType,
		),
		AddressCity:        customer.CurrentAddress.Town,
		AddressSubdivision: customer.CurrentAddress.StateProvinceCode,
		AddressPostalCode:  customer.CurrentAddress.PostCode,
		AddressCountryCode: customer.CurrentAddress.CountryAlpha2,
		PhysicalDocs:       []SubDocument{},
	}

	for _, commonDocument := range customer.Documents {
		if commonDocument.Front != nil && commonDocument.Front.Data != nil {
			document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
				DocumentType:  mapDocumentType(commonDocument.Metadata.Type),
				DocumentValue: base64.StdEncoding.EncodeToString(commonDocument.Front.Data),
			})
		}
	}

	return []Document{
		document,
	}
}

func mapCustomerGender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return "M"
	case common.Female:
		return "F"
	default:
		return "O"
	}
}

func mapDocumentType(documentType common.DocumentType) string {
	switch documentType {
	case common.IDCardEng, common.DriversEng, common.PassportEng:
		return "GOVT_ID"
	case common.IDCard, common.Drivers, common.Passport:
		return "GOVT_ID_INT"
	case common.Selfie:
		return "SELFIE"
	case common.UtilityBill, common.ResidencePermit:
		return "PROOF_OF_ADDRESS"
	case common.ProofOfIncome:
		return "PROOF_OF_INCOME"
	case common.ProofOfAccount:
		return "PROOF_OF_ACCOUNT"
	case common.ACHAuthorization:
		return "AUTHORIZATION"
	case common.BackgroundCheck:
		return "BK_CHECK"
	case common.SSN:
		return "SSN_CARD"
	case common.EINDocument:
		return "EIN_DOC"
	case common.W9Document:
		return "W9_DOC"
	case common.W8Document:
		return "W8_DOC"
	case common.W2Document:
		return "W2_DOC"
	case common.VoidedCheck:
		return "VOIDED_CHECK"
	case common.ArticlesOfIncorporation:
		return "AOI"
	case common.BylawsDocument:
		return "BYLAWS_DOC"
	case common.LetterOfEngagement:
		return "LOE"
	case common.CIPDoc:
		return "CIP_DOC"
	case common.SubscriptionAgreement:
		return "SUBSCRIPTION_AGREEMENT"
	case common.PromissoryNote:
		return "PROMISSORY_NOTE"
	case common.Contract, common.Agreement:
		return "LEGAL_AGREEMENT"
	case common.RegGG:
		return "REG_GG"
	case common.DBADoc:
		return "DBA_DOC"
	case common.DepositAgreement:
		return "DEPOSIT_AGGREEMENT"
	case common.Other:
		fallthrough
	default:
		return "OTHER"
	}
}
