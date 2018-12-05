package verification

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"modulus/kyc/common"
)

// MapCustomerToCreateUserRequest constructs and returns user creation request.
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

func mapCustomerDocuments(customer common.UserData) []Documents {
	document := Documents{
		OwnerName:          fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
		Email:              customer.Email,
		PhoneNumber:        customer.Phone,
		IPAddress:          customer.IPaddress,
		EntityType:         mapCustomerGender(customer.Gender),
		EntityScope:        "Not Known",
		DayOfBirth:         time.Time(customer.DateOfBirth).Day(),
		MonthOfBirth:       int(time.Time(customer.DateOfBirth).Month()),
		YearOfBirth:        time.Time(customer.DateOfBirth).Year(),
		AddressStreet:      customer.CurrentAddress.StreetAddress(),
		AddressCity:        customer.CurrentAddress.Town,
		AddressSubdivision: customer.CurrentAddress.StateProvinceCode,
		AddressPostalCode:  customer.CurrentAddress.PostCode,
		AddressCountryCode: customer.CurrentAddress.CountryAlpha2,
	}

	if (&common.IDCard{}) != customer.IDCard {
		document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("IDCard"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if (&common.Selfie{}) != customer.Selfie {
		document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("Selfie"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if customer.VideoAuth != nil {
		document.PhysicalDocs = append(document.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("VideoAuth"),
			DocumentValue: "data:" + customer.VideoAuth.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
		})
	}

	return []Documents{
		document,
	}
}

// MapUserToOauth returns OAuth token obtaining request.
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

// MapDocumentsToCreateUserRequest constructs and returns CreateDocumentsRequest object.
func MapDocumentsToCreateUserRequest(customer common.UserData) CreateDocumentsRequest {

	request := CreateDocumentsRequest{
		Documents: Documents{
			OwnerName:          fmt.Sprintf("%s %s %s", customer.FirstName, customer.MiddleName, customer.LastName),
			Email:              customer.Email,
			PhoneNumber:        customer.Phone,
			IPAddress:          customer.IPaddress,
			EntityType:         mapCustomerGender(customer.Gender),
			EntityScope:        "Not Known",
			DayOfBirth:         time.Time(customer.DateOfBirth).Day(),
			MonthOfBirth:       int(time.Time(customer.DateOfBirth).Month()),
			YearOfBirth:        time.Time(customer.DateOfBirth).Year(),
			AddressStreet:      customer.CurrentAddress.StreetAddress(),
			AddressCity:        customer.CurrentAddress.Town,
			AddressSubdivision: customer.CurrentAddress.StateProvinceCode,
			AddressPostalCode:  customer.CurrentAddress.PostCode,
			AddressCountryCode: customer.CountryAlpha2,
		},
	}

	if customer.IDCard != nil {
		request.Documents.PhysicalDocs = append(request.Documents.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("IDCard"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if customer.Selfie != nil {
		request.Documents.PhysicalDocs = append(request.Documents.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("Selfie"),
			DocumentValue: "data:image/png;base64," + base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
			//DocumentValue: "data:image/png;base64,SUQs==",
		})
	}

	if customer.VideoAuth != nil {
		request.Documents.PhysicalDocs = append(request.Documents.PhysicalDocs, SubDocument{
			DocumentType:  mapDocumentType("VideoAuth"),
			DocumentValue: "data:" + customer.VideoAuth.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
		})
	}

	return request
}

func mapDocumentType(documentType string) string {
	switch documentType {
	case "IDCard", "DriverLicense", "DriverLicenseTranslation", "Passport", "SNILS":
		return "GOVT_ID_INT"
	case "Selfie":
		return "SELFIE"
	case "VideoAuth":
		return "VIDEO_AUTHORIZATION"
	case "UtilityBill":
		return "PROOF_OF_ADDRESS"
	case "Contract", "Agreement":
		return "LEGAL_AGREEMENT"
	default:
		return "OTHER"
	}
}

// MapResponseError extracts and returns errors from the error response.
func MapResponseError(responseBytes []byte) (result, err error) {

	response := &ResponseError{}
	if err := json.Unmarshal(responseBytes, response); err != nil {
		log.Printf("Error decoding SynapseFi response: %v", err)
		return nil, err
	}
	log.Printf("SynapseFi response status: %v", response.Status)

	errMsg := response.Error[appLanguage]

	return errors.New(errMsg), nil
}
