package applicants

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/strings"
)

func MapCommonCustomerToApplicant(customer common.UserData) ApplicantInfo {
	addresses := MapCommonAddressesToApplicantAddresses(customer)

	return ApplicantInfo{
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		MiddleName:     strings.Pointerize(customer.MiddleName),
		LegalName:      strings.Pointerize(customer.LegalName),
		Gender:         strings.Pointerize(MapGender(customer.Gender)),
		DateOfBirth:    strings.Pointerize(customer.DateOfBirth.Format("2006-01-02")),
		PlaceOfBirth:   strings.Pointerize(customer.PlaceOfBirth),
		CountryOfBirth: strings.Pointerize(customer.CountryOfBirthAlpha2),
		StateOfBirth:   strings.Pointerize(customer.StateOfBirth),
		Country:        strings.Pointerize(customer.CountryAlpha2),
		Nationality:    strings.Pointerize(customer.Nationality),
		Phone:          strings.Pointerize(customer.Phone),
		MobilePhone:    strings.Pointerize(customer.MobilePhone),
		Addresses:      addresses,
	}
}

func MapCommonAddressesToApplicantAddresses(customer common.UserData) []Address {
	commonAddresses := customer.SupplementalAddresses
	currentAddress := MapCommonAddressToApplicantAddress(customer.CurrentAddress)
	if commonAddresses != nil && len(commonAddresses) > 0 {
		addresses := make([]Address, 0)
		for _, commmonAddress := range commonAddresses {
			if address := MapCommonAddressToApplicantAddress(commmonAddress); address != nil {
				addresses = append(addresses, *address)

			}
		}

		if currentAddress != nil {
			addresses = append(addresses, *currentAddress)
		}

		return addresses
	} else {
		if currentAddress != nil {
			return []Address{
				*currentAddress,
			}
		}
		return nil
	}
}

func MapCommonAddressToApplicantAddress(address common.Address) *Address {
	if address == (common.Address{}) {
		return nil
	}
	return &Address{
		Country:        strings.Pointerize(address.CountryAlpha2),
		PostCode:       strings.Pointerize(address.PostCode),
		Town:           strings.Pointerize(address.Town),
		Street:         strings.Pointerize(address.Street),
		SubStreet:      strings.Pointerize(address.SubStreet),
		State:          strings.Pointerize(address.State),
		BuildingName:   strings.Pointerize(address.BuildingName),
		FlatNumber:     strings.Pointerize(address.FlatNumber),
		BuildingNumber: strings.Pointerize(address.BuildingNumber),
		StartDate:      strings.Pointerize(address.StartDate.Format("2006-01-02")),
		EndDate:        strings.Pointerize(address.EndDate.Format("2006-01-02")),
	}
}

func MapGender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return "Male"
	case common.Female:
		return "Female"
	default:
		return ""
	}
}
