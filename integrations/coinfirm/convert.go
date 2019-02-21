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
	driverLicense
	driverLicenseTranslation
	idcard
	snils
	utilityBill
	companyBoard
	companyRegistration
)

type commonDocType int

// prepareCustomerData prepares customer data for KYC process.
func prepareCustomerData(customer *common.UserData) (details model.ParticipantDetails, docfiles []model.File) {
	if customer.IsCompany {
		details, docfiles = prepareCompanyData(customer)
		return
	}
	details, docfiles = prepareIndividualData(customer)
	return
}

// prepareIndividualData prepares an individual for KYC process.
func prepareIndividualData(customer *common.UserData) (details model.ParticipantDetails, docfiles []model.File) {
	details.UserIP = customer.IPaddress
	details.Type = model.Individual
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
	details.CountryAlpha3 = common.CountryAlpha2ToAlpha3[customer.CountryAlpha2]
	details.Postcode = customer.CurrentAddress.PostCode
	details.City = customer.CurrentAddress.Town
	details.Street = customer.CurrentAddress.Street
	details.BirthDate = customer.DateOfBirth.Format(model.DateFormat)

	details.IDNumber, docfiles = prepareIndividualDocuments(customer)

	return
}

// prepareCompanyData prepares a company for KYC process.
func prepareCompanyData(customer *common.UserData) (details model.ParticipantDetails, docfiles []model.File) {
	details.UserIP = customer.IPaddress
	details.Type = model.Corporate
	details.Email = customer.Email
	details.Website = customer.Website
	details.CountryAlpha3 = common.CountryAlpha2ToAlpha3[customer.CountryAlpha2]
	details.Postcode = customer.CurrentAddress.PostCode
	details.City = customer.CurrentAddress.Town
	details.Street = customer.CurrentAddress.Street

	details.IDNumber, docfiles = prepareCompanyDocuments(customer)

	return
}

// prepareIndividualDocuments processes documents provided by an individual and prepares them for a verification.
func prepareIndividualDocuments(customer *common.UserData) (idnum string, docfiles []model.File) {
	if customer.Passport != nil && customer.Passport.Image != nil {
		filetype := commonDocTypeToFileType(passport)
		ext := extFromContentType(customer.Passport.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID {
				idnum = customer.Passport.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
			})
		}
	}
	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		filetype := commonDocTypeToFileType(driverLicense)
		ext := extFromContentType(customer.DriverLicense.FrontImage.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID && len(idnum) == 0 {
				idnum = customer.DriverLicense.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data),
			})
		}
	}
	if customer.DriverLicenseTranslation != nil && customer.DriverLicenseTranslation.FrontImage != nil {
		filetype := commonDocTypeToFileType(driverLicenseTranslation)
		ext := extFromContentType(customer.DriverLicenseTranslation.FrontImage.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID && len(idnum) == 0 {
				idnum = customer.DriverLicenseTranslation.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.DriverLicenseTranslation.FrontImage.Data),
			})
		}
	}
	if customer.IDCard != nil && customer.IDCard.Image != nil {
		filetype := commonDocTypeToFileType(idcard)
		ext := extFromContentType(customer.IDCard.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID && len(idnum) == 0 {
				idnum = customer.IDCard.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data),
			})
		}
	}
	if customer.SNILS != nil && customer.SNILS.Image != nil {
		filetype := commonDocTypeToFileType(snils)
		ext := extFromContentType(customer.SNILS.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID && len(idnum) == 0 {
				idnum = customer.SNILS.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.SNILS.Image.Data),
			})
		}
	}
	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		filetype := commonDocTypeToFileType(utilityBill)
		ext := extFromContentType(customer.UtilityBill.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.UtilityBill.Image.Data),
			})
		}
	}

	return
}

// prepareCompanyDocuments processes documents provided by a company and prepares them for a verification.
func prepareCompanyDocuments(customer *common.UserData) (idnum string, docfiles []model.File) {
	if customer.Passport != nil && customer.Passport.Image != nil {
		filetype := commonDocTypeToFileType(passport)
		ext := extFromContentType(customer.Passport.Image.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID {
				idnum = customer.Passport.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.Passport.Image.Data),
			})
		}
	}
	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		filetype := commonDocTypeToFileType(driverLicense)
		ext := extFromContentType(customer.DriverLicense.FrontImage.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			if filetype == model.FileID && len(idnum) == 0 {
				idnum = customer.DriverLicense.Number
			}
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data),
			})
		}
	}
	if customer.CompanyBoard != nil {
		filetype := commonDocTypeToFileType(companyBoard)
		ext := extFromContentType(customer.CompanyBoard.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.CompanyBoard.Data),
			})
		}
	}
	if customer.CompanyRegistration != nil {
		filetype := commonDocTypeToFileType(companyRegistration)
		ext := extFromContentType(customer.CompanyRegistration.ContentType)
		if len(ext) > 0 && model.IsAcceptedFileExt(ext) && len(filetype) > 0 {
			docfiles = append(docfiles, model.File{
				Type:       filetype,
				Extension:  ext,
				DataBase64: base64.StdEncoding.EncodeToString(customer.CompanyRegistration.Data),
			})
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
	case model.Low, model.Medium:
		res.Status = common.Approved
	case model.High, model.Fail:
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
	case companyBoard:
		filetype = model.FileBoard
	case companyRegistration:
		filetype = model.FileRegister
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
