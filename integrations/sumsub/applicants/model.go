package applicants

type CreateApplicantRequest struct {
	Email string        `json:"email,omitempty"`
	Info  ApplicantInfo `json:"info"`
}

type CreateApplicantResponse struct {
	ID           string `json:"id"`
	CreatedAt    string `json:"createdAt"`
	InspectionID string `json:"inspectionId"`
	JobID        string `json:"jobId"`
	Email        string `json:"email"`
	Info         ApplicantInfo
	Error
}

type Error struct {
	Code        *int    `json:"code"`
	Description *string `json:"description"`
}

type ApplicantInfo struct {
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	MiddleName     string    `json:"middleName,omitempty"`
	LegalName      string    `json:"legalName,omitempty"`
	Gender         string    `json:"gender,omitempty"`
	DateOfBirth    string    `json:"dob,omitempty"`
	PlaceOfBirth   string    `json:"placeOfBirth,omitempty"`
	CountryOfBirth string    `json:"countryOfBirth,omitempty"`
	StateOfBirth   string    `json:"stateOfBirth,omitempty"`
	Country        string    `json:"country,omitempty"`
	Nationality    string    `json:"nationality,omitempty"`
	Phone          string    `json:"phone,omitempty"`
	MobilePhone    string    `json:"mobilePhone,omitempty"`
	Addresses      []Address `json:"addresses,omitempty"`
}

type Address struct {
	Country        string `json:"country,omitempty"`
	PostCode       string `json:"postCode,omitempty"`
	Town           string `json:"town,omitempty"`
	Street         string `json:"street,omitempty"`
	SubStreet      string `json:"subStreet,omitempty"`
	State          string `json:"state,omitempty"`
	BuildingName   string `json:"buildingName,omitempty"`
	FlatNumber     string `json:"flatNumber,omitempty"`
	BuildingNumber string `json:"buildingNumber,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	EndDate        string `json:"endDate,omitempty"`
}
