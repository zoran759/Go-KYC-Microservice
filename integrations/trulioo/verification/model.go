package verification

import "modulus/kyc/integrations/trulioo/configuration"

// StartVerificationRequest represents the model for verification requests.
type StartVerificationRequest struct {
	AcceptTruliooTermsAndConditions bool
	ConfigurationName               string
	ConsentForDataSources           configuration.Consents
	CountryCode                     string
	DataFields                      DataFields
}

// DataFields represents the model for the data fields of the request.
type DataFields struct {
	PersonInfo      *PersonInfo
	Location        *Location
	Communication   *Communication
	Passport        *Passport
	DriverLicence   *DriverLicence
	Document        *Document
	Business        *Business
	NationalIds     []NationalID
	CountrySpecific map[CountryCode]CountrySpecific
}

// CountryCode represents a country ISO 3166-1 alpha-2 code.
type CountryCode = string

// PersonInfo represents the model for the personal info of the customer.
type PersonInfo struct {
	FirstGivenName   string              `json:"FirstGivenName,omitempty"`
	MiddleName       string              `json:"MiddleName,omitempty"`
	FirstSurName     string              `json:"FirstSurName,omitempty"`
	SecondSurname    string              `json:"SecondSurname,omitempty"`
	ISOLatin1Name    string              `json:"ISOLatin1Name,omitempty"`
	MinimumAge       int                 `json:"MinimumAge,omitempty"`
	DayOfBirth       int                 `json:"DayOfBirth,omitempty"`
	MonthOfBirth     int                 `json:"MonthOfBirth,omitempty"`
	YearOfBirth      int                 `json:"YearOfBirth,omitempty"`
	Gender           string              `json:"Gender,omitempty"`
	AdditionalFields *PIAdditionalFields `json:"AdditionalFields,omitempty"`
}

// PIAdditionalFields represents the model for additional fields of the PersonInfo.
type PIAdditionalFields struct {
	FullName string `json:"FullName,omitempty"`
}

// Location represents the model for the address of the customer.
type Location struct {
	BuildingNumber    string            `json:"BuildingNumber,omitempty"`
	BuildingName      string            `json:"BuildingName,omitempty"`
	UnitNumber        string            `json:"UnitNumber,omitempty"`
	StreetName        string            `json:"StreetName,omitempty"`
	StreetType        string            `json:"StreetType,omitempty"`
	City              string            `json:"City,omitempty"`
	Suburb            string            `json:"Suburb,omitempty"`
	County            string            `json:"County,omitempty"`
	StateProvinceCode string            `json:"StateProvinceCode,omitempty"`
	Country           string            `json:"Country,omitempty"`
	PostalCode        string            `json:"PostalCode,omitempty"`
	POBox             string            `json:"POBox,omitempty"`
	AdditionalFields  *AdditionalFields `json:"AdditionalFields,omitempty"`
}

// AdditionalFields represents the model for additional fields of the Location.
type AdditionalFields struct {
	Address1 string `json:"Address1,omitempty"`
}

// Communication represents the model for communication data of the customer.
type Communication struct {
	MobileNumber string `json:"MobileNumber,omitempty"`
	Telephone    string `json:"Telephone,omitempty"`
	Telephone2   string `json:"Telephone2,omitempty"`
	EmailAddress string `json:"EmailAddress,omitempty"`
}

// Passport represents the model for the passport of the customer.
type Passport struct {
	Number        string `json:"Number,omitempty"`
	Mrz1          string `json:"Mrz1,omitempty"`
	Mrz2          string `json:"Mrz2,omitempty"`
	YearOfExpiry  int    `json:"YearOfExpiry,omitempty"`
	MonthOfExpiry int    `json:"MonthOfExpiry,omitempty"`
	DayOfExpiry   int    `json:"DayOfExpiry,omitempty"`
}

// DriverLicence represents the model for the driver licence of the customer.
type DriverLicence struct {
	Number string `json:"Number,omitempty"`
}

// Document represents the model for the given document data of the customer.
type Document struct {
	DocumentFrontImage string
	DocumentBackImage  string
	LivePhoto          string
	DocumentType       string
}

