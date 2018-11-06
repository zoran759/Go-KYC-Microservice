package verification

// Request represents the verification request.
type Request struct {
	CustomerInformation CustomerInformation
	UserNumber          string `json:"user_number,omitempty"`
	RegDate             string `json:"reg_date,omitempty"`
	RegIPAddress        string `json:"reg_ip_address,omitempty"`
}

// CustomerInformation - main user info
type CustomerInformation struct {
	FirstName  string `json:"first_name,omitempty"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Email      string `json:"email,omitempty"`
	Address1   string `json:"address1,omitempty"`
	Address2   string `json:"address2,omitempty"`
	City       string `json:"city,omitempty"`
	Province   string `json:"province,omitempty"`
	PostalCode string `json:"postal_code,omitempty"`
	Country    string `json:"country,omitempty"`
	Phone1     string `json:"phone1,omitempty"`
	Phone2     string `json:"phone2,omitempty"`
	Dob        string `json:"dob,omitempty"`
	Gender     string `json:"gender,omitempty"`
	UserName   string `json:"user_name,omitempty"`
}

// Data represents the verification data.
type Data struct {
	FaceImage   string `json:"face_image,omitempty"`
	FrontImage  string `json:"document_front_image,omitempty"`
	BackImage   string `json:"document_back_image,omitempty"`
	UtilityBill string `json:"document_address_image,omitempty"`
}

// Response represents the response of the verification API.
type Response struct {
	Status int    `json:"status"`
	ID     string `json:"id"`
	Score  int    `json:"score"`
	Rec    string `json:"rec"`
}
