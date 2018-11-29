package common

// TooManyRequests defines the error code returned when KYC status check requests send too frequently.
const TooManyRequests = "429"

// CheckCustomerRequest represents the request for the CheckCustomer handler.
type CheckCustomerRequest struct {
	Provider KYCProvider
	UserData *UserData
}

// CheckStatusRequest represents the status check request payload of the CheckStatus handler.
type CheckStatusRequest struct {
	Provider    KYCProvider
	ReferenceID string
}

// ErrorResponse represents the error response payload from the service.
type ErrorResponse struct {
	Error string
}

// KYCResponse represents the response for the CheckCustomer and the CheckStatus handlers.
type KYCResponse struct {
	Result *Result
	Error  string
}

// Result represents the verification result for the KYCResponse.
type Result struct {
	Status      string
	Details     *Details
	ErrorCode   string
	StatusCheck *KYCStatusCheck
}

// Details defines additional details about the verification result.
type Details struct {
	Finality string
	Reasons  []string
}

// ResultFromKYCResult converts KYC verification result into the API representation.
func ResultFromKYCResult(kycResult KYCResult) (result *Result) {
	result = &Result{}

	result.Status = KYCStatus2Status[kycResult.Status]
	if kycResult.Details != nil {
		result.Details = &Details{
			Finality: KYCFinality2Finality[kycResult.Details.Finality],
			Reasons:  kycResult.Details.Reasons,
		}
	}
	result.ErrorCode = kycResult.ErrorCode
	result.StatusCheck = kycResult.StatusCheck

	return
}
