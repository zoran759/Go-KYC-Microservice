package shuftipro

import (
	"encoding/base64"
	"fmt"
	"time"

	"modulus/kyc/common"

	"github.com/google/uuid"
)

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

// ConsentType represents a type of consent.
type ConsentType string

// DocumentType represents a supported document type.
type DocumentType string

// Request represents a verification request.
type Request struct {
	Reference        string            `json:"reference"`
	CountryAlpha2    string            `json:"country"`
	Language         string            `json:"language,omitempty"`
	Email            string            `json:"email"`
	CallbackURL      string            `json:"callback_url"`
	VerificationMode string            `json:"verification_mode,omitempty"`
	Face             *Face             `json:"face,omitempty"`
	Document         *Document         `json:"document,omitempty"`
	Address          *Address          `json:"address,omitempty"`
	Consent          *Consent          `json:"consent,omitempty"`
	BackgroundChecks *BackgroundChecks `json:"background_checks,omitempty"`
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
	Name            *Name          `json:"name,omitempty"`
	DateOfBirth     string         `json:"dob,omitempty"`
	Number          string         `json:"document_number,omitempty"`
	IssueDate       string         `json:"issue_date,omitempty"`
	ExpiryDate      string         `json:"expiry_date,omitempty"`
}

// Address represents an address.
type Address struct {
	Proof          string         `json:"proof"`
	SupportedTypes []DocumentType `json:"supported_types"`
	FullAddress    string         `json:"full_address"`
	Name           *Name          `json:"name,omitempty"`
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
func (c Client) NewRequest(customer *common.UserData) (r *Request, err error) {
	r = &Request{}

	id := uuid.New()
	r.Reference = fmt.Sprintf("%x", id[:])
	r.CountryAlpha2 = customer.CountryAlpha2
	r.Email = customer.Email
	r.CallbackURL = c.callbackURL

	if !time.Time(customer.DateOfBirth).IsZero() {
		r.BackgroundChecks = &BackgroundChecks{
			Name: Name{
				FirstName:  customer.FirstName,
				MiddleName: customer.MiddleName,
				LastName:   customer.LastName,
			},
			DateOfBirth: customer.DateOfBirth.Format(dateFormat),
		}
	}

	r.Face = getFace(customer)
	r.Document = getDocument(customer)
	r.Address = getAddress(customer)

	return
}

// getFace returns customer selfie in API suitable format or nil if it fails.
func getFace(customer *common.UserData) *Face {
	if customer.Selfie == nil || customer.Selfie.Image == nil || len(customer.Selfie.Image.Data) == 0 {
		return nil
	}
	return &Face{
		Proof: toBase64(customer.Selfie.Image),
	}
}

// getDocument returns customer document in API suitable format or nil if it fails.
func getDocument(customer *common.UserData) *Document {
	d := &Document{
		Name: &Name{
			FirstName:  customer.FirstName,
			MiddleName: customer.MiddleName,
			LastName:   customer.LastName,
		},
		DateOfBirth: customer.DateOfBirth.Format(dateFormat),
	}
	if customer.Passport != nil && customer.Passport.Image != nil && len(customer.Passport.Image.Data) > 0 {
		d.Proof = toBase64(customer.Passport.Image)
		d.SupportedTypes = []DocumentType{Passport}
		d.Number = customer.Passport.Number
		d.IssueDate = customer.Passport.IssuedDate.Format(dateFormat)
		d.ExpiryDate = customer.Passport.ValidUntil.Format(dateFormat)
		return d
	}
	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil && len(customer.DriverLicense.FrontImage.Data) > 0 {
		d.Proof = toBase64(customer.DriverLicense.FrontImage)
		d.SupportedTypes = []DocumentType{DrivingLicense}
		d.Number = customer.DriverLicense.Number
		d.IssueDate = customer.DriverLicense.IssuedDate.Format(dateFormat)
		d.ExpiryDate = customer.DriverLicense.ValidUntil.Format(dateFormat)
		if customer.DriverLicense.BackImage != nil && len(customer.DriverLicense.BackImage.Data) > 0 {
			d.AdditionalProof = toBase64(customer.DriverLicense.BackImage)
		}
		return d
	}
	if customer.IDCard != nil && customer.IDCard.Image != nil && len(customer.IDCard.Image.Data) > 0 {
		d.Proof = toBase64(customer.IDCard.Image)
		d.SupportedTypes = []DocumentType{IDcard}
		d.Number = customer.IDCard.Number
		d.IssueDate = customer.IDCard.IssuedDate.Format(dateFormat)
		d.ExpiryDate = customer.IDCard.ValidUntil.Format(dateFormat)
		return d
	}
	if customer.CreditCard != nil && customer.CreditCard.Image != nil && len(customer.CreditCard.Image.Data) > 0 {
		d.Proof = toBase64(customer.CreditCard.Image)
		d.SupportedTypes = []DocumentType{CreditOrDebitCard}
		d.Number = customer.CreditCard.Number
		d.ExpiryDate = customer.CreditCard.ValidUntil.Format(dateFormat)
		return d
	}
	if customer.DebitCard != nil && customer.DebitCard.Image != nil && len(customer.DebitCard.Image.Data) > 0 {
		d.Proof = toBase64(customer.DebitCard.Image)
		d.SupportedTypes = []DocumentType{CreditOrDebitCard}
		d.Number = customer.DebitCard.Number
		d.ExpiryDate = customer.DebitCard.ValidUntil.Format(dateFormat)
		return d
	}
	// While we have no document selection functionality this kludge, borrowed from IdentityMind implementation, will serve us.
	if customer.Document != nil && customer.Document.Image != nil && len(customer.Document.Image.Data) > 0 {
		d.Proof = toBase64(customer.Document.Image)
		d.Number = customer.Document.Number
		d.IssueDate = customer.Document.IssuedDate.Format(dateFormat)
		d.ExpiryDate = customer.Document.ValidUntil.Format(dateFormat)
		switch customer.Document.Type {
		case common.IDCardType:
			d.SupportedTypes = []DocumentType{IDcard}
		case common.PassportType:
			d.SupportedTypes = []DocumentType{Passport}
		case common.DriverLicenseType:
			d.SupportedTypes = []DocumentType{DrivingLicense}
		case common.CreditCardType:
			fallthrough
		case common.DebitCardType:
			d.SupportedTypes = []DocumentType{CreditOrDebitCard}
		}
		return d
	}

	return nil
}

// getAddress returns customer address in API suitable format or nil if it fails.
func getAddress(customer *common.UserData) *Address {
	addr := customer.CurrentAddress.String()
	if len(addr) == 0 {
		return nil
	}

	a := &Address{
		FullAddress: addr,
		Name: &Name{
			FirstName:  customer.FirstName,
			MiddleName: customer.MiddleName,
			LastName:   customer.LastName,
		},
	}
	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil && len(customer.UtilityBill.Image.Data) > 0 {
		a.Proof = toBase64(customer.UtilityBill.Image)
		a.SupportedTypes = []DocumentType{UtilityBill}
		return a
	}
	if customer.IDCard != nil && customer.IDCard.Image != nil && len(customer.IDCard.Image.Data) > 0 {
		a.Proof = toBase64(customer.IDCard.Image)
		a.SupportedTypes = []DocumentType{IDcard}
		return a
	}

	return nil
}

// toBase64 is a helper function encoding image data into API suitable format.
func toBase64(img *common.DocumentFile) string {
	return "data:" + img.ContentType + ";base64," + base64.StdEncoding.EncodeToString(img.Data)
}
