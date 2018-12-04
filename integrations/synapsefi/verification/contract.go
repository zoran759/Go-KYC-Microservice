package verification

// Config represents service config.
type Config struct {
	Host         string
	ClientID     string
	ClientSecret string
}

// Verification describes the verification interface.
type Verification interface {
	CreateUser(request CreateUserRequest) (*UserResponse, error)
	GetOauthKey(userID string, request CreateOauthRequest) (*OauthResponse, error)
	AddDocument(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error)
	GetUser(userID string) (*UserResponse, error)
}

// Mock represents the service mock.
type Mock struct {
	CreateUserFn  func(request CreateUserRequest) (*UserResponse, error)
	AddDocumentFn func(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error)
	GetOauthKeyFn func(userID string, request CreateOauthRequest) (*OauthResponse, error)
	GetUserFn     func(userID string) (*UserResponse, error)
}

// CreateUser implements the Verification interface for the Mock.
func (mock Mock) CreateUser(request CreateUserRequest) (*UserResponse, error) {
	return mock.CreateUserFn(request)
}

// AddDocument implements the Verification interface for the Mock.
func (mock Mock) AddDocument(userID string, userOAuth string, request CreateDocumentsRequest) (*UserResponse, error) {
	return mock.AddDocumentFn(userID, userOAuth, request)
}

// GetOauthKey implements the Verification interface for the Mock.
func (mock Mock) GetOauthKey(userID string, request CreateOauthRequest) (*OauthResponse, error) {
	return mock.GetOauthKeyFn(userID, request)
}

// GetUser implements the Verification interface for the Mock.
func (mock Mock) GetUser(userID string) (*UserResponse, error) {
	return mock.GetUserFn(userID)
}
