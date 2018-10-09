package configuration

type Config struct {
	Host  string
	Token string
}

type Configuration interface {
	Consents(countryAlpha2 string) (Consents, *int, error)
}

type Mock struct {
	ConsentsFn func(countryAlpha2 string) (Consents, *int, error)
}

func (mock Mock) Consents(countryAlpha2 string) (Consents, *int, error) {
	return mock.ConsentsFn(countryAlpha2)
}
