package complyadvantage

import "modulus/kyc/common"

// service represents the service.
type service struct {
	host string
	key  string
}

// New returns a new verification service object.
func New(c Config) common.CustomerChecker {
	return service{
		host: c.Host,
		key:  c.APIkey,
	}
}

// CheckCustomer implements CustomerChecker interface for the ComplyAdvantage.
func (s service) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	// TODO: implememnt this.

	return
}
