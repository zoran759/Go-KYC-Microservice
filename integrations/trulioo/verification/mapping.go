package verification

import (
	"modulus/kyc/common"
	"time"
)

func MapCustomerToDataFields(customer common.UserData) DataFields {
	return DataFields{
		PersonInfo:    mapCustomerToPersonalInfo(customer),
		Location:      mapCustomerAddressToLocation(customer.CurrentAddress),
		Communication: mapCustomerToCommunication(customer),
		Document:      mapCustomerDocument(customer),
		Business:      mapCustomerBusiness(customer.Business),
	}
}

func mapCustomerToPersonalInfo(customer common.UserData) *PersonInfo {
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

func mapCustomerToCommunication(customer common.UserData) *Communication {
	if customer.MobilePhone == "" && customer.Email == "" && customer.Phone == "" {
		return nil
	}
	return &Communication{
		MobileNumber: customer.MobilePhone,
		EmailAddress: customer.Email,
		Telephone:    customer.Phone,
	}
}

func mapCustomerDocument(customer common.UserData) *Document {
	if customer.Documents != nil && len(customer.Documents) > 0 {
		document := Document{}

		for _, commonDocument := range customer.Documents {
			if document.DocumentType != "" && document.LivePhoto != nil {
				break
			} else if document.LivePhoto == nil && commonDocument.Metadata.Type == common.Selfie {
				document.LivePhoto = commonDocument.Front.Data
			} else if document.DocumentType == "" {
				if mappedType := mapCustomerDocumentType(commonDocument.Metadata.Type); mappedType != "" {
					document.DocumentType = mappedType

					if commonDocument.Front != nil {
						document.DocumentFrontImage = commonDocument.Front.Data
					}
					if commonDocument.Back != nil {
						document.DocumentBackImage = commonDocument.Back.Data
					}
				}
			}
		}

		if document.DocumentType == "" {
			return nil
		}

		return &document
	} else {
		return nil
	}
}

func mapCustomerDocumentType(documentType common.DocumentType) string {
	switch documentType {
	case common.Passport:
		return "Passport"
	case common.Drivers:
		return "DrivingLicence"
	case common.IDCard:
		return "IdentityCard"
	case common.ResidencePermit:
		return "ResidencePermit"
	default:
		return ""
	}
}

func mapCustomerBusiness(business common.Business) *Business {
	if business == (common.Business{}) {
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
