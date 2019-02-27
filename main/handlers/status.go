package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"modulus/kyc/common"
	"modulus/kyc/main/config/providers"
	"modulus/kyc/main/license"
)

// CheckStatus handles requests for a status check.
func CheckStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !license.Valid() {
		writeErrorResponse(w, http.StatusForbidden, license.ErrNoValidLicense)
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeErrorResponse(w, http.StatusMethodNotAllowed, errors.New("used method not allowed for this endpoint"))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	if len(body) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("empty request"))
		return
	}
	log.Println("CheckStatus Request: ", string(body))
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

	service, err1 := providers.GetPlatform(req.Provider)
	if err1 != nil {
		writeErrorResponse(w, http.StatusNotFound, err1)
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
	log.Println("CheckCustomer Response: ", string(resp))
	w.Write(resp)
}
