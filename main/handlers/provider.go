package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"modulus/kyc/common"
	"modulus/kyc/main/handlers/providers"
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
	if len(r.Form) == 0 {
		json.NewEncoder(w).Encode(providerList())
		return
	}

	name := r.Form.Get("name")
	if len(name) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing provider name in the request"))
		return
	}

	res := isProviderImplementedResp{
		Implemented: name == string(common.Example) || common.KYCProviders[common.KYCProvider(name)],
	}

	json.NewEncoder(w).Encode(res)
}

// providerList forms the list of implemented providers.
func providerList() (list providers.ProviderList) {
	for p := range common.KYCProviders {
		list = append(list, p)
	}
	sort.Sort(list)
	return
}
