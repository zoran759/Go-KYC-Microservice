package thomsonreuters

import (
	"net/http"
	"testing"

	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var tomson = ThomsonReuters{
	scheme: "https",
	host:   "rms-world-check-one-api-pilot.thomsonreuters.com",
	path:   "/v1/",
	key:    "c7863652-3d05-4f02-8bf7-40ebb70fe17b",
	secret: "KXT8Pkj5n0Ttm4OSfD31x3Au4zf+2QqSbZIXBFoWq1oi7eGWh0k0dkqSdXmSmy15QcWyob7S/ENIdviedBCLRA==",
}

func TestGetRootGroups(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, status, err := tomson.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)
}

func TestGetGroup(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, status, err := tomson.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.HasChildren {
			assert.NotEmpty(g.Children)
			gID = g.Children[0].ID
			break
		}
	}

	assert.NotEmpty(gID)

	group, status, err := tomson.getGroup(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(group)
}

func TestGetCaseTemplate(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, status, err := tomson.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		gID = g.ID
		break
	}

	assert.NotEmpty(gID)

	ctr, status, err := tomson.getCaseTemplate(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(ctr)
}

func TestGetProviders(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	provs, status, err := tomson.getProviders()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(provs)
}

func TestGetResolutionToolkits(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, status, err := tomson.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		gID = g.ID
		break
	}

	assert.NotEmpty(gID)

	rtks, status, err := tomson.getResolutionToolkits(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(rtks)
}

func TestGetActiveUsers(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	users, status, err := tomson.getActiveUsers()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(users)
}

func TestPerformSynchronousScreening(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	groups, status, err := tomson.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		gID = g.ID
		break
	}

	assert.NotEmpty(gID)

	newcase := model.NewCase{
		GroupID: gID,
	}

	src, status, err := tomson.performSynchronousScreening(newcase)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(src)
}
