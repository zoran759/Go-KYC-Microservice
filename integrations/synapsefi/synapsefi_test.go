package synapsefi

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/synapsefi/verification"
)

func TestNew(t *testing.T) {
	service := New(Config{
		TimeoutThreshold: 1000,
	})

	assert.Equal(t, int64(1000), service.timeoutThreshold)
}

func TestSynapseFI_CheckCustomerValid(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: Valid,
					},
				}, nil
			},
		},
	}

	result, details, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.Nil(t, details) {
		assert.Equal(t, common.Approved, result)
	}
}

func TestSynapseFI_CheckCustomerInvalid(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: Invalid,
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status:       Invalid,
								},
							},
						},
					},
				}, nil
			},
		},
	}

	result, details, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, details) {
		assert.Equal(t, common.Denied, result)
		assert.Equal(t, common.Unknown, details.Finality)
		assert.Equal(t, []string{
			"TYPE:" + Invalid,
		}, details.Reasons)
	}
}

func TestSynapseFI_CheckCustomerPoll(t *testing.T) {
	numberOfTImesPolled := 0

	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: "SUBMITTED",
					},
				}, nil
			},
			GetUserFn: func(userID string) (*verification.UserResponse, error) {
				if numberOfTImesPolled == 0 {
					numberOfTImesPolled++
					return &verification.UserResponse{
						DocumentStatus: verification.DocumentStatus{
							PhysicalDoc: "SUBMITTED",
						},
					}, nil
				} else if numberOfTImesPolled == 1 {
					numberOfTImesPolled++
					return nil, errors.New("test_error")
				}

				return &verification.UserResponse{
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: Invalid,
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status:       Invalid,
								},
							},
						},
					},
				}, nil
			},
		},
		timeoutThreshold: 400,
	}

	result, details, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, details) {
		assert.Equal(t, common.Denied, result)
		assert.Equal(t, common.Unknown, details.Finality)
		assert.Equal(t, []string{
			"TYPE:" + Invalid,
		}, details.Reasons)
	}
}

func TestSynapseFI_CheckCustomerPollTimeout(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: "SUBMITTED",
					},
				}, nil
			},
			GetUserFn: func(userID string) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: Invalid,
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status:       Invalid,
								},
							},
						},
					},
				}, nil
			},
		},
	}

	result, details, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result)
	assert.Nil(t, details)
}

func TestSynapseFI_CheckCustomerError(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return nil, errors.New("test_error")
			},
		},
	}

	result, details, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result)
	assert.Nil(t, details)

	result, details, err = service.CheckCustomer(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Error, result)
	assert.Nil(t, details)
}
