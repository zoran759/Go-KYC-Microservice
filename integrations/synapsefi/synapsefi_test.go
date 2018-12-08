package synapsefi

import (
	"errors"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"

	"github.com/stretchr/testify/assert"
)

var testTime = time.Now().Unix()

type Mock struct {
	CreateUserFn      func(verification.User) (*verification.Response, *string, error)
	AddPhysicalDocsFn func(string, string, string, []verification.SubDocument) (*string, error)
	GetUserFn         func(string) (*verification.Response, *string, error)
}

func (m Mock) CreateUser(user verification.User) (*verification.Response, *string, error) {
	return m.CreateUserFn(user)
}

func (m Mock) AddPhysicalDocs(userID string, rtoken string, docsID string, physdocs []verification.SubDocument) (*string, error) {
	return m.AddPhysicalDocsFn(userID, rtoken, docsID, physdocs)
}

func (m Mock) GetUser(refID string) (*verification.Response, *string, error) {
	return m.GetUserFn(refID)
}

func TestNew(t *testing.T) {
	config := Config{
		Host:         "host",
		ClientID:     "client_id",
		ClientSecret: "client_secret",
	}

	service := SynapseFI{
		verification: verification.NewService(config),
	}

	testservice := New(config)

	assert.Equal(t, service, testservice)
}

func TestSynapseFI_CheckCustomerValid(t *testing.T) {
	service := SynapseFI{
		verification: Mock{
			CreateUserFn: func(verification.User) (*verification.Response, *string, error) {
				return &verification.Response{
					ID: "test_id",
					Documents: []verification.ResponseDocument{
						verification.ResponseDocument{
							ID:              "rdid",
							PermissionScope: "UNVERIFIED",
							VirtualDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
									ID:          "vid",
									Type:        "SSN",
									LastUpdated: testTime,
									Status:      "SUBMITTED",
								},
							},
							PhysicalDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
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
				}, nil, nil
			},
			AddPhysicalDocsFn: func(string, string, string, []verification.SubDocument) (*string, error) {
				return nil, nil
			},
			GetUserFn: func(string) (*verification.Response, *string, error) {
				return &verification.Response{
					ID: "test_id",
					Documents: []verification.ResponseDocument{
						verification.ResponseDocument{
							ID:              "rdid",
							PermissionScope: "SEND|RECEIVE|1000|DAILY",
							VirtualDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
									ID:          "vid",
									Type:        "SSN",
									LastUpdated: testTime,
									Status:      "SUBMITTED|VALID",
								},
							},
							PhysicalDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
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
				}, nil, nil
			},
		},
	}

	assert := assert.New(t)

	result, err := service.CheckCustomer(&common.UserData{
		IDCard: &common.IDCard{
			Number:        "123456789",
			CountryAlpha2: "US",
			IssuedDate:    common.Time(time.Date(1968, 6, 30, 0, 0, 0, 0, time.UTC)),
			Image: &common.DocumentFile{
				Filename:    "ssn.jpg",
				ContentType: "image/jpeg",
				Data:        []byte("fake content"),
			},
		},
	})

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

	result, err = service.CheckStatus(result.StatusCheck.ReferenceID)

	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestSynapseFI_CheckCustomerInvalid(t *testing.T) {
	service := SynapseFI{
		verification: Mock{
			CreateUserFn: func(verification.User) (*verification.Response, *string, error) {
				return &verification.Response{
					ID: "test_id",
					Documents: []verification.ResponseDocument{
						verification.ResponseDocument{
							ID:              "rdid",
							PermissionScope: "UNVERIFIED",
							VirtualDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
									ID:          "vid",
									Type:        "SSN",
									LastUpdated: testTime,
									Status:      "SUBMITTED",
								},
							},
							PhysicalDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
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
				}, nil, nil
			},
			AddPhysicalDocsFn: func(string, string, string, []verification.SubDocument) (*string, error) {
				return nil, nil
			},
			GetUserFn: func(string) (*verification.Response, *string, error) {
				return &verification.Response{
					ID: "test_id",
					Documents: []verification.ResponseDocument{
						verification.ResponseDocument{
							ID:              "rdid",
							PermissionScope: "UNVERIFIED",
							VirtualDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
									ID:          "vid",
									Type:        "SSN",
									LastUpdated: testTime,
									Status:      "SUBMITTED|INVALID",
								},
							},
							PhysicalDocs: []verification.ResponseSubDocument{
								verification.ResponseSubDocument{
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
				}, nil, nil
			},
		},
	}

	assert := assert.New(t)

	result, err := service.CheckCustomer(&common.UserData{
		IDCard: &common.IDCard{
			Number:        "123456789",
			CountryAlpha2: "US",
			IssuedDate:    common.Time(time.Date(1968, 6, 30, 0, 0, 0, 0, time.UTC)),
			Image: &common.DocumentFile{
				Filename:    "ssn.jpg",
				ContentType: "image/jpeg",
				Data:        []byte("fake content"),
			},
		},
	})

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

	result, err = service.CheckStatus(result.StatusCheck.ReferenceID)

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

func TestSynapseFI_CheckCustomerError(t *testing.T) {
	service := SynapseFI{
		verification: Mock{
			CreateUserFn: func(verification.User) (*verification.Response, *string, error) {
				code := "404"

				return nil, &code, errors.New("test_error")
			},
			AddPhysicalDocsFn: func(string, string, string, []verification.SubDocument) (*string, error) {
				return nil, nil
			},
		},
	}

	assert := assert.New(t)

	result, err := service.CheckCustomer(&common.UserData{})
	assert.Error(err)
	assert.EqualError(err, "failed to get document's number from customer documents or no document was supplied")
	assert.Equal(common.Error, result.Status)
	assert.Nil(result.Details)
	assert.Empty(result.ErrorCode)
	assert.Nil(result.StatusCheck)

	result, err = service.CheckCustomer(&common.UserData{
		IDCard: &common.IDCard{
			Number:        "123456789",
			CountryAlpha2: "US",
			IssuedDate:    common.Time(time.Date(1968, 6, 30, 0, 0, 0, 0, time.UTC)),
			Image: &common.DocumentFile{
				Filename:    "ssn.jpg",
				ContentType: "image/jpeg",
				Data:        []byte("fake content"),
			},
		},
	})

	assert.Error(err)
	assert.EqualError(err, "test_error")
	assert.Equal(common.Error, result.Status)
	assert.Nil(result.Details)
	assert.Equal("404", result.ErrorCode)
	assert.Nil(result.StatusCheck)

	result, err = service.CheckCustomer(nil)
	assert.Error(err)
	assert.EqualError(err, "no customer supplied")
	assert.Equal(common.Error, result.Status)
	assert.Nil(result.Details)
	assert.Empty(result.ErrorCode)
	assert.Nil(result.StatusCheck)
}
