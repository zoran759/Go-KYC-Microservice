package verification

import (
	"gitlab.com/modulusglobal/kyc/http"
	"crypto/sha256"
	"encoding/json"
	"github.com/gofrs/uuid"
	"net/url"
)

type service struct {
	config Config
}

func NewService(config Config) Verification {
	return service{
		config: config,
	}
}

func (service service) Verify(request Request) (*Response, error) {
	servicesBytes, err := json.Marshal(request.VerificationServices)
	if err != nil {
		return nil, err
	}

	dataBytes, err := json.Marshal(request.VerificationData)
	if err != nil {
		return nil, err
	}

	reference, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	signatureString :=
		service.config.ClientID +
			request.Country +
			request.Email +
			request.PhoneNumber +
			service.config.RedirectURL +
			string(dataBytes) +
			string(servicesBytes) +
			reference.String() +
			service.config.SecretKey

	hashedSignatureString := sha256.Sum256([]byte(signatureString))

	form := url.Values{}

	form.Add("client_id", service.config.ClientID)
	form.Add("reference", reference.String())
	form.Add("email", request.Email)
	form.Add("phone_number", request.PhoneNumber)
	form.Add("country", request.Country)
	form.Add("redirect_url", service.config.RedirectURL)
	form.Add("verification_services", string(servicesBytes))
	form.Add("verification_data", string(dataBytes))
	form.Add("signature", string(hashedSignatureString[:]))

	_, responseBytes, err := http.Post(
		service.config.Host,
		http.Headers{
			"Content-Type": "x-www-form-urlencoded",
		},
		[]byte(form.Encode()),
	)
	if err != nil {
		return nil, err
	}

	response := new(Response)
	if err := json.Unmarshal(responseBytes, response); err != nil {
		return nil, err
	}

	return response, nil
}
