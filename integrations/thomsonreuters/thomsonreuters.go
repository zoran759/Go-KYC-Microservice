package thomsonreuters

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"modulus/kyc/common"
)

var _ common.KYCPlatform = ThomsonReuters{}

// ThomsonReuters represents the Thomson Reuters API client.
type ThomsonReuters struct {
	scheme string
	host   string
	path   string
	key    string
	secret string
}

// New constructs a new ThomsonReuters client.
func New(c Config) ThomsonReuters {
	u, err := url.Parse(c.Host)
	if err != nil {
		log.Println("During constructing new Thomson Reuters client:", err)
		return ThomsonReuters{}
	}
	if len(u.Scheme) == 0 || len(u.Host) == 0 {
		log.Println("During constructing new Thomson Reuters client: malformed Host format")
		return ThomsonReuters{}
	}

	if !strings.HasSuffix(u.Path, "/") {
		u.Path = u.Path + "/"
	}

	return ThomsonReuters{
		scheme: u.Scheme,
		host:   u.Host,
		path:   u.Path,
		key:    c.APIkey,
		secret: c.APIsecret,
	}
}

// CheckCustomer implements KYCPlatform interface for Thomson Reuters.
func (tr ThomsonReuters) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("customer data is nil")
		return
	}

	gID, code, err := tr.getGroupID()
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	template, code, err := tr.getCaseTemplate(gID)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	newcase := newCase(template, customer)

	src, code, err := tr.performSynchronousScreening(newcase)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	result, err = toResult(src)

	return
}

// CheckStatus implements KYCPlatform interface for Thomson Reuters.
func (tr ThomsonReuters) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	err = errors.New("Thomson Reuters doesn't support a verification status check")
	return
}
