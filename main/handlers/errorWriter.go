package handlers

import (
	"encoding/json"
	"fmt"
	"modulus/kyc/common"
	"net/http"
)

// serviceError represents an error that might happen during creating the KYC provider service.
type serviceError struct {
	status  int
	message string
}

// Error implements the error interface for the serviceError.
func (e serviceError) Error() string {
	return e.message
}

func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	errorResponse := common.ErrorResponse{
		Error: err.Error(),
	}

	resp, err := json.Marshal(errorResponse)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(status)
	w.Write(resp)
}
