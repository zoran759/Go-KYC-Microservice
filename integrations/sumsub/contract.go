package sumsub

// Config defines configuration for the service.
type Config struct {
	Host   string
	APIKey string
}

// Different values of a verification result.
const (
	RedScore     = "RED"
	YellowScore  = "YELLOW"
	GreenScore   = "GREEN"
	ErrorScore   = "ERROR"
	IgnoredScore = "IGNORED"
)

// Different types of the rejection response.
const (
	FinalRejectType = "FINAL"
	RetryRejectType = "RETRY"
)
