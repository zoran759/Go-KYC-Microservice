package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm"
	"modulus/kyc/integrations/example"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/synapsefi"
	"modulus/kyc/main/config"
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

	service, err1 := createStatusChecker(req.Provider)
	if err1 != nil {
		writeErrorResponse(w, err1.status, err1)
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

// createStatusChecker returns the KYCPlatform object for the specified provider or an error if occurred.
func createStatusChecker(provider common.KYCProvider) (service common.KYCPlatform, err *serviceError) {
	if provider == common.Example {
		service = example.Example{}
		return
	}

	if !common.KYCProviders[provider] {
		err = &serviceError{
			status:  http.StatusNotFound,
			message: fmt.Sprintf("unknown KYC provider in the request: %s", provider),
		}
		return
	}

	cfg, ok := config.Cfg[string(provider)]
	if !ok {
		err = &serviceError{
			status:  http.StatusInternalServerError,
			message: fmt.Sprintf("missing config for %s", provider),
		}
		return
	}

	switch provider {
	case common.ComplyAdvantage, common.IDology, common.ThomsonReuters, common.Trulioo:
		err = &serviceError{
			status:  http.StatusUnprocessableEntity,
			message: fmt.Sprintf("%s doesn't support status polling", provider),
		}
	case common.Coinfirm:
		service = coinfirm.New(coinfirm.Config{
			Host:     cfg["Host"],
			Email:    cfg["Email"],
			Password: cfg["Password"],
			Company:  cfg["Company"],
		})
	case common.IdentityMind:
		service = identitymind.New(identitymind.Config{
			Host:     cfg["Host"],
			Username: cfg["Username"],
			Password: cfg["Password"],
		})
	case common.Jumio:
		service = jumio.New(jumio.Config{
			BaseURL: cfg["BaseURL"],
			Token:   cfg["Token"],
			Secret:  cfg["Secret"],
		})
	case common.ShuftiPro:
		service = shuftipro.New(shuftipro.Config{
			Host:        cfg["Host"],
			SecretKey:   cfg["SecretKey"],
			ClientID:    cfg["ClientID"],
			CallbackURL: cfg["CallbackURL"],
		})
	case common.SumSub:
		service = sumsub.New(sumsub.Config{
			Host:   cfg["Host"],
			APIKey: cfg["APIKey"],
		})
	case common.SynapseFI:
		service = synapsefi.New(synapsefi.Config{
			Host:         cfg["Host"],
			ClientID:     cfg["ClientID"],
			ClientSecret: cfg["ClientSecret"],
		})
	default:
		err = &serviceError{
			status:  http.StatusUnprocessableEntity,
			message: fmt.Sprintf("KYC provider not implemented yet: %s", provider),
		}
	}

	return
}
