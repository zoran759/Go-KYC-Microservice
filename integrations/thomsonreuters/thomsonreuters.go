package thomsonreuters

import (
	"errors"
	"fmt"

	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// CheckCustomer implements CustomerChecker interface for Thomson Reuters.
func (s service) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	gID, code, err := s.getGroupID()
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	template, code, err := s.getCaseTemplate(gID)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	toolkits, code, err := s.getResolutionToolkits(gID)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	newcase := createNewCase(template, customer)

	src, code, err := s.performSynchronousScreening(newcase)
	if err != nil {
		if code != nil {
			result.ErrorCode = fmt.Sprintf("%d", *code)
		}
		return
	}

	result, err = toResult(toolkits, src)

	return
}

// getGroupID returns group id.
func (s service) getGroupID() (groupID string, code *int, err error) {
	groups, code, err := s.getRootGroups()
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

// createNewCase constructs a new case for a synchronous screening.
func createNewCase(template model.CaseTemplateResponse, customer *common.UserData) (newcase model.NewCase) {
	// TODO: implement this.

	return
}

// toResult processes the screening result collection and generates the verification result.
func toResult(toolkits model.ResolutionToolkits, src model.ScreeningResultCollection) (result common.KYCResult, err error) {
	// TODO: implement this.

	return
}
