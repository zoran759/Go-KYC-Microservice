package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm"
	"modulus/kyc/integrations/complyadvantage"
	"modulus/kyc/integrations/example"
	"modulus/kyc/integrations/identitymind"
	"modulus/kyc/integrations/idology"
	"modulus/kyc/integrations/jumio"
	"modulus/kyc/integrations/shuftipro"
	"modulus/kyc/integrations/sumsub"
	"modulus/kyc/integrations/synapsefi"
	"modulus/kyc/integrations/thomsonreuters"
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
	requestTime := time.Now().UnixNano()
	log.Println("CheckCustomer Request time: ", requestTime)
	ioutil.WriteFile("requests/"+fmt.Sprint(requestTime)+".json", body, 0644)

	if len(req.Provider) == 0 {
		writeErrorResponse(w, http.StatusBadRequest, errors.New("missing KYC provider id in the request"))
		return
	}

	service, err1 := createCustomerChecker(req.Provider)
	if err1 != nil {
		log.Println("CheckCustomer Error: ", err1)
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
	log.Println("CheckCustomer Response: ", string(resp))
	w.Write(resp)
}

// createCustomerChecker returns the KYCPlatform object for the specified provider or an error if occurred.
func createCustomerChecker(provider common.KYCProvider) (service common.KYCPlatform, err *serviceError) {
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
	case common.Coinfirm:
		service = coinfirm.New(coinfirm.Config{
			Host:     cfg["Host"],
			Email:    cfg["Email"],
			Password: cfg["Password"],
			Company:  cfg["Company"],
		})
	case common.ComplyAdvantage:
		fuzziness, err1 := strconv.ParseFloat(cfg["Fuzziness"], 32)
		if err1 != nil {
			err = &serviceError{
				status:  http.StatusInternalServerError,
				message: fmt.Sprintf("%s config error: %s", provider, err1),
			}
			return
		}
		service = complyadvantage.New(complyadvantage.Config{
			Host:      cfg["Host"],
			APIkey:    cfg["APIkey"],
			Fuzziness: float32(fuzziness),
		})
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
			RedirectURL: cfg["RedirectURL"],
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
	case common.ThomsonReuters:
		service = thomsonreuters.New(thomsonreuters.Config{
			Host:      cfg["Host"],
			APIkey:    cfg["APIkey"],
			APIsecret: cfg["APIsecret"],
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
