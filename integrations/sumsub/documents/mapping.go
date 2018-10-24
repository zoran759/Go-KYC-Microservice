package documents

import (
	"modulus/kyc/common"
)

// MapCommonCustomerDocuments converts input documents into the format acceptable by the API.
func MapCommonCustomerDocuments(customer common.UserData) (documents []Document) {
	if customer.Passport != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "PASSPORT",
				Country:      common.CountryAlpha2ToAlpha3[customer.Passport.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.Passport.IssuedDate.Format("2006-01-02"),
				ValidUntil:   customer.Passport.ValidUntil.Format("2006-01-02"),
				Number:       customer.Passport.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.Passport.Image != nil {
			doc.File = File{
				Data:        customer.Passport.Image.Data,
				Filename:    customer.Passport.Image.Filename,
				ContentType: customer.Passport.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	}

	if customer.IDCard != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "ID_CARD",
				Country:      common.CountryAlpha2ToAlpha3[customer.IDCard.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.IDCard.IssuedDate.Format("2006-01-02"),
				Number:       customer.IDCard.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.IDCard.Image != nil {
			doc.File = File{
				Data:        customer.IDCard.Image.Data,
				Filename:    customer.IDCard.Image.Filename,
				ContentType: customer.IDCard.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	}

	if customer.SNILS != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "SNILS",
				Country:      "RUS",
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.SNILS.IssuedDate.Format("2006-01-02"),
				Number:       customer.SNILS.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.SNILS.Image != nil {
			doc.File = File{
				Data:        customer.SNILS.Image.Data,
				Filename:    customer.SNILS.Image.Filename,
				ContentType: customer.SNILS.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	}

	if customer.DriverLicense != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "DRIVERS",
				Country:      common.CountryAlpha2ToAlpha3[customer.DriverLicense.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.DriverLicense.IssuedDate.Format("2006-01-02"),
				ValidUntil:   customer.DriverLicense.ValidUntil.Format("2006-01-02"),
				Number:       customer.DriverLicense.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.DriverLicense.FrontImage != nil {
			doc.File = File{
				Data:        customer.DriverLicense.FrontImage.Data,
				Filename:    customer.DriverLicense.FrontImage.Filename,
				ContentType: customer.DriverLicense.FrontImage.ContentType,
			}
		}
		if customer.DriverLicense.BackImage != nil {
			doc.Metadata.DocumentSubType = FrontSide
			backdoc := Document{
				Metadata: doc.Metadata,
				File: File{
					Data:        customer.DriverLicense.BackImage.Data,
					Filename:    customer.DriverLicense.BackImage.Filename,
					ContentType: customer.DriverLicense.BackImage.ContentType,
				},
			}
			backdoc.Metadata.DocumentSubType = BackSide
			documents = append(documents, backdoc)
		}
		documents = append(documents, doc)
	}

	if customer.DriverLicenseTranslation != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "DRIVERS_TRANSLATION",
				Country:      common.CountryAlpha2ToAlpha3[customer.DriverLicenseTranslation.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.DriverLicenseTranslation.IssuedDate.Format("2006-01-02"),
				ValidUntil:   customer.DriverLicenseTranslation.ValidUntil.Format("2006-01-02"),
				Number:       customer.DriverLicenseTranslation.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.DriverLicenseTranslation.FrontImage != nil {
			doc.File = File{
				Data:        customer.DriverLicenseTranslation.FrontImage.Data,
				Filename:    customer.DriverLicenseTranslation.FrontImage.Filename,
				ContentType: customer.DriverLicenseTranslation.FrontImage.ContentType,
			}
		}
		if customer.DriverLicenseTranslation.BackImage != nil {
			doc.Metadata.DocumentSubType = FrontSide
			backdoc := Document{
				Metadata: doc.Metadata,
				File: File{
					Data:        customer.DriverLicenseTranslation.BackImage.Data,
					Filename:    customer.DriverLicenseTranslation.BackImage.Filename,
					ContentType: customer.DriverLicenseTranslation.BackImage.ContentType,
				},
			}
			backdoc.Metadata.DocumentSubType = BackSide
			documents = append(documents, backdoc)
		}
		documents = append(documents, doc)
	}

	if customer.CreditCard != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "BANK_CARD",
				Country:      common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				ValidUntil:   customer.CreditCard.ValidUntil.Format("2006-01-02"),
				Number:       customer.CreditCard.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.CreditCard.Image != nil {
			doc.File = File{
				Data:        customer.CreditCard.Image.Data,
				Filename:    customer.CreditCard.Image.Filename,
				ContentType: customer.CreditCard.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	} else if customer.DebitCard != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "BANK_CARD",
				Country:      common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				ValidUntil:   customer.DebitCard.ValidUntil.Format("2006-01-02"),
				Number:       customer.DebitCard.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.DebitCard.Image != nil {
			doc.File = File{
				Data:        customer.DebitCard.Image.Data,
				Filename:    customer.DebitCard.Image.Filename,
				ContentType: customer.DebitCard.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	}

	if customer.UtilityBill != nil && customer.UtilityBill.Image != nil {
		documents = append(documents, Document{
			Metadata: Metadata{
				DocumentType: "UTILITY_BILL",
				Country:      common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
			File: File{
				Data:        customer.UtilityBill.Image.Data,
				Filename:    customer.UtilityBill.Image.Filename,
				ContentType: customer.UtilityBill.Image.ContentType,
			},
		})
	}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		documents = append(documents, Document{
			Metadata: Metadata{
				DocumentType: "SELFIE",
				Country:      common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
			File: File{
				Data:        customer.Selfie.Image.Data,
				Filename:    customer.Selfie.Image.Filename,
				ContentType: customer.Selfie.Image.ContentType,
			},
		})
	}

	if customer.Avatar != nil && customer.Avatar.Image != nil {
		documents = append(documents, Document{
			Metadata: Metadata{
				DocumentType: "PROFILE_IMAGE",
				Country:      common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
			File: File{
				Data:        customer.Avatar.Image.Data,
				Filename:    customer.Avatar.Image.Filename,
				ContentType: customer.Avatar.Image.ContentType,
			},
		})
	}

	if customer.Other != nil {
		doc := Document{
			Metadata: Metadata{
				DocumentType: "OTHER",
				Country:      common.CountryAlpha2ToAlpha3[customer.Other.CountryAlpha2],
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   customer.Other.IssuedDate.Format("2006-01-02"),
				ValidUntil:   customer.Other.ValidUntil.Format("2006-01-02"),
				Number:       customer.Other.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			},
		}
		if customer.Other.Image != nil {
			doc.File = File{
				Data:        customer.Other.Image.Data,
				Filename:    customer.Other.Image.Filename,
				ContentType: customer.Other.Image.ContentType,
			}
		}
		documents = append(documents, doc)
	}

	if len(documents) == 0 {
		return nil
	}

	return
}
