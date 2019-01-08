package model

import "strings"

// ErrorEntity represents error.
type ErrorEntity struct {
	Error    string `json:"error"`
	Cause    string `json:"cause"`
	ObjectID string `json:"objectId"`
}

// Errors represents Error response payload.
type Errors []ErrorEntity

// Error implements the error interface for the Errors.
func (e Errors) Error() string {
	b := strings.Builder{}

	for _, err := range e {
		if b.Len() > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(err.Error)
		b.WriteString(" (")
		b.WriteString(err.Cause)
		b.WriteString(")")
		if len(err.ObjectID) > 0 {
			b.WriteByte(' ')
			b.WriteString(err.ObjectID)
		}
	}

	return b.String()
}
