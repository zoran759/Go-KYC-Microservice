package sumsub

type Config struct {
	Host             string
	APIKey           string
	TimeoutThreshold int64
}

const RedScore = "RED"
const YellowScore = "YELLOW"
const GreenScore = "GREEN"
const ErrorScore = "ERROR"
const IgnoredScore = "IGNORED"

const FinalRejectType = "FINAL"
const RetryRejectTYpe = "RETRY"

const CompleteStatus = "completed"
