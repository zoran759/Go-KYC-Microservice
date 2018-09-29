package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"modulus/kyc/common"
	"modulus/kyc/common/kycErrors"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/trulioo"
	"net/http"
)

// Handler for the CustomerHandler function
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

	// Prepare result variables
	res := new(common.CheckCustomerResponse)

	// Handle request depending on given KYC Provider
	var service common.CustomerChecker
	switch req.Provider {

	// Shufti Pro integration
	case common.Shuftipro:

		// Initialize the service
		service = shuftipro.New(shuftipro.Config{
			Host:        "https://api.shuftipro.com",
			ClientID:    "ac93f3a0fee5afa2d9399d5d0f257dc92bbde89b1e48452e1bfac3c5c1dc99db",
			SecretKey:   "lq34eOTxDe1e6G8a1P7Igqo5YK3ABCDF",
			RedirectURL: "https://api.shuftipro.com",
		})

	// Sum & Substance integration
	case common.SumSubstance:

		// Initialize the service
		service = sumsub.New(sumsub.Config{
			Host:             "test_host",
			APIKey:           "test_key",
			TimeoutThreshold: 123456,
		})

	// Trulioo integration
	case common.Trulioo:

		// Initialize the service
		service = trulioo.New(trulioo.Config{
			Host:         "host",
			NAPILogin:    "login",
			NAPIPassword: "password",
		})

		/*
			// IDology integration : Functionality isn't finished yet, the CheckCustomer method isn't implemented
			case common.IDology:

				// Initialize the service
				service = idology.New(idology.Config{
					Host         :"host",
					Username    :"username",
					Password :"password",
					UseSummaryResult: true, // or false
				})
		*/

	// KYC provider is not recognized
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: kycErrors.InvalidKYCProvider.Error()})
	}

	// Make a request to the KYC provider
	res.KYCResult, res.DetailedKYCResult, err = service.CheckCustomer(&req.UserData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	// Debug output
	log.Printf("Res: %#v\n", res.KYCResult)
	log.Printf("detailedRes: %#v\n", res.DetailedKYCResult)

	// Send the response over HTTP
	json.NewEncoder(w).Encode(res)
}
