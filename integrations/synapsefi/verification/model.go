package verification

type CreateUserRequest struct {
	Logins       []Login    `json:"logins"`
	PhoneNumbers []string   `json:"phone_numbers"`
	LegalNames   []string   `json:"legal_names"`
	Documents    []Document `json:"documents"`
	Extra        Extra      `json:"extra"`
}

type Login struct {
	Email string `json:"email"`
	Scope string `json:"scope"`
}

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
	PhysicalDocs       []SubDocument `json:"physical_docs"`
}

type SubDocument struct {
	DocumentType  string `json:"document_type"`
	DocumentValue string `json:"document_value"`
}

type Extra struct {
	CIPTag     int  `json:"cip_tag"`
	IsBusiness bool `json:"is_business"`
}

type UserResponse struct {
	ID             string             `json:"_id"`
	DocumentStatus DocumentStatus     `json:"doc_status"`
	Documents      []ResponseDocument `json:"documents"`
}

type ResponseDocument struct {
	PhysicalDocs []ResponseSubDocument `json:"physical_docs"`
}

type ResponseSubDocument struct {
	DocumentType string `json:"document_type"`
	Status       string `json:"status"`
}

type DocumentStatus struct {
	PhysicalDoc string `json:"physical_doc"`
}