// Business represents the model for a business.
type Business struct {
	BusinessName                string `json:"BusinessName,omitempty"`
	BusinessRegistrationNumber  string `json:"BusinessRegistrationNumber,omitempty"`
	DayOfIncorporation          int    `json:"DayOfIncorporation,omitempty"`
	MonthOfIncorporation        int    `json:"MonthOfIncorporation,omitempty"`
	YearOfIncorporation         int    `json:"YearOfIncorporation,omitempty"`
	JurisdictionOfIncorporation string `json:"JurisdictionOfIncorporation,omitempty"`
}

// NationalID represents the model for a national ID document.
type NationalID struct {
	Number          string `json:"Number,omitempty"`
	Type            string `json:"Type,omitempty"`
	CityOfIssue     string `json:"CityOfIssue,omitempty"`
	ProvinceOfIssue string `json:"ProvinceOfIssue,omitempty"`
}

// CountrySpecific represents the model for country specific data fields.
type CountrySpecific struct {
	PassportNumber           string `json:"PassportNumber,omitempty"`
	PassportMRZLine1         string `json:"PassportMRZLine1,omitempty"`
	PassportMRZLine2         string `json:"PassportMRZLine2,omitempty"`
	PassportYearOfExpiry     string `json:"PassportYearOfExpiry,omitempty"`
	PassportMonthOfExpiry    string `json:"PassportMonthOfExpiry,omitempty"`
	PassportDayOfExpiry      string `json:"PassportDayOfExpiry,omitempty"`
	PassportCountry          string `json:"PassportCountry,omitempty"`
	DriverLicenceNumber      string `json:"DriverLicenceNumber,omitempty"`
	DriverLicenceVerNumber   string `json:"DriverLicenceVersionNumber,omitempty"`
	DriverLicenceState       string `json:"DriverLicenceState,omitempty"`
	HouseExtension           string `json:"HouseExtension,omitempty"`
	CityOfBirth              string `json:"CityOfBirth,omitempty"`
	StateOfBirth             string `json:"StateOfBirth,omitempty"`
	CountryOfBirth           string `json:"CountryOfBirth,omitempty"`
	MaidenName               string `json:"MaidenName,omitempty"`
	BankAccountNumber        string `json:"BankAccountNumber,omitempty"`
	ItIDDocumentType         string `json:"ItIdDocumentType,omitempty"`
	ItIDDocumentNumber       string `json:"ItIdDocumentNumber,omitempty"`
	ItIDDocumentYearOfIssue  string `json:"ItIdDocumentYearOfIssue,omitempty"`
	ItIDDocumentMonthOfIssue string `json:"ItIdDocumentMonthOfIssue,omitempty"`
	ItIDDocumentDayOfIssue   string `json:"ItIdDocumentDayOfIssue,omitempty"`
	FloorNumber              string `json:"FloorNumber,omitempty"`
	PassportSerie            string `json:"PassportSerie,omitempty"`
	InternalPassportNumber   string `json:"InternalPassportNumber,omitempty"`
	YearOfIssue              string `json:"YearOfIssue,omitempty"`
	MonthOfIssue             string `json:"MonthOfIssue,omitempty"`
	DayOfIssue               string `json:"DayOfIssue,omitempty"`
	NameOnCard               string `json:"NameOnCard,omitempty"`
	SerialNumber             string `json:"SerialNumber,omitempty"`
	VehicleRegistrationPlate string `json:"VehicleRegistrationPlate,omitempty"`
	Address2                 string `json:"Address2,omitempty"`
	WorkTelephone            string `json:"WorkTelephone,omitempty"`
	HouseRegistrationNumber  string `json:"HouseRegistrationNumber,omitempty"`
}

// Response represents the model for a verification response.
type Response struct {
	TransactionID string
	UploadedDt    string
	CountryCode   string
	Record        Record
	Errors        Errors
	ErrorCode     *int `json:"-"`
}

// Record represents the model for the verification record.
type Record struct {
	TransactionRecordID string
	RecordStatus        string
	DatasourceResults   []DatasourceResult
	Errors              Errors
	Rule                RecordRule
}

// DatasourceResult represents the model for the datasource result.
type DatasourceResult struct {
	DatasourceStatus string
	DatasourceName   string
	DatasourceFields []DatasourceField
	Errors           Errors
}

// DatasourceField represents the model for the datasource field.
type DatasourceField struct {
	FieldName string
	Status    string
}

// RecordRule represents the model for the record rule.
type RecordRule struct {
	RuleName string
	Note     string
}

// Error represents a verification error.
type Error struct {
	Code    string
	Message string
}

// Errors represents errors returned by the API.
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
