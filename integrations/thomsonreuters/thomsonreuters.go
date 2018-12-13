package thomsonreuters

import (
	"log"
	"net/url"
	"strings"

	"modulus/kyc/common"
)

// service represents the service.
type service struct {
	scheme string
	host   string
	path   string
	key    string
	secret string
}

// New constructs a new service object.
func New(c Config) ThomsonReuters {
	u, err := url.Parse(c.Host)
	if err != nil {
		log.Println("During constructing new Thomson Reuters service:", err)
		return service{}
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	return service{
		scheme: u.Scheme,
		host:   u.Host,
		path:   u.Path,
		key:    c.APIkey,
		secret: c.APIsecret,
	}
}

// CheckCustomer implements CustomerChecker interface for Thomson Reuters.
func (s service) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	// TODO: implement this.

	return
}

func (s service) id() string {
	return "Thomson Reuters"
}
