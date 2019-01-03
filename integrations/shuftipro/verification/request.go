package verification

import "modulus/kyc/common"

/*
It is important to note here that each service module is independent of other and each one of them is activated according to the nature of request received from you. There are a total of six services which include face, document, address, consent, phone and background_checks.

All verification services are optional. You can provide Shufti Pro a single service or mixture of several services for verifications. All keys are optional too but values are required with the given keys.
*/

// Valid date format.
const dateFormat = "2006-01-02"

// List of ConsentType values.
const (
	Handwritten ConsentType = "handwritten"
	Printed     ConsentType = "printed"
)

// List of DocumentType values.
const (
	Passport           DocumentType = "passport"
	IDcard             DocumentType = "id_card"
	DrivingLicense     DocumentType = "driving_license"
	CreditOrDebitCard  DocumentType = "credit_or_debit_card"
	UtilityBill        DocumentType = "utility_bill"
	BankStatement      DocumentType = "bank_statement"
	RentAgreement      DocumentType = "rent_agreement"
	EmployerLetter     DocumentType = "employer_letter"
	InsuranceAgreement DocumentType = "insurance_agreement"
	TaxBill            DocumentType = "tax_bill"
)

// Supported document types for document verification.
var docTypes = map[DocumentType]bool{
	Passport:          true,
	IDcard:            true,
	DrivingLicense:    true,
	CreditOrDebitCard: true,
}

// Supported types for address verification.
var addressDocTypes = map[DocumentType]bool{
	IDcard:             true,
	Passport:           true,
	DrivingLicense:     true,
	UtilityBill:        true,
	BankStatement:      true,
	RentAgreement:      true,
	EmployerLetter:     true,
	InsuranceAgreement: true,
	TaxBill:            true,
}

// ConsentType represents a type of consent.
type ConsentType string

// DocumentType represents a supported document type.
type DocumentType string

// Request represents a verification request.
type Request struct {
	Reference        string           `json:"reference"`
	Country          string           `json:"country"`
	Language         string           `json:"language,omitempty"`
	Email            string           `json:"email"`
	CallbackURL      string           `json:"callback_url"`
	VerificationMode string           `json:"verification_mode,omitempty"`
	Face             Face             `json:"face,omitempty"`
	Document         Document         `json:"document,omitempty"`
	Address          Address          `json:"address,omitempty"`
	Consent          Consent          `json:"consent,omitempty"`
	BackgroundChecks BackgroundChecks `json:"background_checks,omitempty"`
}

// Face represents the face of the customer.
type Face struct {
	Proof string `json:"proof"`
}

// Document represents a document.
type Document struct {
	Proof           string         `json:"proof"`
	AdditionalProof string         `json:"additional_proof,omitempty"`
	SupportedTypes  []DocumentType `json:"supported_types"`
	Name            Name           `json:"name,omitempty"`
	DateOfBirth     string         `json:"dob,omitempty"`
	DocumentNumber  string         `json:"document_number,omitempty"`
	IssueDate       string         `json:"issue_date,omitempty"`
	ExpiryDate      string         `json:"expiry_date,omitempty"`
}

// Address represents an address.
type Address struct {
	Proof          string         `json:"proof"`
	SupportedTypes []DocumentType `json:"supported_types"`
	FullAddress    string         `json:"full_address"`
	Name           Name           `json:"name,omitempty"`
}

// Consent represents a consent.
type Consent struct {
	Proof          string        `json:"proof"`
	SupportedTypes []ConsentType `json:"supported_types"`
	Text           string        `json:"text"`
}

// BackgroundChecks represents AML based background checks based on this information.
type BackgroundChecks struct {
	Name        Name   `json:"name"`
	DateOfBirth string `json:"dob"`
}

// Name represents a customer name.
type Name struct {
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
}

// NewRequest constructs new request object from the input data.
func NewRequest(customer *common.UserData) (r *Request, err error) {
	// TODO: implement this.

	return
}
