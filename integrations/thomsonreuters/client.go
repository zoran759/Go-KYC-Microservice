package thomsonreuters

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	"modulus/kyc/http"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// getRootGroups retrieves all the top-level groups with their immediate descendants.
func (tomson ThomsonReuters) getRootGroups() (groups model.Groups, code *int, err error) {
	path := "groups"

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
func (tomson ThomsonReuters) getGroup(groupID string) (group model.Group, code *int, err error) {
	path := "groups/" + groupID

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
func (tomson ThomsonReuters) getCaseTemplate(groupID string) (caseTemplate model.CaseTemplateResponse, code *int, err error) {
	path := "groups/" + groupID + "/caseTemplate"

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
func (tomson ThomsonReuters) getProviders() (providers model.ProviderDetails, code *int, err error) {
	path := "reference/providers"

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
// used to construct a valid resolution request(tomson) on the results for a Case belonging to the given Group groupId.
func (tomson ThomsonReuters) getResolutionToolkits(groupID string) (resToolkits model.ResolutionToolkits, code *int, err error) {
	path := "groups/" + groupID + "/resolutionToolkits"

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
func (tomson ThomsonReuters) getActiveUsers() (users model.Users, code *int, err error) {
	path := "users"

	headers := tomson.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tomson.scheme+"://"+tomson.host+tomson.path+path, headers)
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
func (tomson ThomsonReuters) performSynchronousScreening(newcase model.NewCase) (rescol model.ScreeningResultCollection, code *int, err error) {
	path := "cases/screeningRequest"

	payload, err := json.Marshal(newcase)
	if err != nil {
		return
	}

	headers := tomson.createHeaders(mPOST, path, payload)

	status, resp, err := http.Post(tomson.scheme+"://"+tomson.host+tomson.path+path, headers, payload)
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
