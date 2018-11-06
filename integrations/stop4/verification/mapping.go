package verification

import (
	"modulus/kyc/common"
	"time"
)

/*

 */
// MapCustomerToVerificationRequest maps the values of common data to the service specific values.
func MapCustomerToVerificationRequest(customer common.UserData) Request {
	request := Request{
		RegIPAddress: "127.0.0.1",
		RegDate:      time.Now().Format("2006-01-02"),

		CustomerInformation: CustomerInformation{
			FirstName:  customer.FirstName,
			MiddleName: customer.MiddleName,
			LastName:   customer.LastName,
			Email:      customer.Email,
			Address1:   customer.CurrentAddress.String(),
			Phone1:     customer.Phone,
			Phone2:     customer.MobilePhone,
			Dob:        customer.DateOfBirth.Format("2006-01-02"),
		},
	}

	return request
}
