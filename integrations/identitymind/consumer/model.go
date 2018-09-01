package consumer

const maxBillingStreetLength = 100

// KYCRequestData defines the model for Individual KYC Data.
type KYCRequestData struct {
	// Required. Account Name for the user. Maximum length is 60 characters.
	AccountName string `json:"man"`
	// Email address for the user. Maximum length is 60 characters.
	Email string `json:"tea"`
	// OAuth service that authenticated user. For example “google” or “facebook”.
	OAuthService string `json:"soc"`
	// Customer’s IP address.
	IP string `json:"ip"`
	// Device fingerprint blob.
	DeviceFingerprintBlob string `json:"dfp"`
	// Device fingerprint type.
	DeviceFingerprintType DeviceFingerprintType `json:"dft"`
	// Billing First Name.
	BillingFirstName string `json:"bfn"`
	// Billing Middle Name.
	BillingMiddleName string `json:"bmn"`
	// Billing Last Name.
	BillingLastName string `json:"bln"`
	// Billing Street. Include house number, street name and apartment number. Maximum length is 100 characters.
	BillingStreet string `json:"bsn"`
	// Billing Country. ISO 3166 Country code of the Billing Address of the transaction, encoded as a String. Default is “US”.
	BillingCountryAlpha2 string `json:"bco"`
	// Billing Zip / Postal Code.
	BillingPostalCode string `json:"bz"`
	// Billing City.
	BillingCity string `json:"bc"`
	// Billing State.
	BillingState string `json:"bs"`
	// Billing Gender. M, F or Empty.
	BillingGender string `json:"bgd"`
	// Billing Neighborhood.
	BillingNeighborhood string `json:"bnbh"`
	// Shipping First Name.
	ShippingFirstName string `json:"sfn"`
	// Shipping Last Name.
	ShippingLastName string `json:"sln"`
	// Shipping Street. Include house number, street name and apartment number.
	ShippingStreet string `json:"ssn"`
	// Shipping Country. ISO 3166 Country code of the Billing Address of the transaction, encoded as a String. Default is “US”.
	ShippingCountry string `json:"sco"`
	// Shipping Zip / Postal Code.
	ShippingPostalCode string `json:"sz"`
	// Shipping City.
	ShippingCity string `json:"sc"`
	// Shipping State.
	ShippingState string `json:"ss"`
	// Customer longitude.
	CustomerLongitude string `json:"clong"`
	// Customer latitude.
	CustomerLatitude string `json:"clat"`
	// Customer Browser Language.
	CustomerBrowserLanguage string `json:"blg"`
	// Affiliate Id. The client specific identifier for the affiliate that generated this transaction.
	AffiliateID string `json:"aflid"`
	// The signup/affiliate creation date of the affiliate associated with this transaction. Either a ISO8601 encoded string or a unix timestamp.
	AffiliateCreationDate string `json:"aflsd"`
	// Customer primary phone number.
	CustomerPrimaryPhone string `json:"phn"`
	// Customer mobile phone number.
	CustomerMobilePhone string `json:"pm"`
	// Customer work phone number.
	CustomerWorkPhone string `json:"pw"`
	// Transaction Time in UTC. Encoded either as a Unix timestamp, or ISO8601 string. Do not include milliseconds. I chose the string.
	TransactionTime string `json:"tti"`
	// Transaction Identifier. If not provided, an id will be allocated internally by IDM.
	TransactionIdentifier string `json:"tid"`
	// Credit Card unique identifier (Hash). IdentityMind will supply procedure to generate hash. NOTE: The hash must be of the full card number, not a masked or tokenized representation.
	CreditCardUIDHash string `json:"pccn"`
	// A masked or tokenized version of the credit card number. IdentityMind will supply procedure to generate token.
	CreditCardNumberToken string `json:"pcct"`
	// The type of the card.
	CardType CardType `json:"pcty"`
	// Generic payment account unique identifier (Hash). This is used when IdentityMind does not natively support the payment type. NOTE: The hash must be of the full account number, not a masked or tokenized representation.
	GenericPaymentAccountUIDHash string `json:"phash"`
	// A masked or tokenized version of the account token.
	AccountToken string `json:"ptoken"`
	// ACH account unique identifier (Hash). NOTE: The hash must be of the full account number, not a masked or tokenized representation.
	ACHAccountUIDHash string `json:"pach"`
	// A virtual currency address for the funding source. For example the Bitcoin P2PKH address.
	VirtualCurrencyAddress string `json:"pbc"`
	// The policy profile to be used to evaluate this transaction. Prior to IDMRisk 1.18 this was encoded in the smna and smid fields.
	Profile string `json:"profile"`
	// Deprecated.
	Smna string `json:"smna"`
	// Free form memo field for client use.
	Memo string `json:"memo"`
	// Merchant Identifier. Used when a reseller is proxying requests for their merchant’s. Please contact IdentityMind support for further details of the usage of this field.
	MerchantID string `json:"m"`
	// List of Source Digital Currency Addresses.
	SrcDigitalCurrencyAddresses []string `json:"sdcad"`
	// List of Destination Digital Currency Addresses.
	DstDigitalCurrencyAddresses []string `json:"ddcad"`
	// Digital Currency Transaction Hash.
	DigitalCurrencyTransactionHash string `json:"dcth"`
	// An array of tags to be applied to the transaction.
	Tags []string `json:"tags"`
	// Required if using Document Verification, the document front side image data, Base64 encoded. If provided this will override the configured “Jumio client integration”. 5MB maximum size.
	ScanData string `json:"scanData"`
	// If using Document Verification, the document back side image data, Base64 encoded. 5MB maximum size.
	BacksideImageData string `json:"backsideImageData"`
	// If using Document Verification, a serialized JSON array of face image data, Base64 encoded. 5MB maximum size.
	FaceImages []string `json:"faceImages"`
	// Stage of application being processed. An integer between 1 and 5. If not provided, defaults to 1.
	Stage int `json:"stage"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, this parameter should match the exact application ID of the merchant.
	MerchantAid string `json:"merchantAid"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, whether the individual provides a personal guarantee of debt.
	PersonalGuarantee bool `json:"personalguarantee"`
	// If this individual is linked to a merchant (business) as one of the owners of the business, the percentage of ownership.
	Ownership float64 `json:"ownership"`
	// Title of the applicant.
	Title string `json:"title"`
	// Required if using Document Verification, the country in which the document was issued in.
	DocumentCountry string `json:"docCountry"`
	// Required if using Document Verification, the Type of the Document - Passport (PP) | Driver’s Licence (DL) | Government issued Identity Card (ID) |Residence Permit (RP) | Utility Bill (UB).
	DocumentType DocumentType `json:"docType"`
	// Applicant’s date of birth encoded as an ISO8601 string.
	DateOfBirth string `json:"dob"`
	// The applicant’s social security number or national identification number. It is a structed string defined as [ISO-3166-1 (alpha-2)]:[national id].For example “US:123456789” represents a United States Social Security Number. For backwards compatibility if no country code is provided then the identifier is assumed to be a US SSN.
	ApplicantSSN string `json:"assn"`
	// Last 4 digits of the applicant’s social security number or national identification number. If you wish to display the assn4l on the UI, both assn4l and assn values must be present in this request.
	ApplicantSSNLast4 string `json:"assnl4"`
	// Deprecated.
	Smid string `json:"smid"`
	// AVS (Address Verification System) Result value from the Gateway.
	AVSResult string `json:"avs_result"`
}

// TODO plan:
// - write func for ApplicantSSN preparation.
// - review model.
