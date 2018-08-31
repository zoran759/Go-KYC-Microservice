package verification

type Config struct {
	Host         string
	ClientID     string
	ClientSecret string
}

type Verification interface {
	CreateUser(request CreateUserRequest) (*UserResponse, error)
	GetUser(userID string) (*UserResponse, error)
}
