package common

// CustomerChecker represents a KYC verificator.
// CheckCustomer verifies the given UserData using a specified KYC provider's API.
// It returns verification result as the KYCResult and an error if occurred.
type CustomerChecker interface {
	CheckCustomer(*UserData) (KYCResult, error)
}

// StatusChecker represents a KYC verification status checker.
// CheckStatus checks the status of the KYC verification using the data provided in KYCStatusCheck.
// It returns verification result as the KYCResult and an error if occurred.
type StatusChecker interface {
	CheckStatus(string) (KYCResult, error)
}
