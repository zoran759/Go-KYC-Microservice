package handlers

import (
	"encoding/json"
	"net/http"

	"modulus/kyc/common"
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

// writeErrorResponse writes the error response to the connection using the specified HTTP status code and error object.
func writeErrorResponse(w http.ResponseWriter, status int, err error) {
	errorResponse := common.ErrorResponse{
		Error: err.Error(),
	}

	resp, _ := json.Marshal(errorResponse)

	w.WriteHeader(status)
	w.Write(resp)
}
