package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"modulus/kyc/common"
	"modulus/kyc/main/handlers/providers"
)

// CheckStatus handles requests for a status check.
func CheckStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if len(body) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty request"))
		return
	}

	req := common.CheckStatusRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if len(req.Provider) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing KYC provider id in the request"))
		return
	}
	if len(req.ReferenceID) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing verification id in the request"))
		return
	}

	service := providers.Provider(req.Provider)
	if service == nil {
		writeErrorResponse(w, http.StatusNotFound, fmt.Errorf("unknown KYC provider in the request: %s", req.Provider))
		return
	}

	response := common.KYCResponse{}

	result, err := service.CheckStatus(req.ReferenceID)
	if err != nil {
		response.Error = err.Error()
	}

	response.Result = common.ResultFromKYCResult(result)

	resp, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(resp)
}
