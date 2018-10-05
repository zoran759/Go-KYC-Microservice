package common

// CustomerChecker represents a KYC verificator.
type CustomerChecker interface {
	CheckCustomer(customer *UserData) (KYCResult, error)
}

// StatusChecker represents a KYC verification status checker.
// A CheckStatus checks the status of the KYC verification using
// KYCProvider param for KYC provider's name for which the check should be performed
// and second param as the submission id returned from a KYC provider API.
type StatusChecker interface {
	CheckStatus(KYCProvider, string) (KYCResult, error)
}
