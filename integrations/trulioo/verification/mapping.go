package verification

import (
	"encoding/base64"
	"modulus/kyc/common"
	"time"
)

// MapCustomerToDataFields converts input customer data to the format acceptable by the API.
func MapCustomerToDataFields(customer *common.UserData) DataFields {
	return DataFields{
		PersonInfo:    mapCustomerToPersonalInfo(customer),
		Location:      mapCustomerAddressToLocation(customer.CurrentAddress),
		Communication: mapCustomerToCommunication(customer),
		Document:      mapCustomerDocument(customer),
		Business:      mapCustomerBusiness(customer.Business),
	}
}

func mapCustomerToPersonalInfo(customer *common.UserData) *PersonInfo {
	dateOfBirth := time.Time(customer.DateOfBirth)

	return &PersonInfo{
		FirstGivenName: customer.FirstName,
		MiddleName:     customer.MiddleName,
		FirstSurName:   customer.PaternalLastName,
		SecondSurname:  customer.LastName,
		ISOLatin1Name:  customer.LatinISO1Name,
		DayOfBirth:     dateOfBirth.Day(),
		MonthOfBirth:   int(dateOfBirth.Month()),
		YearOfBirth:    dateOfBirth.Year(),
		Gender:         mapGender(customer.Gender),
	}
}

func mapCustomerAddressToLocation(address common.Address) *Location {
	if address == (common.Address{}) {
		return nil
	}

	return &Location{
		BuildingNumber:    address.BuildingNumber,
		BuildingName:      address.BuildingName,
		UnitNumber:        address.FlatNumber,
		StreetName:        address.Street,
		StreetType:        address.StreetType,
		City:              address.Town,
		Suburb:            address.Suburb,
		County:            address.County,
		Country:           address.CountryAlpha2,
		StateProvinceCode: address.StateProvinceCode,
		PostalCode:        address.PostCode,
		POBox:             address.PostOfficeBox,
	}
}

func mapCustomerToCommunication(customer *common.UserData) *Communication {
	if customer.MobilePhone == "" && customer.Email == "" && customer.Phone == "" {
		return nil
	}
	return &Communication{
		MobileNumber: customer.MobilePhone,
		EmailAddress: customer.Email,
		Telephone:    customer.Phone,
	}
}

func mapCustomerDocument(customer *common.UserData) (document *Document) {
	document = &Document{}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		document.LivePhoto = base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data)
	}

	if customer.Passport != nil && customer.Passport.Image != nil {
		document.DocumentType = "Passport"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.Passport.Image.Data)
		return
	}

	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		document.DocumentType = "DrivingLicence"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data)
		if customer.DriverLicense.BackImage != nil {
			document.DocumentBackImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.BackImage.Data)
		}
		return
	}

	if customer.IDCard != nil && customer.IDCard.Image != nil {
		document.DocumentType = "IdentityCard"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data)
		return
	}

	if customer.ResidencePermit != nil && customer.ResidencePermit.Image != nil {
		document.DocumentType = "ResidencePermit"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.ResidencePermit.Image.Data)
		return
	}

	return nil
}

func mapCustomerBusiness(business *common.Business) *Business {
	if business == nil {
		return nil
	}

	incorporationDate := time.Time(business.IncorporationDate)

	return &Business{
		BusinessName:                business.Name,
		BusinessRegistrationNumber:  business.RegistrationNumber,
		DayOfIncorporation:          incorporationDate.Day(),
		MonthOfIncorporation:        int(incorporationDate.Month()),
		YearOfIncorporation:         incorporationDate.Year(),
		JurisdictionOfIncorporation: business.IncorporationJurisdiction,
	}
}

func mapGender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return "M"
	case common.Female:
		return "F"
	default:
		return ""
	}
}
