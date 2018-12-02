package verification

// Config represents service configuration.
type Config struct {
	Host   string
	APIKey string
}

// Verification represents KYC verification interface.
type Verification interface {
	CheckApplicantStatus(applicantID string) (string, *ReviewResult, error)
	RequestApplicantCheck(applicantID string) error
}

// Mock represents the service mock.
type Mock struct {
	StartVerificationFn    func(applicantID string) (bool, *int, error)
	CheckApplicantStatusFn func(applicantID string) (string, *ReviewResult, error)
}

// CheckApplicantStatus implements Verification interface for the Mock.
func (mock Mock) CheckApplicantStatus(applicantID string) (string, *ReviewResult, error) {
	return mock.CheckApplicantStatusFn(applicantID)
}

// RequestApplicantCheck implements Verification interface for the Mock.
func (mock Mock) RequestApplicantCheck(applicantID string) error {
	return nil
}
