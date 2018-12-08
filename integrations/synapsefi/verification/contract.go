package verification

// Config represents service config.
type Config struct {
	Host         string
	ClientID     string
	ClientSecret string
	fingerprint  string
}

// Verification describes the verification interface.
type Verification interface {
	CreateUser(User) (*Response, *string, error)
	AddPhysicalDocs(string, string, PhysicalDocs) (*string, error)
	GetUser(string) (*Response, *string, error)
}
