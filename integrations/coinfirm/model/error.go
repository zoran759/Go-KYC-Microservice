package model

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Err string `json:"error"`
}

// Error implements the error interface for ErrorResponse.
func (e ErrorResponse) Error() string {
	return e.Err
}
