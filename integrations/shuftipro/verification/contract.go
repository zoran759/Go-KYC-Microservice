package verification

type Config struct {
	Host        string
	SecretKey   string
	ClientID    string
	RedirectURL string
}

type Verification interface {
	Verify(request Request) (*Response, error)
	CheckStatus(reference string) (*Response, error)
}

type Mock struct {
	VerifyFn      func(request Request) (*Response, error)
	CheckStatusFn func(reference string) (*Response, error)
}

func (mock Mock) Verify(request Request) (*Response, error) {
	return mock.VerifyFn(request)
}

func (mock Mock) CheckStatus(reference string) (*Response, error) {
	return mock.CheckStatusFn(reference)
}
