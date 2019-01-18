package model

// List of status values.
const (
	// Participant has been invited to join to KYC process by email, however we have not received filled KYC form from participant.
	Empty Status = "empty"
	// Participant provided filled KYC form and can be evaluated.
	New Status = "new"
	// The KYC form which was submitted by participant is under evaluation by Coinfirm analysts.
	InProgress Status = "inprogress"
	// Data provided by participant is incomplete or does not meet the requirements set in KYC form (e.g. expired proof of residency).
	// Participant is supposed submit corrected KYC following the email notification.
	Incomplete Status = "incomplete"
	// Coinfirm analysts evaluated the risk associated to participant as low.
	Low Status = "low"
	// Coinfirm analysts evaluated the risk associated to participant as medium.
	Medium Status = "medium"
	// Coinfirm analysts evaluated the risk associated to participant as high.
	High Status = "high"
	// Coinfirm analysts evaluated the risk associated to participant as unacceptable
	Fail Status = "fail"
)

// Status represents the current participant Status in the Coinfirm KYC flow.
type Status string

// StatusResponse represents the response on the request of the current participant status.
type StatusResponse struct {
	CurrentStatus  Status `json:"current_status"`
	ICOOwnerStatus string `json:"ico_owner_status"`
}
