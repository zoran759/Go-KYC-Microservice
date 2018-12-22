package model

// List of CaseEntityType values.
const (
	IndividualCET   CaseEntityType = "INDIVIDUAL"
	OrganisationCET CaseEntityType = "ORGANISATION"
	VesselCET       CaseEntityType = "VESSEL"
	UnspecifiedCET  CaseEntityType = "UNSPECIFIED"
)

// CaseTemplateResponse represents the Case template, containing metadata required by the client to construct a valid Case.
type CaseTemplateResponse struct {
	GroupID                   string                             `json:"groupId"`
	GroupScreeningType        GroupScreeningType                 `json:"groupScreeningType"`
	MandatoryProviderTypes    []ProviderType                     `json:"mandatoryProviderTypes"`
	CustomFields              []FieldDefinition                  `json:"customFields"`
	SecondaryFieldsByProvider map[string]SecondaryFieldsByEntity `json:"secondaryFieldsByProvider"`
}

// NewCase defines Case data that can be sent when creating a new Case.
type NewCase struct {
	GroupID         string         `json:"groupId"`
	ID              string         `json:"caseId,omitempty"`
	EntityType      CaseEntityType `json:"entityType"`
	Name            string         `json:"name"`
	ProviderTypes   []ProviderType `json:"providerTypes"`
	CustomFields    []Field        `json:"customFields,omitempty"`
	SecondaryFields []Field        `json:"secondaryFields,omitempty"`
}

// CaseEntityType represents the case entity type enumeration.
type CaseEntityType string
