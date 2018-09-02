package synapsefi

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"
)

func mapResponseToResult(response verification.UserResponse) (common.KYCResult, *common.DetailedKYCResult, error) {
	if response.DocumentStatus.PhysicalDoc == Valid {
		return common.Approved, nil, nil
	} else if response.DocumentStatus.PhysicalDoc == Invalid {
		details := common.DetailedKYCResult{
			Finality: common.Unknown,
		}

		for _, document := range response.Documents[0].PhysicalDocs {
			if document.Status != Valid {
				details.Reasons = append(
					details.Reasons,
					fmt.Sprintf("%s:%s",
						document.DocumentType,
						document.Status,
					),
				)
			}
		}

		return common.Denied, &details, nil
	} else if response.DocumentStatus.PhysicalDoc == MissingOrInvalid {
		return common.Error, nil, errors.New("There are no documents provided, or they are invalid")
	}

	return common.Error, nil, errors.New("Unknown status: " + response.DocumentStatus.PhysicalDoc)
}
