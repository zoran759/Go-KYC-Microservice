package expectid

import (
	"fmt"
	"net/url"
	"time"

	"github.com/achiku/xml"
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/http"
)

// verify sends a vefirication request into the API.
// It expects an url-encoded request body as the param.
// It returns a response from the API or the error if occured.
func (c *Client) verify(requestBody string) (resp *Response, err error) {
	headers := http.Headers{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	_, response, err := http.Post(c.config.Host, headers, []byte(requestBody))
	if err != nil {
		return
	}

	resp = &Response{}
	if err = xml.Unmarshal(response, resp); err != nil {
		return
	}

	return
}

// makeRequestBody returns url-encoded request body.
// It expects a customer data as the param.
func (c *Client) makeRequestBody(customer *common.UserData) string {
	v := url.Values{}

	// FIXME: probably, at this moment some fields have no corresponding value in the common data and some are I don't understand to what values they're related.

	// Required. IDology API username (128 bytes).
	v.Set("username", c.config.Username)
	// Required. IDology API password (255 bytes).
	v.Set("password", c.config.Password)
	// Required. First Name.
	v.Set("firstName", customer.FirstName)
	// Required. Last Name.
	v.Set("lastName", customer.LastName)
	// Required. Street address.
	v.Set("address", customer.AddressString)

	// Conditional. City. City and State required if enabled.
	if len(customer.CurrentAddress.Town) > 0 {
		v.Set("city", customer.CurrentAddress.Town)
	}
	//  Conditional. State(2). City and State required if enabled.
	if len(customer.CurrentAddress.State) > 0 {
		v.Set("state", customer.CurrentAddress.State)
	}
	// Conditional. 5-digit zip code (5). Zip Code required if enabled.
	if len(customer.CurrentAddress.PostCode) == 5 {
		v.Set("zip", customer.CurrentAddress.PostCode)
	}

	// "invoice" - Optional. Your invoice or order number.
	// "amount" - Optional. Order amount.
	// "shipping"  - Optional. Shipping amount.
	// "tax" - Optional. Tax amount.
	// "total" - Optional. Total amount(sum of the above).

	// "idType" - Optional. Type of ID provided.
	// "idIssuer" - Optional. Issuing agency of ID.
	// "idNumber" - Optional. Number on ID.
	// "paymentMethod" - Optional. Payment method.

	// "ssnLast4" - Optional. Last 4 digits of SSN (4) . Results improve with the addition of this field.
	// "ssn" - Optional. Full ssn (9).

	// Optional. Month of Birth (2). Results improve with the addition of this field.
	v.Set("dobMonth", time.Time(customer.DateOfBirth).Month().String())
	// Optional. Year of Birth (4). Results improve with the addition of this field. YOB is the minimum DOB information accepted by IDology.
	v.Set("dobYear", fmt.Sprintf("%d", time.Time(customer.DateOfBirth).Year()))

	// "ipAddress" - Optional. IP Address . Include periods in the address,for example - 11.111.111.11

	// Optional. Email address.
	if len(customer.Email) > 0 {
		v.Set("email", customer.Email)
	}

	// Optional. Phone number (10).
	if len(customer.Phone) == 10 {
		v.Set("telephone", customer.Phone)
	} else if len(customer.MobilePhone) == 10 {
		v.Set("telephone", customer.Phone)
	}

	// "sku" - Optional. SKU.
	// "uid" - Optional. User ID (External application).
	// "altAddress" - Optional. Alternate street address. Submit a secondary address for verification, such as a shipping address.
	// "altCity" - Optional. Alternate city.
	// "altState" - Optional. Alternate state(2) State in abbreviated format. i.e. for Georgia send "GA".
	// "altZip" - Optional. Alternate 5-digit zip code (5).
	// "purchaseDate" - Optional. Card purchase date.
	// "captureQueryId" - Optional. This <id-number>should be submitted if data is being sent from the ExpectID Scan Onboard product. This field should be populated with the <id-number> provided in the API Response from the Scan Onboard polling service. This will couple calls to the ExpectID IQ API service to results from Scan Onboard.
	// "score" - Optional. score.
	// "c_custom_field_1" - Optional. Custom field 1.
	// "c_custom_field_2" - Optional. Custom field 2.
	// "c_custom_field_3" - Optional. Custom field 3.

	return v.Encode()
}
