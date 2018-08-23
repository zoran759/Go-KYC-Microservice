package idology

import (
	"gitlab.com/modulusglobal/kyc/common"
)

// Config holds configuration settings for the service which interacts with IDology API.
type Config struct {
	Host,
	Username,
	Password string
}

// CustomerChecker defines the customer verification interface for IDology.
type CustomerChecker interface {
	CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error)
}
