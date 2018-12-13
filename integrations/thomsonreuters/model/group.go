package model

// Group represents a group accessible by user.
type Group struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	ParentID    string     `json:"parentId"`
	HasChildren bool       `json:"hasChildren"`
	Children    []Group    `json:"children"`
	Status      StatusType `json:"status"`
}

// Groups represents a list of groups.
type Groups []Group
