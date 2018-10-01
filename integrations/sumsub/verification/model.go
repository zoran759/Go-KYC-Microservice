package verification

type StartVerificationResponse struct {
	OK int `json:"ok"`
	Error
}

type Error struct {
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}

type ApplicantStatusResponse struct {
	ID           string       `json:"id"`
	InspectionID string       `json:"inspectionId"`
	JobID        string       `json:"jobId"`
	CreateDate   string       `json:"createDate"`
	ReviewDate   string       `json:"reviewDate"`
	ReviewResult ReviewResult `json:"reviewResult"`
	ReviewStatus string       `json:"reviewStatus"`
	ApplicantID  string       `json:"applicantId"`
	Error
}

type ReviewResult struct {
	ReviewAnswer     string   `json:"reviewAnswer"`
	Label            string   `json:"label"`
	RejectLabels     []string `json:"rejectLabels"`
	ReviewRejectType string   `json:"reviewRejectType"`
	ErrorCode        int      `json:"-"`
}
