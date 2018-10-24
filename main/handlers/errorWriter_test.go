package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServiceError(t *testing.T) {
	e := serviceError{
		message: "Test error message",
	}

	assert.Equal(t, "Test error message", e.Error())
}

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	writeErrorResponse(w, http.StatusBadRequest, errors.New("Test error"))

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, `{"Error":"Test error"}`, string(body))

	w = httptest.NewRecorder()

	writeErrorResponse(w, http.StatusInternalServerError, errors.New(""))

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, `{"Error":""}`, string(body))
}
