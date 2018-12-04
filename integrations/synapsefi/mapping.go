package synapsefi

import (
	"fmt"

	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"

	"github.com/pkg/errors"
)

func mapResponseToResult(response *verification.UserResponse) (result common.KYCResult, err error) {
	if response.DocumentStatus.PhysicalDoc == DocStatusValid {
		result.Status = common.Approved
		return result, err

	} else if response.DocumentStatus.PhysicalDoc == DocStatusInvalid {
		result.Status = common.Denied
		result.Details = &common.KYCDetails{
			Finality: common.Unknown,
		}

		for _, document := range response.Documents[0].PhysicalDocs {
			if document.Status != DocStatusValid {
				result.Details.Reasons = append(
					result.Details.Reasons,
					fmt.Sprintf("%s:%s",
						document.DocumentType,
						document.Status,
					),
				)
			}
		}

		return result, err

	} else if response.DocumentStatus.PhysicalDoc == DocStatusMissingOrInvalid {
		err = errors.New("There are no documents provided, or they are invalid")
		return result, err

	}

	err = errors.New("Unknown status: " + response.DocumentStatus.PhysicalDoc)
	return result, err
}
