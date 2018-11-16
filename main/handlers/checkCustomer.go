package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"modulus/kyc/common"
	"modulus/kyc/integrations/example"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/idology"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/trulioo"
	"modulus/kyc/main/config"
)

// CheckCustomer handles requests for KYC verifications.
func CheckCustomer(w http.ResponseWriter, r *http.Request) {
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

	req := common.CheckCustomerRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if len(req.Provider) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing KYC provider id in the request"))
		return
	}

	service, err1 := createCustomerChecker(req.Provider)
	if err1 != nil {
		writeErrorResponse(w, err1.status, err1)
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

	w.Write(resp)
}

// createCustomerChecker returns the CustomerChecker object for the specified provider or an error if occurred.
func createCustomerChecker(provider common.KYCProvider) (service common.CustomerChecker, err *serviceError) {
	if provider == common.Example {
		service = &example.Service{}
		return
	}

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
	case common.IdentityMind:
		service = identitymind.New(identitymind.Config{
			Host:     cfg["Host"],
			Username: cfg["Username"],
			Password: cfg["Password"],
		})
	case common.IDology:
		useSummaryResult, err1 := strconv.ParseBool(cfg["UseSummaryResult"])
		if err1 != nil {
			err = &serviceError{
				status:  http.StatusInternalServerError,
				message: fmt.Sprintf("%s config error: %s", provider, err1),
			}
			return
		}
		service = idology.New(idology.Config{
			Host:             cfg["Host"],
			Username:         cfg["Username"],
			Password:         cfg["Password"],
			UseSummaryResult: useSummaryResult,
		})
	case common.ShuftiPro:
		service = shuftipro.New(shuftipro.Config{
			Host:        cfg["Host"],
			SecretKey:   cfg["SecretKey"],
			ClientID:    cfg["ClientID"],
			RedirectURL: cfg["RedirectURL"],
		})
	case common.SumSub:
		service = sumsub.New(sumsub.Config{
			Host:   cfg["Host"],
			APIKey: cfg["APIKey"],
		})
	case common.Trulioo:
		service = trulioo.New(trulioo.Config{
			Host:         cfg["Host"],
			NAPILogin:    cfg["NAPILogin"],
			NAPIPassword: cfg["NAPIPassword"],
		})
	default:
		err = &serviceError{
			status:  http.StatusUnprocessableEntity,
			message: fmt.Sprintf("KYC provider not implemented yet: %s", provider),
		}
	}

	return
}
