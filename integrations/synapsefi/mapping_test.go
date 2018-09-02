package synapsefi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"
)

func Test_mapResponseToResult(t *testing.T) {
	approvedResponse := verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: Valid,
		},
	}

	status, details, err := mapResponseToResult(approvedResponse)
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, status)
		assert.Nil(t, details)
	}

	deniedResponse := verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: Invalid,
		},
		Documents: []verification.ResponseDocument{
			{
				PhysicalDocs: []verification.ResponseSubDocument{
					{
						DocumentType: "PASSPORT",
						Status:       Invalid,
					},
					{
						DocumentType: "SELFIE",
						Status:       Valid,
					},
					{
						DocumentType: "TYPE",
						Status:       Invalid,
					},
				},
			},
		},
	}

	status, details, err = mapResponseToResult(deniedResponse)
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, status)
		if assert.NotNil(t, details) {
			assert.Equal(t, common.Unknown, details.Finality)
			assert.Equal(t, []string{
				"PASSPORT:" + Invalid,
				"TYPE:" + Invalid,
			}, details.Reasons)
		}
	}

	missingResponse := verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: MissingOrInvalid,
		},
	}

	status, details, err = mapResponseToResult(missingResponse)
	if assert.Error(t, err) && assert.Equal(t, common.Error, status) && assert.Nil(t, details) {
		assert.Equal(t, "There are no documents provided, or they are invalid", err.Error())
	}

	unexpectedResponse := verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: "UnexpectedStatus",
		},
	}

	status, details, err = mapResponseToResult(unexpectedResponse)
	if assert.Error(t, err) && assert.Equal(t, common.Error, status) && assert.Nil(t, details) {
		assert.Equal(t, "Unknown status: UnexpectedStatus", err.Error())
	}
}
