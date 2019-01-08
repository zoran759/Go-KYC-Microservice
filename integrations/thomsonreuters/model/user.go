package model

// Users represents a list of active users (customers) in the Thomson Reuters API client’s account.
type Users []UserSummary

// UserSummary represents a user (customer) in the Thomson Reuters API client’s account.
type UserSummary struct {
	UserID    string     `json:"userId"`
	Email     string     `json:"email"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	FullName  string     `json:"fullName"`
	Status    StatusType `json:"status"`
}
