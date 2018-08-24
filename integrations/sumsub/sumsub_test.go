package sumsub

import (
	"testing"

	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/applicants"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/documents"
	"gitlab.com/lambospeed/kyc/integrations/sumsub/verification"
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{
		SupplementalAddresses: []common.Address{
			{
				CountryAlpha2: "CountryAlpha2",
				PostCode:      "code",
				Town:          "Possum Springs",
			},
		},
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:    "PASSPORT",
					SubType: "FRONT_SIDE",
					Country: "RUSSIA",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.jpeg",
					ContentType: "image/jpeg",
					Data:        []byte{123, 23, 21, 2, 233},
				},
			},
		},
	})
	if assert.NoError(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Approved, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Unclear, status)
		assert.Equal(t, common.DetailedKYCResult{
			Finality: common.Final,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Denied, status)
		assert.Equal(t, common.DetailedKYCResult{
			Finality: common.NonFinal,
			Reasons: []string{
				"INCOMPLETE_DOCUMENT",
				"WRONG_USER_REGION",
			},
		}, *result)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Error, status)
		assert.Equal(t, common.DetailedKYCResult{
			Finality: common.Unknown,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.NotNil(t, result) {
		assert.Equal(t, common.Error, status)
		assert.Equal(t, common.DetailedKYCResult{
			Finality: common.Final,
			Reasons: []string{
				"ID_INVALID",
			},
		}, *result)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{
		Documents: []common.Document{
			{
				Front: &common.DocumentFile{},
			},
		},
	})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
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

	status, result, err := sumsubService.CheckCustomer(nil)
	if assert.Error(t, err) && assert.Nil(t, result) {
		assert.Equal(t, common.Error, status)
	}
}
