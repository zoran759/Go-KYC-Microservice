package thomsonreuters

import (
	"encoding/json"
	"errors"
	"log"
	stdhttp "net/http"
	"net/url"
	"strings"

	"modulus/kyc/common"
	"modulus/kyc/http"
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

// getRootGroups retrieves all the top-level groups with their immediate descendants.
func (s service) getRootGroups() (groups model.Groups, code *int, err error) {
	path := "groups"

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &groups)

	return
}

// getGroup retrieves a specified group including its immediate descendants.
func (s service) getGroup(groupID string) (group model.Group, code *int, err error) {
	path := "groups/" + groupID

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &group)

	return

}

// getCaseTemplate retrieves the CaseTemplate for the given Group.
func (s service) getCaseTemplate(groupID string) (caseTemplate model.CaseTemplateResponse, code *int, err error) {
	path := "groups/" + groupID + "/caseTemplate"

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &caseTemplate)

	return
}

// getProviders retrieves a list of all available providers and their sources.
func (s service) getProviders() (providers model.ProviderDetails, code *int, err error) {
	path := "reference/providers"

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &providers)

	return
}

// getResolutionToolkits retrieves the ResolutionToolkits for the given Group for all enabled provider types,
// used to construct a valid resolution request(s) on the results for a Case belonging to the given Group groupId
func (s service) getResolutionToolkits(groupID string) (resToolkits model.ResolutionToolkits, code *int, err error) {
	path := "groups/" + groupID + "/resolutionToolkits"

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &resToolkits)

	return
}

// getActiveUsers retrieves a list of active users (customers) in the Thomson Reuters API client’s account.
func (s service) getActiveUsers() (users model.Users, code *int, err error) {
	path := "users"

	headers := s.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(s.scheme+s.host+s.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &users)

	return
}

// performSynchronousScreening performs a synchronous screening for a given case.
// The returned result collection contains the regular case result details plus identity documents and important events.
func (s service) performSynchronousScreening(newcase model.NewCase) (rescol model.ScreeningResultCollection, code *int, err error) {
	path := "cases/screeningRequest"

	payload, err := json.Marshal(newcase)
	if err != nil {
		return
	}

	headers := s.createHeaders(mPOST, path, payload)

	status, resp, err := http.Post(s.scheme+s.host+s.path+path, headers, payload)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, errs)
		if err != nil {
			err = errors.New("http error")
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &rescol)

	return
}
