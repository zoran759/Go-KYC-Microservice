package model

// AuthRequest represents the authentication request.
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents the response on the authentication request.
type AuthResponse struct {
	Token string `json:"token"`
}
