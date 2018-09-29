package sumsub

import (
	"testing"

	"modulus/kyc/common"
	"modulus/kyc/integrations/sumsub/applicants"
	"modulus/kyc/integrations/sumsub/documents"
	"modulus/kyc/integrations/sumsub/verification"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	sumsubService := New(Config{
		Host:             "test_host",
		APIKey:           "test_key",
		TimeoutThreshold: 123456,
	})

	assert.Equal(t, int64(123456), sumsubService.timeoutThreshold)
}

func TestSumSub_CheckCustomerGreen(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return CompleteStatus, &verification.ReviewResult{
					ReviewAnswer: GreenScore,
				}, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{
		SupplementalAddresses: []common.Address{
			{
				CountryAlpha2: "CountryAlpha2",
				PostCode:      "code",
				Town:          "Possum Springs",
			},
		},
		Passport: &common.Passport{
			CountryAlpha2: "RU",
			Image: &common.DocumentFile{
				Filename:    "passport.jpeg",
				ContentType: "image/jpeg",
				Data:        []byte{123, 23, 21, 2, 233},
			},
		},
	})
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Approved, result.Status)
	}
}

func TestSumSub_CheckCustomerYellow(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return CompleteStatus, &verification.ReviewResult{
					ReviewAnswer: YellowScore,
					Label:        "TEST_LABEL",
					RejectLabels: []string{
						"ID_INVALID",
					},
					ReviewRejectType: "FINAL",
				}, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result.Details) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.Equal(t, common.KYCDetails{
			Finality: common.Final,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result.Details)
	}
}

func TestSumSub_CheckCustomerRed(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return CompleteStatus, &verification.ReviewResult{
					ReviewAnswer: RedScore,
					Label:        "TEST_LABEL",
					RejectLabels: []string{
						"INCOMPLETE_DOCUMENT",
						"WRONG_USER_REGION",
					},
					ReviewRejectType: "RETRY",
				}, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result.Details) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Equal(t, common.KYCDetails{
			Finality: common.NonFinal,
			Reasons: []string{
				"INCOMPLETE_DOCUMENT",
				"WRONG_USER_REGION",
			},
		}, *result.Details)
	}
}

func TestSumSub_CheckCustomerError(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return CompleteStatus, &verification.ReviewResult{
					ReviewAnswer: ErrorScore,
					Label:        "TEST_LABEL",
					RejectLabels: []string{
						"ID_INVALID",
					},
					ReviewRejectType: "EXTERNAL",
				}, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
		assert.Equal(t, common.KYCDetails{
			Finality: common.Unknown,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result.Details)
	}
}

func TestSumSub_CheckCustomerIgnored(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return CompleteStatus, &verification.ReviewResult{
					ReviewAnswer: IgnoredScore,
					Label:        "TEST_LABEL",
					RejectLabels: []string{
						"ID_INVALID",
					},
					ReviewRejectType: "FINAL",
				}, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
		assert.Equal(t, common.KYCDetails{
			Finality: common.Final,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result.Details)
	}
}

func TestSumSub_CheckCustomerTimeout(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "pending", nil, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerErrorTimeout(t *testing.T) {
	checkApplicantInvoked := false

	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return true, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				if !checkApplicantInvoked {
					checkApplicantInvoked = true
					return "", nil, errors.New("Check applicant error")
				}

				return "pending", nil, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerNotStartedUnknownReasons(t *testing.T) {
	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return false, nil
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerNotStartedError(t *testing.T) {
	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, error) {
				return false, errors.New("Unable to start a check")
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerDocumentUploadError(t *testing.T) {
	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, error) {
				return &documents.Metadata{}, errors.New("Bad document")
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{
		Other: &common.Other{
			Image: &common.DocumentFile{},
		},
	})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerCreateApplicantError(t *testing.T) {
	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return nil, errors.New("test_error")
			},
		},
	}

	result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}

func TestSumSub_CheckCustomerNoApplicantError(t *testing.T) {
	sumsubService := &SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return nil, errors.New("test_error")
			},
		},
	}

	result, err := sumsubService.CheckCustomer(nil)
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Error, result.Status)
	}
}
