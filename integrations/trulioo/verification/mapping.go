package verification

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/numbers"
	"gitlab.com/lambospeed/kyc/strings"
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
		FirstGivenName: strings.Pointerize(customer.FirstName),
		MiddleName:     strings.Pointerize(customer.MiddleName),
		FirstSurName:   strings.Pointerize(customer.PaternalLastName),
		SecondSurname:  strings.Pointerize(customer.LastName),
		ISOLatin1Name:  strings.Pointerize(customer.LatinISO1Name),
		DayOfBirth:     numbers.PointerizeInt(dateOfBirth.Day()),
		MonthOfBirth:   numbers.PointerizeInt(int(dateOfBirth.Month())),
		YearOfBirth:    numbers.PointerizeInt(dateOfBirth.Year()),
		Gender:         strings.Pointerize(mapGender(customer.Gender)),
	}
}

func mapCustomerAddressToLocation(address common.Address) *Location {
	if address == (common.Address{}) {
		return nil
	}

	return &Location{
		BuildingNumber:    strings.Pointerize(address.BuildingNumber),
		BuildingName:      strings.Pointerize(address.BuildingName),
		UnitNumber:        strings.Pointerize(address.FlatNumber),
		StreetName:        strings.Pointerize(address.Street),
		StreetType:        strings.Pointerize(address.StreetType),
		City:              strings.Pointerize(address.Town),
		Suburb:            strings.Pointerize(address.Suburb),
		County:            strings.Pointerize(address.County),
		Country:           strings.Pointerize(address.CountryAlpha2),
		StateProvinceCode: strings.Pointerize(address.StateProvinceCode),
		PostalCode:        strings.Pointerize(address.PostCode),
		POBox:             strings.Pointerize(address.PostOfficeBox),
	}
}

func mapCustomerToCommunication(customer common.UserData) *Communication {
	if customer.MobilePhone == "" && customer.Email == "" && customer.Phone == "" {
		return nil
	}
	return &Communication{
		MobileNumber: strings.Pointerize(customer.MobilePhone),
		EmailAddress: strings.Pointerize(customer.Email),
		Telephone:    strings.Pointerize(customer.Phone),
	}
}

func mapCustomerDocument(customer common.UserData) *Document {
	if customer.Documents != nil && len(customer.Documents) > 0 {
		commonDocument := customer.Documents[0]

		document := Document{
			DocumentType: commonDocument.Metadata.Type,
		}

		if commonDocument.Front != nil {
			document.DocumentFrontImage = commonDocument.Front.Data
		}
		if commonDocument.Back != nil {
			document.DocumentBackImage = commonDocument.Back.Data
		}

		return &document
	} else {
		return nil
	}
}

func mapCustomerBusiness(business common.Business) *Business {
	if business == (common.Business{}) {
		return nil
	}

	incorporationDate := time.Time(business.IncorporationDate)

	return &Business{
		BusinessName:                strings.Pointerize(business.Name),
		BusinessRegistrationNumber:  strings.Pointerize(business.RegistrationNumber),
		DayOfIncorporation:          numbers.PointerizeInt(incorporationDate.Day()),
		MonthOfIncorporation:        numbers.PointerizeInt(int(incorporationDate.Month())),
		YearOfIncorporation:         numbers.PointerizeInt(incorporationDate.Year()),
		JurisdictionOfIncorporation: strings.Pointerize(business.IncorporationJurisdiction),
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
