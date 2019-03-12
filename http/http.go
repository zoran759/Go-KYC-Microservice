package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

// defaultHTTPTimeout holds the default value for a HTTP request timeout.
// Currently, the timeout value for requests isn't configurable.
var defaultHTTPTimeout = time.Minute * 3

// Headers represents a HTTP request headers.
type Headers map[string]string

// Post sends a HTTP POST request to the endpoint using the specified headers and body.
// It returns a HTTP status code or zero in the case of an error, a response body or an error if occurred.
func Post(endpoint string, headers Headers, body []byte) (int, []byte, error) {
	return Request(http.MethodPost, endpoint, headers, body)
}

// Get sends a HTTP GET request to the endpoint using the specified headers.
func Get(endpoint string, headers Headers) (int, []byte, error) {
	return Request(http.MethodGet, endpoint, headers, []byte{})
}

// Patch sends a HTTP PATCH request to the endpoint using the specified headers and body.
// It returns a HTTP status code or zero in the case of an error, a response body or an error if occurred.
func Patch(endpoint string, headers Headers, body []byte) (int, []byte, error) {
	return Request(http.MethodPatch, endpoint, headers, body)
}

// Request sends a HTTP request to the endpoint using the specified method and headers.
// The body will be used as the request body.
func Request(method string, endpoint string, headers Headers, body []byte) (int, []byte, error) {
	request, err := http.NewRequest(method, endpoint, bytes.NewReader(body))

	if err != nil {
		return 0, nil, err
	}

	for header, value := range headers {
		request.Header.Set(header, value)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultHTTPTimeout)
	defer cancel()

	response, err := http.DefaultClient.Do(request.WithContext(ctx))

	if err != nil {
		return 0, nil, err
	}

	return extractCodeAndBodyFromResponse(response)
}

// extractCodeAndBodyFromResponse extracts the content from a HTTP response.
// It returns a HTTP status code, a response body content or an error if occurred.
func extractCodeAndBodyFromResponse(response *http.Response) (int, []byte, error) {
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return 0, nil, err
	}

	return response.StatusCode, responseBody, nil
}
