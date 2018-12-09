package verification

import (
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

var testTime = time.Now().Unix()

func TestToKYCResult(t *testing.T) {
	assert := assert.New(t)

	// Approved result.
	r := &Response{
		ID: "test_id",
		Documents: []ResponseDocument{
			ResponseDocument{
				ID:              "rdid",
				PermissionScope: "SEND|RECEIVE|1000|DAILY",
				VirtualDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "vid",
						Type:        "SSN",
						LastUpdated: testTime,
						Status:      "SUBMITTED|VALID",
					},
				},
				PhysicalDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "phid",
						Type:        "SSN_CARD",
						LastUpdated: testTime,
						Status:      "SUBMITTED|VALID",
					},
				},
			},
		},
		Permission:   "SEND|RECEIVE|1000|DAILY",
		RefreshToken: "rtoken",
	}

	result, err := r.ToKYCResult()

	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}

	// Denied result.
	r = &Response{
		ID: "test_id",
		Documents: []ResponseDocument{
			ResponseDocument{
				ID:              "rdid",
				PermissionScope: "UNVERIFIED",
				VirtualDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "vid",
						Type:        "SSN",
						LastUpdated: testTime,
						Status:      "SUBMITTED|INVALID",
					},
				},
				PhysicalDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "phid",
						Type:        "SSN_CARD",
						LastUpdated: testTime,
						Status:      "SUBMITTED|VALID",
					},
				},
			},
		},
		Permission:   "UNVERIFIED",
		RefreshToken: "rtoken",
	}

	result, err = r.ToKYCResult()

	if assert.NoError(err) {
		assert.Equal(common.Denied, result.Status)
		if assert.NotNil(result.Details) {
			assert.Equal(common.Unknown, result.Details.Finality)
			assert.Len(result.Details.Reasons, 2)
			assert.Equal("Docs set permission: UNVERIFIED", result.Details.Reasons[0])
			assert.Equal("Virtual doc | type: SSN | status: SUBMITTED|INVALID", result.Details.Reasons[1])
		}
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}

	// Denied result.
	r = &Response{
		ID: "test_id",
		Documents: []ResponseDocument{
			ResponseDocument{
				ID:              "rdid",
				PermissionScope: "UNVERIFIED",
				VirtualDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "vid",
						Type:        "SSN",
						LastUpdated: testTime,
						Status:      "SUBMITTED|VALID",
					},
				},
				PhysicalDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "phid",
						Type:        "SSN_CARD",
						LastUpdated: testTime,
						Status:      "SUBMITTED|INVALID",
					},
				},
			},
		},
		Permission:   "UNVERIFIED",
		RefreshToken: "rtoken",
	}

	result, err = r.ToKYCResult()

	if assert.NoError(err) {
		assert.Equal(common.Denied, result.Status)
		if assert.NotNil(result.Details) {
			assert.Equal(common.Unknown, result.Details.Finality)
			assert.Len(result.Details.Reasons, 2)
			assert.Equal("Docs set permission: UNVERIFIED", result.Details.Reasons[0])
			assert.Equal("Physical doc | type: SSN_CARD | status: SUBMITTED|INVALID", result.Details.Reasons[1])
		}
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}

	// Status check required result.
	r = &Response{
		ID: "test_id",
		Documents: []ResponseDocument{
			ResponseDocument{
				ID:              "rdid",
				PermissionScope: "UNVERIFIED",
				VirtualDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "vid",
						Type:        "SSN",
						LastUpdated: testTime,
						Status:      "SUBMITTED",
					},
				},
				PhysicalDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "phid",
						Type:        "SSN_CARD",
						LastUpdated: testTime,
						Status:      "SUBMITTED",
					},
				},
			},
		},
		Permission:   "UNVERIFIED",
		RefreshToken: "rtoken",
	}

	result, err = r.ToKYCResult()

	if assert.NoError(err) {
		assert.Equal(common.Unclear, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)

		if assert.NotNil(result.StatusCheck) {
			assert.Equal(common.SynapseFI, result.StatusCheck.Provider)
			assert.Equal("test_id", result.StatusCheck.ReferenceID)
			assert.NotZero(result.StatusCheck.LastCheck)
		}
	}

	// Error result due to missing docs.
	r = &Response{
		ID:           "test_id",
		Permission:   "UNVERIFIED",
		RefreshToken: "rtoken",
	}

	result, err = r.ToKYCResult()

	if assert.Error(err) {
		assert.EqualError(err, "documents for verification are missing, please, supply one")
		assert.Equal(common.Error, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}

	// Full coverage (more than one doc with valid status in the response).
	r = &Response{
		ID: "test_id",
		Documents: []ResponseDocument{
			ResponseDocument{
				ID:              "rdid",
				PermissionScope: "UNVERIFIED",
				VirtualDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "vid",
						Type:        "SSN",
						LastUpdated: testTime,
						Status:      "SUBMITTED|INVALID",
					},
				},
				PhysicalDocs: []ResponseSubDocument{
					ResponseSubDocument{
						ID:          "phid",
						Type:        "SSN_CARD",
						LastUpdated: testTime,
						Status:      "SUBMITTED|VALID",
					},
				},
			},
			ResponseDocument{
				ID:              "verified",
				PermissionScope: "SEND|RECEIVE|1000|DAILY",
			},
		},
		Permission:   "UNVERIFIED",
		RefreshToken: "rtoken",
	}

	result, err = r.ToKYCResult()

	if assert.NoError(err) {
		assert.Equal(common.Denied, result.Status)
		if assert.NotNil(result.Details) {
			assert.Equal(common.Unknown, result.Details.Finality)
			assert.Len(result.Details.Reasons, 2)
			assert.Equal("Docs set permission: UNVERIFIED", result.Details.Reasons[0])
			assert.Equal("Virtual doc | type: SSN | status: SUBMITTED|INVALID", result.Details.Reasons[1])
		}
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}

}
