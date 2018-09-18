package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

func TestPost(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"http://testpost.io/ping",
		func(request *http.Request) (*http.Response, error) {
			defer request.Body.Close()

			responseBody, err := ioutil.ReadAll(request.Body)

			if err != nil {
				return httpmock.NewStringResponse(http.StatusInternalServerError, "failure"), nil
			}

			resp := httpmock.NewBytesResponse(200, responseBody)

			return resp, nil
		},
	)
	httpmock.RegisterResponder(
		http.MethodPost,
		"http://testpost.io/failure",
		httpmock.NewErrorResponder(errors.New("failure")),
	)

	mockBody := `{"test": "test"}`

	mockBodyBytes, _ := json.Marshal(mockBody)

	status, responseBody, err := Post("http://testpost.io/ping", Headers{"Content-Type": "application/json"}, mockBodyBytes)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, mockBodyBytes, responseBody)
	}

	status, responseBody, err = Post("http://testpost.io/failure", Headers{"Content-Type": "application/json"}, mockBodyBytes)

	if assert.Error(t, err) {
		assert.Equal(t, 0, status)
		assert.Nil(t, responseBody)
		assert.Equal(t, "Post http://testpost.io/failure: failure", err.Error())
	}
}

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"http://testpost.io/ping",
		func(request *http.Request) (*http.Response, error) {

			resp := httpmock.NewStringResponse(200, "GetRequest")

			return resp, nil
		},
	)
	httpmock.RegisterResponder(
		http.MethodGet,
		"http://testpost.io/failure",
		httpmock.NewErrorResponder(errors.New("failure")),
	)

	status, responseBody, err := Get("http://testpost.io/ping", Headers{})

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, []byte("GetRequest"), responseBody)
	}

	status, responseBody, err = Get("http://testpost.io/failure", Headers{})

	if assert.Error(t, err) {
		assert.Equal(t, 0, status)
		assert.Nil(t, responseBody)
		assert.Equal(t, "Get http://testpost.io/failure: failure", err.Error())
	}
}

func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Second * 70)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	_, _, err := Get(ts.URL, Headers{})

	assert.Error(t, err)
}
