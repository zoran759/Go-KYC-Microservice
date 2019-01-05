package thomsonreuters

import (
	"net/http"
	"testing"

	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetGroupID(t *testing.T) {
	assert := assert.New(t)

	tr := New(Config{
		Host:      "https://fakehost/v1/",
		APIkey:    "key",
		APIsecret: "secret",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))

	groups, status, err := tr.getRootGroups()

	assert.NoError(err)
	assert.Nil(status)
	assert.Len(groups, 1)

	group := groups[0]
	assert.Equal("0a3687d0-65b4-1cc3-9975-f20b0000066f", group.ID)
	assert.Equal("CriptoHub S.A. - API (P)", group.Name)
	assert.Empty(group.ParentID)
	assert.True(group.HasChildren)
	assert.Equal(model.ActiveStatus, group.Status)
	assert.Len(group.Children, 1)
}
