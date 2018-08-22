package http

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"testing"
)

func Test_escapeQuotes(t *testing.T) {
	testString := `"test"`

	escapedString := escapeQuotes(testString)

	assert.Equal(t, "\\\"test\\\"", escapedString)
}

func TestCreateFormFile(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := CreateFormFile(writer, "file", "filename.png", "image/png")
	assert.NoError(t, err)
	assert.NotNil(t, part)

	part, err = CreateFormFile(nil, "", "", "")
	assert.Error(t, err)
	assert.Nil(t, part)
}
