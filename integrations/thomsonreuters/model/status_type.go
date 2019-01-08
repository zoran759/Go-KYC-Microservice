package model

// StatusType represents the enumeration of status types.
type StatusType string

// List of StatusType values.
const (
	ActiveStatus   StatusType = "ACTIVE"
	InactiveStatus StatusType = "INACTIVE"
	DeletedStatus  StatusType = "DELETED"
)
