package verification

import (
	"github.com/gofrs/uuid"
)

// RegistrationRequest represents the verification request.
type RegistrationRequest struct {
	AuthRequest
	NewUserRequest
}

type AuthRequest struct {
	MerchantID	string	`json:"merchant_id"`
	Password 	string	`json:"password"`
}

type NewUserRequest struct {
	CustomerInformation CustomerInformation
	UserNumber          uuid.UUID `json:"user_number,omitempty"`
	RegDate             string `json:"reg_date,omitempty"`
	RegIPAddress        string `json:"reg_ip_address,omitempty"`
}

type CustomerInformation struct {
	FirstName  	CustomerInformationField
	MiddleName 	CustomerInformationField
	LastName   	CustomerInformationField
	Email      	CustomerInformationField
	Address1   	CustomerInformationField
	Address2   	CustomerInformationField
	City       	CustomerInformationField
	Province   	CustomerInformationField
	PostalCode 	CustomerInformationField
	Country    	CustomerInformationField
	Phone1     	CustomerInformationField
	Phone2     	CustomerInformationField
	Dob        	CustomerInformationField
	Gender     	CustomerInformationField
	//IdValues	CustomerInformationField
}

type CustomerInformationField struct {
	FieldName string
	FieldVal string
}

type CustomerInformationDoctypeField struct {
	Type int
	Value string
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
	Status  int    	  	`json:"status"`
	Details string
	ID      string 		`json:"id"`
	Score   int    		`json:"score"`
	Rec     string 		`json:"rec"`
}

// Unified error response
type ResponseError struct {
	Error		map[string]string
	ErrorCode	string
	HttpCode	string
	Status 		bool
}