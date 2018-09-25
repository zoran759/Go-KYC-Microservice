package jumio

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.com/lambospeed/kyc/common"
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

// Request timings recommendation according to https://github.com/Jumio/implementation-guides/blob/master/netverify/netverify-retrieval-api.md#usage.
var timings = [10]time.Duration{
	40 * time.Second,
	60 * time.Second,
	100 * time.Second,
	160 * time.Second,
	240 * time.Second,
	340 * time.Second,
	460 * time.Second,
	600 * time.Second,
	760 * time.Second,
	940 * time.Second,
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
	PersonalNumber string `json:"personalNumber"`
}

// populateFields populate the fields of the request object with input data.
func (r *Request) populateFields(customer *common.UserData) (err error) {
	if err = r.populateDocumentFields(customer.Documents); err != nil {
		return
	}

	r.MerchantIDScanReference = "Modulus"
	r.FirstName = customer.FirstName
	r.LastName = customer.LastName
	r.DOB = customer.DateOfBirth.Format(dateFormat)

	return
}

// populateDocumentFields processes customer documents and populate the fields relate to a document with input data.
func (r *Request) populateDocumentFields(documents []common.Document) (err error) {
	docs := map[common.DocumentType]int{}
	for i, doc := range documents {
		switch doc.Metadata.Type {
		case common.IDCard:
			docs[common.IDCard] = i
		case common.Passport:
			docs[common.Passport] = i
		case common.Drivers:
			docs[common.Drivers] = i
		case common.Selfie:
			if doc.Front != nil {
				if !acceptableImageMimeType[doc.Front.ContentType] {
					err = fmt.Errorf("unacceptable selfie image format: %s", doc.Front.ContentType)
					return
				}
				r.FaceImage, err = toBase64(doc.Front.Data)
				if err != nil {
					err = fmt.Errorf("during encoding selfi image data: %s", err)
					return
				}
				r.FaceImageMimeType = doc.Front.ContentType
			}
		}
	}

	if len(docs) == 0 {
		err = errors.New("missing acceptable document for the verification (anyone of passport, driving license or id card)")
		return
	}

	tryDocument := func(doc *common.Document) (err error) {
		if !acceptableImageMimeType[doc.Front.ContentType] {
			err = fmt.Errorf("unacceptable %s front image format: %s", docTypeToName[doc.Metadata.Type], doc.Front.ContentType)
			return
		}
		r.FrontsideImage, err = toBase64(doc.Front.Data)
		if err != nil {
			err = fmt.Errorf("during encoding %s front image: %s", docTypeToName[doc.Metadata.Type], err)
			return
		}
		r.FrontsideImageMimeType = doc.Front.ContentType
		if doc.Back != nil {
			if !acceptableImageMimeType[doc.Back.ContentType] {
				err = fmt.Errorf("unacceptable %s back image format: %s", docTypeToName[doc.Metadata.Type], doc.Back.ContentType)
				return
			}
			r.BacksideImage, err = toBase64(doc.Back.Data)
			if err != nil {
				err = fmt.Errorf("during encoding %s back image: %s", docTypeToName[doc.Metadata.Type], err)
				return
			}
			r.BacksideImageMimeType = doc.Back.ContentType
		}
		r.Country = common.CountryName2ToAlpha3[strings.ToUpper(doc.Metadata.Country)]
		r.IDType = documentTypeToIDType[doc.Metadata.Type]
		r.Expiry = doc.Metadata.ValidUntil.Format(dateFormat)
		r.Number = doc.Metadata.Number

		return
	}

	i, ok := docs[common.Passport]
	if ok && documents[i].Front != nil {
		err = tryDocument(&documents[i])
		if err == nil {
			r.USState = documents[i].Metadata.Country
		}
		return
	}
	i, ok = docs[common.Drivers]
	if ok && documents[i].Front != nil {
		err = tryDocument(&documents[i])
		if err == nil {
			c := strings.ToUpper(documents[i].Metadata.Country)
			if c == "USA" || c == "UNITED STATES OF AMERICA" {
				r.USState = documents[i].Metadata.StateCode
			}
		}
		return
	}
	i, ok = docs[common.IDCard]
	if ok && documents[i].Front != nil {
		err = tryDocument(&documents[i])
		if err == nil {
			r.USState = documents[i].Metadata.Country
		}
	}

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
