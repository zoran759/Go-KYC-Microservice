package verification

import (
	"modulus/kyc/common"
	"fmt"
	"time"
	"errors"
	"log"
	"encoding/json"
	"encoding/base64"
)

func MapCustomerToCreateUserRequest(customer common.UserData, isSimpleMode bool) CreateUserRequest {

	request := CreateUserRequest{
		Logins: []Login{
			{
				Email: customer.Email,
				Scope: "READ_AND_WRITE",
			},
		},
		LegalNames: []string{
			fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
		},
		Extra: Extra{
			CIPTag:     1,
			IsBusiness: false,
		},
	}

	phoneNumbers := make([]string, 0)
	if customer.Phone != "" {
		phoneNumbers = append(phoneNumbers, customer.Phone)
	}
	if customer.MobilePhone != "" {
		phoneNumbers = append(phoneNumbers, customer.MobilePhone)
	}

	request.PhoneNumbers = phoneNumbers
	if isSimpleMode && customer.IDCard != nil && customer.Selfie != nil {
		request.Documents = mapCustomerDocuments(customer)
	}

	return request
}

func mapCustomerDocuments(customer common.UserData) []Document  {
	document := Document{
		PhysicalDocs:       []SubDocument{},
		OwnerName:    fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
		Email:        customer.Email,
		PhoneNumber:  customer.Phone,
		IPAddress:    "127.0.0.1",
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
		AddressCountryCode: customer.CountryAlpha2,
	}

	if (&common.IDCard{}) != customer.IDCard {
		document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
			DocumentType: mapDocumentType("IDCard"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if (&common.Selfie{}) != customer.Selfie {
		document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
			DocumentType: mapDocumentType("Selfie"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	return []Document{
		document,
	}

}

func MapUserToOauth(refreshToken string) CreateOauthRequest {
	return CreateOauthRequest{
		RefreshToken: refreshToken,
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

func MapDocumentsToCreateUserRequest(customer common.UserData) CreateDocumentsRequest {

	request := CreateDocumentsRequest{
		Documents: Document{
			OwnerName:    fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
			Email:        customer.Email,
			PhoneNumber:  customer.Phone,
			IPAddress:    "127.0.0.1",
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
			AddressCountryCode: customer.CountryAlpha2,
			PhysicalDocs:       []SubDocument{},
		},
	}

	if customer.IDCard != nil {
		request.Documents.PhysicalDocs = append(request.Documents.PhysicalDocs, SubDocument{
			DocumentType: mapDocumentType("IDCard"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if customer.Selfie != nil {
		request.Documents.PhysicalDocs = append(request.Documents.PhysicalDocs, SubDocument{
			DocumentType: mapDocumentType("Selfie"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	return request
}

func mapDocumentType(documentType string) string {
	switch documentType {
	case "IDCard", "DriverLicense", "DriverLicenseTranslation", "Passport", "SNILS", "DocumentPhoto":
		return "GOVT_ID_INT"
	case "Selfie":
		return "SELFIE"
	case "UtilityBill", "ResidencePermit":
		return "PROOF_OF_ADDRESS"
	case "Contract", "Agreement":
		return "LEGAL_AGREEMENT"
	default:
		return "OTHER"
	}
}

func MapResponseError(responseBytes []byte) (result error, err error) {

	response := &ResponseError{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		log.Printf("Error decoding SynapseFi response: %v", err)
		return nil, err
	}
	log.Printf("SynapseFi response status: %v", response.Status);

	errMsg := response.Error[AppLanguage]

	return errors.New(errMsg), nil
}

