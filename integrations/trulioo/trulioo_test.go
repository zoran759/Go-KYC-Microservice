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

	assert := assert.New(t)
	if assert.NoError(err) {
		assert.NotNil(result.Details)
		assert.Len(result.Details.Reasons, 2)
		assert.Equal("400 Test error", result.Details.Reasons[0])
		assert.Equal("500 Another test error", result.Details.Reasons[1])
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
	if assert.Error(err) && assert.Nil(result.Details) {
		assert.Equal("Test error1;Another test error2;", err.Error())
	}

	service.verification = verification.Mock{
		VerifyFn: func(countryAlpha2 string, consents configuration.Consents, fields verification.DataFields) (*verification.Response, error) {
			return nil, errors.New("test error")
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(err) && assert.Nil(result.Details) {
		assert.Equal("test error", err.Error())
	}

	service.configuration = configuration.Mock{
		ConsentsFn: func(countryAlpha2 string) (configuration.Consents, *int, error) {
			return nil, nil, errors.New("test error2")
		},
	}

	result, err = service.CheckCustomer(&common.UserData{})
	if assert.Error(err) && assert.Nil(result.Details) {
		assert.Equal("test error2", err.Error())
	}

	result, err = service.CheckCustomer(nil)
	if assert.Error(err) && assert.Nil(result.Details) {
		assert.Equal("No customer supplied", err.Error())
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

func TestAEValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "Amir",
		LastName:      "Saliba",
		LatinISO1Name: "Amir Saliba",
		DateOfBirth:   common.Time(time.Date(1986, 10, 24, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "AE",
		Passport: &common.Passport{
			Number:     "506079867",
			Mrz1:       "P<ARESALIBA<<AMIR<<<<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "5060798672ARE8610245M24051827847968855030570",
			ValidUntil: common.Time(time.Date(2024, 5, 18, 0, 0, 0, 0, time.UTC)),
		},
		NationalID: &common.NationalID{
			Number: "784-7968-8550305-0",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestAUValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:            "John",
		LastName:             "Smith",
		MiddleName:           "Henry",
		Email:                "testpersonAU@gdctest.com",
		DateOfBirth:          common.Time(time.Date(1983, 3, 5, 0, 0, 0, 0, time.UTC)),
		CountryOfBirthAlpha2: "AU",
		CountryAlpha2:        "AU",
		Phone:                "0398968785",
		CurrentAddress: common.Address{
			Suburb:            "Doncaster",
			Street:            "Lawford",
			StreetType:        "st",
			BuildingNumber:    "10",
			FlatNumber:        "3",
			PostCode:          "3108",
			StateProvinceCode: "VIC",
		},
		Passport: &common.Passport{
			Number:     "N1236548",
			Mrz1:       "P<SAGSMITH<<JOHN<HENRY<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "N1236548<1AUS8303052359438740809<<<<<<<<<<54",
			ValidUntil: common.Time(time.Date(2018, 12, 5, 0, 0, 0, 0, time.UTC)),
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestCNValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:         "言明",
		LastName:          "胡",
		LatinISO1Name:     "SUSAN MITCHELL",
		DateOfBirth:       common.Time(time.Date(1976, 6, 5, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2:     "CN",
		Phone:             "15022246786",
		BankAccountNumber: "1000000000000456",
		Passport: &common.Passport{
			Number:     "G86453123",
			Mrz1:       "P<CHNSUSAN<<MITCHELL<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "C8456324<1CNL8303052359438740809<<<<<<<<<<54",
			ValidUntil: common.Time(time.Date(2020, 6, 12, 0, 0, 0, 0, time.UTC)),
		},
		NationalID: &common.NationalID{
			Number: "440861896421345987",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestGBValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "Julia",
		LastName:      "Audi",
		MiddleName:    "Ronald",
		Gender:        1,
		DateOfBirth:   common.Time(time.Date(1979, 10, 26, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "GB",
		Phone:         "+44865413985",
		MobilePhone:   "+448654139123",
		UKNHSNumber:   "1634567897",
		UKNINumber:    "BL261079C",
		CurrentAddress: common.Address{
			Town:           "Aylesbury",
			Street:         "Chiltern",
			StreetType:     "Court",
			BuildingName:   "Beck",
			BuildingNumber: "12",
			PostCode:       "HP22 6EP",
		},
		Passport: &common.Passport{
			Number:     "54846031",
			Mrz1:       "P<SAGAUDI<<JULIA<<<<<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "99003853<1CZE1101018M1207046110101111<<<<<94",
			ValidUntil: common.Time(time.Date(2020, 1, 7, 0, 0, 0, 0, time.UTC)),
		},
		DriverLicense: &common.DriverLicense{
			Number: "AUDI9710269JR9AB",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestKRValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "한",
		LastName:      "이친",
		LatinISO1Name: "Lee Chin Han",
		DateOfBirth:   common.Time(time.Date(1992, 3, 14, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "KR",
		Passport: &common.Passport{
			Number:        "JH56851",
			Mrz1:          "P<KORHAN<<LEE<CHIN<<<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:          "JH56851<<7KOR9203147M2111038<<<<<<<<<<<<<<04",
			CountryAlpha2: "KR",
			ValidUntil:    common.Time(time.Date(2021, 11, 3, 0, 0, 0, 0, time.UTC)),
		},
		IDCard: &common.IDCard{
			Number: "AM125D",
		},
		DriverLicense: &common.DriverLicense{
			Number: "경기9852342287",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestMXValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:        "Jose",
		LastName:         "Santano",
		MaternalLastName: "Martinez",
		DateOfBirth:      common.Time(time.Date(1978, 12, 12, 0, 0, 0, 0, time.UTC)),
		StateOfBirth:     "Sonora",
		CountryAlpha2:    "MX",
		Phone:            "6251140504",
		CurrentAddress: common.Address{
			Town:     "Hermosillo",
			PostCode: "83010",
		},
		Passport: &common.Passport{
			Number:     "S85416687",
			Mrz1:       "P<SAGSANTANO<MARTINEZ<<JOSE<<<<<<<<<<<<<<<<<",
			Mrz2:       "99003853<1CZE1101018M1207046110101111<<<<<94",
			ValidUntil: common.Time(time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)),
		},
		NationalID: &common.NationalID{
			Number: "HEGG560427MVZRRL05",
		},
		SocialService: &common.SocialService{
			Number: "HEGG560427ABC",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestMYValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:            "Wu",
		LastName:             "Boon Chai",
		FullName:             "Boon Chai Wu",
		DateOfBirth:          common.Time(time.Date(1987, 11, 29, 0, 0, 0, 0, time.UTC)),
		CountryOfBirthAlpha2: "MY",
		StateOfBirth:         "SRW",
		Gender:               1,
		CountryAlpha2:        "MY",
		Phone:                "0122323273",
		Location: &common.Location{
			Latitude:  "45.661",
			Longitude: "-111.067",
		},
		CurrentAddress: common.Address{
			Town:     "Miri",
			PostCode: "35500",
		},
		Passport: &common.Passport{
			Number:     "S8441773",
			Mrz1:       "P<SAGWU<<BOON<<CHAI<<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "S8441773<1MYS8303052359438740809<<<<<<<<<<54",
			ValidUntil: common.Time(time.Date(2022, 10, 20, 0, 0, 0, 0, time.UTC)),
		},
		NationalID: &common.NationalID{
			Number: "871129134567",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestNZValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "Snow",
		LastName:      "Huntsman",
		MiddleName:    "White",
		DateOfBirth:   common.Time(time.Date(1976, 3, 6, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "NZ",
		Phone:         "078475332",
		VehicleRegistrationPlate: "ABC123",
		CurrentAddress: common.Address{
			Town:           "Auckland",
			Suburb:         "Mt Roskill",
			Street:         "Carr",
			StreetType:     "Road",
			BuildingNumber: "50",
			FlatNumber:     "2",
			PostCode:       "1041",
		},
		Passport: &common.Passport{
			Number:     "A7894634",
			Mrz1:       "P<NZLWHITE<<SNOW<<<<<<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "N8456324<1NZL8303052359438740809<<<<<<<<<<54",
			ValidUntil: common.Time(time.Date(2020, 6, 4, 0, 0, 0, 0, time.UTC)),
		},
		DriverLicense: &common.DriverLicense{
			Number:  "8465341",
			Version: "3",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestRUValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "Борис",
		LastName:      "Иванов",
		MiddleName:    "Сергеевич",
		DateOfBirth:   common.Time(time.Date(1922, 12, 30, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "RU",
		Phone:         "8311234567",
		CurrentAddress: common.Address{
			Town:           "Октябрьский",
			Street:         "Советская",
			BuildingNumber: "9",
			FlatNumber:     "34",
			PostCode:       "606123",
		},
		Passport: &common.Passport{
			Number:     "2594123456",
			IssuedDate: common.Time(time.Date(2005, 8, 16, 0, 0, 0, 0, time.UTC)),
		},
		SocialService: &common.SocialService{
			Number: "12345678901",
		},
		TaxID: &common.TaxID{
			Number: "123456789012",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}

func TestUSValidTestEntity(t *testing.T) {
	assert := assert.New(t)

	service := New(Config{
		Host:         "https://api.globaldatacompany.com",
		NAPILogin:    "modulus.dev",
		NAPIPassword: "p9LF(m~CEKam*@88RHKDJ",
	})

	customer := &common.UserData{
		FirstName:     "Justin",
		LastName:      "Williams",
		MiddleName:    "Mark",
		Email:         "testpersonUS@gdctest.com",
		DateOfBirth:   common.Time(time.Date(1988, 8, 4, 0, 0, 0, 0, time.UTC)),
		CountryAlpha2: "US",
		Phone:         "802 660 9697",
		CurrentAddress: common.Address{
			Town:              "New York",
			Street:            "9th",
			StreetType:        "Avenue",
			BuildingNumber:    "420",
			FlatNumber:        "18",
			PostCode:          "10001",
			StateProvinceCode: "NY",
		},
		Passport: &common.Passport{
			Number:     "S85416687",
			Mrz1:       "P<USAWILLIAMS<<JUSTIN<<<<<<<<<<<<<<<<<<<<<<<",
			Mrz2:       "99003853<1USA1101018M1207046110101111<<<<<94",
			IssuedDate: common.Time(time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC)),
		},
		DriverLicense: &common.DriverLicense{
			Number: "0812319884104",
		},
		SocialService: &common.SocialService{
			Number: "000568791",
		},
	}

	result, err := service.CheckCustomer(customer)
	if assert.NoError(err) {
		assert.Equal(common.Approved, result.Status)
		assert.Nil(result.Details)
		assert.Empty(result.ErrorCode)
		assert.Nil(result.StatusCheck)
	}
}
