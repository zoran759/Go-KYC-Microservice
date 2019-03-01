package model

// List of ContactDetailType values.
const (
	Email ContactDetailType = "EMAIL"
	Fax   ContactDetailType = "FAX"
	Voice ContactDetailType = "VOICE"
	URL   ContactDetailType = "URL"
)

// List of CountryLinkType values.
const (
	Nationality  CountryLinkType = "NATIONALITY"
	OperatesIn   CountryLinkType = "OPERATESIN"
	POB          CountryLinkType = "POB"
	POD          CountryLinkType = "POD"
	RegisteredIn CountryLinkType = "REGISTEREDIN"
	Resident     CountryLinkType = "RESIDENT"
	VesselFlag   CountryLinkType = "VESSELFLAG"
	Location     CountryLinkType = "LOCATION"
)

// List of DetailType values.
const (
	Biography      DetailType = "BIOGRAPHY"
	Funding        DetailType = "FUNDING"
	Identification DetailType = "IDENTIFICATION"
	Note           DetailType = "NOTE"
	Reports        DetailType = "REPORTS"
	Regulation     DetailType = "REGULATION"
	SanctionDT     DetailType = "SANCTION"
	UnknownDT      DetailType = "UNKNOWN"
)

// List of EntityUpdateCategory values.
const (
	C1         EntityUpdateCategory = "C1"
	C2         EntityUpdateCategory = "C2"
	C3         EntityUpdateCategory = "C3"
	C4         EntityUpdateCategory = "C4"
	UnknownEUC EntityUpdateCategory = "UNKNOWN"
)

// List of EventType values.
const (
	Birth EventType = "BIRTH"
	Death EventType = "DEATH"
)

// List of FieldResult values.
const (
	Matched    FieldResult = "MATCHED"
	NotMatched FieldResult = "NOT_MATCHED"
	UnknownFR  FieldResult = "UNKNOWN"
)

// List of MatchStrength values.
const (
	Weak   MatchStrength = "WEAK"
	Medium MatchStrength = "MEDIUM"
	Strong MatchStrength = "STRONG"
	Exact  MatchStrength = "EXACT"
)

// List of NameType values.
const (
	Primary       NameType = "PRIMARY"
	Aka           NameType = "AKA"
	AkaEnhanced   NameType = "AKAENHANCED"
	Fka           NameType = "FKA"
	Dba           NameType = "DBA"
	Maiden        NameType = "MAIDEN"
	LangVariation NameType = "LANG_VARIATION"
	Previous      NameType = "PREVIOUS"
	VehicleID     NameType = "VEHICLE_ID"
	LowQualityAka NameType = "LOW_QUALITY_AKA"
	NativeAka     NameType = "NATIVE_AKA"
)

// List of ProfileActionType values.
const (
	CivilAction        ProfileActionType = "CIVIL_ACTION"
	CriminalConviction ProfileActionType = "CRIMINAL_CONVICTION"
	Enforcement        ProfileActionType = "ENFORCEMENT"
	Sanction           ProfileActionType = "SANCTION"
)

// List of ProfileEntityType values.
const (
	CountryPET      ProfileEntityType = "COUNTRY"
	IndividualPET   ProfileEntityType = "INDIVIDUAL"
	OrganisationPET ProfileEntityType = "ORGANISATION"
	VesselPET       ProfileEntityType = "VESSEL"
)

// ScreeningResultCollection contains an array of Watchlist Screening Result objects from a synchronous screening operation.
type ScreeningResultCollection struct {
	CaseID  string                     `json:"caseId"`
	Results []WatchlistScreeningResult `json:"results"`
}

// WatchlistScreeningResult represents the result found after performing a synchronous screening of a Case.
// This contains the abstract case screening result details plus includes identity documents and important events.
type WatchlistScreeningResult struct {
	ResultID              string                 `json:"resultId"`
	ReferenceID           string                 `json:"referenceId"`
	SubmittedTerm         string                 `json:"submittedTerm"`
	MatchedTerm           string                 `json:"matchedTerm"`
	MatchedNameType       NameType               `json:"matchedNameType"`
	MatchStrength         MatchStrength          `json:"matchStrength"`
	PrimaryName           string                 `json:"primaryName"`
	Gender                string                 `json:"gender"`
	ProviderType          ProviderType           `json:"providerType"`
	Category              string                 `json:"category"`
	Categories            []string               `json:"categories"`
	CountryLinks          []CountryLink          `json:"countryLinks"`
	Events                []Event                `json:"events"`
	SecondaryFieldResults []SecondaryFieldResult `json:"secondaryFieldResults"`
	IdentityDocuments     []IdentityDocument     `json:"identityDocuments"`
	Sources               []string               `json:"sources"`
	CreationDate          string                 `json:"creationDate"`
	ModificationDate      string                 `json:"modificationDate"`
}

