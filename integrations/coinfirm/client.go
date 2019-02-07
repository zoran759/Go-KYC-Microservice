package coinfirm

import (
	"encoding/json"
	"errors"
	stdhttp "net/http"

	"modulus/kyc/http"
	"modulus/kyc/integrations/coinfirm/model"
)

// newAuthToken requests the API for a new user token required to access nearly all endpoints.
func (c Coinfirm) newAuthToken(headers http.Headers) (token model.AuthResponse, status *int, err error) {
	authreq := model.AuthRequest{
		Email:    c.config.Email,
		Password: c.config.Password,
	}

	body, err := json.Marshal(authreq)
	if err != nil {
		return
	}

	code, resp, err := http.Post(c.config.Host+"/auth/login", headers, body)
	if err != nil {
		return
	}

	if code != stdhttp.StatusOK {
		status = &code
		eresp := &model.ErrorResponse{}
		if err = json.Unmarshal(resp, eresp); err != nil {
			err = errors.New("http error")
			return
		}
		err = eresp
		return
	}

	err = json.Unmarshal(resp, &token)

	return
}

// newParticipant requests the API to add new participant without data.
func (c Coinfirm) newParticipant(headers http.Headers, nParticipant model.NewParticipant) (participant model.NewParticipantResponse, status *int, err error) {
	body, err := json.Marshal(nParticipant)
	if err != nil {
		return
	}

	code, resp, err := http.Request(stdhttp.MethodPut, c.config.Host+"/kyc/customers/"+c.config.Company, headers, body)
	if err != nil {
		return
	}

	if code != stdhttp.StatusOK {
		status = &code
		eresp := &model.ErrorResponse{}
		if err = json.Unmarshal(resp, eresp); err != nil {
			err = errors.New("http error")
			return
		}
		err = eresp
		return
	}

	err = json.Unmarshal(resp, &participant)

	return
}

// sendParticipantDetails sends individual participant data to the API.
func (c Coinfirm) sendParticipantDetails(headers http.Headers, pID string, details model.ParticipantDetails) (status *int, err error) {
	body, err := json.Marshal(details)
	if err != nil {
		return
	}

	code, resp, err := http.Request(stdhttp.MethodPut, c.config.Host+"/kyc/forms/"+c.config.Company+"/"+pID, headers, body)
	if err != nil {
		return
	}

	if code != stdhttp.StatusCreated {
		status = &code
		eresp := &model.ErrorResponse{}
		if err = json.Unmarshal(resp, eresp); err != nil {
			err = errors.New("http error")
			return
		}
		err = eresp
	}

	return
}

// sendDocFile sends a document file to the API to add it to KYC process.
func (c Coinfirm) sendDocFile(headers http.Headers, pID string, docfile *model.File) (status *int, err error) {
	body, err := json.Marshal(docfile)
	if err != nil {
		return
	}

	code, resp, err := http.Post(c.config.Host+"/kyc/files/"+c.config.Company+"/"+pID, headers, body)
	if err != nil {
		return
	}

	if code != stdhttp.StatusCreated {
		status = &code
		eresp := &model.ErrorResponse{}
		if err = json.Unmarshal(resp, eresp); err != nil {
			err = errors.New("http error")
			return
		}
		err = eresp
	}

	return
}

// getParticipantCurrentStatus requests the current participant status in KYC flow from the API.
func (c Coinfirm) getParticipantCurrentStatus(headers http.Headers, pID string) (status model.StatusResponse, code *int, err error) {
	rcode, resp, err := http.Get(c.config.Host+"/kyc/status/"+c.config.Company+"/"+pID, headers)
	if err != nil {
		return
	}

	if rcode != stdhttp.StatusOK {
		code = &rcode
		eresp := &model.ErrorResponse{}
		if err = json.Unmarshal(resp, eresp); err != nil {
			err = errors.New("http error")
			return
		}
		err = eresp
		return
	}

	err = json.Unmarshal(resp, &status)

	return
}
