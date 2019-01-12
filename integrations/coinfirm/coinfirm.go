package coinfirm

import (
	"modulus/kyc/common"
	"modulus/kyc/http"
)

// Coinfirm represents the Coinfirm API client.
type Coinfirm struct {
	host     string
	email    string
	password string
	company  string
	token    string
	headers  http.Headers
}

// New constructs a new Coinfirm API client instance.
// It accepts Config object as the config param.
func New(c Config) *Coinfirm {
	return &Coinfirm{
		host:     c.Host,
		email:    c.Email,
		password: c.Password,
		company:  c.Company,
		headers: http.Headers{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	}
}

// CheckCustomer implements CustomerChecker interface for the Coinfirm.
func (c *Coinfirm) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	// TODO: implement this.

	return
}

// CheckStatus implements StatusChecker interface for the Coinfirm.
func (c *Coinfirm) CheckStatus(pID string) (res common.KYCResult, err error) {
	// TODO: implement this.

	return
}
