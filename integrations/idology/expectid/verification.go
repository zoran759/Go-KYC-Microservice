package expectid

import (
	"net/url"

	"github.com/achiku/xml"
	"gitlab.com/modulusglobal/kyc/common"
	"gitlab.com/modulusglobal/kyc/http"
)

func (c *client) verify(requestBody string) (resp *Response, err error) {
	// TODO: implement this.
	headers := http.Headers{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	status, response, err := http.Post(c.Host, headers, []byte(requestBody))
	if err != nil {
		return
	}
	// TODO: status check.
	_ = status

	resp = &Response{}
	if err = xml.Unmarshal(response, resp); err != nil {
		return
	}

	return
}

func (c *client) makeRequestBody(customer *common.UserData) string {
	// TODO: implement this.
	values := url.Values{}

	return values.Encode()
}
