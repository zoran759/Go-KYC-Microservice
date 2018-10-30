package synapsefi

import (
	"testing"
	"errors"
	"github.com/stretchr/testify/assert"
	"modulus/kyc/integrations/synapsefi/verification"
	"modulus/kyc/common"
)

func TestNewShort(t *testing.T) {
	service := New(Config{
	})

	assert.Equal(t, int64(3600), service.timeoutThreshold)
	assert.Equal(t, "simple", service.kycFlow)
}

func TestNewFull(t *testing.T) {
	service := New(Config{
		TimeoutThreshold: 1000,
		KYCFlow: "complex",
	})

	assert.Equal(t, int64(1000), service.timeoutThreshold)
	assert.Equal(t, "complex", service.kycFlow)
}

func TestSynapseFI_CheckCustomerValid(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusValid,
					},
				}, nil
			},
		},
		kycFlow: "simple",
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, result.Status)
	}
}

func TestSynapseFI_CheckCustomerInvalid(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusInvalid,
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status: DocStatusInvalid,
								},
							},
						},
					},
				}, nil
			},
		},
		kycFlow: "simple",
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Equal(t, common.Unknown, result.Details.Finality)
		assert.Equal(t, []string{
			"TYPE:" + DocStatusInvalid,
		}, result.Details.Reasons)
	}
}

func TestSynapseFI_CheckCustomerPoll(t *testing.T) {

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
						PhysicalDoc: "SUBMITTED|INVALID",
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status: "SUBMITTED|INVALID",
								},
							},
						},
					},
				}, nil
			},
		},
		timeoutThreshold: 400,
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result.Details) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Equal(t, common.Unknown, result.Details.Finality)
		assert.Equal(t, []string{
			"TYPE:" + DocStatusInvalid,
		}, result.Details.Reasons)
	}
}

//func TestSynapseFI_CheckCustomerWithoutDocuments(t *testing.T) {
//	service := SynapseFI{
//		verification: verification.Mock{
//			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
//				return &verification.UserResponse{
//					ID: "someid",
//				}, nil
//			},
//		},
//	}
//
//	result, err := service.CheckCustomer(&common.UserData{})
//	if assert.NoError(t, err) {
//		assert.Equal(t, common.Denied, result.Status)
//		assert.Equal(t, common.Unknown, result.Details.Finality)
//		assert.Equal(t, []string{
//			DocStatusMissingOrInvalid,
//		}, result.Details.Reasons)
//	}
//}

func TestSynapseFI_CheckCustomerError(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return nil, errors.New("test_error")
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)

	result, err = service.CheckCustomer(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)
}

func TestSynapseFI_CheckCustomerValidComplexFlow(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusValid,
					},
				}, nil
			},
			GetOauthKeyFn: func(userID string, request verification.CreateOauthRequest) (*verification.OauthResponse, error) {
				return &verification.OauthResponse{
					ID: "someid",
					OAuthKey: "somekey",
					RefreshToken: "sometoken",
					ExpiresAt: "1498297390",
				}, nil
			},
			AddDocumentFn: func(userID string, userOAuth string, request verification.CreateDocumentsRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusValid,
					},
				}, nil
			},
		},
		kycFlow: "complex",
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Approved, result.Status)
	}
}

func TestSynapseFI_CheckCustomerValidComplexFlowInvalid(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusValid,
					},
				}, nil
			},
			GetOauthKeyFn: func(userID string, request verification.CreateOauthRequest) (*verification.OauthResponse, error) {
				return &verification.OauthResponse{
					ID: "someid",
					OAuthKey: "somekey",
					RefreshToken: "sometoken",
					ExpiresAt: "1498297390",
				}, nil
			},
			AddDocumentFn: func(userID string, userOAuth string, request verification.CreateDocumentsRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusInvalid,
					},
					Documents: []verification.ResponseDocument{
						{
							PhysicalDocs: []verification.ResponseSubDocument{
								{
									DocumentType: "TYPE",
									Status: DocStatusInvalid,
								},
							},
						},
					},
				}, nil
			},
		},
		kycFlow: "complex",
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Equal(t, common.Unknown, result.Details.Finality)
		assert.Equal(t, []string{
			"TYPE:" + DocStatusInvalid,
		}, result.Details.Reasons)
	}
}

func TestSynapseFI_CheckCustomerComplexCreateUserError(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return nil, errors.New("test_error")
			},
		},
		kycFlow: "complex",
	}

	result, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)

	result, err = service.CheckCustomer(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)
}

func TestSynapseFI_CheckCustomerComplexOAuthError(t *testing.T) {
	service := SynapseFI{
		verification: verification.Mock{
			CreateUserFn: func(request verification.CreateUserRequest) (*verification.UserResponse, error) {
				return &verification.UserResponse{
					ID: "someid",
					DocumentStatus: verification.DocumentStatus{
						PhysicalDoc: DocStatusValid,
					},
				}, nil
			},
			GetOauthKeyFn: func(userID string, request verification.CreateOauthRequest) (*verification.OauthResponse, error) {
				return nil, errors.New("test_error")
			},
		},
		kycFlow: "complex",
	}


	result, err := service.CheckCustomer(&common.UserData{})
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)

	result, err = service.CheckCustomer(nil)
	assert.Error(t, err)
	assert.Equal(t, common.Error, result.Status)
	assert.Nil(t, result.Details)
}

