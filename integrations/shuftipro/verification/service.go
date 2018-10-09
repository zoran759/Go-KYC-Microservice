package verification

import (
	"crypto/sha256"
	"encoding/hex"
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
			reference.String() +
			string(dataBytes) +
			string(servicesBytes) +
			service.config.SecretKey

	hasher := sha256.New()
	hasher.Write([]byte(signatureString))
	hashedSignatureString := hasher.Sum(nil)

	form := url.Values{}
	form.Add("client_id", service.config.ClientID)
	form.Add("country", request.Country)
	form.Add("email", request.Email)
	form.Add("phone_number", request.PhoneNumber)
	form.Add("redirect_url", service.config.RedirectURL)
	form.Add("reference", reference.String())
	form.Add("verification_data", string(dataBytes))
	form.Add("verification_services", string(servicesBytes))
	form.Add("signature", hex.EncodeToString(hashedSignatureString[:]))

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
