package jumio

// ErrorResponse represents a response with an error.
type ErrorResponse struct {
	Message    string `json:"message"`
	HTTPStatus string `json:"httpStatus"`
	RequestURI string `json:"requestURI"`
}

// Error implements the error interface for the ErrorResponse.
func (e ErrorResponse) Error() string {
	return "HTTP status: " + e.HTTPStatus + " | " + e.Message
}
