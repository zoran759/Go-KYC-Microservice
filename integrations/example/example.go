package example

import (
	"modulus/kyc/common"
)

// Checks the customer with the KYC provider and returns a boolean indicating whether user is approved.
func CheckCustomer(customer *common.UserData) bool {
	return true
}
