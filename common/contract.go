package common

// KYCPlatform describes KYC provider platform.
//
// * CheckCustomer verifies the given UserData using a specified KYC provider's API.
// * CheckStatus checks the status of the existing KYC verification.
type KYCPlatform interface {
	CheckCustomer(customer *UserData) (KYCResult, error)
	CheckStatus(referenceID string) (KYCResult, error)
}
