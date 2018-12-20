package synapsefi

import (
	"errors"
	"flag"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/synapsefi/verification"

	"github.com/stretchr/testify/assert"
)

var (
	useSandbox = flag.Bool("sandbox", false, "activate sandbox testing")
	testTime   = time.Now().Unix()
)

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

func TestAgainstSandbox(t *testing.T) {
	if !*useSandbox {
		t.Skip("use '-sandbox' command line flag to activate sandbox testing")
	}

	assert := assert.New(t)

	customer := &common.UserData{
		LegalName:   "John Doe",
		Email:       "john.doe@example.com",
		IPaddress:   "192.168.0.116",
		Gender:      common.Male,
		DateOfBirth: common.Time(time.Date(1960, 7, 24, 0, 0, 0, 0, time.UTC)),
		Phone:       "+912111222333",
		CurrentAddress: common.Address{
			CountryAlpha2:     "US",
			Town:              "San Francisco",
			Street:            "Mission Street",
			BuildingNumber:    "101",
			PostCode:          "94105",
			StateProvinceCode: "CA",
		},
		IDCard: &common.IDCard{
			Number:        "777772222",
			CountryAlpha2: "US",
			Image: &common.DocumentFile{
				Filename:    "ssn.jpg",
				ContentType: "image/jpeg",
				Data:        []byte("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"),
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "selfie.jpg",
				ContentType: "image/jpeg",
				Data:        []byte("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"),
			},
		},
	}

	service := New(Config{
		Host:         "https://uat-api.synapsefi.com/v3.1/",
		ClientID:     "client_id_9tCAZNrlj3gOeUYGIKucqEb0pQmx6zy1W2VBasX7",
		ClientSecret: "client_secret_UWVu3Y5E82Amai9q1Tk0PGlXLhyZHf0rNCneSpos",
	})

	result, err := service.CheckCustomer(customer)

	if assert.NoError(err) {
		assert.Equal(common.Unclear, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)

		if assert.NotNil(result.StatusCheck) {
			assert.Equal(common.SynapseFI, result.StatusCheck.Provider)
			assert.NotEmpty(result.StatusCheck.ReferenceID)
			assert.NotZero(result.StatusCheck.LastCheck)
		}
	}

	time.Sleep(5 * time.Second)

	result, err = service.CheckStatus(result.StatusCheck.ReferenceID)

	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}
