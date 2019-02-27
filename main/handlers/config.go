package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"modulus/kyc/main/config"
)

// ConfigResponse represents a configuration change response.
type ConfigResponse struct {
	Updated bool
	Errors  []string
}

// UpdateConfig handles requests for config updates.
func UpdateConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

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

	req := config.Config{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	updated, errs := config.Update(req)
	if updated {
		err = config.Save()
	}

	if err != nil {
		errs = append(errs, err.Error())
	}
	if errs == nil {
		errs = []string{}
	}

	resp := ConfigResponse{
		Updated: updated,
		Errors:  errs,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Config:", err)
	}
}
