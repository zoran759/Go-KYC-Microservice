package thomsonreuters

import (
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// newCase constructs a new case for the synchronous screening.
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
// We will use exact matching as denied result of verification.
func toResult(src model.ScreeningResultCollection) (result common.KYCResult, err error) {
	for _, r := range src.Results {
		if r.MatchStrength == model.Exact {
			if !matchesExactly(r.SecondaryFieldResults) {
				continue
			}

			result.Status = common.Denied
			reasons := []string{
				"Case ID: " + src.CaseID,
				"Matched Term: " + r.MatchedTerm,
				"Category: " + r.Category,
			}
			result.Details = &common.KYCDetails{
				Reasons: reasons,
			}
			return
		}
	}

	result.Status = common.Approved

	return
}

// matchesExactly checks if all secondary field results match the case.
func matchesExactly(secondaryFieldResults []model.SecondaryFieldResult) bool {
	for _, r := range secondaryFieldResults {
		if r.FieldResult == model.NotMatched {
			return false
		}
	}
	return true
}
