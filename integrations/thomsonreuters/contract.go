package thomsonreuters

import "modulus/kyc/common"

// Config represents the service config.
type Config struct {
	Host      string
	APIkey    string
	APIsecret string
}

// ThomsonReuters describes verification service interface for Thomson Reuters.
type ThomsonReuters interface {
	CheckCustomer(customer *common.UserData) (common.KYCResult, error)
	// Prevents inappropriate interface usage.
	id() string
}
