package jumio

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/http"
)

// service defines the model for the Jumio performNetverify API.
type service struct {
	baseURL     string
	credentials string
}

// New constructs new service object to use with the Jumio performNetverify API.
func New(config Config) common.CustomerChecker {
	return &service{
		baseURL:     config.BaseURL,
		credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Token+":"+config.Secret)),
	}
}

// CheckCustomer implements customer verification using the Jumio performNetverify API.
func (s *service) CheckCustomer(customer *common.UserData) (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	result = common.Error

	if customer == nil {
		err = errors.New("no customer supplied")
		return
	}

	req := &Request{}
	if err = req.populateFields(customer); err != nil {
		err = fmt.Errorf("invalid customer data: %s", err)
		return
	}

	response, err := s.sendRequest(req)
	if err != nil {
		err = fmt.Errorf("during sending request: %s", err)
		return
	}

	scanDetails, err := s.retrieveDetails(response.JumioIDScanReference)
	if err != nil {
		err = fmt.Errorf("during retrieving result: %s", err)
		return
	}

	result, details, err = scanDetails.toResult()

	return
}

// sendRequest sends a vefirication request into the API.
// It returns a response from the API or the error if occured.
func (s *service) sendRequest(request *Request) (response *Response, err error) {
	body, err := json.Marshal(request)
	if err != nil {
		return
	}

	headers := http.Headers{
		"Accept":         accept,
		"Content-Type":   contentType,
		"Content-Length": fmt.Sprintf("%d", len(body)),
		"User-Agent":     userAgent,
		"Authorization":  s.credentials,
	}

	_, resp, err := http.Post(s.baseURL+performNetverifyEndpoint, headers, body)
	if err != nil {
		return
	}

	response = &Response{}
	err = json.Unmarshal(resp, response)

	return
}

func (s *service) retrieveDetails(scanref string) (response *DetailsResponse, err error) {
	if len(scanref) == 0 {
		err = errors.New("empty Jumio’s reference number of an existing scan")
		return
	}

	_, err = s.retrieveScanStatus(scanref)
	if err != nil {
		return
	}

	response, err = s.retrieveScanDetails(scanref)

	return
}

// retrieveScanStatus retrieves the status of an Jumio scan.
func (s *service) retrieveScanStatus(scanref string) (status ScanStatus, err error) {
	headers := http.Headers{
		"Accept":        accept,
		"User-Agent":    userAgent,
		"Authorization": s.credentials,
	}
	endpoint := s.baseURL + scanStatusEndpoint + scanref

	for _, d := range timings {
		time.Sleep(d)

		_, resp, err1 := http.Get(endpoint, headers)
		if err1 != nil {
			err = err1
			return
		}

		response := &StatusResponse{}
		err1 = json.Unmarshal(resp, response)
		if err1 != nil {
			err = err1
			return
		}

		if response.Status != PendingStatus {
			status = response.Status
			return
		}
	}

	err = errors.New("Jumio scan status is 'pending' after 10 allowed attempts - the retrieval aborted")

	return
}

// retrieveScanDetails retrieves details of an Jumio scan.
func (s *service) retrieveScanDetails(scanref string) (response *DetailsResponse, err error) {
	headers := http.Headers{
		"Accept":        accept,
		"User-Agent":    userAgent,
		"Authorization": s.credentials,
	}

	_, resp, err := http.Get(fmt.Sprintf(s.baseURL+scanDetailsEndpoint, scanref), headers)
	if err != nil {
		return
	}

	response = &DetailsResponse{}
	err = json.Unmarshal(resp, response)

	return
}
