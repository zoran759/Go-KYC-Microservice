package thomsonreuters

import (
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	"modulus/kyc/http"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// getRootGroups retrieves all the top-level groups with their immediate descendants.
func (tr ThomsonReuters) getRootGroups() (groups model.Groups, code *int, err error) {
	path := "groups"

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		err = fmt.Errorf("during fetching top level groups: %s", err)
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
		if err != nil || len(errs) == 0 {
			err = fmt.Errorf("during fetching top level groups: http error %d", status)
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &groups)

	return
}

// getGroup retrieves a specified group including its immediate descendants.
func (tr ThomsonReuters) getGroup(groupID string) (group model.Group, code *int, err error) {
	path := "groups/" + groupID

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		err = fmt.Errorf("during fetching the group with id %s: %s", groupID, err)
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
		if err != nil {
			err = fmt.Errorf("during fetching the group with id %s: http error %d", groupID, status)
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &group)

	return
}

// getCaseTemplate retrieves the CaseTemplate for the given Group.
func (tr ThomsonReuters) getCaseTemplate(groupID string) (caseTemplate model.CaseTemplateResponse, code *int, err error) {
	path := "groups/" + groupID + "/caseTemplate"

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		err = fmt.Errorf("during fetching a case template for the group with id %s: %s", groupID, err)
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
		if err != nil {
			err = fmt.Errorf("during fetching a case template for the group with id %s: http error %d", groupID, status)
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &caseTemplate)

	return
}

/*
This method will likely return a significant amount of data (more than 1 Mb), well, I guess we ain't gonna need it.

// getProviders retrieves a list of all available providers and their sources.
func (tr ThomsonReuters) getProviders() (providers model.ProviderDetails, code *int, err error) {
	path := "reference/providers"

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
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
*/

/*
Currently, we don't use this.

// getResolutionToolkits retrieves the ResolutionToolkits for the given Group for all enabled provider types,
// used to construct a valid resolution request(tr) on the results for a Case belonging to the given Group groupId.
func (tr ThomsonReuters) getResolutionToolkits(groupID string) (resToolkits model.ResolutionToolkits, code *int, err error) {
	path := "groups/" + groupID + "/resolutionToolkits"

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
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
*/

/*
Currently, we don't use this.

// getActiveUsers retrieves a list of active users (customers) in the Thomson Reuters API client’s account.
func (tr ThomsonReuters) getActiveUsers() (users model.Users, code *int, err error) {
	path := "users"

	headers := tr.createHeaders(mGET, path, nil)

	status, resp, err := http.Get(tr.scheme+"://"+tr.host+tr.path+path, headers)
	if err != nil {
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
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
*/

// performSynchronousScreening performs a synchronous screening for a given case.
// The returned result collection contains the regular case result details plus identity documents and important events.
func (tr ThomsonReuters) performSynchronousScreening(newcase model.NewCase) (rescol model.ScreeningResultCollection, code *int, err error) {
	path := "cases/screeningRequest"

	payload, err := json.Marshal(newcase)
	if err != nil {
		return
	}

	headers := tr.createHeaders(mPOST, path, payload)

	status, resp, err := http.Post(tr.scheme+"://"+tr.host+tr.path+path, headers, payload)
	if err != nil {
		err = fmt.Errorf("during performing synchronous screening: %s", err)
		return
	}

	if status != stdhttp.StatusOK {
		code = &status
		errs := model.Errors{}
		err = json.Unmarshal(resp, &errs)
		if err != nil {
			err = fmt.Errorf("during performing synchronous screening: http error %d", status)
			return
		}
		err = errs

		return
	}

	err = json.Unmarshal(resp, &rescol)

	return
}
