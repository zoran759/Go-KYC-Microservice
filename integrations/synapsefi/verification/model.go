package verification

// ResponseError represents error response from the API.
type ResponseError struct {
	Error     map[string]string `json:"error"`
	ErrorCode string            `json:"error_code"`
	HTTPCode  string            `json:"http_code"`
	Status    bool              `json:"success"`
}

// CreateDocumentsRequest represents KYC documents addition request.
type CreateDocumentsRequest struct {
	Documents Document `json:"documents"`
}

// CreateUserRequest represents user creation request.
type CreateUserRequest struct {
	Logins       []Login    `json:"logins"`
	PhoneNumbers []string   `json:"phone_numbers"`
	LegalNames   []string   `json:"legal_names"`
	Documents    []Document `json:"documents"`
	Extra        Extra      `json:"extra,omitempty"`
}

// CreateOauthRequest represents OAuth token obtaining request.
type CreateOauthRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// OauthResponse represents OAuth response.
type OauthResponse struct {
	ID           string `json:"user_id"`
	OAuthKey     string `json:"oauth_key"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

// Login represents login data of the customer in a user creation request.
type Login struct {
	Email string `json:"email"`
	Scope string `json:"scope"`
}

// UserResponse represents user creation response.
type UserResponse struct {
	ID             string             `json:"_id"`
	DocumentStatus DocumentStatus     `json:"doc_status"`
	Documents      []ResponseDocument `json:"documents"`
	RefreshToken   string             `json:"refresh_token"`
}

// UserResponseClient represents ?..
type UserResponseClient struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Extra represents extra info about the customer.
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
	PhysicalDocs       []SubDocument `json:"physical_docs,omitempty"`
}

// SubDocument represents sub-document object for KYC.
type SubDocument struct {
	DocumentType  string `json:"document_type"`
	DocumentValue string `json:"document_value"`
}

// ResponseDocument represents documents info in the API response.
type ResponseDocument struct {
	PhysicalDocs []ResponseSubDocument `json:"physical_docs"`
}

// ResponseSubDocument represents sub-document in the API response.
type ResponseSubDocument struct {
	DocumentType string `json:"document_type"`
	DocumentID   string `json:"id"`
	Status       string `json:"status"`
}

// DocumentStatus represents document verification result.
type DocumentStatus struct {
	PhysicalDoc string `json:"physical_doc"`
}
