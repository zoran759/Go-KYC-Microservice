package verification

import "gitlab.com/modulusglobal/kyc/integrations/trulioo/configuration"

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
	FirstGivenName *string
	MiddleName     *string
	FirstSurName   *string
	SecondSurname  *string
	ISOLatin1Name  *string
	DayOfBirth     *int
	MonthOfBirth   *int
	YearOfBirth    *int
	Gender         *string
}

type Location struct {
	BuildingNumber    *string
	BuildingName      *string
	UnitNumber        *string
	StreetName        *string
	StreetType        *string
	City              *string
	Suburb            *string
	County            *string
	StateProvinceCode *string
	Country           *string
	PostalCode        *string
	POBox             *string
}

type Communication struct {
	MobileNumber *string
	Telephone    *string
	EmailAddress *string
}

type Document struct {
	DocumentFrontImage []byte
	DocumentBackImage  []byte
	LivePhoto          []byte
	DocumentType       string
}

type Business struct {
	BusinessName                *string
	BusinessRegistrationNumber  *string
	DayOfIncorporation          *int
	MonthOfIncorporation        *int
	YearOfIncorporation         *int
	JurisdictionOfIncorporation *string
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
