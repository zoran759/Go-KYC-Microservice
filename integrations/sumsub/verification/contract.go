package verification

// Config represents service configuration.
type Config struct {
	Host   string
	APIKey string
}

// Verificator represents KYC verification checker.
type Verificator interface {
	CheckApplicantStatus(applicantID string) (string, *ReviewResult, error)
}

// Mock represents the service mock.
type Mock struct {
	StartVerificationFn    func(applicantID string) (bool, *int, error)
	CheckApplicantStatusFn func(applicantID string) (string, *ReviewResult, error)
}

// CheckApplicantStatus implememnts Verificator interface for the Mock.
func (mock Mock) CheckApplicantStatus(applicantID string) (string, *ReviewResult, error) {
	return mock.CheckApplicantStatusFn(applicantID)
}
