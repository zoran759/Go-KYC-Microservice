package verification

import (
	"modulus/kyc/common"
	"time"
	"errors"
	"golang.org/x/net/html"
	"strings"
)

// MapCustomerToVerificationRequest maps the values of common data to the service specific values.
func MapCustomerToVerificationRequest(customer common.UserData) RegistrationRequest {
	request := RegistrationRequest{}

	request.RegIPAddress = "127.1.2.3"
	request.RegDate = time.Now().Format("2006-01-02")

	request.CustomerInformation = CustomerInformation{}

	request.CustomerInformation.FirstName =
			CustomerInformationField{
				FieldName:  "customer_information[first_name]",
				FieldVal: customer.FirstName,
			}
	request.CustomerInformation.MiddleName =
		CustomerInformationField{
			FieldName:  "customer_information[middle_name]",
			FieldVal: customer.MiddleName,
		}
	request.CustomerInformation.LastName =
		CustomerInformationField{
			FieldName:  "customer_information[last_name]",
			FieldVal: customer.LastName,
		}
	request.CustomerInformation.Email =
		CustomerInformationField{
			FieldName:  "customer_information[email]",
			FieldVal: customer.Email,
		}
	request.CustomerInformation.Country =
		CustomerInformationField{
			FieldName:  "customer_information[country]",
			FieldVal: customer.CountryAlpha2,
		}
	request.CustomerInformation.PostalCode =
		CustomerInformationField{
			FieldName:  "customer_information[postal_code]",
			FieldVal: customer.CurrentAddress.PostCode,
		}
	request.CustomerInformation.Province =
		CustomerInformationField{
			FieldName:  "customer_information[province]",
			FieldVal: customer.CurrentAddress.StateProvinceCode,
		}
	request.CustomerInformation.Gender =
		CustomerInformationField{
			FieldName:  "customer_information[gender]",
			FieldVal: MapGender(customer.Gender),
		}
	request.CustomerInformation.Address1 =
		CustomerInformationField{
			FieldName:  "customer_information[address1]",
			FieldVal: customer.CurrentAddress.String(),
		}
	request.CustomerInformation.Address2 =
		CustomerInformationField{
			FieldName:  "customer_information[address2]",
			FieldVal: "",
		}
	request.CustomerInformation.Phone1 =
		CustomerInformationField{
			FieldName:  "customer_information[phone1]",
			FieldVal: customer.Phone,
		}
	request.CustomerInformation.Phone2 =
		CustomerInformationField{
			FieldName:  "customer_information[phone2]",
			FieldVal: customer.MobilePhone,
		}
	request.CustomerInformation.Dob =
		CustomerInformationField{
			FieldName:  "customer_information[dob]",
			FieldVal: customer.DateOfBirth.Format("2006-01-02"),
		}

	return request
}

func MapResponseError(responseStatus int, responseBytes []byte) (result error, err error) {
	message := extractErrorFromHtml(string(responseBytes))
	return errors.New(message), nil
}

func extractErrorFromHtml(responseBody string) (result string) {
	result = "Unknown message"

	doc, err := html.Parse(strings.NewReader(string(responseBody)))
	if err != nil {
		return result
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			result = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return result
}

func mapDocumentType(documentType string) int {
	switch documentType {
	case "IDCard":
		return 3
	case "DriverLicense", "DriverLicenseTranslation":
		return 4
	case "ResidencePermit":
		return 12
	default:
		return 1
	}
}

func MapGender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return "M"
	case common.Female:
		return "F"
	default:
		return ""
	}
}

