package sumsub

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/http"
	"modulus/kyc/integrations/sumsub/applicants"
	"modulus/kyc/integrations/sumsub/documents"
	"modulus/kyc/integrations/sumsub/verification"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	testImageUpload = flag.Bool("use-images", false, "test document images uploading")
	testPassport    []byte
	testSelfie      []byte
)

func TestNew(t *testing.T) {
	sumsubService := New(Config{
		Host:   "test_host",
		APIKey: "test_key",
	})

	assert.NotNil(t, sumsubService)
}

func TestSumSub_CheckCustomerGreen(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{ID: "test id"}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "completed", &verification.ReviewResult{
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
		assert.Equal(t, common.Unclear, result.Status)
		assert.NotNil(t, result.StatusCheck)
		assert.Equal(t, common.SumSub, result.StatusCheck.Provider)
		assert.Equal(t, "test id", result.StatusCheck.ReferenceID)
		assert.NotZero(t, time.Time(result.StatusCheck.LastCheck))
	}

	result, err = sumsubService.CheckStatus(result.StatusCheck.ReferenceID)
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Nil(t, result.StatusCheck)
		assert.Equal(t, common.Approved, result.Status)
	}
}

func TestSumSub_CheckCustomerYellow(t *testing.T) {
	sumsubService := SumSub{
		applicants: applicants.Mock{
			CreateApplicantFn: func(email string, applicant applicants.ApplicantInfo) (*applicants.CreateApplicantResponse, error) {
				return &applicants.CreateApplicantResponse{ID: "test id"}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "completed", &verification.ReviewResult{
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
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.NotNil(t, result.StatusCheck)
		assert.Equal(t, common.SumSub, result.StatusCheck.Provider)
		assert.Equal(t, "test id", result.StatusCheck.ReferenceID)
		assert.NotZero(t, time.Time(result.StatusCheck.LastCheck))
	}

	result, err = sumsubService.CheckStatus(result.StatusCheck.ReferenceID)
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
				return &applicants.CreateApplicantResponse{ID: "test id"}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "completed", &verification.ReviewResult{
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
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.NotNil(t, result.StatusCheck)
		assert.Equal(t, common.SumSub, result.StatusCheck.Provider)
		assert.Equal(t, "test id", result.StatusCheck.ReferenceID)
		assert.NotZero(t, time.Time(result.StatusCheck.LastCheck))
	}

	result, err = sumsubService.CheckStatus(result.StatusCheck.ReferenceID)
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
				return &applicants.CreateApplicantResponse{ID: "test id"}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "completed", &verification.ReviewResult{
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
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.NotNil(t, result.StatusCheck)
		assert.Equal(t, common.SumSub, result.StatusCheck.Provider)
		assert.Equal(t, "test id", result.StatusCheck.ReferenceID)
		assert.NotZero(t, time.Time(result.StatusCheck.LastCheck))
	}

	result, err = sumsubService.CheckStatus(result.StatusCheck.ReferenceID)
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
				return &applicants.CreateApplicantResponse{ID: "test id"}, nil
			},
		},
		documents: documents.Mock{
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
			},
			CheckApplicantStatusFn: func(applicantID string) (string, *verification.ReviewResult, error) {
				return "completed", &verification.ReviewResult{
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
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.NotNil(t, result.StatusCheck)
		assert.Equal(t, common.SumSub, result.StatusCheck.Provider)
		assert.Equal(t, "test id", result.StatusCheck.ReferenceID)
		assert.NotZero(t, time.Time(result.StatusCheck.LastCheck))
	}

	result, err = sumsubService.CheckStatus(result.StatusCheck.ReferenceID)
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
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
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
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return true, nil, nil
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
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return false, nil, nil
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
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, nil
			},
		},
		verification: verification.Mock{
			StartVerificationFn: func(applicantID string) (bool, *int, error) {
				return false, nil, errors.New("Unable to start a check")
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
			UploadDocumentFn: func(applicantID string, document documents.Document) (*documents.Metadata, *int, error) {
				return &documents.Metadata{}, nil, errors.New("Bad document")
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

func TestSumSubImageUpload(t *testing.T) {
	if !*testImageUpload {
		t.Skip("use '-use-images' flag to activate images uploading test")
	}

	assert := assert.New(t)

	if !assert.NotEmpty(testPassport, "testPassport must contain the content of the image data file 'passport.jpg'") {
		return
	}
	if !assert.NotEmpty(testSelfie, "testSelfie must contain the content of the image data file 'selfie.png'") {
		return
	}

	customer := &common.UserData{
		FirstName:     "John",
		LastName:      "Doe",
		Email:         "john.doe@mail.com",
		Gender:        common.Male,
		DateOfBirth:   common.Time(time.Date(1975, 06, 15, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "US",
		Nationality:   "US",
		CurrentAddress: common.Address{
			CountryAlpha2:     "US",
			State:             "Georgia",
			Town:              "Albany",
			Street:            "PeachTree Avenue",
			StreetType:        "Avenue",
			BuildingNumber:    "7315",
			FlatNumber:        "13",
			PostCode:          "31707",
			StateProvinceCode: "GA",
			StartDate:         common.Time(time.Date(1975, 06, 20, 0, 0, 0, 0, time.UTC)),
		},
		Passport: &common.Passport{
			Number:        "0123456789",
			CountryAlpha2: "US",
			State:         "GA",
			IssuedDate:    common.Time(time.Date(2015, 06, 20, 0, 0, 0, 0, time.UTC)),
			ValidUntil:    common.Time(time.Date(2025, 06, 19, 0, 0, 0, 0, time.UTC)),
			Image: &common.DocumentFile{
				Filename:    "passport.jpg",
				ContentType: "image/jpeg",
				Data:        testPassport,
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "selfie.png",
				ContentType: "image/png",
				Data:        testSelfie,
			},
		},
	}

	config := Config{
		Host:   "https://test-api.sumsub.com",
		APIKey: "GKTBNXNEPJHCXY",
	}

	service := New(config)

	result, err := service.CheckCustomer(customer)

	if !assert.Nil(err) {
		return
	}
	if !assert.NotNil(result.StatusPolling, "status polling data has to be provided") {
		return
	}

	applicantID := result.StatusPolling.CustomerID
	t.Log("Received applicant id:", applicantID)

	// Simulate approved result of the verification.
	_, _, err = http.Post(
		fmt.Sprintf(config.Host+"/resources/applicants/%s/status/testCompleted?key=%s", applicantID, config.APIKey),
		http.Headers{
			"Content-Type": "application/json",
		},
		[]byte(`{"reviewAnswer":"GREEN","rejectLabels":[]}`))

	if !assert.Nil(err) {
		return
	}

	result, err = service.CheckStatus(applicantID)

	if !assert.Nil(err) {
		return
	}
	assert.Equal(common.Approved, result.Status)
	assert.Nil(result.Details)
	assert.Empty(result.ErrorCode)
	assert.Nil(result.StatusPolling)

	// Get back the downloaded documents.
	type doc struct {
		IDDocType string `json:"idDocType"`
		Country   string `json:"country"`
		ImageID   int    `json:"imageId"`
	}
	type docs struct {
		Status struct {
			InspectionID string `json:"inspectionId"`
		} `json:"status"`
		DocumentStatus []doc `json:"documentStatus"`
	}

	_, body, err := http.Get(fmt.Sprintf(config.Host+"/resources/applicants/%s/state?key=%s", applicantID, config.APIKey), nil)

	if !assert.Nil(err) {
		return
	}
	if !assert.NotEmpty(body) {
		return
	}

	docsFromAPI := docs{}
	err = json.Unmarshal(body, &docsFromAPI)
	if !assert.Nil(err) {
		return
	}

	// Get back the downloaded document images.
	type img struct {
		doctype string
		data    []byte
	}

	docImages := []img{}
	for _, d := range docsFromAPI.DocumentStatus {
		_, body, err := http.Get(fmt.Sprintf(config.Host+"/resources/inspections/%s/resources/%d?key=%s", docsFromAPI.Status.InspectionID, d.ImageID, config.APIKey), nil)
		if !assert.Nil(err) {
			return
		}
		if !assert.NotEmpty(body) {
			return
		}
		docImages = append(docImages, img{
			doctype: d.IDDocType,
			data:    body,
		})
	}

	if !assert.NotEmpty(docImages) {
		return
	}

	for _, docImg := range docImages {
		switch docImg.doctype {
		case "PASSPORT":
			assert.Equal(customer.Passport.Image.Data, docImg.data)
		case "SELFIE":
			assert.Equal(customer.Selfie.Image.Data, docImg.data)
		}
	}
}

func init() {
	testPassport, _ = ioutil.ReadFile("../../test_data/passport.jpg")
	testSelfie, _ = ioutil.ReadFile("../../test_data/selfie.png")
}
