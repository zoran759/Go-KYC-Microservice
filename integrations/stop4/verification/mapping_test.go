package verification

import (
	"testing"
	"time"
	"modulus/kyc/common"
	"github.com/stretchr/testify/assert"
	"errors"
)

func TestMapCustomerToVerificationRequest(t *testing.T) {
	testTime := common.Time(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))

	customer := common.UserData{
		FirstName:            "FirstName",
		PaternalLastName:     "PaternalLastName",
		LastName:             "LastName",
		MiddleName:           "MiddleName",
		LegalName:            "LegalName",
		LatinISO1Name:        "LATIN",
		Email:                "Email",
		Gender:               common.Male,
		DateOfBirth:          testTime,
		PlaceOfBirth:         "PlaceOfBirth",
		CountryOfBirthAlpha2: "CountryOfBirth",
		StateOfBirth:         "StateOfBirth",
		CountryAlpha2:        "CountryAlpha2",
		Nationality:          "Nationality",
		Phone:                "Phone",
		MobilePhone:          "MobilePhone",
		CurrentAddress: common.Address{
			CountryAlpha2:     "Country1",
			County:            "County1",
			State:             "State1",
			Town:              "Town1",
			Suburb:            "Suburb1",
			Street:            "Street1",
			StreetType:        "StreetType1",
			SubStreet:         "SubStreet1",
			BuildingName:      "BuildingName1",
			BuildingNumber:    "BuildingNumber1",
			FlatNumber:        "FlatNumber1",
			PostCode:          "PostCode1",
			StateProvinceCode: "SPC1",
			PostOfficeBox:     "POB1",
			StartDate:         testTime,
			EndDate:           testTime,
		},
		Business: &common.Business{
			Name:                      "BusinessName",
			RegistrationNumber:        "RegNumber",
			IncorporationDate:         testTime,
			IncorporationJurisdiction: "IncorporationJurisdiction",
		},
		IDCard: &common.IDCard{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			Number:        "Number",
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Selfie: &common.Selfie{
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		UtilityBill: &common.UtilityBill{
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Passport: &common.Passport{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			Number:        "Number",
			Image: &common.DocumentFile{
				Filename:    "Filename",
				ContentType: "ContentType",
				Data:        []byte{1, 2, 3, 4, 5, 6, 7},
			},
		},
		Other: &common.Other{
			CountryAlpha2: "Country",
			IssuedDate:    testTime,
			ValidUntil:    testTime,
			Number:        "Number",
		},
	}

	verificationRequest := MapCustomerToVerificationRequest(customer)

	assert.Equal(t, customer.CurrentAddress.String(), verificationRequest.CustomerInformation.Address1.FieldVal)
	assert.Equal(t, customer.DateOfBirth.Format("2006-01-02"), verificationRequest.CustomerInformation.Dob.FieldVal)
	assert.Equal(t, customer.Email, verificationRequest.CustomerInformation.Email.FieldVal)
	assert.Equal(t, customer.FirstName, verificationRequest.CustomerInformation.FirstName.FieldVal)
	assert.Equal(t, customer.LastName, verificationRequest.CustomerInformation.LastName.FieldVal)
	assert.Equal(t, customer.MiddleName, verificationRequest.CustomerInformation.MiddleName.FieldVal)
	assert.Equal(t, customer.Phone, verificationRequest.CustomerInformation.Phone1.FieldVal)
	assert.Equal(t, "M", verificationRequest.CustomerInformation.Gender.FieldVal)
	assert.Equal(t, customer.MobilePhone, verificationRequest.CustomerInformation.Phone2.FieldVal)
	assert.Equal(t, customer.CurrentAddress.StateProvinceCode, verificationRequest.CustomerInformation.Province.FieldVal)
	assert.Equal(t, "127.1.2.3", verificationRequest.RegIPAddress)
	assert.Equal(t, "PO BOX POB1 County1 Town1 SPC1 PostCode1 ", verificationRequest.CustomerInformation.Address1.FieldVal)
}

func TestMapGender(t *testing.T) {
	assert.Equal(t, "M", MapGender(common.Male))
	assert.Equal(t, "F", MapGender(common.Female))
	assert.Empty(t, MapGender(common.Gender(10)))
}

func TestMapErrorCode(t *testing.T) {
	assert.Equal(t, "Registration IP Address is not present", MapErrorCode(-4))
	assert.Equal(t, "Wrong customer email format", MapErrorCode(-13))
	assert.Equal(t, "Merchant mismatch", MapErrorCode(-104))
	assert.Equal(t, "Unknown error", MapErrorCode(-1925))
}

func TestMapResponseError(t *testing.T) {
	status := 404

	body := `<!doctype html><html lang="en">
	<meta charset="utf-8">
	<title>apiary.io—404—No Resource Found!</title>
	<body>
	<div id="bg"></div>
	<div id="message">
	<h1>The resource you're looking for doesn't exist. <br>Please check the <a href="https://coreservices2.docs.apiary.io/">API documentation</a>.</h1>
	<p></p>
	</div>
	</body>
	</html>`

	htmlResponse, err := MapResponseError(status, []byte(body))
	assert.Equal(t, errors.New("The resource you're looking for doesn't exist. "), htmlResponse)
	assert.Nil(t, err)

	body = `<!doctype html><html lang="en">
	<meta charset="utf-8">
	<title>apiary.io—404—No Resource Found!</title>
	<body>
	<div id="bg"></div>
	<div id="message">
	<h3><h1>The resource you're looking for doesn't exist. <br>Please check the <a href="https://coreservices2.docs.apiary.io/">API documentation</a>.</h1></h3>
	<p></p>
	</div>
	</body>
	</html>`

	htmlResponse, err = MapResponseError(status, []byte(body))
	assert.Equal(t, errors.New("The resource you're looking for doesn't exist. "), htmlResponse)
	assert.Nil(t, err)

	body = `}`
	htmlResponse, err = MapResponseError(status, []byte(body))
	assert.Equal(t, errors.New("Unknown message"), htmlResponse)
	assert.Nil(t, err)
}

func TestMapResponseErrorFailed(t *testing.T) {
	status := 404

	body := `<!doctype html><html lang="en">
	<meta charset="utf-8">
	<title>apiary.io—404—No Resource Found!</title>
	<body>
	<div id="bg"></div>
	<div id="message">
	<h2>The resource you're looking for doesn't exist. <br>Please check the <a href="https://coreservices2.docs.apiary.io/">API documentation</a>.</h2>
	<p></p>
	</div>
	</body>
	</html>`

	htmlResponse, err := MapResponseError(status, []byte(body))

	assert.Equal(t, errors.New("Unknown message"), htmlResponse)
	assert.Nil(t, err)
}

func TestMapDocumentType(t *testing.T) {
	docType := mapDocumentType("IDCard")
	assert.Equal(t, 3, docType)

	docType = mapDocumentType("")
	assert.Equal(t, 1, docType)

	docType = mapDocumentType("DriverLicenseTranslation")
	assert.Equal(t, 4, docType)
}
