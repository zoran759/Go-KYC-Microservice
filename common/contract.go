package common

// KYCPlatform describes a KYC/AML platform that can be used for KYC/AML evaluation.
//
// * CheckCustomer verifies the given UserData using a specified platform's API.
//   It returns verification result as the KYCResult and an error if occurred.
//
// * CheckStatus checks the status of an existing verification using the reference provided.
//   It returns verification result as the KYCResult and an error if occurred.
type KYCPlatform interface {
	CheckCustomer(customer *UserData) (KYCResult, error)
	CheckStatus(referenceID string) (KYCResult, error)
}

// CustomerChecker represents a KYC verificator.
// CheckCustomer verifies the given UserData using a specified KYC provider's API.
// It returns verification result as the KYCResult and an error if occurred.
type CustomerChecker interface {
	CheckCustomer(customer *UserData) (KYCResult, error)
}

// StatusChecker represents a KYC verification status checker.
// CheckStatus checks the status of the KYC verification using the data provided in KYCStatusCheck.
// It returns verification result as the KYCResult and an error if occurred.
type StatusChecker interface {
	CheckStatus(referenceID string) (KYCResult, error)
}

// CredentialsChecker describes a validator of KYC provider configuration credentials.
// It tries to use the credentials in a real request to the provider API to see
// whether the credentials are valid and returns a bool result or an error if occured.
type CredentialsChecker interface {
	CheckCredentials() (bool, error)
}
