package verification

import "gitlab.com/modulusglobal/kyc/integrations/trulioo/configuration"

type Config struct {
	Host  string
	Token string
}

type Verification interface {
	Verify(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*VerificationResponse, error)
}

type Mock struct {
	VerifyFn func(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*VerificationResponse, error)
}

func (mock Mock) Verify(countryAlpha2 string, consents configuration.Consents, fields DataFields) (*VerificationResponse, error) {
	return mock.VerifyFn(countryAlpha2, consents, fields)
}
