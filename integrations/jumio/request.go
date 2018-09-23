package jumio

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/http"
)

const maxImageDataLength = 5 << 20

var headers = http.Headers{
	"Accept":       "application/json",
	"Content-Type": "application/json",
	"User-Agent":   "Modulus Exchange/v1.0",
}

var acceptableImageMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
}

// Request defines the model for a request body for the performNetverify API request.
type Request struct {
	// Required. Max. length 100.
	MerchantIDScanReference string `json:"merchantIdScanReference"`
	// Required. Max. 5MB & <8000 pixels per side. Base64 encoded image of ID front side.
	FrontsideImage string `json:"frontsideImage"`
	// Max. 5MB & <8000 pixels per side. Base64 encoded image of face. Mandatory if Face match enabled.
	FaceImage string `json:"faceImage,omitempty"`
	// Required. Possible countries: ISO 3166-1 alpha-3 country code; XKX (Kosovo).
	Country string `json:"country"`
	// Required. Type of the document used for the verification.
	IDType IDType `json:"idType"`
	// Mime type of front side image.
	FrontsideImageMimeType string `json:"frontsideImageMimeType,omitempty"`
	// Mime type of face image.
	FaceImageMimeType string `json:"faceImageMimeType,omitempty"`
	// Max. 5MB & <8000 pixels per side. Base64 encoded image of ID back side.
	BacksideImage string `json:"backsideImage,omitempty"`
	// Mime type of back side image.
	BacksideImageMimeType string `json:"backsideImageMimeType,omitempty"`
	// Defines fields which will be extracted during the ID verification. If a field is not listed in this parameter, it will not be processed for this transaction, regardless of customer portal settings. Max. length 100.
	EnabledFields string `json:"enabledFields,omitempty"`
	// Your reporting criteria for each scan. Max. length 100.
	MerchantReportingCriteria string `json:"merchantReportingCriteria,omitempty"`
	// Identification of the customer must not contain sensitive data like PII (Personally Identificable Information) or account login. Max. length 100.
	CustomerID string `json:"customerId,omitempty"`
	// Callback URL for the confirmation after the verification is completed. Max. length 255.
	CallbackURL string `json:"callbackUrl,omitempty"`
	// First name of the customer. Max. length 100.
	FirstName string `json:"firstName,omitempty"`
	// Last name of the customer. Max. length 100.
	LastName string `json:"lastName,omitempty"`
	// Possible values if idType = PASSPORT or ID_CARD: US state code (2 chars) or Alpha-2 country code. If idType = DRIVING_LICENSE: US state code.
	USState string `json:"usState,omitempty"`
	// Date of expiry in the format YYYY-MM-DD.
	Expiry string `json:"expiry,omitempty"`
	// Identification number of the document. Max. length 100.
	Number string `json:"number,omitempty"`
	// Date of birth in the format YYYY-MM-DD.
	DOB string `json:"dob,omitempty"`
	// Possible values:
	// • onFinish (default): Callback is only sent after the whole verification;
	// • onAllSteps: Additional callback is sent when the images are received.
	// Max. length 255.
	CallbackGranularity string `json:"callbackGranularity,omitempty"`
	// FIXME: What's the difference between this and the Identification number?
	// Personal number of the document. Max. length 14.
	PersonalNumber string `json:"personalNumber"`
}

// populateFields populate the fields of the request object with input data.
func (r *Request) populateFields(customer *common.UserData) (err error) {
	// TODO: implement this.

	return
}