// CountryLink represents a country link for the Screening Result.
type CountryLink struct {
	Country     Country         `json:"country"`
	CountryText string          `json:"countryText"`
	Type        CountryLinkType `json:"type"`
}

// Country represents a country.
type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// CountryLinkType represents the country link type enumeration.
type CountryLinkType string

// Event represents an event.
type Event struct {
	Type             EventType `json:"type"`
	Address          Address   `json:"address"`
	AllegedAddresses []Address `json:"allegedAddresses"`
	Year             int       `json:"year"`
	Month            int       `json:"month"`
	Day              int       `json:"day"`
	FullDate         string    `json:"fullDate"`
}

// Address represents an address.
type Address struct {
	Country  Country `json:"country"`
	Region   string  `json:"region"`
	City     string  `json:"city"`
	Street   string  `json:"street"`
	PostCode string  `json:"postCode"`
}

// EventType represents the event type enumeration.
type EventType string

// IdentityDocument represents an identity document from the Result.
type IdentityDocument struct {
	Type         string                       `json:"type"`
	Number       string                       `json:"number"`
	Issuer       string                       `json:"issuer"`
	IssueDate    string                       `json:"issueDate"`
	ExpiryDate   string                       `json:"expiryDate"`
	LocationType IdentityDocumentLocationType `json:"locationType"`
	Entity       Entity                       `json:"entity"`
}

// Entity represents a specific entity - usually corresponding to an individual or an organisation
// - that can appear as a potential screening result against a Case.
// Synonyms for this concept include "Profile".
type Entity struct {
	ID                    string               `json:"entityId"`
	ExternalImportID      string               `json:"externalImportId"`
	Active                bool                 `json:"active"`
	Type                  ProfileEntityType    `json:"entityType"`
	Provider              Provider             `json:"provider"`
	Category              string               `json:"category"`
	SubCategory           string               `json:"subCategory"`
	Names                 []Name               `json:"names"`
	Sources               []ProviderSource     `json:"sources"`
	SourceURIs            []string             `json:"sourceUris"`
	Details               []Detail             `json:"details"`
	Actions               []ActionDetail       `json:"actions"`
	Addresses             []Address            `json:"addresses"`
	Associates            []Associate          `json:"associates"`
	Contacts              []ContactDetail      `json:"contacts"`
	CountryLinks          []CountryLink        `json:"countryLinks"`
	PreviousCountryLinks  []CountryLink        `json:"previousCountryLinks"`
	IdentityDocuments     []IdentityDocument   `json:"identityDocuments"`
	Images                []Image              `json:"images"`
	Files                 []File               `json:"files"`
	Weblinks              []Weblink            `json:"weblinks"`
	Description           string               `json:"description"`
	SourceDescription     string               `json:"sourceDescription"`
	CreationDate          string               `json:"creationDate"`
	ModificationDate      string               `json:"modificationDate"`
	LastAdjunctChangeDate string               `json:"lastAdjunctChangeDate"`
	DeletionDate          string               `json:"deletionDate"`
	Comments              string               `json:"comments"`
	UpdateCategory        EntityUpdateCategory `json:"updateCategory"`
	UpdatedDates          EntityUpdatedDates   `json:"updatedDates"`
}

// ActionDetail represents an action detail.
type ActionDetail struct {
	ActionID        string            `json:"actionId"`
	ActionType      ProfileActionType `json:"actionType"`
	Title           string            `json:"title"`
	Text            string            `json:"text"`
	Files           []File            `json:"files"`
	PublicationType string            `json:"publicationType"`
	Published       string            `json:"published"`
	Reference       string            `json:"reference"`
	Source          ProviderSource    `json:"source"`
	StartDate       string            `json:"startDate"`
	EndDate         string            `json:"endDate"`
	Comment         string            `json:"comment"`
}

// ProfileActionType represents the action detail type enumeration.
type ProfileActionType string

// File represents a file.
type File struct {
	Caption string   `json:"caption"`
	URI     string   `json:"uri"`
	Tags    []string `json:"tags"`
}

// ContactDetail represents a contact detail.
type ContactDetail struct {
	ContactDetailType ContactDetailType `json:"contactDetailType"`
	Detail            string            `json:"detail"`
	Country           Country           `json:"country"`
}

