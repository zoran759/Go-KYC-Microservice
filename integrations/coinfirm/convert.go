package coinfirm

import (
	"encoding/base64"
	"errors"
	"mime"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/coinfirm/model"
)

const (
	unspecified commonDocType = iota
	passport
	idcard
	snils
	utilityBill
	driverLicense
	driverLicenseTranslation
)

type commonDocType int

// prepareCustomerData prepares customer data for KYC process.
func prepareCustomerData(customer *common.UserData) (details model.ParticipantDetails, docfile *model.File) {
	details.UserIP = customer.IPaddress
	details.Type = "individual"
	details.FirstName = customer.FirstName
	details.LastName = customer.LastName
	details.MiddleName = customer.MiddleName
	details.Email = customer.Email
	details.Nationality = common.CountryAlpha2ToAlpha3[customer.Nationality]
	if len(customer.Phone) > 0 {
		details.Phone = customer.Phone
	} else {
		details.Phone = customer.MobilePhone
	}
	details.Country = common.CountryAlpha2ToAlpha3[customer.CountryAlpha2]
	details.Postcode = customer.CurrentAddress.PostCode
	details.City = customer.CurrentAddress.Town
	details.Street = customer.CurrentAddress.Street
	details.BirthDate = customer.DateOfBirth.Format(model.DateFormat)

	details.IDNumber, docfile = prepareCustomerDocument(customer)

	return
}

// prepareCustomerDocument processes provided customer documents and prepares one of them in a particular order.
func prepareCustomerDocument(customer *common.UserData) (docnum string, docfile *model.File) {
	if customer.Passport != nil && customer.Passport.Image != nil {
		filetype := commonDocTypeToFileType(passport)
		ext := extFromContentType(customer.Passport.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.Passport.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
			}
			return
		}
	}
	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		filetype := commonDocTypeToFileType(driverLicense)
		ext := extFromContentType(customer.DriverLicense.FrontImage.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.DriverLicense.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data),
			}
			return
		}
	}
	if customer.DriverLicenseTranslation != nil && customer.DriverLicenseTranslation.FrontImage != nil {
		filetype := commonDocTypeToFileType(driverLicenseTranslation)
		ext := extFromContentType(customer.DriverLicenseTranslation.FrontImage.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.DriverLicenseTranslation.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.DriverLicenseTranslation.FrontImage.Data),
			}
			return
		}
	}
	if customer.IDCard != nil && customer.IDCard.Image != nil {
		filetype := commonDocTypeToFileType(idcard)
		ext := extFromContentType(customer.IDCard.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.IDCard.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			}
			return
		}
	}
	if customer.SNILS != nil && customer.SNILS.Image != nil {
		filetype := commonDocTypeToFileType(snils)
		ext := extFromContentType(customer.SNILS.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.SNILS.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data),
			}
			return
		}
	}
	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		filetype := commonDocTypeToFileType(utilityBill)
		ext := extFromContentType(customer.UtilityBill.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docnum = customer.UtilityBill.Number
			docfile = &model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data),
			}
		}
	}

	return
}

// toResult processes the current status check result and generates the verification result.
func toResult(pID string, status model.StatusResponse) (res common.KYCResult, err error) {
	switch status.CurrentStatus {
	case model.New, model.InProgress:
		res.Status = common.Unclear
		res.StatusCheck = &common.KYCStatusCheck{
			Provider:    common.Coinfirm,
			ReferenceID: pID,
			LastCheck:   time.Now(),
		}
	case model.Incomplete:
		err = errors.New("data provided by participant is incomplete or does not meet the requirements set in KYC form")
	case model.Low:
		res.Status = common.Approved
	case model.Medium, model.High, model.Fail:
		s := string(status.CurrentStatus)
		if status.CurrentStatus == model.Fail {
			s = "unacceptable"
		}
		res.Status = common.Denied
		res.Details = &common.KYCDetails{
			Reasons: []string{
				"Coinfirm analysts evaluated the risk associated to participant as " + s,
			},
		}
	default:
		err = errors.New("unexpected status value: " + string(status.CurrentStatus))
	}

	return
}

// commonDocTypeToFileType is a helper function to decide what API file type kind the document is.
func commonDocTypeToFileType(doctype commonDocType) (filetype model.FileType) {
	if doctype == unspecified {
		return
	}

	switch doctype {
	case passport, idcard, snils, driverLicense, driverLicenseTranslation:
		filetype = model.FileID
	case utilityBill:
		filetype = model.FileAddress
	}

	return
}

// extFromContentType is a helper function that tries to figure out file extension from its mime type.
func extFromContentType(contentType string) (ext string) {
	if len(contentType) == 0 {
		return
	}

	exts, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return
	}

	if len(exts) > 0 {
		ext = exts[0][1:]
	}

	return
}
