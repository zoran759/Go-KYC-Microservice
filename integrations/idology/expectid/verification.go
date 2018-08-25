package expectid

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"time"

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
	err = xml.Unmarshal(response, resp)

	return
}

// makeRequestBody returns url-encoded request body.
// It expects a customer data as the param.
func (c *Client) makeRequestBody(customer *common.UserData) string {
	// IDology recommends submitting all of the fields (even blank values)
	// for optimal performance and full change management functionality.

	// FIXME: probably, at this moment some fields have no corresponding value in the common data and some are I don't understand to what values they're related.

	v := url.Values{}

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
	} else {
		v.Set("city", "")
	}
	//  Conditional. State(2). City and State required if enabled.
	if len(customer.CurrentAddress.StateProvinceCode) > 0 {
		v.Set("state", customer.CurrentAddress.StateProvinceCode)
	} else {
		v.Set("state", "")
	}
	// Conditional. 5-digit zip code (5). Zip Code required if enabled.
	if len(customer.CurrentAddress.PostCode) == 5 {
		v.Set("zip", customer.CurrentAddress.PostCode)
	} else {
		v.Set("zip", "")
	}
	// Optional. Your invoice or order number.
	v.Set("invoice", "")
	// Optional. Order amount.
	v.Set("amount", "")
	// Optional. Shipping amount.
	v.Set("shipping", "")
	// Optional. Tax amount.
	v.Set("tax", "")
	// Optional. Total amount(sum of the above).
	v.Set("total", "")
	// Optional. Type of ID provided.
	v.Set("idType", "")
	// Optional. Issuing agency of ID.
	v.Set("idIssuer", "")
	// Optional. Number on ID.
	v.Set("idNumber", "")
	// Optional. Payment method.
	v.Set("paymentMethod", "")
	// "ssnLast4" - Optional. Last 4 digits of SSN (4). Results improve with the addition of this field.
	v.Set("ssnLast4", "")
	// "ssn" - Optional. Full ssn (9).
	v.Set("ssn", "")
	for _, d := range customer.Documents {
		if d.Metadata.Type == common.IDCard {
			ssnLast4 := d.Metadata.Number[len(d.Metadata.Number)-4:]
			if len(ssnLast4) == 4 {
				v.Set("ssnLast4", ssnLast4)
			}
			v.Set("ssn", d.Metadata.Number)

			break
		}
	}
	// Optional. Month of Birth (2). Results improve with the addition of this field.
	// Optional. Year of Birth (4). Results improve with the addition of this field. YOB is the minimum DOB information accepted by IDology.
	if !time.Time(customer.DateOfBirth).IsZero() {
		v.Set("dobMonth",
			fmt.Sprintf("%2d", time.Time(customer.DateOfBirth).Month()),
		)
		v.Set("dobYear", fmt.Sprintf("%d", time.Time(customer.DateOfBirth).Year()))
	} else {
		v.Set("dobMonth", "")
		v.Set("dobYear", "")
	}
	// Optional. IP Address . Include periods in the address,for example - 11.111.111.11
	v.Set("ipAddress", "")
	// Optional. Email address.
	if len(customer.Email) > 0 {
		v.Set("email", customer.Email)
	} else {
		v.Set("email", "")
	}
	// Optional. Phone number (10).
	if len(customer.Phone) == 10 {
		v.Set("telephone", customer.Phone)
	} else if len(customer.MobilePhone) == 10 {
		v.Set("telephone", customer.MobilePhone)
	} else {
		v.Set("telephone", "")
	}
	// Optional. SKU.
	v.Set("sku", "")
	// Optional. User ID (External application).
	v.Set("uid", "")
	// Optional. Alternate street address. Submit a secondary address for verification, such as a shipping address.
	v.Set("altAddress", "")
	// Optional. Alternate city.
	v.Set("altCity", "")
	// Optional. Alternate state(2) State in abbreviated format. i.e. for Georgia send "GA".
	v.Set("altState", "")
	// Optional. Alternate 5-digit zip code (5).
	v.Set("altZip", "")
	// Optional. Card purchase date.
	v.Set("purchaseDate", "")
	// Optional. This <id-number>should be submitted if data is being sent from the ExpectID Scan Onboard product. This field should be populated with the <id-number> provided in the API Response from the Scan Onboard polling service. This will couple calls to the ExpectID IQ API service to results from Scan Onboard.
	v.Set("captureQueryId", "")
	// Optional. score.
	v.Set("score", "")
	// Optional. Custom field 1.
	v.Set("c_custom_field_1", "")
	// Optional. Custom field 2.
	v.Set("c_custom_field_2", "")
	// Optional. Custom field 3.
	v.Set("c_custom_field_3", "")

	return v.Encode()
}
