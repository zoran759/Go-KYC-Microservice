package verification

import "fmt"

// ErrorResponse represents error response from the API.
type ErrorResponse struct {
	Text     map[string]string `json:"error"`
	Code     string            `json:"error_code"`
	HTTPCode string            `json:"http_code"`
	Success  bool              `json:"success"`
}

// Error implements error interface for the ErrorResponse.
func (er ErrorResponse) Error() string {
	return fmt.Sprintf("http status: %s | error code: %s | error: %s", er.HTTPCode, er.Code, er.Text[appLanguage])
}
