package documents

import (
	"gitlab.com/modulusglobal/kyc/common"
	"gitlab.com/modulusglobal/kyc/strings"
)

func MapCommonCustomerDocuments(customer common.UserData) []Document {
	if customer.Documents != nil && len(customer.Documents) > 0 {
		documents := make([]Document, 0)

		for _, commonDocument := range customer.Documents {
			if commonDocument.Front == nil {
				continue
			}
			metadata := Metadata{
				DocumentType: commonDocument.Metadata.Type,
				Country:      commonDocument.Metadata.Country,
				FirstName:    strings.Pointerize(customer.FirstName),
				MiddleName:   strings.Pointerize(customer.MiddleName),
				LastName:     strings.Pointerize(customer.LastName),
				DateIssued:   strings.Pointerize(commonDocument.Metadata.DateIssued.Format("2006-01-02")),
				ValidUntil:   strings.Pointerize(commonDocument.Metadata.ValidUntil.Format("2006-01-02")),
				Number:       strings.Pointerize(commonDocument.Metadata.Number),
				DateOfBirth:  strings.Pointerize(customer.DateOfBirth.Format("2006-01-02")),
				PlaceOfBirth: strings.Pointerize(customer.PlaceOfBirth),
			}

			if commonDocument.Back != nil {
				// can't take address of a const, so use this workaround
				frontSide := FrontSide
				backSide := BackSide

				metadata.DocumentSubType = &frontSide
				backMetadata := metadata
				backMetadata.DocumentSubType = &backSide

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
