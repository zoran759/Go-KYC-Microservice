package synapsefi

import (
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"

	"github.com/pkg/errors"
)

var _ common.KYCPlatform = SynapseFI{}

// SynapseFI represents the verification service.
type SynapseFI struct {
	verification verification.Verification
}

// New constructs and returns the new verification service object.
func New(config Config) SynapseFI {
	return SynapseFI{
		verification: verification.NewService(config),
	}
}

// CheckCustomer implements KYCPlatform interface for the SynapseFI.
func (service SynapseFI) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	user := verification.MapCustomerToUser(customer)

	if len(user.Documents[0].VirtualDocs) == 0 {
		err = errors.New("failed to get document's number from customer documents or no document was supplied")
		return
	}

	physDocs := verification.MapCustomerToPhysicalDocs(customer)

	if len(physDocs) == 0 {
		err = errors.New("failed to get document's content or no document was supplied")
		return
	}

	response, code, err := service.verification.CreateUser(user)
	if err != nil {
		if code != nil {
			result.ErrorCode = *code
		}
		return
	}

	if len(response.Documents) == 0 {
		err = errors.New("failed to get documents id for the user " + response.ID)
		return
	}

	code, err = service.verification.AddPhysicalDocs(response.ID, response.RefreshToken, response.Documents[0].ID, physDocs)
	if err != nil {
		if code != nil {
			result.ErrorCode = *code
		}
		return
	}

	result.Status = common.Unclear
	result.StatusCheck = &common.KYCStatusCheck{
		Provider:    common.SynapseFI,
		ReferenceID: response.ID,
		LastCheck:   time.Now(),
	}

	return
}

// CheckStatus implements KYCPlatform interface for the SynapseFI.
func (service SynapseFI) CheckStatus(refID string) (result common.KYCResult, err error) {
	resp, code, err := service.verification.GetUser(refID)
	if err != nil {
		if code != nil {
			result.ErrorCode = *code
		}
		return
	}

	result, err = resp.ToKYCResult()

	return
}
