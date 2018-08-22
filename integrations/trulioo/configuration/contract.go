package configuration

type Config struct {
	Host  string
	Token string
}

type Configuration interface {
	Consents(countryAlpha2 string) (Consents, error)
}

type Mock struct {
	ConsentsFn func(countryAlpha2 string) (Consents, error)
}

func (mock Mock) Consents(countryAlpha2 string) (Consents, error) {
	return mock.ConsentsFn(countryAlpha2)
}
