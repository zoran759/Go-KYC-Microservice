package handlers

import (
	"encoding/json"
	"errors"
	"modulus/kyc/common"
	"net/http"
)

type isProviderImplementedResp struct {
	Implemented bool
}

// IsProviderImplemented handles requests for check whether the provider specified in the request is implemented.
func IsProviderImplemented(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := r.ParseForm(); err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	name := r.Form.Get("name")
	if len(name) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("provider name was not specified"))
		return
	}

	res := isProviderImplementedResp{
		Implemented: name == string(common.Example) || common.KYCProviders[common.KYCProvider(name)],
	}

	json.NewEncoder(w).Encode(res)
}
