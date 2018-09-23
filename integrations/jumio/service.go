package jumio

import (
	"encoding/base64"

	"gitlab.com/lambospeed/kyc/common"
)

// service defines the model for the Jumio performNetverify API.
type service struct {
	host        string
	credentials string
}

// New constructs new service object to use with the Jumio performNetverify API.
func New(config Config) common.CustomerChecker {
	return &service{
		host:        config.Host,
		credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
	}
}

// CheckCustomer implements customer verification using the Jumio performNetverify API.
func (s *service) CheckCustomer(customer *common.UserData) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	// TODO: implement this.

	return
}
