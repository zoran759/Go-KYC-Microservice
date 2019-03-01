package thomsonreuters

import (
	"errors"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// getGroupID returns group id.
func (tr ThomsonReuters) getGroupID() (groupID string, code *int, err error) {
	/*
	 * It's hard to determine what group we require for verification because
	 * there's no "standard" classification of groups by a usage purpose or anything else.
	 * Currently, we will use the first *active* top or root group.
	 * As I see, the solution for this problem is to introduce other select criteria
	 * related on some attribute's unique value or their combination (group name, id, etc).
	 */
	groups, code, err := tr.getRootGroups()
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
