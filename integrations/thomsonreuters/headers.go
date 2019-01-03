package thomsonreuters

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"time"

	"modulus/kyc/http"
)

const (
	content = "application/json"

	mGET  httpMethod = "get "
	mHEAD httpMethod = "head "
	mPOST httpMethod = "post "
)

type httpMethod string

// createHeaders creates HTTP headers required to perform request.
func (tomson ThomsonReuters) createHeaders(method httpMethod, endpoint string, payload []byte) http.Headers {
	date := time.Now().Format(time.RFC3339)

	dataToSign := bytes.Buffer{}

	dataToSign.WriteString("(request-target): ")
	dataToSign.WriteString(string(method))
	dataToSign.WriteString(tomson.path)
	dataToSign.WriteString(endpoint)
	dataToSign.WriteByte('\n')
	dataToSign.WriteString("host: ")
	dataToSign.WriteString(tomson.host)
	dataToSign.WriteByte('\n')
	dataToSign.WriteString("date: ")
	dataToSign.WriteString(date)

	if method == mPOST {
		dataToSign.WriteByte('\n')
		dataToSign.WriteString("content-type: ")
		dataToSign.WriteString(content)
		dataToSign.WriteByte('\n')
		dataToSign.WriteString("content-length: ")
		dataToSign.WriteString(fmt.Sprintf("%d", len(payload)))

		if len(payload) > 0 {
			dataToSign.WriteByte('\n')
			dataToSign.Write(payload)
		}
	}

	// TODO: remove after debug.
	log.Print("Data to sign:\n\n", dataToSign.String(), "\n\n")

	mac := hmac.New(sha256.New, []byte(tomson.secret))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(dataToSign.Bytes()))

	aheader := strings.Builder{}
	aheader.WriteString(`Signature keyId="`)
	aheader.WriteString(tomson.key)
	aheader.WriteString(`",algorithm="hmac-sha256",headers="(request-target) host date`)

	if method == mPOST {
		aheader.WriteString(" content-type content-length")
	}

	aheader.WriteString(`",signature="`)
	aheader.WriteString(signature)
	aheader.WriteByte('"')

	headers := http.Headers{
		"Date":          date,
		"Authorization": aheader.String(),
	}

	if method == mPOST {
		headers["Content-Type"] = content
		headers["Content-Length"] = fmt.Sprintf("%d", len(payload))
	}

	// TODO: remove after debug.
	log.Println("Signature:", signature)
	log.Println("Authorization header:", aheader.String())

	return headers
}
