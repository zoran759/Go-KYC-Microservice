package common

// CustomerChecker represents a KYC verificator.
type CustomerChecker interface {
	CheckCustomer(customer *UserData) (KYCResult, error)
}

// StatusChecker represents a KYC verification status checker.
// A CheckStatus checks the status of the KYC verification using
// input param as the submission id returned from a KYC provider API.
type StatusChecker interface {
	CheckStatus(string) (KYCResult, error)
}
