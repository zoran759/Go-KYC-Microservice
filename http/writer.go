package http

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"
)

// quoteEscaper holds the replacer function for special characters in a string.
var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// escapeQuotes escapes special characters in a string.
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

// CreateFormFile returns a writer that writes data from a file to the form field.
func CreateFormFile(writer *multipart.Writer, fieldname, filename string, contentType string) (io.Writer, error) {
	if writer == nil {
		return nil, errors.New("no writer supplied")
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", contentType)
	return writer.CreatePart(h)
}
