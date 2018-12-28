package thomsonreuters

import (
	"net/http"
	"testing"

	"modulus/kyc/integrations/thomsonreuters/model"

	"github.com/stretchr/testify/assert"
)

// TODO: missing valid credentials.
var s = service{
	scheme: "https",
	host:   "rms-world-check-one-api-pilot.thomsonreuters.com",
	path:   "/v1/",
	key:    "key",
	secret: "secret",
}

func TestNew(t *testing.T) {
	assert := assert.New(t)

	// Test URL parsing error.
	svc := New(Config{
		Host:      "::",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s := svc.(service)

	assert.Empty(s)

	// Test malformed Host.
	svc = New(Config{
		Host:      "host",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s = svc.(service)

	assert.Empty(s)

	// Test valid config.
	svc = New(Config{
		Host:      "https://rms-world-check-one-api-pilot.thomsonreuters.com/v1",
		APIkey:    "key",
		APIsecret: "secret",
	})
	s = svc.(service)

	assert.NotEmpty(s)
	assert.Equal("https", s.scheme)
	assert.Equal("rms-world-check-one-api-pilot.thomsonreuters.com", s.host)
	assert.Equal("/v1/", s.path)
	assert.Equal("key", s.key)
	assert.Equal("secret", s.secret)
}

func TestGetRootGroups(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	groups, status, err := s.getRootGroups()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(groups)
}

func TestGetGroup(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	groups, status, err := s.getRootGroups()

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

	group, status, err := s.getGroup(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(group)
}

func TestGetCaseTemplate(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	groups, status, err := s.getRootGroups()

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

	ctr, status, err := s.getCaseTemplate(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(ctr)
}

func TestGetProviders(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	provs, status, err := s.getProviders()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(provs)
}

func TestGetResolutionToolkits(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	groups, status, err := s.getRootGroups()

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

	rtks, status, err := s.getResolutionToolkits(gID)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(rtks)
}

func TestGetActiveUsers(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	users, status, err := s.getActiveUsers()

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(users)
}

func TestPerformSynchronousScreening(t *testing.T) {
	// TODO: implement this.
	assert := assert.New(t)

	groups, status, err := s.getRootGroups()

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

	src, status, err := s.performSynchronousScreening(newcase)

	assert.NoError(err)
	assert.Equal(http.StatusOK, *status)
	assert.NotEmpty(src)
}
