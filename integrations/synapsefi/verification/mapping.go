package verification

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"modulus/kyc/common"
)

// MapCustomerToUser constructs and returns user object with virtual documents if any from customer data.
func MapCustomerToUser(customer *common.UserData) User {
	user := User{
		Logins: []Login{
			Login{
				Email: customer.Email,
				Scope: "READ_AND_WRITE",
			},
		},
		LegalNames: []string{customer.LegalName},
		Extra: Extra{
			CIPTag:     1,
			IsBusiness: false,
		},
	}

	if len(customer.Phone) > 0 {
		user.PhoneNumbers = append(user.PhoneNumbers, customer.Phone)
	}
	if len(customer.MobilePhone) > 0 {
		user.PhoneNumbers = append(user.PhoneNumbers, customer.MobilePhone)
	}

	user.Documents = mapCustomerToVirtualDocs(customer)

	return user
}

func mapCustomerToVirtualDocs(customer *common.UserData) []Document {
	doc := Document{
		OwnerName:          customer.LegalName,
		Email:              customer.Email,
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

	switch {
	case len(customer.Phone) > 0:
		doc.PhoneNumber = customer.Phone
	case len(customer.MobilePhone) > 0:
		doc.PhoneNumber = customer.MobilePhone
	}

	if customer.IDCard != nil {
		if customer.IDCard.CountryAlpha2 == "US" {
			doc.VirtualDocs = append(doc.VirtualDocs, SubDocument{
				Type:  "SSN",
				Value: customer.IDCard.Number,
			})
		} else {
			doc.VirtualDocs = append(doc.VirtualDocs, SubDocument{
				Type:  "PERSONAL_IDENTIFICATION",
				Value: customer.IDCard.Number,
			})
		}
	}

	if customer.Passport != nil {
		doc.VirtualDocs = append(doc.VirtualDocs, SubDocument{
			Type:  "PASSPORT",
			Value: customer.Passport.Number,
		})
	}

	if customer.DriverLicense != nil {
		doc.VirtualDocs = append(doc.VirtualDocs, SubDocument{
			Type:  "DRIVERS_LICENSE",
			Value: customer.DriverLicense.Number,
		})
	}

	return []Document{doc}
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

// MapCustomerToPhysicalDocs constructs and returns physical documents from customer data.
func MapCustomerToPhysicalDocs(customer *common.UserData) (docs []SubDocument) {
	if customer.Passport != nil && customer.Passport.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("Passport"),
			Value: "data:" + customer.Passport.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
		})
	}

	if customer.IDCard != nil && customer.IDCard.Image != nil {
		doctype := mapDocType("IDCard")
		if customer.IDCard.CountryAlpha2 == "US" {
			doctype = "SSN_CARD"
		}
		docs = append(docs, SubDocument{
			Type:  doctype,
			Value: "data:" + customer.IDCard.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
		})
	}

	if customer.SNILS != nil && customer.SNILS.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("SNILS"),
			Value: "data:" + customer.SNILS.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data),
		})
	}

	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("DriverLicense"),
			Value: "data:" + customer.DriverLicense.FrontImage.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data),
		})
	}

	if customer.DriverLicenseTranslation != nil && customer.DriverLicenseTranslation.FrontImage != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("DriverLicenseTranslation"),
			Value: "data:" + customer.DriverLicenseTranslation.FrontImage.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.DriverLicenseTranslation.FrontImage.Data),
		})
	}

	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("UtilityBill"),
			Value: "data:" + customer.UtilityBill.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data),
		})
	}

	if customer.Agreement != nil && customer.Agreement.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("Agreement"),
			Value: "data:" + customer.Agreement.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.Agreement.Image.Data),
		})
	}

	if customer.Contract != nil && customer.Contract.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("Contract"),
			Value: "data:" + customer.Contract.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.Contract.Image.Data),
		})
	}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("Selfie"),
			Value: "data:" + customer.Selfie.Image.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data),
		})
	}

	if customer.VideoAuth != nil {
		docs = append(docs, SubDocument{
			Type:  mapDocType("VideoAuth"),
			Value: "data:" + customer.VideoAuth.ContentType + ";base64," + base64.StdEncoding.EncodeToString(customer.VideoAuth.Data),
		})
	}

	return
}

func mapDocType(docType string) string {
	switch docType {
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

// MapErrorResponse extracts and returns errors from the error response.
func MapErrorResponse(response []byte) (code *string, err error) {
	eresp := &ErrorResponse{}

	if err = json.Unmarshal(response, eresp); err != nil {
		return
	}

	code = &eresp.Code
	err = eresp

	return
}
