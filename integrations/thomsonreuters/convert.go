package thomsonreuters

import (
	"modulus/kyc/common"
	"modulus/kyc/integrations/thomsonreuters/model"
)

// newCase constructs a new case for a synchronous screening.
func newCase(template model.CaseTemplateResponse, customer *common.UserData) (newcase model.NewCase) {
	// TODO: implement this.

	return
}

// toResult processes the screening result collection and generates the verification result.
func toResult(toolkits model.ResolutionToolkits, src model.ScreeningResultCollection) (result common.KYCResult, err error) {
	// TODO: implement this.

	return
}
