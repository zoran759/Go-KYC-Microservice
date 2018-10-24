package verification

type Config struct {
	Host         string
	ClientID     string
	ClientSecret string
}

type Verification interface {
	CreateUser(request CreateUserRequest) (*UserResponse, error)
	GetOauthKey(userID string, request CreateOauthRequest) (*OauthResponse, error)
	AddDocument(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error)
	GetUser(userID string) (*UserResponse, error)
}
type Mock struct {
	CreateUserFn func(request CreateUserRequest) (*UserResponse, error)
	AddDocumentFn func(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error)
	GetOauthKeyFn func(userID string, request CreateOauthRequest) (*OauthResponse, error)
	GetUserFn    func(userID string) (*UserResponse, error)
}

func (mock Mock) CreateUser(request CreateUserRequest) (*UserResponse, error) {
	return mock.CreateUserFn(request)
}

func (mock Mock) AddDocument(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error) {
	return mock.AddDocumentFn(userID, userOAuth, request)
}

func (mock Mock) GetOauthKey(userID string, request CreateOauthRequest) (*OauthResponse, error) {
	return mock.GetOauthKeyFn(userID, request)
}

func (mock Mock) GetUser(userID string) (*UserResponse, error) {
	return mock.GetUserFn(userID)
}
