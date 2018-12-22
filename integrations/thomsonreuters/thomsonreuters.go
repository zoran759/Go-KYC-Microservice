package thomsonreuters

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"
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
func New(c Config) common.CustomerChecker {
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
	gID, code, err := s.getGroupID()
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	template, code, err := s.getCaseTemplate(gID)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	toolkits, code, err := s.getResolutionToolkits(gID)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	_ = template
	_ = toolkits

	return
}

// getGroupID returns group id.
func (s service) getGroupID() (groupID string, code *int, err error) {
	groups, code, err := s.getRootGroups()
	if err != nil {
		return
	}

	// Obtain id of the first active root group.
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		groupID = g.ID
		break
	}

	if len(groupID) == 0 {
		err = errors.New("the verification prerequisites error: no active root group")
	}

	return
}
