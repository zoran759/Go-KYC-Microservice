package common

// CustomerChecker defines the interface for a KYC verificator.
type CustomerChecker interface {
	CheckCustomer(customer *UserData) (KYCResult, *DetailedKYCResult, error)
}
