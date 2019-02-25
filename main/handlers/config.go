package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"modulus/kyc/main/config"
)

// ConfigRequest represents a configuration change request.
type ConfigRequest struct {
	Config config.Config
}

// ConfigResponse represents a configuration change response.
type ConfigResponse struct {
	Errors []string `json:",omitempty"`
}

// UpdateConfig handles requests for config updates.
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
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

	req := ConfigRequest{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	updated, errs := config.Update(req.Config)

	if updated {
		config.Save()
	}

	resp := ConfigResponse{
		Errors: errs,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Config:", err)
	}
}
