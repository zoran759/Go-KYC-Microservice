package verification

type Request struct {
	Email                string
	PhoneNumber          string
	Country              string
	VerificationServices Services
	VerificationData     Data
}

type Services struct {
	DocumentType       string `json:"document_type,omitempty"`
	DocumentIDNumber   string `json:"document_id_no,omitempty"`
	DocumentExpiryDate string `json:"document_expiry_date,omitempty"`
	Address            string `json:"address,omitempty"`
	FirstName          string `json:"first_name,omitempty"`
	MiddleName         string `json:"middle_mame,omitempty"`
	LastName           string `json:"last_name,omitempty"`
	DateOfBirth        string `json:"dob,omitempty"`
	CardFirst6Digits   string `json:"card_first_6_digits,omitempty"`
	CardLast4Digits    string `json:"card_last_4_digits"`
}

type Data struct {
	FaceImage   string `json:"face_image,omitempty"`
	FrontImage  string `json:"front_image,omitempty"`
	BackImage   string `json:"back_image,omitempty"`
	UtilityBill string `json:"document_address_image,omitempty"`
}

type Response struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Reference  string `json:"reference"`
	Signature  string `json:"signature"`
}
