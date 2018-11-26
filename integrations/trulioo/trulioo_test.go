package trulioo

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/integrations/trulioo/configuration"
	"modulus/kyc/integrations/trulioo/verification"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var testImageUpload = flag.Bool("use-images", false, "test document images uploading")

func TestNew(t *testing.T) {
	_ = New(Config{})
}

func TestTrulioo_CheckCustomerNoMatch(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
				return configuration.Consents{}, nil, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
				return &verification.Response{
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

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Denied, result.Status)
		assert.Len(t, result.Details.Reasons, 2)
		assert.Equal(t, common.Unknown, result.Details.Finality)

		assert.Equal(t, []string{
			"Datasource Name has status: status; field statuses: Field name : status; Field name2 : status2; error: test error;",
			"Datasource Name1 has status: status1; field statuses: Field name3 : status3; Field name4 : status; error: test error2;",
		}, result.Details.Reasons)
	}
}

func TestTrulioo_CheckCustomerUnclear(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
				return configuration.Consents{}, nil, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
				return &verification.Response{
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

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) {
		assert.Equal(t, common.Unclear, result.Status)
		assert.Len(t, result.Details.Reasons, 2)
		assert.Equal(t, common.Unknown, result.Details.Finality)

		assert.Equal(t, []string{
			"Datasource Name has status: status; field statuses: Field name : status; Field name2 : status2; error: test error;",
			"Datasource Name1 has status: status1; field statuses: Field name3 : status3; Field name4 : status; error: test error2;",
		}, result.Details.Reasons)
	}
}

func TestTrulioo_CheckCustomerApproved(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
				return configuration.Consents{}, nil, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
				return &verification.Response{
					Record: verification.Record{
						RecordStatus: Match,
					},
				}, nil
			},
		},
	}

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.NoError(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, common.Approved, result.Status)
	}
}

func TestTrulioo_CheckCustomerError(t *testing.T) {
	service := Trulioo{
		configuration: configuration.Mock{
			ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
				return configuration.Consents{}, nil, nil
			},
		},
		verification: verification.Mock{
			VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
				return &verification.Response{
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

	result, err := service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, "Test error;Another test error;", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
			return &verification.Response{
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

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, "Test error1;Another test error2;", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
			return nil, errors.New("test error")
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, "test error", err.Error())
	}

	service.configuration = configuration.Mock{
		ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
			return nil, nil, errors.New("test error2")
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, "test error2", err.Error())
	}

	result, err = service.CheckCustomer(nil)
	if assert.Error(t, err) && assert.Nil(t, result.Details) {
		assert.Equal(t, "No customer supplied", err.Error())
	}
}

