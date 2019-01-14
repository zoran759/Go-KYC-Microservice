package model

// DateFormat is the date representation form appropriate for the API.
const DateFormat = "2006-01-02"

// NewParticipant represents a new participant without data.
// A new participant creation has to be requested first then send details if succeed.
type NewParticipant struct {
	Email           string          `json:"email"`
	CryptoAddresses []CryptoAddress `json:"addresses,omitempty"`
	UUID            string          `json:"uuid,omitempty"`
}

// NewParticipantResponse represents the response on the request of adding of the new participant.
type NewParticipantResponse struct {
	UUID string `json:"uuid"`
}

// ParticipantDetails represents an individual participant data.
// They have to be send after successful creation of the new participant.
type ParticipantDetails struct {
	UserIP        string   `json:"user_ip,omitempty"`
	Type          string   `json:"type"`
	Pep           bool     `json:"pep,omitempty"`
	FirstName     string   `json:"first_name"`
	LastName      string   `json:"last_name"`
	MiddleName    string   `json:"middle_name,omitempty"`
	Email         string   `json:"email"`
	Nationality   string   `json:"nationality"`
	IDNumber      string   `json:"id_number"`
	Phone         string   `json:"phone"`
	PhoneVerified bool     `json:"phone_verified"`
	Country       string   `json:"country"`
	Postcode      string   `json:"postcode"`
	City          string   `json:"city,omitempty"`
	Street        string   `json:"street,omitempty"`
	BirthDate     string   `json:"bdate"`
	FileFundsText string   `json:"fileFundsText,omitempty"`
	Sow           *Sow     `json:"sow,omitempty"`
	Custom        []string `json:"custom,omitempty"`
}

// Sow represents sow.
type Sow struct {
	BusinessActivities bool `json:"sow_business_activities"`
	StockSales         bool `json:"sow_stock_sales"`
	RealEstateSale     bool `json:"sow_real_estate_sale"`
	Donation           bool `json:"sow_donation"`
	Inherited          bool `json:"sow_inherited"`
	CryptoTrading      bool `json:"sow_crypto_trading"`
	ICOContribution    bool `json:"sow_ico_contribution"`
	Other              bool `json:"sow_other"`
}
