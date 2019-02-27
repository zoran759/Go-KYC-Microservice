package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/main/config/providers"
	"modulus/kyc/main/license"
)

// CheckCustomer handles requests for KYC verifications.
func CheckCustomer(w http.ResponseWriter, r *http.Request) {
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

	req := common.CheckCustomerRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	requestTime := time.Now().UnixNano()
	log.Println("CheckCustomer Request time: ", requestTime)
	ioutil.WriteFile("requests/"+fmt.Sprint(requestTime)+".json", body, 0644)

	if len(req.Provider) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing KYC provider id in the request"))
		return
	}

	service, err1 := providers.GetPlatform(req.Provider)
	if err1 != nil {
		log.Println("CheckCustomer Error: ", err1)
		writeErrorResponse(w, http.StatusNotFound, err1)
		return
	}

	response := common.KYCResponse{}

	result, err := service.CheckCustomer(req.UserData)
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
