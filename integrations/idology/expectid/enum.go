package expectid

// SummaryResultKey defines a message identifier for summary-result.
type SummaryResultKey string

// Possible SummaryResultKey values.
const (
	Success SummaryResultKey = "id.success"
	Failure SummaryResultKey = "id.failure"
	Partial SummaryResultKey = "id.partial"
)

// SummaryResultMessage defines the message for summary-result.
type SummaryResultMessage string

// Possible SummaryResultMessage values.
const (
	PassMessage SummaryResultMessage = "Pass"
	FailMessage SummaryResultMessage = "Fail"
)

// ResultKey defines ExpectID result.
// This value indicates if the information submitted matched the information located,
// subject to IDology’s default tolerances and error-correcting logic (if enabled)
// on common typos in addresses or names.
type ResultKey string

// Possible ResultKey values.
const (
	Match           ResultKey = "result.match"
	NoMatch         ResultKey = "result.no.match"
	MatchRestricted ResultKey = "result.match.restricted"
)
