package handlers

import (
	"encoding/json"
	"github.com/pkg/errors"
	"modulus/kyc/integrations/ciphertrace"
	"modulus/kyc/main/config"
	"net/http"
)

// Check txHash for BTC and ETH.
func CipherTraceCheck(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Coin   string `json:"coin"`
		TxHash string `json:"txHash"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		err = &serviceError{
			status:  http.StatusBadRequest,
			message: err.Error(),
		}
		return
	}
	cfg, ok := config.Cfg["CipherTrace"]
	if !ok {
		err = &serviceError{
			status:  http.StatusInternalServerError,
			message: "missing config for CipherTrace",
		}
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	service := ciphertrace.NewCipherService(cfg["URL"], cfg["Key"], cfg["Username"])
	switch req.Coin {
	case "BTC":
		res, err, code := service.GtAddressRiskInfo(req.TxHash)
		if err != nil {
			writeErrorResponse(w, code, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	case "ETH":
		res, err, code := service.GtAddressRiskInfoETH(req.TxHash)
		if err != nil {
			writeErrorResponse(w, code, err)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	default:
		writeErrorResponse(w, http.StatusBadRequest, errors.New("not supported coin"))
	}
}
