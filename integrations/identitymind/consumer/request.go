package consumer

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"modulus/kyc/common"

	"github.com/google/uuid"
)

const (
	maxBillingStreetLength = 100
	maxAccountNameLength   = 60
	maxEmailLength         = 60
	maxImageDataLength     = 5 << 20
)

// KYCRequestData defines the model for Individual KYC Data.
type KYCRequestData struct {
	// Required. Account Name for the user. Maximum length is 60 characters.
	AccountName string `json:"man"`
	// Email address for the user. Maximum length is 60 characters.
	Email string `json:"tea,omitempty"`
	// OAuth service that authenticated user. For example “google” or “facebook”.
	OAuthService string `json:"soc,omitempty"`
	// Customer’s IP address.
	IP string `json:"ip,omitempty"`
	// Device fingerprint blob.
	DeviceFingerprintBlob string `json:"dfp,omitempty"`
	// Device fingerprint type.
	DeviceFingerprintType DeviceFingerprintType `json:"dft,omitempty"`
	// Billing First Name.
	BillingFirstName string `json:"bfn,omitempty"`
	// Billing Middle Name.
	BillingMiddleName string `json:"bmn,omitempty"`
	// Billing Last Name.
	BillingLastName string `json:"bln,omitempty"`
	// Billing Street. Include house number, street name and apartment number. Maximum length is 100 characters.
	BillingStreet string `json:"bsn,omitempty"`
	// Billing Country. ISO 3166 Country code of the Billing Address of the transaction, encoded as a String. Default is “US”.
	BillingCountryAlpha2 string `json:"bco,omitempty"`
	// Billing Zip / Postal Code.
	BillingPostalCode string `json:"bz,omitempty"`
	// Billing City.
	BillingCity string `json:"bc,omitempty"`
	// Billing State.
	BillingState string `json:"bs,omitempty"`
	// Billing Gender. M, F or Empty.
	BillingGender string `json:"bgd"`
	// Billing Neighborhood.
	BillingNeighborhood string `json:"bnbh,omitempty"`
	// Shipping First Name.
	ShippingFirstName string `json:"sfn,omitempty"`
	// Shipping Last Name.
	ShippingLastName string `json:"sln,omitempty"`
	// Shipping Street. Include house number, street name and apartment number.
	ShippingStreet string `json:"ssn,omitempty"`
	// Shipping Country. ISO 3166 Country code of the Billing Address of the transaction, encoded as a String. Default is “US”.
	ShippingCountry string `json:"sco,omitempty"`
	// Shipping Zip / Postal Code.
	ShippingPostalCode string `json:"sz,omitempty"`
	// Shipping City.
	ShippingCity string `json:"sc,omitempty"`
	// Shipping State.
	ShippingState string `json:"ss,omitempty"`
	// Customer longitude.
	CustomerLongitude string `json:"clong,omitempty"`
	// Customer latitude.
	CustomerLatitude string `json:"clat,omitempty"`
	// Customer Browser Language.
	CustomerBrowserLanguage string `json:"blg,omitempty"`
	// Affiliate Id. The client specific identifier for the affiliate that generated this transaction.
	AffiliateID string `json:"aflid,omitempty"`
	// The signup/affiliate creation date of the affiliate associated with this transaction. Either a ISO8601 encoded string or a unix timestamp.
	AffiliateCreationDate string `json:"aflsd,omitempty"`
	// Customer primary phone number.
	CustomerPrimaryPhone string `json:"phn,omitempty"`
	// Customer mobile phone number.
	CustomerMobilePhone string `json:"pm,omitempty"`
	// Customer work phone number.
	CustomerWorkPhone string `json:"pw,omitempty"`
	// Transaction Time in UTC. Encoded either as a Unix timestamp, or ISO8601 string. Do not include milliseconds. I chose the string.
	// FIXME: I can't imagine how to implement this "Object"...
	// TransactionTime string `json:"tti,omitempty"`
	// Transaction Identifier. If not provided, an id will be allocated internally by IDM.
	TxID string `json:"tid,omitempty"`
	// Credit Card unique identifier (Hash). IdentityMind will supply procedure to generate hash. NOTE: The hash must be of the full card number, not a masked or tokenized representation.
	CreditCardUIDHash string `json:"pccn,omitempty"`
	// A masked or tokenized version of the credit card number. IdentityMind will supply procedure to generate token.
	CreditCardNumberToken string `json:"pcct,omitempty"`
	// The type of the card.
	CardType CardType `json:"pcty,omitempty"`
	// Generic payment account unique identifier (Hash). This is used when IdentityMind does not natively support the payment type. NOTE: The hash must be of the full account number, not a masked or tokenized representation.
	GenericPaymentAccountUIDHash string `json:"phash,omitempty"`
	// A masked or tokenized version of the account token.
	AccountToken string `json:"ptoken,omitempty"`
	// ACH account unique identifier (Hash). NOTE: The hash must be of the full account number, not a masked or tokenized representation.
	ACHAccountUIDHash string `json:"pach,omitempty"`
	// A virtual currency address for the funding source. For example the Bitcoin P2PKH address.
	VirtualCurrencyAddress string `json:"pbc,omitempty"`
	// The policy profile to be used to evaluate this transaction.
	// If no profile is specified, the DEFAULT profile is used. Prior to IDMRisk 1.19, this was encoded in the smid field.
	Profile string `json:"profile,omitempty"`
	// Free form memo field for client use.
	Memo string `json:"memo,omitempty"`
	// Merchant Identifier. Used when a reseller is proxying requests for their merchant’s. Please contact IdentityMind support for further details of the usage of this field.
	MerchantID string `json:"m,omitempty"`
	// List of Source Digital Currency Addresses.
	SrcDigitalCurrencyAddresses []string `json:"sdcad,omitempty"`
	// List of Destination Digital Currency Addresses.
	DstDigitalCurrencyAddresses []string `json:"ddcad,omitempty"`
	// Digital Currency Transaction Hash.
	DigitalCurrencyTransactionHash string `json:"dcth,omitempty"`
	// An array of tags to be applied to the transaction.
	Tags []string `json:"tags,omitempty"`
	// Required if using Document Verification, the document front side image data, Base64 encoded. If provided this will override the configured “Jumio client integration”. 5MB maximum size.
	ScanData string `json:"scanData,omitempty"`
	// If using Document Verification, the document back side image data, Base64 encoded. 5MB maximum size.
	BacksideImageData string `json:"backsideImageData,omitempty"`
	// If using Document Verification, a serialized JSON array of face image data, Base64 encoded. 5MB maximum size.
	FaceImages []string `json:"faceImages,omitempty"`
	// Stage of application being processed. An integer between 1 and 5. If not provided, defaults to 1.
	Stage int `json:"stage,omitempty"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, this parameter should match the exact application ID of the merchant.
	MerchantAid string `json:"merchantAid,omitempty"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, whether the individual provides a personal guarantee of debt.
	PersonalGuarantee bool `json:"personalguarantee,omitempty"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, the percentage of ownership.
	Ownership float64 `json:"ownership,omitempty"`
	// Title of the applicant.
	Title string `json:"title,omitempty"`
	// Required if using Document Verification, the country in which the document was issued in.
	DocumentCountry string `json:"docCountry,omitempty"`
	// Required if using Document Verification, the Type of the Document.
	DocumentType DocumentType `json:"docType,omitempty"`
	// Applicant’s date of birth encoded as an ISO8601 string.
	DateOfBirth string `json:"dob,omitempty"`
	// The applicant’s social security number or national identification number. It is a structed string defined as [ISO-3166-1 (alpha-2)]:[national id].For example “US:123456789” represents a United States Social Security Number. For backwards compatibility if no country code is provided then the identifier is assumed to be a US SSN.
	ApplicantSSN string `json:"assn,omitempty"`
	// Last 4 digits of the applicant’s social security number or national identification number. If you wish to display the assn4l on the UI, both assn4l and assn values must be present in this request.
	ApplicantSSNLast4 string `json:"assnl4,omitempty"`
	// AVS (Address Verification System) Result value from the Gateway.
	AVSResult string `json:"avs_result,omitempty"`
}

// setApplicantSSN sets ApplicantSSN field in KYCRequestData in the required format.
func (r *KYCRequestData) setApplicantSSN(ssn *common.IDCard) {
	if len(ssn.CountryAlpha2) == 0 || len(ssn.Number) == 0 {
		return
	}

	r.ApplicantSSN = ssn.CountryAlpha2 + ":" + ssn.Number
}

// populateFields populate the fields of the request object with input data.
func (r *KYCRequestData) populateFields(customer *common.UserData) (err error) {
	if len(customer.AccountName) == 0 {
		customer.AccountName = uuid.New().String()
	}
	if len(customer.AccountName) > maxAccountNameLength {
		err = fmt.Errorf("account length %d exceeded limit of %d symbols", len(customer.AccountName), maxAccountNameLength)
		return
	}
	if len(customer.Email) > maxEmailLength {
		err = fmt.Errorf("email length %d exceeded limit of %d symbols", len(customer.Email), maxEmailLength)
		return
	}
	billingStreet := customer.CurrentAddress.HouseStreetApartment()
	if len(billingStreet) > maxBillingStreetLength {
		err = fmt.Errorf("street address length %d exceeded limit of %d symbols", len(billingStreet), maxBillingStreetLength)
		return
	}

	r.AccountName = customer.AccountName
	r.Email = customer.Email
	r.IP = customer.IPaddress
	r.BillingFirstName = customer.FirstName
	r.BillingMiddleName = customer.MiddleName
	r.BillingLastName = customer.LastName
	r.BillingStreet = billingStreet
	r.BillingCountryAlpha2 = customer.CountryAlpha2
	r.BillingPostalCode = customer.CurrentAddress.PostCode
	r.BillingCity = customer.CurrentAddress.Town
	r.BillingState = customer.CurrentAddress.State
	r.BillingGender = gender2API[customer.Gender]
	if customer.Location != nil {
		r.CustomerLongitude = customer.Location.Longitude
		r.CustomerLatitude = customer.Location.Latitude

	}
	r.CustomerPrimaryPhone = customer.Phone
	r.CustomerMobilePhone = customer.MobilePhone
	r.DateOfBirth = customer.DateOfBirth.Format("2006-01-02")

	// Documents processing.
	err = r.populateDocumentFields(customer)

	return
}

// populateDocumentFields processes customer documents and populate the fields required for Document Verification.
func (r *KYCRequestData) populateDocumentFields(customer *common.UserData) (err error) {
	if customer.IDCard != nil {
		r.setApplicantSSN(customer.IDCard)
		if len(customer.IDCard.Number) > 4 {
			r.ApplicantSSNLast4 = customer.IDCard.Number[len(customer.IDCard.Number)-4:]
		}
	}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		face, e := toBase64(customer.Selfie.Image)
		if e != nil {
			err = fmt.Errorf("during encoding selfi image data: %s", e)
			return
		}
		r.FaceImages = append(r.FaceImages, face)
	}

	if customer.Passport != nil && customer.Passport.Image != nil {
		r.ScanData, err = toBase64(customer.Passport.Image)
		if err != nil {
			err = fmt.Errorf("during encoding passport image: %s", err)
			return
		}
		r.DocumentCountry = customer.Passport.CountryAlpha2
		r.DocumentType = Passport

		return
	}

	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		r.ScanData, err = toBase64(customer.DriverLicense.FrontImage)
		if err != nil {
			err = fmt.Errorf("during encoding driver licence front image: %s", err)
			return
		}
		if customer.DriverLicense.BackImage != nil {
			r.BacksideImageData, err = toBase64(customer.DriverLicense.BackImage)
			if err != nil {
				err = fmt.Errorf("during encoding driver licence back image: %s", err)
				return
			}
		}
		r.DocumentCountry = customer.DriverLicense.CountryAlpha2
		r.DocumentType = DriverLicence

		return
	}

	if customer.IDCard != nil && customer.IDCard.Image != nil {
		r.ScanData, err = toBase64(customer.IDCard.Image)
		if err != nil {
			err = fmt.Errorf("during encoding id card image: %s", err)
			return
		}
		r.DocumentCountry = customer.IDCard.CountryAlpha2
		r.DocumentType = GovernmentIssuedIDCard

		return
	}

	if customer.SNILS != nil && customer.SNILS.Image != nil {
		r.ScanData, err = toBase64(customer.SNILS.Image)
		if err != nil {
			err = fmt.Errorf("during encoding SNILS image: %s", err)
			return
		}
		r.DocumentCountry = "RU"
		r.DocumentType = GovernmentIssuedIDCard

		return
	}

	if customer.ResidencePermit != nil && customer.ResidencePermit.Image != nil {
		r.ScanData, err = toBase64(customer.ResidencePermit.Image)
		if err != nil {
			err = fmt.Errorf("during encoding residence permit image: %s", err)
			return
		}
		r.DocumentCountry = customer.ResidencePermit.CountryAlpha2
		r.DocumentType = ResidencePermit

		return
	}

	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		r.ScanData, err = toBase64(customer.UtilityBill.Image)
		if err != nil {
			err = fmt.Errorf("during encoding utility bill image: %s", err)
			return
		}
		r.DocumentCountry = customer.UtilityBill.CountryAlpha2
		r.DocumentType = UtilityBill
	}

	return
}

// createRequestBody creates request body from the object data.
func (r *KYCRequestData) createRequestBody() (body []byte, err error) {
	body, err = json.Marshal(r)

	return
}

// toBase64 returns base64-encoded representation of the data.
func toBase64(src *common.DocumentFile) (dst string, err error) {
	if len(src.Data) == 0 {
		return
	}

	if base64.StdEncoding.EncodedLen(len(src.Data)) > maxImageDataLength {
		err = errors.New("too large image file")
		return
	}

	dst = src.ContentType + ";base64," + base64.StdEncoding.EncodeToString(src.Data)

	return
}
