package verification

import (
	"encoding/json"
	"log"
	"modulus/kyc/http"
	"net/url"

	"github.com/gofrs/uuid"
)

type service struct {
	config Config
}

// NewService constructs a new verification service object.
func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

func (service service) Verify(request Request) (*Response, error) {
	customerInformationBytes, err := json.Marshal(request.CustomerInformation)
	if err != nil {
		return nil, err
	}

	reference, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("merchant_id", service.config.MerchantID)
	form.Add("password", service.config.Password)
	form.Add("reg_ip_address", request.RegIPAddress)
	form.Add("reg_date", request.RegDate)
	form.Add("user_number", reference.String())
	form.Add("customer_information", string(customerInformationBytes))

	code, responseBytes, err := http.Post(
		service.config.Host,
		http.Headers{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		[]byte(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	log.Println(form.Encode())
	log.Println(string(responseBytes))
	log.Println(code)
	response := new(Response)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	return response, nil
}