func TestImagesUpload(t *testing.T) {
	if !*testImageUpload {
		t.Skip("use '-use-images' flag to activate images uploading test")
	}

	assert := assert.New(t)

	host := "https://api.globaldatacompany.com"
	nAPILogin := "modulus.dev"
	nAPIPassword := "p9LF(m~CEKam*@88RHKDJ"

	r, err := http.NewRequest(http.MethodGet, host+"/configuration/v1/testentities/Identity%20Verification/US", nil)
	if !assert.NoError(err) {
		return
	}

	r.SetBasicAuth(nAPILogin, nAPIPassword)

	resp, err := http.DefaultClient.Do(r)
	if !assert.NoError(err) {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if !assert.NoError(err) || len(body) == 0 {
		return
	}

	testEntities := []map[string]interface{}{}
	err = json.Unmarshal(body, &testEntities)
	if !assert.NoError(err) {
		return
	}
	if !assert.Len(testEntities, 3, "testEntities length should be equal to 3") {
		return
	}

	customer, err := fillCustomer(testEntities[2])
	if !assert.NoError(err) {
		return
	}

	service := New(Config{
		Host:         host,
		NAPILogin:    nAPILogin,
		NAPIPassword: nAPIPassword,
	})

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func fillCustomer(testEntity map[string]interface{}) (customer *common.UserData, err error) {
	customer = &common.UserData{}

	if personInfoI, ok := testEntity["PersonInfo"]; ok {
		personInfo, ok := personInfoI.(map[string]interface{})
		if ok {
			if firstName, ok := personInfo["FirstGivenName"].(string); ok {
				if firstName == "J" {
					firstName = "Justin"
				}
				customer.FirstName = firstName
			}
			if middleName, ok := personInfo["MiddleName"].(string); ok {
				customer.MiddleName = middleName
			}
			if firstSurName, ok := personInfo["FirstSurName"].(string); ok {
				customer.LastName = firstSurName
			}
			if dayOfBirth, ok := personInfo["DayOfBirth"].(float64); ok {
				if monthOfBirth, ok := personInfo["MonthOfBirth"].(float64); ok {
					if yearOfBirth, ok := personInfo["YearOfBirth"].(float64); ok {
						dob, err := time.Parse("2006-01-02", fmt.Sprintf("%04.f-%02.f-%02.f", yearOfBirth, monthOfBirth, dayOfBirth))
						if err != nil {
							return nil, err
						}
						customer.DateOfBirth = common.Time(dob)
					}
				}
			}
			if gender, ok := personInfo["Gender"].(string); ok {
				switch gender {
				case "M":
					customer.Gender = common.Male
				case "F":
					customer.Gender = common.Female
				}
			}
		}
	}

	if locationI, ok := testEntity["Location"]; ok {
		location, ok := locationI.(map[string]interface{})
		if ok {
			if buildingNumber, ok := location["BuildingNumber"].(string); ok {
				customer.CurrentAddress.BuildingNumber = buildingNumber
			}
			if unitNumber, ok := location["UnitNumber"].(string); ok {
				customer.CurrentAddress.FlatNumber = unitNumber
			}
			if streetName, ok := location["StreetName"].(string); ok {
				customer.CurrentAddress.Street = streetName
			}
			if streetType, ok := location["StreetType"].(string); ok {
				customer.CurrentAddress.StreetType = streetType
			}
			if city, ok := location["City"].(string); ok {
				customer.CurrentAddress.Town = city
			}
			if suburb, ok := location["Suburb"].(string); ok {
				customer.CurrentAddress.Suburb = suburb
			}
			if stateProvinceCode, ok := location["StateProvinceCode"].(string); ok {
				customer.CurrentAddress.StateProvinceCode = stateProvinceCode
			}
			if postalCode, ok := location["PostalCode"].(string); ok {
				customer.CurrentAddress.PostCode = postalCode
			}
		}
	}

	if communicationI, ok := testEntity["Communication"]; ok {
		communication, ok := communicationI.(map[string]interface{})
		if ok {
			if telephone, ok := communication["Telephone"].(string); ok {
				customer.Phone = telephone
			}
			if email, ok := communication["EmailAddress"].(string); ok {
				customer.Email = email
			}
		}
	}

	if driversI, ok := testEntity["DriverLicence"]; ok {
		drivers, ok := driversI.(map[string]interface{})
		if ok {
			if number, ok := drivers["Number"].(string); ok {
				customer.DriverLicense = &common.DriverLicense{
					Number: number,
				}
			}
		}
	}

	passportI, ok := testEntity["Passport"]
	if !ok {
		err = errors.New("no passport in the test entity")
		return
	}
	passport, ok := passportI.(map[string]interface{})
	if !ok {
		err = errors.New("unknown model for a passport in the test entity")
		return
	}

	customer.Passport = &common.Passport{}
	mrz, ok := passport["Mrz1"].(string)
	if ok {
		customer.Passport.Mrz1 = mrz
	}
	mrz, ok = passport["Mrz2"].(string)
	if ok {
		customer.Passport.Mrz2 = mrz
	}
	number, ok := passport["Number"].(string)
	if ok {
		customer.Passport.Number = number
	}
	if dayOfExpiry, ok := passport["DayOfExpiry"].(float64); ok {
		if monthOfExpiry, ok := passport["MonthOfExpiry"].(float64); ok {
			if yearOfExpiry, ok := passport["YearOfExpiry"].(float64); ok {
				expDate, err := time.Parse("2006-01-02", fmt.Sprintf("%04.f-%02.f-%02.f", yearOfExpiry, monthOfExpiry, dayOfExpiry))
				if err != nil {
					return nil, err
				}
				customer.Passport.ValidUntil = common.Time(expDate)
			}
		}
	}

	passportImage, err := ioutil.ReadFile("../../test_data/passport.jpg")
	if err != nil {
		return nil, err
	}
	customer.Passport.Image = &common.DocumentFile{
		Filename:    "passport.jpg",
		ContentType: "image/jpeg",
		Data:        passportImage,
	}

	customer.CountryAlpha2 = "US"

	return
}
