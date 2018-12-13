package model

// UserSummary represents a user (customer) in the Thomson Reuters API client’s account.
type UserSummary struct {
	UserID    string     `json:"userId"`
	Email     string     `json:"email"`
	FullName  string     `json:"fullName"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Status    StatusType `json:"status"`
}

// Users represents a list of active users (customers) in the Thomson Reuters API client’s account.
type Users []UserSummary
