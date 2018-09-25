package jumio

// Status represents the verification scan status.
type Status string

// Possible values of Status.
const (
	PendingStatus Status = "PENDING"
	DoneStatus    Status = "DONE"
	FailedStatus  Status = "FAILED"
)
