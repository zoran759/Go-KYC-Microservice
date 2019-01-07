package thomsonreuters

import (
	"encoding/json"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/stretchr/testify/assert"
)

func TestNewCase(t *testing.T) {
	assert := assert.New(t)

	customer := &common.UserData{
		FullName:             "John Martin Doe",
		Gender:               common.Male,
		DateOfBirth:          common.Time(time.Date(1965, 3, 8, 0, 0, 0, 0, time.UTC)),
		CountryOfBirthAlpha2: "GN",
		CountryAlpha2:        "MY",
		Nationality:          "HT",
	}

	template := model.CaseTemplateResponse{}
	err := json.Unmarshal([]byte(caseTemplateResponse), &template)

	assert.NoError(err)

	newcase := newCase(template, customer)

	assert.Equal(template.GroupID, newcase.GroupID)
	assert.Empty(newcase.ID)
	assert.Equal(model.IndividualCET, newcase.EntityType)
	assert.Equal(customer.FullName, newcase.Name)
	assert.Len(newcase.ProviderTypes, 1)
	assert.Equal(model.WatchList, newcase.ProviderTypes[0])
	assert.Empty(newcase.CustomFields)
	assert.Len(newcase.SecondaryFields, 5)
	assert.Equal("SFCT_1", newcase.SecondaryFields[0].TypeID)
	assert.Equal(model.Gender(customer.Gender), newcase.SecondaryFields[0].Value)
	assert.Equal("SFCT_2", newcase.SecondaryFields[1].TypeID)
	assert.Equal("1965-03-08", newcase.SecondaryFields[1].DateTimeValue)
	assert.Equal("SFCT_3", newcase.SecondaryFields[2].TypeID)
	assert.Equal(common.CountryAlpha2ToAlpha3[customer.CountryAlpha2], newcase.SecondaryFields[2].Value)
	assert.Equal("SFCT_4", newcase.SecondaryFields[3].TypeID)
	assert.Equal(common.CountryAlpha2ToAlpha3[customer.CountryOfBirthAlpha2], newcase.SecondaryFields[3].Value)
	assert.Equal("SFCT_5", newcase.SecondaryFields[4].TypeID)
	assert.Equal(common.CountryAlpha2ToAlpha3[customer.Nationality], newcase.SecondaryFields[4].Value)
}

func TestToResultApproved(t *testing.T) {
	assert := assert.New(t)

	src := model.ScreeningResultCollection{
		CaseID: "32737c50-0058-4f28-a0fa-01776aba71e4",
		Results: []model.WatchlistScreeningResult{
			model.WatchlistScreeningResult{
				ResultID:        "0a3687cf-673a-1553-9a06-c6dc00d5378b",
				SubmittedTerm:   "Сергей Васильевич Сарбаш",
				MatchedTerm:     "САРБАШ,Сергей Васильевич",
				MatchedNameType: model.NativeAka,
				MatchStrength:   model.Exact,
				PrimaryName:     "Sergey SARBASH",
				Gender:          model.Male,
				ProviderType:    model.WatchList,
				Category:        "LEGAL",
				SecondaryFieldResults: []model.SecondaryFieldResult{
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "MALE",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "MALE",
						SubmittedValue: "MALE",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							DateTimeValue: "1967-06-12",
						},
						FieldResult:            model.NotMatched,
						MatchedDateTimeValue:   "1967-06-12",
						SubmittedDateTimeValue: "1975-09-21",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						FieldResult:    model.UnknownFR,
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
				},
				CreationDate:     "2019-01-04T23:03:29.980Z",
				ModificationDate: "2019-01-04T23:03:29.980Z",
			},
		},
	}

	res, err := toResult(src)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestToResultDenied(t *testing.T) {
	assert := assert.New(t)

	src := model.ScreeningResultCollection{
		CaseID: "24da33ec-9ad9-463c-9ef7-9e0dce1bfcbb",
		Results: []model.WatchlistScreeningResult{
			model.WatchlistScreeningResult{
				ResultID:        "0a3687d0-673a-15cf-9a06-ae7c00d3929c",
				SubmittedTerm:   "Сергей Владимирович Железняк",
				MatchedTerm:     "Сергей Владимирович Железняк",
				MatchedNameType: model.NativeAka,
				MatchStrength:   model.Exact,
				PrimaryName:     "Sergei Vladimirovich ZHELEZNYAK",
				Gender:          model.Male,
				ProviderType:    model.WatchList,
				Category:        "POLITICAL INDIVIDUAL",
				SecondaryFieldResults: []model.SecondaryFieldResult{
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "MALE",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "MALE",
						SubmittedValue: "MALE",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							DateTimeValue: "1970-07-30",
						},
						FieldResult:            model.Matched,
						MatchedDateTimeValue:   "1970-07-30",
						SubmittedDateTimeValue: "1970-07-30",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						FieldResult:    model.UnknownFR,
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
				},
				CreationDate:     "2019-01-04T21:17:00.013Z",
				ModificationDate: "2019-01-04T21:17:00.013Z",
			},
			model.WatchlistScreeningResult{
				ResultID:        "0a3687d0-673a-15cf-9a06-ae7c00d3923a",
				SubmittedTerm:   "Сергей Владимирович Железняк",
				MatchedTerm:     "Sergey ZHELEZNYAK",
				MatchedNameType: model.Primary,
				MatchStrength:   model.Strong,
				PrimaryName:     "Sergey ZHELEZNYAK",
				Gender:          model.Male,
				ProviderType:    model.WatchList,
				Category:        "CRIME - FINANCIAL",
				SecondaryFieldResults: []model.SecondaryFieldResult{
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "MALE",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "MALE",
						SubmittedValue: "MALE",
					},
					model.SecondaryFieldResult{
						FieldResult:            model.UnknownFR,
						SubmittedDateTimeValue: "1970-07-30",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						FieldResult:    model.UnknownFR,
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						FieldResult:    model.UnknownFR,
						SubmittedValue: "RUS",
					},
					model.SecondaryFieldResult{
						Field: model.Field{
							Value: "RUS",
						},
						FieldResult:    model.Matched,
						MatchedValue:   "RUS",
						SubmittedValue: "RUS",
					},
				},
				CreationDate:     "2019-01-04T21:17:00.013Z",
				ModificationDate: "2019-01-04T21:17:00.013Z",
			},
		},
	}

	res, err := toResult(src)

	assert.NoError(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 3)
	assert.Equal("Case ID: 24da33ec-9ad9-463c-9ef7-9e0dce1bfcbb", res.Details.Reasons[0])
	assert.Equal("Matched Term: Сергей Владимирович Железняк", res.Details.Reasons[1])
	assert.Equal("Category: POLITICAL INDIVIDUAL", res.Details.Reasons[2])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
