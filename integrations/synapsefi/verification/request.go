package verification

// OAuthRequest represents OAuth token obtaining request.
type OAuthRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// OAuthResponse represents response on OAuth token request.
type OAuthResponse struct {
	ID           string `json:"user_id"`
	OAuthKey     string `json:"oauth_key"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// User represents the user to verify with KYC.
type User struct {
	Logins       []Login    `json:"logins"`
	PhoneNumbers []string   `json:"phone_numbers"`
	LegalNames   []string   `json:"legal_names"`
	Documents    []Document `json:"documents"`
	Extra        Extra      `json:"extra,omitempty"`
}

// Login represents login data for onboarding the user.
type Login struct {
	Email string `json:"email"`
	Scope string `json:"scope"`
}

// Extra represents extra info of the user.
type Extra struct {
	CIPTag     int  `json:"cip_tag"`
	IsBusiness bool `json:"is_business"`
}

// Document represents customer's KYC Document.
type Document struct {
	OwnerName          string        `json:"name"`
	Email              string        `json:"email"`
	PhoneNumber        string        `json:"phone_number"`
	IPAddress          string        `json:"ip"`
	EntityType         string        `json:"entity_type"`
	EntityScope        string        `json:"entity_scope"`
	DayOfBirth         int           `json:"day"`
	MonthOfBirth       int           `json:"month"`
	YearOfBirth        int           `json:"year"`
	AddressStreet      string        `json:"address_street"`
	AddressCity        string        `json:"address_city"`
	AddressSubdivision string        `json:"address_subdivision"`
	AddressPostalCode  string        `json:"address_postal_code"`
	AddressCountryCode string        `json:"address_country_code"`
	VirtualDocs        []SubDocument `json:"virtual_docs,omitempty"`
	PhysicalDocs       []SubDocument `json:"physical_docs,omitempty"`
}

// SubDocument represents sub-document object for KYC document.
type SubDocument struct {
	Type  string `json:"document_type"`
	Value string `json:"document_value"`
}

// PhysicalDocs represents physical KYC documents to add to a user.
type PhysicalDocs struct {
	Documents []Document `json:"documents"`
}
