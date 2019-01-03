package thomsonreuters

import (
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/stretchr/testify/assert"
)

func TestNewCase(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	template := model.CaseTemplateResponse{}
	customer := &common.UserData{}

	newcase := newCase(template, customer)

	assert.Equal(template.GroupID, newcase.GroupID)
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
