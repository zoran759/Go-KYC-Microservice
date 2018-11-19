package jumio

import (
	"encoding/base64"
	"errors"
	"fmt"

	"modulus/kyc/common"

	"github.com/google/uuid"
)

const (
	performNetverifyEndpoint = "/performNetverify"
	scanStatusEndpoint       = "/scans/"
	scanDetailsEndpoint      = "/scans/%s/data"

	maxImageDataLength = 5 << 20
	dateFormat         = "2006-01-02"

	accept      = "application/json"
	contentType = "application/json"
	userAgent   = "Modulus Exchange/v1.0"
)

var acceptableImageMimeType = map[string]bool{
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
	// FIXME: I didn't figure out how to use it?
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
	PersonalNumber string `json:"personalNumber,omitempty"`
}

// populateFields populate the fields of the request object with input data.
func (r *Request) populateFields(customer *common.UserData) (err error) {
	if err = r.populateDocumentFields(customer); err != nil {
		return
	}

	r.MerchantIDScanReference = uuid.New().String()
	r.FirstName = customer.FirstName
	r.LastName = customer.LastName
	r.DOB = customer.DateOfBirth.Format(dateFormat)

	return
}

// populateDocumentFields processes customer documents and populate the fields relate to a document with input data.
func (r *Request) populateDocumentFields(customer *common.UserData) (err error) {
	if customer.Selfie != nil && customer.Selfie.Image != nil {
		if !acceptableImageMimeType[customer.Selfie.Image.ContentType] {
			err = fmt.Errorf("unacceptable selfie image format: %s", customer.Selfie.Image.ContentType)
			return
		}
		r.FaceImage, err = toBase64(customer.Selfie.Image.Data)
		if err != nil {
			err = fmt.Errorf("during encoding selfi image data: %s", err)
			return
		}
		r.FaceImageMimeType = customer.Selfie.Image.ContentType
	}

	if customer.Passport != nil && customer.Passport.Image != nil {
		if !acceptableImageMimeType[customer.Passport.Image.ContentType] {
			err = fmt.Errorf("unacceptable passport image format: %s", customer.Passport.Image.ContentType)
			return
		}
		r.FrontsideImage, err = toBase64(customer.Passport.Image.Data)
		if err != nil {
			err = fmt.Errorf("during encoding passport image: %s", err)
			return
		}
		r.FrontsideImageMimeType = customer.Passport.Image.ContentType
		r.Country = common.CountryAlpha2ToAlpha3[customer.Passport.CountryAlpha2]
		r.IDType = Passport
		r.Expiry = customer.Passport.ValidUntil.Format(dateFormat)
		r.Number = customer.Passport.Number
		if customer.Passport.CountryAlpha2 == "US" {
			r.USState = customer.Passport.State
		}

		return
	}

	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		if !acceptableImageMimeType[customer.DriverLicense.FrontImage.ContentType] {
			err = fmt.Errorf("unacceptable driver license front image format: %s", customer.DriverLicense.FrontImage.ContentType)
			return
		}
		r.FrontsideImage, err = toBase64(customer.DriverLicense.FrontImage.Data)
		if err != nil {
			err = fmt.Errorf("during encoding driver license front image: %s", err)
			return
		}
		r.FrontsideImageMimeType = customer.DriverLicense.FrontImage.ContentType
		if customer.DriverLicense.BackImage != nil {
			if !acceptableImageMimeType[customer.DriverLicense.BackImage.ContentType] {
				err = fmt.Errorf("unacceptable driver license back image format: %s", customer.DriverLicense.BackImage.ContentType)
				return
			}
			r.BacksideImage, err = toBase64(customer.DriverLicense.BackImage.Data)
			if err != nil {
				err = fmt.Errorf("during encoding driver license back image: %s", err)
				return
			}
			r.BacksideImageMimeType = customer.DriverLicense.BackImage.ContentType
		}
		r.Country = common.CountryAlpha2ToAlpha3[customer.DriverLicense.CountryAlpha2]
		r.IDType = DrivingLicense
		r.Expiry = customer.DriverLicense.ValidUntil.Format(dateFormat)
		r.Number = customer.DriverLicense.Number
		if customer.DriverLicense.CountryAlpha2 == "US" {
			r.USState = customer.DriverLicense.State
		}

		return
	}

	if customer.IDCard != nil && customer.IDCard.Image != nil {
		if !acceptableImageMimeType[customer.IDCard.Image.ContentType] {
			err = fmt.Errorf("unacceptable id card image format: %s", customer.IDCard.Image.ContentType)
			return
		}
		r.FrontsideImage, err = toBase64(customer.IDCard.Image.Data)
		if err != nil {
			err = fmt.Errorf("during encoding id card image: %s", err)
			return
		}
		r.FrontsideImageMimeType = customer.IDCard.Image.ContentType
		r.Country = common.CountryAlpha2ToAlpha3[customer.IDCard.CountryAlpha2]
		r.IDType = IDCard
		r.Number = customer.IDCard.Number

		return
	}

	if customer.SNILS != nil && customer.SNILS.Image != nil {
		if !acceptableImageMimeType[customer.SNILS.Image.ContentType] {
			err = fmt.Errorf("unacceptable SNILS image format: %s", customer.SNILS.Image.ContentType)
			return
		}
		r.FrontsideImage, err = toBase64(customer.SNILS.Image.Data)
		if err != nil {
			err = fmt.Errorf("during encoding SNILS image: %s", err)
			return
		}
		r.FrontsideImageMimeType = customer.SNILS.Image.ContentType
		r.Country = "RUS"
		r.IDType = IDCard
		r.Number = customer.SNILS.Number

		return
	}

	err = errors.New("missing acceptable document for the verification (anyone of passport, driving license or id card)")
	return
}

// toBase64 returns base64-encoded representation of the data.
func toBase64(src []byte) (dst string, err error) {
	if len(src) == 0 {
		return
	}

	if base64.StdEncoding.EncodedLen(len(src)) > maxImageDataLength {
		err = errors.New("too large image file")
		return
	}

	dst = base64.StdEncoding.EncodeToString(src)

	return
}
