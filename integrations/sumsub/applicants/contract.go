package applicants

type Config struct {
	Host   string
	APIKey string
}

type Applicants interface {
	CreateApplicant(email string, applicant ApplicantInfo) (*CreateApplicantResponse, error)
}

type Mock struct {
	CreateApplicantFn func(email string, applicant ApplicantInfo) (*CreateApplicantResponse, error)
}

func (mock Mock) CreateApplicant(email string, applicant ApplicantInfo) (*CreateApplicantResponse, error) {
	return mock.CreateApplicantFn(email, applicant)
}