func MapErrorCode(code int) string {
	switch(code) {
	case -1:
		return "Username is not present"
	case -2:
		return "Registration Date is not present"
	case -3:
		return "Registration Date incorrect date format"
	case -4:
		return "Registration IP Address is not present"
	case -5:
		return "The Registration IP Address is not valid"
	case -6:
		return "Registration Device Id is not present"
	case -7:
		return "Session id is not present"
	case -8:
		return "Session id is not present"
	case -9:
		return "No customer information"
	case -10:
		return "Customer first name is not present"
	case -11:
		return "Customer last name is not present"
	case -12:
		return "Customer email is not present"
	case -13:
		return "Wrong customer email format"
	case -14:
		return "Customer address 1 is not present"
	case -15:
		return "Customer city is not present"
	case -16:
		return "Customer province is not present"
	case -17:
		return "Customer country is not present"
	case -18:
		return "Customer phone 1 is not present"
	case -19:
		return "Customer day of birth is not present"
	case -20:
		return "Customer day of birth incorrect date format"
	case -21:
		return "Deposit count is not present"
	case -22:
		return "Deposit count is not a integer number"
	case -23:
		return "Withdrawal count is not present"
	case -24:
		return "Withdrawal count is not a integer number"
	case -25:
		return "No deposit limits information"
	case -26:
		return "Payment method type is not present"
	case -27:
		return "Payment method type is not valid"
	case -28:
		return "No payment method information"
	case -29:
		return "BIN is not present"
	case -30:
		return "BIN is not a integer number"
	case -31:
		return "Last 4 digits is not present"
	case -32:
		return "Last 4 digits is not a integer numberv"
	case -33:
		return "Routing number is not present"
	case -34:
		return "Routing number is not a integer number"
	case -35:
		return "Account number is not present"
	case -36:
		return "Account number is not a integer number"
	case -37:
		return "Ecase -wallet id is not present"
	case -38:
		return "Ecase -wallet id is not a integer number"
	case -39:
		return "Deposit limit min is not present"
	case -40:
		return "Deposit limit min is not a integer number"
	case -41:
		return "Deposit limit daily is not present"
	case -42:
		return "Deposit limit daily is not a integer number"
	case -43:
		return "Deposit limit weekly is not present"
	case -44:
		return "Deposit limit weekly is not a integer number"
	case -45:
		return "Deposit limit monthly is not present"
	case -46:
		return "Deposit limit monthly is not a integer number"
	case -47:
		return "Transaction id is not present"
	case -48:
		return "Transaction already exists"
	case -49:
		return "Amount ID is not present"
	case -50:
		return "Amount is not a integer number"
	case -51:
		return "Currency is not present"
	case -52:
		return "Time is not present"
	case -53:
		return "IP is not present"
	case -54:
		return "The IP Address is not valid"
	case -55:
		return "Device id is not present"
	case -56:
		return "Invalid Session ID"
	case -57:
		return "Invalid Internal Transaction Id"
	case -58:
		return "Merchant id is not present"
	case -59:
		return "Password is not present"
	case -60:
		return "Invalid Merchant Login"
	case -61:
		return "Customer email is in blacklist"
	case -62:
		return "Wrong time format"
	case -63:
		return "Evidence"
	case -64:
		return "User is not allowed"
	case -65:
		return "Associated account with evidence"
	case -66:
		return "Current IP Address is not present"
	case -67:
		return "Current IP Address is not valid"
	case -69:
		return "Process transaction merchant_id error"
	case -70:
		return "Process transaction password error"
	case -71:
		return "Process transaction amount error"
	case -72:
		return "Process transaction currency error"
	case -73:
		return "Process transaction transaction_id error"
	case -74:
		return "Process transaction description error"
	case -75:
		return "Process transaction transaction_type error"
	case -76:
		return "Process transaction card_holder error"
	case -77:
		return "Process transaction card_number error"
	case -78:
		return "Process transaction expiration_date error"
	case -79:
		return "Process transaction expiration_year error"
	case -80:
		return "Process transaction cvv error"
	case -81:
		return "Process transaction address_line_1 error"
	case -82:
		return "Process transaction address_line_2 error"
	case -83:
		return "Process transaction city error"
	case -84:
		return "Process transaction state error"
	case -85:
		return "Process transaction postal_code error"
	case -86:
		return "Process transaction country error"
	case -87:
		return "Process transaction customer_email error"
	case -88:
		return "Process transaction customer_phone"
	case -89:
		return "Process transaction ip_address_error"
	case -90:
		return "Process transaction Duplicated Transaction Id error"
	case -91:
		return "Process transaction processor error"
	case -92:
		return "Process transaction Processor not exist error"
	case -93:
		return "Process transaction first_name error"
	case -94:
		return "Process transaction last_name error"
	case -95:
		return "Processor Error"
	case -96:
		return "Quantity is not a number"
	case -97:
		return "Wrong local_time format"
	case -98:
		return "Wrong billing_email format"
	case -99:
		return "Not Telesign information"
	case -100:
		return "Not Telesign phone"
	case -101:
		return "Not Telesign language"
	case -102:
		return "Not Telesign verify code"
	case -103:
		return "Not Telesign valid language"
	case -104:
		return "Merchant mismatch"
	case -105:
		return "Password mismatch"
	case -106:
		return "Customer Information first_name error"
	case -107:
		return "Customer Information last_name error"
	case -108:
		return "Customer Information email error"
	case -109:
		return "Customer Information address1 error"
	case -110:
		return "Customer Information address2 error"
	case -111:
		return "Customer Information city error"
	case -112:
		return "Customer Information province error"
	case -113:
		return "Customer Information postal_code error"
	case -114:
		return "Customer Information country error"
	case -115:
		return "Customer Information phone1 error"
	case -116:
		return "Customer Information phone2 error"
	case -117:
		return "Customer Information DOB error"
	case -118:
		return "Username error"
	case -119:
		return "Usernumber error"
	case -120:
		return "Registration date error"
	case -121:
		return "Registration IP Address error"
	case -122:
		return "Registration DeviceId error"
	case -123:
		return "Bonus Code error"
	case -124:
		return "Bonus Submission Date error"
	case -125:
		return "Bonus Amounr error"
	case -126:
		return "BonusId error"
	case -127:
		return "Status error"
	case -128:
		return "Website error"
	case -129:
		return "How did you hear error"
	case -130:
		return "AffiliateId error"
	case -131:
		return "Id_Type error"
	case -132:
		return "Id_Value error"
	case -133:
		return "Transaction does not exist"
	case -134:
		return "ReferenceId does not exist"
	case -135:
		return "Code entered does not exists"
	default:
		return "Unknown error"
	}
}



