package thomsonreuters

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	"modulus/kyc/http"
	"modulus/kyc/integrations/thomsonreuters/model"
)

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
