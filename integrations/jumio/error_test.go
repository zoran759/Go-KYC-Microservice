package jumio_test

import (
	"modulus/kyc/integrations/jumio"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	assert := assert.New(t)

	err := jumio.ErrorResponse{
		Message:    "Bad Request: frontsideImage invalid",
		HTTPStatus: "400",
		RequestURI: "https://lon.netverify.com/api/netverify/v2/performNetverify",
	}

	assert.Equal("HTTP status: 400 | Bad Request: frontsideImage invalid", err.Error())
}
