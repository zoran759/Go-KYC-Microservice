package complyadvantage

import "modulus/kyc/common"

// Response represents a repsonse from the ComplyAdvantage API.
type Response struct {
	Code    int      `json:"code"`
	Status  string   `json:"status"`
	Content *Content `json:"content"`
}

// Content represents a response content if search succeed.
type Content struct {
	Data Data `json:"data"`
}

// Data represents response payload.
type Data struct {
	ID            int             `json:"id"`
	Ref           string          `json:"ref"`
	SearcherID    int             `json:"searcher_id"`
	AssigneeID    int             `json:"assignee_id"`
	Filters       ResponseFilters `json:"filters"`
	MatchStatus   string          `json:"match_status"`
	RiskLevel     string          `json:"risk_level"`
	SearchTerm    string          `json:"search_term"`
	SubmittedTerm string          `json:"submitted_term"`
	ClientRef     *string         `json:"client_ref"`
	TotalHits     int             `json:"total_hits"`
	UpdatedAt     string          `json:"updated_at"`
	CreatedAt     string          `json:"created_at"`
	Tags          []interface{}   `json:"tags"`
	Limit         int             `json:"limit"`
	Offset        int             `json:"offset"`
	Searcher      Person          `json:"searcher"`
	Assignee      Person          `json:"assignee"`
	Hits          []Hit           `json:"hits"`
}

// ResponseFilters represents filters used for the search.
type ResponseFilters struct {
	BirthYear      int      `json:"birth_year"`
	CountryCodes   []string `json:"country_codes"`
	RemoveDeceased int      `json:"remove_deceased"`
	Types          []string `json:"types"`
	ExactMatch     bool     `json:"exact_match"`
	Fuzziness      int      `json:"fuzziness"`
}

// Person represents a person related to the search.
type Person struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
}

// Hit represents an entity object matching the search term.
type Hit struct {
	Doc               Doc                          `json:"doc"`
	MatchTypes        []string                     `json:"match_types"`
	MatchTypesDetails map[string]MatchTypesDetails `json:"match_types_details"`
	Score             float32                      `json:"score"`
}

// Doc represents a matched doc in the search result.
type Doc struct {
	ID             string                `json:"id"`
	EntityType     string                `json:"entity_type"`
	Name           string                `json:"name"`
	Aka            []Aka                 `json:"aka"`
	Associates     []Associate           `json:"associates"`
	Assets         []Asset               `json:"assets"`
	Fields         []Field               `json:"fields"`
	LastUpdatedUTC string                `json:"last_updated_utc"`
	Media          []Media               `json:"media"`
	Sources        []string              `json:"sources"`
	SourceNotes    map[string]SourceNote `json:"source_notes"`
	Types          []string              `json:"types"`
	Keywords       []interface{}         `json:"keywords"`
}

// Aka represents the search term "aka".
type Aka struct {
	Name string `json:"name"`
}

// Associate represents an associated entity.
type Associate struct {
	Association string `json:"association"`
	Name        string `json:"name"`
	Comment     string `json:"comment"`
}

// Asset represents an asset item of the doc.
type Asset struct {
	Type      string `json:"type"`
	Source    string `json:"source"`
	PublicURL string `json:"public_url"`
}

// Field represents a doc data field.
type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Tag    string `json:"tag"`
	Source string `json:"source"`
}

// Media represents a media associated with the doc.
type Media struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	PdfURL  string `json:"pdf_url"`
	Date    string `json:"date"`
	Snippet string `json:"snippet"`
}

// SourceNote represents a note from the data source (eg. a sanction list).
type SourceNote struct {
	Name              string   `json:"name"`
	URL               string   `json:"url"`
	ListingStartedUTC string   `json:"listing_started_utc"`
	ListingEndedUTC   string   `json:"listing_ended_utc"`
	AMLTypes          []string `json:"aml_types"`
	CountryCodes      []string `json:"country_codes"`
}

// MatchTypesDetails represents match types details for the search term.
type MatchTypesDetails struct {
	Type       string              `json:"type"`
	MatchTypes map[string][]string `json:"match_types"`
}

// toResult processes the response and generates the verification result.
func (r Response) toResult() (result common.KYCResult, err error) {
	// TODO: implement this.

	return
}
