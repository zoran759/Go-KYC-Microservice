package model

// List of FieldValueType values.
const (
	Country FieldValueType = "COUNTRY"
	Gender  FieldValueType = "GENDER"
	Text    FieldValueType = "TEXT"
	Date    FieldValueType = "DATE"
)

// List of GroupScreeningType values.
const (
	CaseManagementAudit GroupScreeningType = "CASE_MANAGEMENT_AUDIT"
	ZeroFootprint       GroupScreeningType = "ZERO_FOOTPRINT"
)

// List of ProviderType values.
const (
	WatchList       ProviderType = "WATCHLIST"
	PassportCheck   ProviderType = "PASSPORT_CHECK"
	ClientWatchList ProviderType = "CLIENT_WATCHLIST"
)

// FieldValueType represents enumerated value formats for a field.
type FieldValueType string

// GroupScreeningType represents the group screening type enumeration.
type GroupScreeningType string

// ProviderType represents the provider type enumeration.
type ProviderType string

// FieldDefinition represents Secondary or Custom Field metadata describing the rules
// for populating the corresponding Field data when creating or updating a Case.
type FieldDefinition struct {
	TypeID         string         `json:"typeId"`
	Label          string         `json:"label"`
	FieldRequired  bool           `json:"fieldRequired"`
	FieldValueType FieldValueType `json:"fieldValueType"`
	RegExp         string         `json:"regExp"`
}

// SecondaryFieldsByEntity represents the wrapper which contains secondary fields grouped by CaseEntityType.
type SecondaryFieldsByEntity struct {
	SecondaryFieldsByEntity map[string][]FieldDefinition `json:"secondaryFieldsByEntity"`
}

// CaseTemplateResponse represents the Case template, containing metadata required by the client to construct a valid Case.
type CaseTemplateResponse struct {
	GroupScreeningType        GroupScreeningType                 `json:"groupScreeningType"`
	GroupID                   string                             `json:"groupId"`
	CustomFields              []FieldDefinition                  `json:"customFields"`
	SecondaryFieldsByProvider map[string]SecondaryFieldsByEntity `json:"secondaryFieldsByProvider"`
	MandatoryProviderTypes    []ProviderType                     `json:"mandatoryProviderTypes"`
}
