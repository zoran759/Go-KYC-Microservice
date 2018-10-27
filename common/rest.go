package common

// CheckCustomerRequest represents the request for the CheckCustomer handler.
type CheckCustomerRequest struct {
	Provider KYCProvider
	UserData *UserData
}

// CheckStatusRequest represents the status check request payload of the CheckStatus handler.
type CheckStatusRequest struct {
	Provider   KYCProvider
	CustomerID string
}

// KYCResponse represents the response for the CheckCustomer and the CheckStatus handlers.
type KYCResponse struct {
	Result *KYCResult
	Error  string
}

// ErrorResponse represents the error response payload from the service.
type ErrorResponse struct {
	Error string
}
