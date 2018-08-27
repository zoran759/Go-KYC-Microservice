package trulioo

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/verification"
)

func TestNew(t *testing.T) {
	_ = New(Config{})
}

func TestTrulioo_CheckCustomerNoMatch(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, error) {
				return configuration.Consents{}, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
				return &verification.VerificationResponse{
					Record: verification.Record{
						RecordStatus: NoMatch,
						DatasourceResults: []verification.DatasourceResult{
							{
								DatasourceStatus: "status",
								DatasourceName:   "Name",
								DatasourceFields: []verification.DatasourceField{
									{
										FieldName: "Field name",
										Status:    "status",
									},
									{
										FieldName: "Field name2",
										Status:    "status2",
									},
								},
								Errors: verification.Errors{
									{
										Code:    "400",
										Message: "test error",
									},
								},
							},
							{
								DatasourceStatus: "status1",
								DatasourceName:   "Name1",
								DatasourceFields: []verification.DatasourceField{
									{
										FieldName: "Field name3",
										Status:    "status3",
									},
									{
										FieldName: "Field name4",
										Status:    "status",
									},
								},
								Errors: verification.Errors{
									{
										Code:    "400",
										Message: "test error2",
									},
								},
							},
							{},
						},
					},
				}, nil
			},
		},
	}

	result, detailedKYCResult, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result)
		assert.Len(t, detailedKYCResult.Reasons, 2)
		assert.Equal(t, common.Unknown, detailedKYCResult.Finality)

		assert.Equal(t, []string{
			"Datasource Name has status: status; field statuses: Field name : status; Field name2 : status2; error: test error;",
			"Datasource Name1 has status: status1; field statuses: Field name3 : status3; Field name4 : status; error: test error2;",
		}, detailedKYCResult.Reasons)
	}
}

func TestTrulioo_CheckCustomerUnclear(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, error) {
				return configuration.Consents{}, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
				return &verification.VerificationResponse{
					Record: verification.Record{
						RecordStatus: "sdfsdf",
						DatasourceResults: []verification.DatasourceResult{
							{
								DatasourceStatus: "status",
								DatasourceName:   "Name",
								DatasourceFields: []verification.DatasourceField{
									{
										FieldName: "Field name",
										Status:    "status",
									},
									{
										FieldName: "Field name2",
										Status:    "status2",
									},
								},
								Errors: verification.Errors{
									{
										Code:    "400",
										Message: "test error",
									},
								},
							},
							{
								DatasourceStatus: "status1",
								DatasourceName:   "Name1",
								DatasourceFields: []verification.DatasourceField{
									{
										FieldName: "Field name3",
										Status:    "status3",
									},
									{
										FieldName: "Field name4",
										Status:    "status",
									},
								},
								Errors: verification.Errors{
									{
										Code:    "400",
										Message: "test error2",
									},
								},
							},
							{},
						},
					},
				}, nil
			},
		},
	}

	result, detailedKYCResult, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Unclear, result)
		assert.Len(t, detailedKYCResult.Reasons, 2)
		assert.Equal(t, common.Unknown, detailedKYCResult.Finality)

		assert.Equal(t, []string{
			"Datasource Name has status: status; field statuses: Field name : status; Field name2 : status2; error: test error;",
			"Datasource Name1 has status: status1; field statuses: Field name3 : status3; Field name4 : status; error: test error2;",
		}, detailedKYCResult.Reasons)
	}
}

func TestTrulioo_CheckCustomerApproved(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, error) {
				return configuration.Consents{}, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
				return &verification.VerificationResponse{
					Record: verification.Record{
						RecordStatus: Match,
					},
				}, nil
			},
		},
	}

	result, detailedKYCResult, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.Nil(t, detailedKYCResult) {
		assert.Equal(t, common.Approved, result)
	}
}

func TestTrulioo_CheckCustomerError(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, error) {
				return configuration.Consents{}, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
				return &verification.VerificationResponse{
					Record: verification.Record{
						RecordStatus: NoMatch,
						Errors: verification.Errors{
							{
								Code:    "400",
								Message: "Test error",
							},
							{
								Code:    "500",
								Message: "Another test error",
							},
						},
					},
				}, nil
			},
		},
	}

	result, detailedKYCResult, err := service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, detailedKYCResult) && assert.Equal(t, common.Error, result) {
		assert.Equal(t, "Test error;Another test error;", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
			return &verification.VerificationResponse{
				Errors: verification.Errors{
					{
						Code:    "400",
						Message: "Test error1",
					},
					{
						Code:    "500",
						Message: "Another test error2",
					},
				},
			}, nil
		},
	}

	result, detailedKYCResult, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, detailedKYCResult) && assert.Equal(t, common.Error, result) {
		assert.Equal(t, "Test error1;Another test error2;", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.VerificationResponse, error) {
			return nil, errors.New("test error")
		},
	}

	result, detailedKYCResult, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, detailedKYCResult) && assert.Equal(t, common.Error, result) {
		assert.Equal(t, "test error", err.Error())
	}

	service.configuration = configuration.Mock{
		ConsentsFn: func(countryAlpha2 string) (configuration.Consents, error) {
			return nil, errors.New("test error2")
		},
	}

	result, detailedKYCResult, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, detailedKYCResult) && assert.Equal(t, common.Error, result) {
		assert.Equal(t, "test error2", err.Error())
	}

	result, detailedKYCResult, err = service.CheckCustomer(nil)
	if assert.Error(t, err) && assert.Nil(t, detailedKYCResult) && assert.Equal(t, common.Error, result) {
		assert.Equal(t, "No customer supplied", err.Error())
	}
}
