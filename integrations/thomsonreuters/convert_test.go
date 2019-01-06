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
	// TODO: implement this.
	assert := assert.New(t)

	toolkits := model.ResolutionToolkits{}
	src := model.ScreeningResultCollection{}

	res, err := toResult(toolkits, src)

	assert.NoError(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestToResultDenied(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	toolkits := model.ResolutionToolkits{}
	src := model.ScreeningResultCollection{}

	res, err := toResult(toolkits, src)

	assert.NoError(err)
	assert.Equal(common.Denied, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestToResultError(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	toolkits := model.ResolutionToolkits{}
	src := model.ScreeningResultCollection{}

	res, err := toResult(toolkits, src)

	assert.NoError(err)
	assert.Equal(common.Error, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}
