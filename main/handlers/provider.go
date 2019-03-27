package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"modulus/kyc/common"
	"modulus/kyc/main/config"
	"modulus/kyc/main/config/providers"
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
		Implemented: config.IsKnownName(name),
	}

	json.NewEncoder(w).Encode(res)
}

// providerList forms the list of implemented providers.
func providerList() (list providers.List) {
	for p := range providers.Providers {
		list = append(list, p)
	}
	for np := range config.NotProviders {
		if np != config.ServiceSection {
			list = append(list, common.KYCProvider(np))
		}
	}

	sort.Sort(list)
	return
}