// ContactDetailType represents the contact detail type enumeration.
type ContactDetailType string

// Detail represents an entity detail.
type Detail struct {
	DetailType DetailType `json:"detailType"`
	Title      string     `json:"title"`
	Text       string     `json:"text"`
}

// DetailType represents the detail type enumeration.
type DetailType string

// ProfileEntityType represents the profile entity type enumeration.
type ProfileEntityType string

// Image represents an image.
type Image struct {
	Caption      string   `json:"caption"`
	Height       int      `json:"height"`
	Width        int      `json:"width"`
	ImageUseCode string   `json:"imageUseCode"`
	URI          string   `json:"uri"`
	Tags         []string `json:"tags"`
}

// Name represents a name.
type Name struct {
	GivenName      string       `json:"givenName"`
	LastName       string       `json:"lastName"`
	FullName       string       `json:"fullName"`
	Prefix         string       `json:"prefix"`
	Suffix         string       `json:"suffix"`
	Type           NameType     `json:"type"`
	LanguageCode   LanguageCode `json:"languageCode"`
	OriginalScript string       `json:"originalScript"`
}

// LanguageCode represents a language code.
type LanguageCode struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

// NameType represents the watchlist profile name type enumeration.
type NameType string

// EntityUpdateCategory represents the entity update category enumeration.
type EntityUpdateCategory string

// EntityUpdatedDates represents entity updated dates.
type EntityUpdatedDates struct {
	AgeUpdated                 string `json:"ageUpdated"`
	AliasesUpdated             string `json:"aliasesUpdated"`
	AlternativeSpellingUpdated string `json:"alternativeSpellingUpdated"`
	AsOfDateUpdated            string `json:"asOfDateUpdated"`
	CategoryUpdated            string `json:"categoryUpdated"`
	CitizenshipsUpdated        string `json:"citizenshipsUpdated"`
	CompaniesUpdated           string `json:"companiesUpdated"`
	DeceasedUpdated            string `json:"deceasedUpdated"`
	DOBsUpdated                string `json:"dobsUpdated"`
	EIUpdated                  string `json:"eiUpdated"`
	EnteredUpdated             string `json:"enteredUpdated"`
	ExternalSourcesUpdated     string `json:"externalSourcesUpdated"`
	FirstNameUpdated           string `json:"firstNameUpdated"`
	ForeignAliasUpdated        string `json:"foreignAliasUpdated"`
	FurtherInformationUpdated  string `json:"furtherInformationUpdated"`
	IDNumbersUpdated           string `json:"idNumbersUpdated"`
	KeywordsUpdated            string `json:"keywordsUpdated"`
	LastNameUpdated            string `json:"lastNameUpdated"`
	LinkedToUpdated            string `json:"linkedToUpdated"`
	LocationsUpdated           string `json:"locationsUpdated"`
	LowQualityAliasesUpdated   string `json:"lowQualityAliasesUpdated"`
	PassportsUpdated           string `json:"passportsUpdated"`
	PlaceOfBirthUpdated        string `json:"placeOfBirthUpdated"`
	PositionUpdated            string `json:"positionUpdated"`
	SSNUpdated                 string `json:"ssnUpdated"`
	SubCategoryUpdated         string `json:"subCategoryUpdated"`
	TitleUpdated               string `json:"titleUpdated"`
	UpdatecategoryUpdated      string `json:"updatecategoryUpdated"`
}

// Weblink represents a web link.
type Weblink struct {
	Caption string   `json:"caption"`
	URI     string   `json:"uri"`
	Tags    []string `json:"tags"`
}

// IdentityDocumentLocationType represents an identity document location type.
type IdentityDocumentLocationType struct {
	Type    string  `json:"type"`
	Name    string  `json:"name"`
	Country Country `json:"country"`
}

// MatchStrength represents the matched result strength enumeration.
type MatchStrength string

// SecondaryFieldResult belongs to a Result and provides
// the details of the comparison of the Case’s secondary field
// with the corresponding field in the reference data.
type SecondaryFieldResult struct {
	TypeID                 string      `json:"typeId"`
	Field                  Field       `json:"field"`
	FieldResult            FieldResult `json:"fieldResult"`
	MatchedValue           string      `json:"matchedValue"`
	MatchedDateTimeValue   string      `json:"matchedDateTimeValue"`
	SubmittedValue         string      `json:"submittedValue"`
	SubmittedDateTimeValue string      `json:"submittedDateTimeValue"`
}

// FieldResult represents the field result enumeration.
type FieldResult string
