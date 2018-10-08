package kycErrors

import "errors"

// Error thrown into HTTP response
type ErrorResponse struct {
	Error     string            `json:"error,omitempty"`
	ErrorData map[string]string `json:"errorData,omitempty"`
}

var (
	InvalidKYCProvider = errors.New("INVALID_KYC_PROVIDER")
)
