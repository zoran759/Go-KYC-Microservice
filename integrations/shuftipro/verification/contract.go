package verification

type Config struct {
	Host        string
	SecretKey   string
	ClientID    string
	RedirectURL string
}

type Verification interface {
	Verify(request Request) (*Response, error)
}

type Mock struct {
	VerifyFn func(request Request) (*Response, error)
}

func (mock Mock) Verify(request Request) (*Response, error) {
	return mock.VerifyFn(request)
}
