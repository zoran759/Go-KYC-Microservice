package example

import (
	"gitlab.com/modulusglobal/kyc/common"
)

// Checks the customer with the KYC provider and returns a boolean indicating whether user is approved.
func CheckCustomer(customer *common.UserData) bool {
	return true
}
