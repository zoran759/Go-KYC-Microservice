package model

// List of ResolutionStatusType values.
const (
	Positive    ResolutionStatusType = "POSITIVE"
	Possible    ResolutionStatusType = "POSSIBLE"
	False       ResolutionStatusType = "FALSE"
	Unspecified ResolutionStatusType = "UNSPECIFIED"
)

// ResolutionStatusType represents the enumeration of resolution status types.
type ResolutionStatusType string

// ResolutionField represents the ID, label and type for a Status, a Risk or a Reason.
type ResolutionField struct {
	ID    string               `json:"id"`
	Label string               `json:"label"`
	Type  ResolutionStatusType `json:"type"`
}

// ResolutionFields describes all resolution fields (statuses, risks and reasons) in detail.
type ResolutionFields struct {
	Reasons  []ResolutionField `json:"reasons"`
	Risks    []ResolutionField `json:"risks"`
	Statuses []ResolutionField `json:"statuses"`
}

// StatusRule represents the rules that should be applied when resolving a Result with a specific Status.
type StatusRule struct {
	ReasonRequired bool     `json:"reasonRequired"`
	Reasons        []string `json:"reasons"`
	RemarkRequired bool     `json:"remarkRequired"`
	Risks          []string `json:"risks"`
}

// ResolutionRules represents a collection of StatusRules keyed by Status ID.
type ResolutionRules map[string]StatusRule

// ResolutionToolkitResponse represents a resolution toolkit settings applicable for a Group.
type ResolutionToolkitResponse struct {
	GroupID          string           `json:"groupId"`
	ResolutionFields ResolutionFields `json:"resolutionFields"`
	ResolutionRules  ResolutionRules  `json:"resolutionRules"`
}

// ResolutionToolkits represents Map of ResolutionToolkits for the given Group keyed by provider type
type ResolutionToolkits map[string]ResolutionToolkitResponse
