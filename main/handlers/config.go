package handlers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"modulus/kyc/main/config"
)

// ConfigUpdateResponse represents a configuration change response.
type ConfigUpdateResponse struct {
	Updated bool
	Errors  []string
}

// ConfigHandler handles requests for config management.
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var (
		resp   interface{}
		status int
		err    error
	)

	switch r.Method {
	case http.MethodGet:
		resp = getConfig()
	case http.MethodPost:
		resp, status, err = updateConfig(r)
	default:
		w.Header().Set("Allow", http.MethodGet+", "+http.MethodPost)
		writeErrorResponse(w, http.StatusMethodNotAllowed, errors.New("used method not allowed for this endpoint"))
		return
	}

	if err != nil {
		writeErrorResponse(w, status, err)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("Config:", err)
	}
}

func updateConfig(r *http.Request) (resp ConfigUpdateResponse, status int, err error) {
	status = http.StatusBadRequest

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = http.StatusInternalServerError
		return
	}
	if len(body) == 0 {
		err = errors.New("empty request")
		return
	}

	req := config.Config{}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return
	}

	updated, errs := config.Update(req)
	if updated {
		err = config.Save()
	}

	if err != nil {
		errs = append(errs, err.Error())
		err = nil
	}
	if errs == nil {
		errs = []string{}
	}

	resp = ConfigUpdateResponse{
		Updated: updated,
		Errors:  errs,
	}

	return
}

func getConfig() config.Config {
	return config.GetConfig()
}
