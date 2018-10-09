package verification

// Error represents an error in the response if occured.
type Error struct {
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}

// ApplicantStatusResponse represents status check response.
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

// ReviewResult represents a review result in the response.
type ReviewResult struct {
	ReviewAnswer     string   `json:"reviewAnswer"`
	Label            string   `json:"label"`
	RejectLabels     []string `json:"rejectLabels"`
	ReviewRejectType string   `json:"reviewRejectType"`
	ErrorCode        int      `json:"-"`
}
