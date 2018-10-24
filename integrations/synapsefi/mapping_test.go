package synapsefi

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"
)

func Test_mapResponseToResult(t *testing.T) {
	approvedResponse := &verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: DocStatusValid,
		},
	}

	status, err := mapResponseToResult(approvedResponse)
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, status.Status)
		assert.Nil(t, status.Details)
	}

	deniedResponse := &verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: DocStatusInvalid,
		},
		Documents: []verification.ResponseDocument{
			{
				PhysicalDocs: []verification.ResponseSubDocument{
					{
						DocumentType: "PASSPORT",
						Status:       DocStatusInvalid,
					},
					{
						DocumentType: "SELFIE",
						Status:       DocStatusValid,
					},
					{
						DocumentType: "TYPE",
						Status:       DocStatusInvalid,
					},
				},
			},
		},
	}

	status, err = mapResponseToResult(deniedResponse)
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, status.Status)
		if assert.NotNil(t, status.Details) {
			assert.Equal(t, common.Unknown, status.Details.Finality)
			assert.Equal(t, []string{
				"PASSPORT:" + DocStatusInvalid,
				"TYPE:" + DocStatusInvalid,
			}, status.Details.Reasons)
		}
	}

	missingResponse := &verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: DocStatusMissingOrInvalid,
		},
	}

	status, err = mapResponseToResult(missingResponse)
	if assert.Error(t, err) && assert.Equal(t, common.Error, status.Status) && assert.Nil(t, status.Details) {
		assert.Equal(t, "There are no documents provided, or they are invalid", err.Error())
	}

	unexpectedResponse := &verification.UserResponse{
		ID: "ID",
		DocumentStatus: verification.DocumentStatus{
			PhysicalDoc: "UnexpectedStatus",
		},
	}

	status, err = mapResponseToResult(unexpectedResponse)
	if assert.Error(t, err) && assert.Equal(t, common.Error, status.Status) && assert.Nil(t, status.Details) {
		assert.Equal(t, "Unknown status: UnexpectedStatus", err.Error())
	}
}
