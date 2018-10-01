package verification

type Config struct {
	Host   string
	APIKey string
}

type Verification interface {
	StartVerification(applicantID string) (bool, *int, error)
	CheckApplicantStatus(applicantID string) (string, *ReviewResult, error)
}

type Mock struct {
	StartVerificationFn    func(applicantID string) (bool, *int, error)
	CheckApplicantStatusFn func(applicantID string) (string, *ReviewResult, error)
}

func (mock Mock) StartVerification(applicantID string) (bool, *int, error) {
	return mock.StartVerificationFn(applicantID)
}

func (mock Mock) CheckApplicantStatus(applicantID string) (string, *ReviewResult, error) {
	return mock.CheckApplicantStatusFn(applicantID)
}
