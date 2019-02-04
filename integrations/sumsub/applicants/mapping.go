package applicants

import (
	"modulus/kyc/common"
)

func MapCommonCustomerToApplicant(customer common.UserData) ApplicantInfo {
	addresses := MapCommonAddressesToApplicantAddresses(customer)

	return ApplicantInfo{
		FirstName:      customer.FirstName,
		LastName:       customer.LastName,
		MiddleName:     customer.MiddleName,
		LegalName:      customer.LegalName,
		Gender:         MapGender(customer.Gender),
		DateOfBirth:    customer.DateOfBirth.Format("2006-01-02"),
		PlaceOfBirth:   customer.PlaceOfBirth,
		CountryOfBirth: common.CountryAlpha2ToAlpha3[customer.CountryOfBirthAlpha2],
		StateOfBirth:   customer.StateOfBirth,
		Country:        common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
		Nationality:    common.CountryAlpha2ToAlpha3[customer.Nationality],
		Phone:          customer.Phone,
		MobilePhone:    customer.MobilePhone,
		Addresses:      addresses,
	}
}

func MapCommonAddressesToApplicantAddresses(customer common.UserData) []Address {
	commonAddresses := customer.SupplementalAddresses
	currentAddress := MapCommonAddressToApplicantAddress(customer.CurrentAddress)
	if commonAddresses != nil && len(commonAddresses) > 0 {
		addresses := make([]Address, 0)
		for _, commonAddress := range commonAddresses {
			if address := MapCommonAddressToApplicantAddress(commonAddress); address != nil {
				addresses = append(addresses, *address)

			}
		}

		if currentAddress != nil {
			addresses = append(addresses, *currentAddress)
		}

		return addresses
	}

	if currentAddress != nil {
		return []Address{*currentAddress}
	}

	return nil
}

func MapCommonAddressToApplicantAddress(address common.Address) *Address {
	if address == (common.Address{}) {
		return nil
	}
	return &Address{
		Country:        common.CountryAlpha2ToAlpha3[address.CountryAlpha2],
		PostCode:       address.PostCode,
		Town:           address.Town,
		Street:         address.Street,
		SubStreet:      address.SubStreet,
		State:          address.State,
		BuildingName:   address.BuildingName,
		FlatNumber:     address.FlatNumber,
		BuildingNumber: address.BuildingNumber,
		StartDate:      address.StartDate.Format("2006-01-02"),
		EndDate:        address.EndDate.Format("2006-01-02"),
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
