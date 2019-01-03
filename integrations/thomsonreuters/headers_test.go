package thomsonreuters

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateHeaders(t *testing.T) {
	tomson := New(Config{
		Host:      "https://rms-world-check-one-api-pilot.thomsonreuters.com/v1/",
		APIkey:    "key",
		APIsecret: "secret",
	})

	endpoint := "groups"

	// Test with GET and no payload.
	headers := tomson.createHeaders(mGET, endpoint, nil)

	assert := assert.New(t)

	assert.Len(headers, 2)
	assert.Contains(headers, "Date")
	assert.Contains(headers, "Authorization")

	date, err := time.Parse(time.RFC3339, headers["Date"])

	assert.NoError(err)
	assert.NotZero(date)

	dataToSign := "(request-target): " + string(mGET) + tomson.path + endpoint + "\nhost: " + tomson.host + "\ndate: " + headers["Date"]
	mac := hmac.New(sha256.New, []byte(tomson.secret))
	signature := base64.StdEncoding.EncodeToString(mac.Sum([]byte(dataToSign)))
	aheader := `Signature keyId="` + tomson.key + `",algorithm="hmac-sha256",headers="(request-target) host date",signature="` + signature + `"`

	assert.Equal(aheader, headers["Authorization"])

	// Test with POST and a payload.
	endpoint = "cases/screeningRequest"
	payload := []byte(`{"name": "John Doe", "providerTypes": ["WATCHLIST"]}`)

	headers = tomson.createHeaders(mPOST, endpoint, payload)

	assert.Len(headers, 4)
	assert.Contains(headers, "Date")
	assert.Contains(headers, "Authorization")
	assert.Contains(headers, "Content-Type")
	assert.Contains(headers, "Content-Length")

	date, err = time.Parse(time.RFC3339, headers["Date"])

	assert.NoError(err)
	assert.NotZero(date)

	dataToSign = "(request-target): " + string(mPOST) + tomson.path + endpoint + "\nhost: " + tomson.host + "\ndate: " + headers["Date"] + "\ncontent-type: " + content + "\ncontent-length: " + fmt.Sprintf("%d\n", len(payload)) + string(payload)
	signature = base64.StdEncoding.EncodeToString(mac.Sum([]byte(dataToSign)))
	aheader = `Signature keyId="` + tomson.key + `",algorithm="hmac-sha256",headers="(request-target) host date content-type content-length",signature="` + signature + `"`

	assert.Equal(aheader, headers["Authorization"])
	assert.Equal(content, headers["Content-Type"])
	assert.Equal(fmt.Sprintf("%d", len(payload)), headers["Content-Length"])

}
