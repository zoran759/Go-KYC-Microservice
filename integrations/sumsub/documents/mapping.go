package documents

import (
	"gitlab.com/lambospeed/kyc/common"
)

func MapCommonCustomerDocuments(customer common.UserData) []Document {
	if customer.Documents != nil && len(customer.Documents) > 0 {
		documents := make([]Document, 0)

		for _, commonDocument := range customer.Documents {
			if commonDocument.Front == nil {
				continue
			}
			metadata := Metadata{
				DocumentType: MapDocumentType(commonDocument.Metadata.Type),
				Country:      commonDocument.Metadata.Country,
				FirstName:    customer.FirstName,
				MiddleName:   customer.MiddleName,
				LastName:     customer.LastName,
				DateIssued:   commonDocument.Metadata.DateIssued.Format("2006-01-02"),
				ValidUntil:   commonDocument.Metadata.ValidUntil.Format("2006-01-02"),
				Number:       commonDocument.Metadata.Number,
				DateOfBirth:  customer.DateOfBirth.Format("2006-01-02"),
				PlaceOfBirth: customer.PlaceOfBirth,
			}

			if commonDocument.Back != nil {
				metadata.DocumentSubType = FrontSide
				backMetadata := metadata
				backMetadata.DocumentSubType = BackSide

				documents = append(documents, Document{
					Metadata: backMetadata,
					File: File{
						Data:        commonDocument.Back.Data,
						Filename:    commonDocument.Back.Filename,
						ContentType: commonDocument.Back.ContentType,
					},
				})
			}

			documents = append(documents, Document{
				Metadata: metadata,
				File: File{
					Data:        commonDocument.Front.Data,
					Filename:    commonDocument.Front.Filename,
					ContentType: commonDocument.Front.ContentType,
				},
			})

		}

		return documents
	} else {
		return nil
	}
}

func MapDocumentType(documentType common.DocumentType) string {
	switch documentType {
	case common.IDCard:
		return "ID_CARD"
	case common.Passport:
		return "PASSPORT"
	case common.Drivers:
		return "DRIVERS"
	case common.BankCard:
		return "BANK_CARD"
	case common.UtilityBill:
		return "UTILITY_BILL"
	case common.SNILS:
		return "SNILS"
	case common.Selfie:
		return "SELFIE"
	case common.ProfileImage:
		return "PROFILE_IMAGE"
	case common.IDDocPhoto:
		return "ID_DOC_PHOTO"
	case common.Agreement:
		return "AGREEMENT"
	case common.Contract:
		return "CONTRACT"
	case common.ResidencePermit:
		return "RESIDENCE_PERMIT"
	case common.EmploymentCertificate:
		return "EMPLOYMENT_CERTIFICATE"
	case common.DriversTranslation:
		return "DRIVERS_TRANSLATION"
	case common.Other:
		fallthrough
	default:
		return "OTHER"
	}
}
