package model

import "modulus/kyc/common"

// List of FieldValueType values.
const (
	CountryFVT FieldValueType = "COUNTRY"
	GenderFVT  FieldValueType = "GENDER"
	TextFVT    FieldValueType = "TEXT"
	DateFVT    FieldValueType = "DATE"
)

// List of GENDER field value type values.
const (
	Male    = "MALE"
	Female  = "FEMALE"
	Unknown = "UNKNOWN"
)

// FieldValueType represents enumerated value formats for a field.
type FieldValueType string

// FieldDefinition represents Secondary or Custom Field metadata describing the rules
// for populating the corresponding Field data when creating or updating a Case.
type FieldDefinition struct {
	Label          string         `json:"label"`
	TypeID         string         `json:"typeId"`
	FieldRequired  bool           `json:"fieldRequired"`
	FieldValueType FieldValueType `json:"fieldValueType"`
	RegExp         string         `json:"regExp"`
}

// SecondaryFieldsByEntity represents the wrapper which contains secondary fields grouped by CaseEntityType.
type SecondaryFieldsByEntity struct {
	SecondaryFieldsByEntity map[string][]FieldDefinition `json:"secondaryFieldsByEntity"`
}

// Field holds the value of a Custom Field or a Secondary Field. Valid value type for this field is specified in the corresponding FieldDefinition of CaseTemplateResponse.
type Field struct {
	TypeID        string `json:"typeId"`
	Value         string `json:"value,omitempty"`
	DateTimeValue string `json:"dateTimeValue,omitempty"`
}

// Gender converts customer's gender to the API acceptable value.
func Gender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return Male
	case common.Female:
		return Female
	}
	return Unknown
}
