package handlers

import (
	"modulus/kyc/common"
	"modulus/kyc/common/kycErrors"
	"modulus/kyc/integrations/shuftipro"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// CheckCustomerHandler represents the handler for the CustomerHandler function.
func CheckCustomerHandler(w http.ResponseWriter, r *http.Request) {

	// Read request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	// Parse request
	var req common.CheckCustomerRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("CustomerHandler request received Customer:%#v\nProvider:%v\n", req.UserData, req.Provider)

	// Prepare result variables
	res := new(common.CheckCustomerResponse)

	// Handle request depending on given KYC Provider
	switch req.Provider {

	case common.Shuftipro:
		{
			//Example Shufti Pro integration
			service := shuftipro.New(shuftipro.Config{
				Host:        "https://api.shuftipro.com",
				ClientID:    "ac93f3a0fee5afa2d9399d5d0f257dc92bbde89b1e48452e1bfac3c5c1dc99db",
				SecretKey:   "lq34eOTxDe1e6G8a1P7Igqo5YK3ABCDF",
				RedirectURL: "https://api.shuftipro.com",
			})

			// Make a request to the KYC provider
			res.KYCResult, err = service.CheckCustomer(&req.UserData)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
				return
			}

			log.Printf("Res: %#v\n", res.KYCResult)
			log.Printf("detailedRes: %#v\n", res.KYCResult.Details)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: kycErrors.InvalidKYCProvider.Error()})
	}

	// Send the response over HTTP
	json.NewEncoder(w).Encode(res)
}
