package configuration

type Consents []string

type Error struct {
	Message string
}

func (error Error) Error() string {
	return error.Message
}
