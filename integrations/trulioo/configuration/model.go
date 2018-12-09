package configuration

// Consents represents consents required for the verification.
type Consents []string

// Error represents error returning in case of a failure of getting consents.
type Error struct {
	Message string
}

// Error() implements the error interface for the Error.
func (error Error) Error() string {
	return error.Message
}
