package applicants

import (
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestMapCommonCustomerToApplicant(t *testing.T) {
	assert := assert.New(t)

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
		Nationality:          "US",
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

	assert.Equal(customer.FirstName, applicantInfo.FirstName)
	assert.Equal(customer.LastName, applicantInfo.LastName)
	assert.Equal(customer.MiddleName, applicantInfo.MiddleName)
	assert.Equal(customer.LegalName, applicantInfo.LegalName)
	assert.Equal("Male", applicantInfo.Gender)
	assert.Equal(customer.DateOfBirth.Format("2006-01-02"), applicantInfo.DateOfBirth)
	assert.Equal(customer.PlaceOfBirth, applicantInfo.PlaceOfBirth)
	assert.Equal("AUS", applicantInfo.CountryOfBirth)
	assert.Equal(customer.StateOfBirth, applicantInfo.StateOfBirth)
	assert.Equal("AUS", applicantInfo.Country)
	assert.Equal("USA", applicantInfo.Nationality)
	assert.Equal(customer.Phone, applicantInfo.Phone)
	assert.Equal(customer.MobilePhone, applicantInfo.MobilePhone)

	assert.Equal(3, len(applicantInfo.Addresses))
	assert.Equal("AUS", applicantInfo.Addresses[0].Country)
	assert.Equal(customer.SupplementalAddresses[0].PostCode, applicantInfo.Addresses[0].PostCode)
	assert.Equal(customer.SupplementalAddresses[0].Town, applicantInfo.Addresses[0].Town)
	assert.Equal(customer.SupplementalAddresses[0].SubStreet, applicantInfo.Addresses[0].SubStreet)
	assert.Equal(customer.SupplementalAddresses[0].State, applicantInfo.Addresses[0].State)
	assert.Equal(customer.SupplementalAddresses[0].BuildingName, applicantInfo.Addresses[0].BuildingName)
	assert.Equal(customer.SupplementalAddresses[0].FlatNumber, applicantInfo.Addresses[0].FlatNumber)
	assert.Equal(customer.SupplementalAddresses[0].BuildingName, applicantInfo.Addresses[0].BuildingName)
	assert.Equal(customer.SupplementalAddresses[0].StartDate.Format("2006-01-02"), applicantInfo.Addresses[0].StartDate)
	assert.Equal(customer.SupplementalAddresses[0].EndDate.Format("2006-01-02"), applicantInfo.Addresses[0].EndDate)

	assert.Equal("AUS", applicantInfo.Addresses[1].Country)
	assert.Equal(customer.SupplementalAddresses[1].PostCode, applicantInfo.Addresses[1].PostCode)
	assert.Equal(customer.SupplementalAddresses[1].Town, applicantInfo.Addresses[1].Town)
	assert.Equal(customer.SupplementalAddresses[1].SubStreet, applicantInfo.Addresses[1].SubStreet)
	assert.Equal(customer.SupplementalAddresses[1].State, applicantInfo.Addresses[1].State)
	assert.Equal(customer.SupplementalAddresses[1].BuildingName, applicantInfo.Addresses[1].BuildingName)
	assert.Equal(customer.SupplementalAddresses[1].FlatNumber, applicantInfo.Addresses[1].FlatNumber)
	assert.Equal(customer.SupplementalAddresses[1].BuildingName, applicantInfo.Addresses[1].BuildingName)
	assert.Equal(customer.SupplementalAddresses[1].StartDate.Format("2006-01-02"), applicantInfo.Addresses[1].StartDate)
	assert.Equal(customer.SupplementalAddresses[1].EndDate.Format("2006-01-02"), applicantInfo.Addresses[1].EndDate)

	assert.Equal("AUS", applicantInfo.Addresses[2].Country)
	assert.Equal(customer.CurrentAddress.PostCode, applicantInfo.Addresses[2].PostCode)
	assert.Equal(customer.CurrentAddress.Town, applicantInfo.Addresses[2].Town)
	assert.Equal(customer.CurrentAddress.SubStreet, applicantInfo.Addresses[2].SubStreet)
	assert.Equal(customer.CurrentAddress.State, applicantInfo.Addresses[2].State)
	assert.Equal(customer.CurrentAddress.BuildingName, applicantInfo.Addresses[2].BuildingName)
	assert.Equal(customer.CurrentAddress.FlatNumber, applicantInfo.Addresses[2].FlatNumber)
	assert.Equal(customer.CurrentAddress.BuildingName, applicantInfo.Addresses[2].BuildingName)
	assert.Equal(customer.CurrentAddress.StartDate.Format("2006-01-02"), applicantInfo.Addresses[2].StartDate)
	assert.Equal(customer.CurrentAddress.EndDate.Format("2006-01-02"), applicantInfo.Addresses[2].EndDate)

	customerWithNoAddresses := common.UserData{
		FirstName: "test",
		LastName:  "testtest",
	}

	applicantInfo = MapCommonCustomerToApplicant(customerWithNoAddresses)

	assert.Nil(applicantInfo.Addresses)

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
	assert.Equal("Female", applicantInfo.Gender)
	if assert.Equal(1, len(applicantInfo.Addresses)) {
		assert.Equal("AUS", applicantInfo.Addresses[0].Country)
		assert.Equal(customerWithOneAddress.CurrentAddress.PostCode, applicantInfo.Addresses[0].PostCode)
		assert.Equal(customerWithOneAddress.CurrentAddress.Town, applicantInfo.Addresses[0].Town)
		assert.Equal(customerWithOneAddress.CurrentAddress.SubStreet, applicantInfo.Addresses[0].SubStreet)
		assert.Equal(customerWithOneAddress.CurrentAddress.State, applicantInfo.Addresses[0].State)
		assert.Equal(customerWithOneAddress.CurrentAddress.BuildingName, applicantInfo.Addresses[0].BuildingName)
		assert.Equal(customerWithOneAddress.CurrentAddress.FlatNumber, applicantInfo.Addresses[0].FlatNumber)
		assert.Equal(customerWithOneAddress.CurrentAddress.BuildingName, applicantInfo.Addresses[0].BuildingName)
		assert.Equal(customerWithOneAddress.CurrentAddress.StartDate.Format("2006-01-02"), applicantInfo.Addresses[0].StartDate)
		assert.Equal(customerWithOneAddress.CurrentAddress.EndDate.Format("2006-01-02"), applicantInfo.Addresses[0].EndDate)
	}
}
