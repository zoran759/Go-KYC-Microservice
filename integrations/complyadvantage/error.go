package complyadvantage

import (
	"encoding/json"
	"fmt"
)

// ErrorResponse represents error response from the API.
type ErrorResponse struct {
	Code    int                    `json:"code"`
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

// Error implements error interface for the ErrorResponse.
func (e ErrorResponse) Error() string {
	text := fmt.Sprintf("%d %s", e.Code, e.Message)

	if len(e.Errors) > 0 {
		details, err := json.Marshal(e.Errors)
		if err != nil {
			return text
		}

		text += " | " + string(details)
	}

	return text
}
