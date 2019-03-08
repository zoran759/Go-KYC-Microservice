package verification

// OldRequest represents the verification request.
type OldRequest struct {
	Email                string
	PhoneNumber          string
	Country              string
	VerificationServices Services
	VerificationData     Data
}

// Services represents the verification services of the Shufti Pro API.
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

// Data represents the verification data.
type Data struct {
	FaceImage   string `json:"face_image,omitempty"`
	FrontImage  string `json:"document_front_image,omitempty"`
	BackImage   string `json:"document_back_image,omitempty"`
	UtilityBill string `json:"document_address_image,omitempty"`
}

// OldResponse represents the response of the verification API.
type OldResponse struct {
	StatusCode string `json:"status_code"`
	Message    string `json:"message"`
	Reference  string `json:"reference"`
	Signature  string `json:"signature"`
}
