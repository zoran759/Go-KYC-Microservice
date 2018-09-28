package verification

import "modulus/kyc/integrations/trulioo/configuration"

type StartVerificationRequest struct {
	AcceptTruliooTermsAndConditions bool
	ConfigurationName               string
	ConsentForDataSources           configuration.Consents
	CountryCode                     string
	DataFields                      DataFields
}

type DataFields struct {
	PersonInfo    *PersonInfo
	Location      *Location
	Communication *Communication
	Document      *Document
	Business      *Business
}

type PersonInfo struct {
	FirstGivenName string `json:"FirstGivenName,omitempty"`
	MiddleName     string `json:"MiddleName,omitempty"`
	FirstSurName   string `json:"FirstSurName,omitempty"`
	SecondSurname  string `json:"SecondSurname,omitempty"`
	ISOLatin1Name  string `json:"ISOLatin1Name,omitempty"`
	DayOfBirth     int    `json:"DayOfBirth,omitempty"`
	MonthOfBirth   int    `json:"MonthOfBirth,omitempty"`
	YearOfBirth    int    `json:"YearOfBirth,omitempty"`
	Gender         string `json:"Gender,omitempty"`
}

type Location struct {
	BuildingNumber    string `json:"BuildingNumber,omitempty"`
	BuildingName      string `json:"BuildingName,omitempty"`
	UnitNumber        string `json:"UnitNumber,omitempty"`
	StreetName        string `json:"StreetName,omitempty"`
	StreetType        string `json:"StreetType,omitempty"`
	City              string `json:"City,omitempty"`
	Suburb            string `json:"Suburb,omitempty"`
	County            string `json:"County,omitempty"`
	StateProvinceCode string `json:"StateProvinceCode,omitempty"`
	Country           string `json:"Country,omitempty"`
	PostalCode        string `json:"PostalCode,omitempty"`
	POBox             string `json:"POBox,omitempty"`
}

type Communication struct {
	MobileNumber string `json:"MobileNumber,omitempty"`
	Telephone    string `json:"Telephone,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
}

type Document struct {
	DocumentFrontImage []byte
	DocumentBackImage  []byte
	LivePhoto          []byte
	DocumentType       string
}

type Business struct {
	BusinessName                string `json:"BusinessName,omitempty"`
	BusinessRegistrationNumber  string `json:"BusinessRegistrationNumber,omitempty"`
	DayOfIncorporation          int    `json:"DayOfIncorporation,omitempty"`
	MonthOfIncorporation        int    `json:"MonthOfIncorporation,omitempty"`
	YearOfIncorporation         int    `json:"YearOfIncorporation,omitempty"`
	JurisdictionOfIncorporation string `json:"JurisdictionOfIncorporation,omitempty"`
}

type VerificationResponse struct {
	TransactionID string
	UploadedDt    string
	CountryCode   string
	Record        Record
	Errors        Errors
}

type Record struct {
	TransactionRecordID string
	RecordStatus        string
	DatasourceResults   []DatasourceResult
	Errors              Errors
	Rule                RecordRule
}

type DatasourceResult struct {
	DatasourceStatus string
	DatasourceName   string
	DatasourceFields []DatasourceField
	Errors           Errors
}

type DatasourceField struct {
	FieldName string
	Status    string
}

type RecordRule struct {
	RuleName string
	Note     string
}

type Error struct {
	Code    string
	Message string
}

type Errors []Error

func (errors Errors) Error() string {
	err := ""
	for _, responseErr := range errors {
		if responseErr.Message != "" {
			err += responseErr.Message + ";"
		}
	}

	return err
}
