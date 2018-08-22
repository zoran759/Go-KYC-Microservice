package applicants

type CreateApplicantRequest struct {
	Email string        `json:"email"`
	Info  ApplicantInfo `json:"info"`
}

type CreateApplicantResponse struct {
	ID           string
	CreatedAt    string
	InspectionID string
	JobID        string
	Email        string
	Info         ApplicantInfo
	Error
}

type Error struct {
	Description *string `json:"description"`
}

type ApplicantInfo struct {
	FirstName      string    `json:"firstName"`
	LastName       string    `json:"lastName"`
	MiddleName     *string   `json:"middleName"`
	LegalName      *string   `json:"legalName"`
	Gender         *string   `json:"gender"`
	DateOfBirth    *string   `json:"dob"`
	PlaceOfBirth   *string   `json:"placeOfBirth"`
	CountryOfBirth *string   `json:"countryOfBirth"`
	StateOfBirth   *string   `json:"stateOfBirth"`
	Country        *string   `json:"country"`
	Nationality    *string   `json:"nationality"`
	Phone          *string   `json:"phone"`
	MobilePhone    *string   `json:"mobilePhone"`
	Addresses      []Address `json:"addresses"`
}

type Address struct {
	Country        *string `json:"country"`
	PostCode       *string `json:"postCode"`
	Town           *string `json:"town"`
	Street         *string `json:"street"`
	SubStreet      *string `json:"subStreet"`
	State          *string `json:"state"`
	BuildingName   *string `json:"buildingName"`
	FlatNumber     *string `json:"flatNumber"`
	BuildingNumber *string `json:"buildingNumber"`
	StartDate      *string `json:"startDate"`
	EndDate        *string `json:"endDate"`
}
