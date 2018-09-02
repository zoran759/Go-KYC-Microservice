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
type Mock struct {
	CreateUserFn func(request CreateUserRequest) (*UserResponse, error)
	GetUserFn    func(userID string) (*UserResponse, error)
}

func (mock Mock) CreateUser(request CreateUserRequest) (*UserResponse, error) {
	return mock.CreateUserFn(request)
}

func (mock Mock) GetUser(userID string) (*UserResponse, error) {
	return mock.GetUserFn(userID)
}
