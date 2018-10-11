package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub"
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

	var service common.StatusChecker

	switch req.Provider {
	case common.IDology:
		fallthrough
	case common.ShuftiPro:
		fallthrough
	case common.Trulioo:
		writeErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("%s doesn't support status polling", req.Provider))
		return
	case common.SumSub:
		service = sumsub.New(sumsub.Config{
			Host:   "https://test-api.sumsub.com",
			APIKey: "GKTBNXNEPJHCXY",
		})
	default:
		writeErrorResponse(w, http.StatusNotFound, fmt.Errorf("unknown KYC provider id in the request: %s", req.Provider))
		return
	}

	response := common.CheckStatusResponse{}

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
