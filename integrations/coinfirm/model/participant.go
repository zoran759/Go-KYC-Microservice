package model

// DateFormat is the date representation form appropriate for the API.
const DateFormat = "2006-01-02"

// List of participantType type values.
const (
	Corporate  participantType = "corporate"
	Individual participantType = "individual"
)

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

// ParticipantDetails represents a participant data.
// They have to be send after successful creation of the new participant.
type ParticipantDetails struct {
	UserIP        string          `json:"user_ip,omitempty"`
	Type          participantType `json:"type,omitempty"`
	Pep           bool            `json:"pep,omitempty"`
	FirstName     string          `json:"first_name,omitempty"`
	LastName      string          `json:"last_name,omitempty"`
	MiddleName    string          `json:"middle_name,omitempty"`
	Email         string          `json:"email,omitempty"`
	Website       string          `json:"website,omitempty"`
	Nationality   string          `json:"nationality,omitempty"`
	IDNumber      string          `json:"id_number,omitempty"`
	Phone         string          `json:"phone,omitempty"`
	PhoneVerified bool            `json:"phone_verified,omitempty"`
	CountryAlpha3 string          `json:"country,omitempty"`
	Postcode      string          `json:"postcode,omitempty"`
	City          string          `json:"city,omitempty"`
	Street        string          `json:"street,omitempty"`
	BirthDate     string          `json:"bdate,omitempty"`
	FileFundsText string          `json:"fileFundsText,omitempty"`
	Sow           *Sow            `json:"sow,omitempty"`
	Beneficials   []Beneficial    `json:"beneficials,omitempty"`
	Custom        []interface{}   `json:"custom,omitempty"`
}

// participantType represents participant type enumeration.
type participantType string

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

// Beneficial represents a beneficial.
type Beneficial struct {
	Name              string `json:"beneficial_name,omitempty"`
	SowRealEstateSale bool   `json:"sow_real_estate_sale,omitempty"`
	Pep               bool   `json:"beneficial_pep,omitempty"`
	NationAlpha3      string `json:"beneficial_nation,omitempty"`
	Address           string `json:"beneficial_adr,omitempty"`
	BirthDate         string `json:"beneficial_bdate,omitempty"`
	Proc              int    `json:"beneficial_proc,omitempty"`
	WealthText        string `json:"beneficial_wealth_text,omitempty"`
}
