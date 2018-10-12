package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/main/config"
	"net/http"
)

// CheckStatusHandler handles requests for a status check.
func CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
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
	if len(req.CustomerID) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing verification id in the request"))
		return
	}

	service, err1 := createStatusChecker(req.Provider)
	if err1 != nil {
		writeErrorResponse(w, err1.status, err1)
		return
	}

	response := common.KYCResponse{}

	result, err := service.CheckStatus(req.CustomerID)
	if err != nil {
		response.Error = err.Error()
	}
	response.Result = &result

	resp, err := json.Marshal(response)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.Write(resp)
}

func createStatusChecker(provider common.KYCProvider) (service common.StatusChecker, err *serviceError) {
	if !common.KYCProviders[provider] {
		err = &serviceError{
			status:  http.StatusNotFound,
			message: fmt.Sprintf("unknown KYC provider in the request: %s", provider),
		}
		return
	}

	cfg, ok := config.KYC[provider]
	if !ok {
		err = &serviceError{
			status:  http.StatusInternalServerError,
			message: fmt.Sprintf("missing config for %s", provider),
		}
		return
	}

	switch provider {
	case common.IDology, common.ShuftiPro, common.Trulioo:
		err = &serviceError{
			status:  http.StatusUnprocessableEntity,
			message: fmt.Sprintf("%s doesn't support status polling", provider),
		}
	case common.SumSub:
		service = sumsub.New(sumsub.Config{
			Host:   cfg["Host"],
			APIKey: cfg["APIKey"],
		})
	default:
		err = &serviceError{
			status:  http.StatusInternalServerError,
			message: fmt.Sprintf("KYC provider not implemented yet: %s", provider),
		}
	}

	return
}
