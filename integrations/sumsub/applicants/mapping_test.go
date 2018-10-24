package applicants

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"modulus/kyc/common"
	"time"
)

func TestMapCommonCustomerToApplicant(t *testing.T) {
	testTime := common.Time(time.Now())
	customer := common.UserData{
		FirstName:            "FirstName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		Email:                "Email",
		Gender:               common.Male,
		DateOfBirth:          testTime,
		PlaceOfBirth:         "PlaceOfBirth",
		CountryOfBirthAlpha2: "AU",
		StateOfBirth:         "StateOfBirth",
		CountryAlpha2:        "AU",
		Nationality:          "Nationality",
		Phone:                "Phone",
		MobilePhone:          "MobilePhone",
		CurrentAddress: common.Address{
			CountryAlpha2:  "AU",
			PostCode:       "PostCode1",
			Town:           "Town1",
			Street:         "Street1",
			SubStreet:      "SubStreet1",
			State:          "State1",
			BuildingName:   "BuildingName1",
			FlatNumber:     "FlatNumber1",
			BuildingNumber: "BuildingNumber1",
			StartDate:      testTime,
			EndDate:        testTime,
		},
		SupplementalAddresses: []common.Address{
			{
				CountryAlpha2:  "AU",
				PostCode:       "PostCode",
				Town:           "Town",
				Street:         "Street",
				SubStreet:      "SubStreet",
				State:          "State",
				BuildingName:   "BuildingName",
				FlatNumber:     "FlatNumber",
				BuildingNumber: "BuildingNumber",
				StartDate:      testTime,
				EndDate:        testTime,
			},
			{
				CountryAlpha2:  "AU",
				PostCode:       "PostCode1",
				Town:           "Town1",
				Street:         "Street1",
				SubStreet:      "SubStreet1",
				State:          "State1",
				BuildingName:   "BuildingName1",
				FlatNumber:     "FlatNumber1",
				BuildingNumber: "BuildingNumber1",
				StartDate:      testTime,
				EndDate:        testTime,
			},
		},
	}

	applicantInfo := MapCommonCustomerToApplicant(customer)

	assert.Equal(t, customer.FirstName, applicantInfo.FirstName)
	assert.Equal(t, customer.LastName, applicantInfo.LastName)
	assert.Equal(t, customer.MiddleName, applicantInfo.MiddleName)
	assert.Equal(t, customer.LegalName, applicantInfo.LegalName)
	assert.Equal(t, "Male", applicantInfo.Gender)
	assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), applicantInfo.DateOfBirth)
	assert.Equal(t, customer.PlaceOfBirth, applicantInfo.PlaceOfBirth)
	assert.Equal(t, "Australia", applicantInfo.CountryOfBirth)
	assert.Equal(t, customer.StateOfBirth, applicantInfo.StateOfBirth)
	assert.Equal(t, "AUS", applicantInfo.Country)
	assert.Equal(t, customer.Nationality, applicantInfo.Nationality)
	assert.Equal(t, customer.Phone, applicantInfo.Phone)
	assert.Equal(t, customer.MobilePhone, applicantInfo.MobilePhone)

	assert.Equal(t, 3, len(applicantInfo.Addresses))
	assert.Equal(t, "Australia", applicantInfo.Addresses[0].Country)
	assert.Equal(t, customer.SupplementalAddresses[0].PostCode, applicantInfo.Addresses[0].PostCode)
	assert.Equal(t, customer.SupplementalAddresses[0].Town, applicantInfo.Addresses[0].Town)
	assert.Equal(t, customer.SupplementalAddresses[0].SubStreet, applicantInfo.Addresses[0].SubStreet)
	assert.Equal(t, customer.SupplementalAddresses[0].State, applicantInfo.Addresses[0].State)
	assert.Equal(t, customer.SupplementalAddresses[0].BuildingName, applicantInfo.Addresses[0].BuildingName)
	assert.Equal(t, customer.SupplementalAddresses[0].FlatNumber, applicantInfo.Addresses[0].FlatNumber)
	assert.Equal(t, customer.SupplementalAddresses[0].BuildingName, applicantInfo.Addresses[0].BuildingName)
	assert.Equal(t, customer.SupplementalAddresses[0].StartDate.Format("2006-01-02"), applicantInfo.Addresses[0].StartDate)
	assert.Equal(t, customer.SupplementalAddresses[0].EndDate.Format("2006-01-02"), applicantInfo.Addresses[0].EndDate)

	assert.Equal(t, "Australia", applicantInfo.Addresses[1].Country)
	assert.Equal(t, customer.SupplementalAddresses[1].PostCode, applicantInfo.Addresses[1].PostCode)
	assert.Equal(t, customer.SupplementalAddresses[1].Town, applicantInfo.Addresses[1].Town)
	assert.Equal(t, customer.SupplementalAddresses[1].SubStreet, applicantInfo.Addresses[1].SubStreet)
	assert.Equal(t, customer.SupplementalAddresses[1].State, applicantInfo.Addresses[1].State)
	assert.Equal(t, customer.SupplementalAddresses[1].BuildingName, applicantInfo.Addresses[1].BuildingName)
	assert.Equal(t, customer.SupplementalAddresses[1].FlatNumber, applicantInfo.Addresses[1].FlatNumber)
	assert.Equal(t, customer.SupplementalAddresses[1].BuildingName, applicantInfo.Addresses[1].BuildingName)
	assert.Equal(t, customer.SupplementalAddresses[1].StartDate.Format("2006-01-02"), applicantInfo.Addresses[1].StartDate)
	assert.Equal(t, customer.SupplementalAddresses[1].EndDate.Format("2006-01-02"), applicantInfo.Addresses[1].EndDate)

	assert.Equal(t, "Australia", applicantInfo.Addresses[2].Country)
	assert.Equal(t, customer.CurrentAddress.PostCode, applicantInfo.Addresses[2].PostCode)
	assert.Equal(t, customer.CurrentAddress.Town, applicantInfo.Addresses[2].Town)
	assert.Equal(t, customer.CurrentAddress.SubStreet, applicantInfo.Addresses[2].SubStreet)
	assert.Equal(t, customer.CurrentAddress.State, applicantInfo.Addresses[2].State)
	assert.Equal(t, customer.CurrentAddress.BuildingName, applicantInfo.Addresses[2].BuildingName)
	assert.Equal(t, customer.CurrentAddress.FlatNumber, applicantInfo.Addresses[2].FlatNumber)
	assert.Equal(t, customer.CurrentAddress.BuildingName, applicantInfo.Addresses[2].BuildingName)
	assert.Equal(t, customer.CurrentAddress.StartDate.Format("2006-01-02"), applicantInfo.Addresses[2].StartDate)
	assert.Equal(t, customer.CurrentAddress.EndDate.Format("2006-01-02"), applicantInfo.Addresses[2].EndDate)

	customerWithNoAddresses := common.UserData{
		FirstName: "test",
		LastName:  "testtest",
	}

	applicantInfo = MapCommonCustomerToApplicant(customerWithNoAddresses)

	assert.Nil(t, applicantInfo.Addresses)

	customerWithOneAddress := common.UserData{
		Gender: common.Female,
		CurrentAddress: common.Address{
			CountryAlpha2:  "AU",
			PostCode:       "PostCode1",
			Town:           "Town1",
			Street:         "Street1",
			SubStreet:      "SubStreet1",
			State:          "State1",
			BuildingName:   "BuildingName1",
			FlatNumber:     "FlatNumber1",
			BuildingNumber: "BuildingNumber1",
			StartDate:      testTime,
			EndDate:        testTime,
		},
	}

	applicantInfo = MapCommonCustomerToApplicant(customerWithOneAddress)
	assert.Equal(t, "Female", applicantInfo.Gender)
	if assert.Equal(t, 1, len(applicantInfo.Addresses)) {
		assert.Equal(t, "Australia", applicantInfo.Addresses[0].Country)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.PostCode, applicantInfo.Addresses[0].PostCode)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.Town, applicantInfo.Addresses[0].Town)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.SubStreet, applicantInfo.Addresses[0].SubStreet)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.State, applicantInfo.Addresses[0].State)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.BuildingName, applicantInfo.Addresses[0].BuildingName)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.FlatNumber, applicantInfo.Addresses[0].FlatNumber)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.BuildingName, applicantInfo.Addresses[0].BuildingName)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.StartDate.Format("2006-01-02"), applicantInfo.Addresses[0].StartDate)
		assert.Equal(t, customerWithOneAddress.CurrentAddress.EndDate.Format("2006-01-02"), applicantInfo.Addresses[0].EndDate)
	}
}
