package thomsonreuters

import (
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// newCase constructs a new case for a synchronous screening.
func newCase(template model.CaseTemplateResponse, customer *common.UserData) (newcase model.NewCase) {
	secfields, ok := template.SecondaryFieldsByProvider["watchlist"]
	if !ok {
		return
	}
	fields := secfields.SecondaryFieldsByEntity["individual"]
	if len(fields) == 0 {
		return
	}

	newcase.GroupID = template.GroupID
	newcase.EntityType = model.IndividualCET
	newcase.ProviderTypes = []model.ProviderType{model.WatchList}
	newcase.Name = customer.Fullname()

	for _, f := range fields {
		switch f.Label {
		case "GENDER":
			newcase.SecondaryFields = append(newcase.SecondaryFields, model.Field{
				TypeID: f.TypeID,
				Value:  model.Gender(customer.Gender),
			})
		case "DATE_OF_BIRTH":
			if time.Time(customer.DateOfBirth).IsZero() {
				continue
			}
			newcase.SecondaryFields = append(newcase.SecondaryFields, model.Field{
				TypeID:        f.TypeID,
				DateTimeValue: customer.DateOfBirth.Format("2006-01-02"),
			})
		case "COUNTRY_LOCATION":
			if len(customer.CountryAlpha2) == 0 {
				continue
			}
			newcase.SecondaryFields = append(newcase.SecondaryFields, model.Field{
				TypeID: f.TypeID,
				Value:  common.CountryAlpha2ToAlpha3[customer.CountryAlpha2],
			})
		case "PLACE_OF_BIRTH":
			if len(customer.CountryOfBirthAlpha2) == 0 {
				continue
			}
			newcase.SecondaryFields = append(newcase.SecondaryFields, model.Field{
				TypeID: f.TypeID,
				Value:  common.CountryAlpha2ToAlpha3[customer.CountryOfBirthAlpha2],
			})
		case "NATIONALITY":
			if len(customer.Nationality) == 0 {
				continue
			}
			newcase.SecondaryFields = append(newcase.SecondaryFields, model.Field{
				TypeID: f.TypeID,
				Value:  common.CountryAlpha2ToAlpha3[customer.Nationality],
			})
		}
	}

	return
}

// toResult processes the screening result collection and generates the verification result.
func toResult(toolkits model.ResolutionToolkits, src model.ScreeningResultCollection) (result common.KYCResult, err error) {
	// TODO: implement this.

	return
}
