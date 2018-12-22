package model

// List of GroupScreeningType values.
const (
	CaseManagementAudit GroupScreeningType = "CASE_MANAGEMENT_AUDIT"
	ZeroFootprint       GroupScreeningType = "ZERO_FOOTPRINT"
)

// Group represents a group accessible by user.
type Group struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ParentID    string     `json:"parentId"`
	Status      StatusType `json:"status"`
	HasChildren bool       `json:"hasChildren"`
	Children    []Group    `json:"children"`
}

// Groups represents a list of groups.
type Groups []Group

// GroupScreeningType represents the group screening type enumeration.
type GroupScreeningType string
