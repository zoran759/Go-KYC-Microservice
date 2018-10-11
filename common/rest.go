package common

// CheckCustomerRequest represents the request for the CheckCustomer handler.
type CheckCustomerRequest struct {
	Provider KYCProvider
	UserData UserData
}

// CheckCustomerResponse represents the response for the CheckCustomer handler.
type CheckCustomerResponse struct {
	// KYC result
	KYCResult KYCResult
}

// CheckStatusRequest represents the status check request payload of the CheckStatus handler.
type CheckStatusRequest struct {
	Provider   KYCProvider
	CustomerID string
}

// CheckStatusResponse represents the response payload for the CheckStatus handler.
type CheckStatusResponse struct {
	Result *KYCResult
	Error  string
}

// ErrorResponse represents the error response payload from the service.
type ErrorResponse struct {
	Error string
}
