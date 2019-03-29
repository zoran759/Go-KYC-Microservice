package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"

	"github.com/stretchr/testify/assert"
)

func TestConfigHandler(t *testing.T) {
	endpoint := "/Config"

	defer os.Remove(config.DefaultDevFile)

	type testCase struct {
		name     string
		method   string
		body     []byte
		response string
		status   int
	}

	testCases := []testCase{
		testCase{
			name:     "Get config",
			method:   http.MethodGet,
			response: `{"IDology":{"Host":"https://web.idologylive.com/api/idiq.svc","Password":"fakepassword","UseSummaryResult":"false","Username":"fakeuser"},"IdentityMind":{"Host":"https://sandbox.identitymind.com/im","Password":"fakepassword","Username":"fakeuser"},"ShuftiPro":{"CallbackURL":"https://api.shuftipro.com","ClientID":"fakeID","Host":"https://api.shuftipro.com","SecretKey":"fakeKey"},"Sum\u0026Substance":{"APIKey":"fakeKey","Host":"https://test-api.sumsub.com"},"Trulioo":{"Host":"https://api.globaldatacompany.com","NAPILogin":"fakelogin","NAPIPassword":"fakepassword"}}` + "\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Update config",
			method:   http.MethodPost,
			body:     []byte(`{"Jumio":{"BaseURL":"https://lon.netverify.com/api/netverify/v2","Token":"4c620e55-3a97-e31e-719f-339b5c9d4b7c","Secret":"RBGPMico1YWETwqLKVKsFGH0JSE7aQfI"}}`),
			response: "{\"Updated\":true,\"Errors\":[]}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Update config with errors",
			method:   http.MethodPost,
			body:     []byte(`{"ComplyAdvantage":{"Host":"https://api.complyadvantage.com","APIkey":"uihgewnwe78jbIUGFIkbssufKUBkzbjkkb"}}`),
			response: "{\"Updated\":true,\"Errors\":[\"missing or empty option 'Fuzziness' for the ComplyAdvantage\"]}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Update empty config",
			method:   http.MethodPost,
			body:     []byte(`{}`),
			response: "{\"Updated\":false,\"Errors\":[\"no config update data provided\"]}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Empty update request",
			method:   http.MethodPost,
			body:     []byte{},
			response: "{\"Error\":\"empty request\"}",
			status:   http.StatusBadRequest,
		},
		testCase{
			name:     "Update wrong config",
			method:   http.MethodPost,
			body:     []byte(`{"Config":"This config format is invalid"}`),
			response: "{\"Error\":\"json: cannot unmarshal string into Go value of type config.Options\"}",
			status:   http.StatusBadRequest,
		},
		testCase{
			name:     "Update config with unknown option",
			method:   http.MethodPost,
			body:     []byte(`{"Jumio":{"BaseURL":"https://lon.netverify.com/api/netverify/v2","Token":"4c620e55-3a97-e31e-719f-339b5c9d4b7c","Secret":"RBGPMico1YWETwqLKVKsFGH0JSE7aQfI","Foobar":"foobar value"}}`),
			response: "{\"Updated\":true,\"Errors\":[\"Jumio: unknown option 'Foobar' was filtered out\"]}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Wrong method",
			method:   http.MethodPatch,
			response: "{\"Error\":\"used method not allowed for this endpoint\"}",
			status:   http.StatusMethodNotAllowed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, endpoint, bytes.NewReader(tc.body))
			w := httptest.NewRecorder()
			handlers.ConfigHandler(w, req)
			assert.Equal(t, tc.response, w.Body.String())
			assert.Equal(t, tc.status, w.Code)
		})
	}
}
