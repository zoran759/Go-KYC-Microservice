package thomsonreuters

import (
	"errors"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// getGroupID returns group id.
func (tomson ThomsonReuters) getGroupID() (groupID string, code *int, err error) {
	groups, code, err := tomson.getRootGroups()
	if err != nil {
		return
	}

	// Obtain id of the first active root group.
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		groupID = g.ID
		break
	}

	if len(groupID) == 0 {
		err = errors.New("the verification prerequisites error: no active root group")
	}

	return
}
