package jumio

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	stdhttp "net/http"
	"time"

	"modulus/kyc/common"
	"modulus/kyc/http"
)

// Verification represents the verification interface of Jumio performNetverify API.
type Verification interface {
	common.CustomerChecker
	common.StatusChecker
}

// service defines the model for the Jumio performNetverify API.
type service struct {
	baseURL     string
	credentials string
}

// New constructs new service object to use with the Jumio performNetverify API.
func New(config Config) Verification {
	return &service{
		baseURL:     config.BaseURL,
		credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
	}
}

// CheckCustomer implements customer verification using the Jumio performNetverify API.
func (s *service) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	req := &Request{}
	if err = req.populateFields(customer); err != nil {
		err = fmt.Errorf("invalid customer data: %s", err)
		return
	}

	response, errorCode, err := s.sendRequest(req)
	if err != nil {
		if errorCode != nil {
			result.ErrorCode = fmt.Sprintf("%d", *errorCode)
		}
		err = fmt.Errorf("during sending request: %s", err)
		return
	}

	result.Status = common.Unclear
	result.StatusCheck = &common.KYCStatusCheck{
		Provider:    common.Jumio,
		ReferenceID: response.JumioIDScanReference,
		LastCheck:   time.Now(),
	}

	return
}

// sendRequest sends a vefirication request into the API.
// It returns a response from the API or the error if occured.
func (s *service) sendRequest(request *Request) (response *Response, errorCode *int, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}

	headers := s.headers()
	headers["Content-Type"] = contentType
	headers["Content-Length"] = fmt.Sprintf("%d", len(body))

	statusCode, resp, err := http.Post(s.baseURL+performNetverifyEndpoint, headers, body)
	if err != nil {
		return
	}
	if statusCode != stdhttp.StatusOK {
		errorCode = &statusCode
		err = errors.New("http error")
		return
	}

	response = &Response{}
	err = json.Unmarshal(resp, response)

	return
}

// CheckStatus implements the StatusChecker interface for Jumio.
func (s *service) CheckStatus(referenceID string) (result common.KYCResult, err error) {
	if len(referenceID) == 0 {
		err = errors.New("empty Jumio’s reference number of an existing scan")
		return
	}

	status, errorCode, err := s.retrieveScanStatus(referenceID)
	if err != nil {
		if errorCode != nil {
			result.ErrorCode = fmt.Sprintf("%d", *errorCode)
		}
		err = fmt.Errorf("during sending request: %s", err)
		return
	}

	switch status {
	case PendingStatus:
		result.Status = common.Unclear
		result.StatusCheck = &common.KYCStatusCheck{
			Provider:    common.Jumio,
			ReferenceID: referenceID,
			LastCheck:   time.Now(),
		}
	case DoneStatus, FailedStatus:
		scanDetails := &DetailsResponse{}
		scanDetails, errorCode, err = s.retrieveScanDetails(referenceID)
		if err != nil {
			if errorCode != nil {
				result.ErrorCode = fmt.Sprintf("%d", *errorCode)
			}
			err = fmt.Errorf("during sending request: %s", err)
			return
		}
		result, err = scanDetails.toResult()
	default:
		err = fmt.Errorf("unknown status of the verification: %s", status)
	}

	return
}

// retrieveScanStatus retrieves the status of an Jumio scan.
func (s *service) retrieveScanStatus(referenceID string) (status ScanStatus, errorCode *int, err error) {
	statusCode, resp, err := http.Get(s.baseURL+scanStatusEndpoint+referenceID, s.headers())
	if err != nil {
		return
	}
	if statusCode != stdhttp.StatusOK {
		errorCode = &statusCode
		err = errors.New("http error")
		return
	}

	response := &StatusResponse{}
	err = json.Unmarshal(resp, response)
	if err != nil {
		return
	}

	status = response.Status

	return
}

// retrieveScanDetails retrieves details of an Jumio scan.
func (s *service) retrieveScanDetails(referenceID string) (response *DetailsResponse, errorCode *int, err error) {
	statusCode, resp, err := http.Get(fmt.Sprintf(s.baseURL+scanDetailsEndpoint, referenceID), s.headers())
	if err != nil {
		return
	}
	if statusCode != stdhttp.StatusOK {
		errorCode = &statusCode
		err = errors.New("http error")
		return
	}

	response = &DetailsResponse{}
	err = json.Unmarshal(resp, response)

	return
}

// headers is a helper that constructs HTTP request headers.
func (s *service) headers() http.Headers {
	return http.Headers{
		"Accept":        accept,
		"User-Agent":    userAgent,
		"Authorization": s.credentials,
	}
}
