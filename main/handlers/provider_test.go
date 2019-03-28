package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsProviderImplemented(t *testing.T) {
	endpoint := "/Provider"

	list, _ := json.Marshal(providerList())
	list = append(list, '\n')

	type testCase struct {
		name     string
		query    string
		response string
		status   int
	}

	testCases := []testCase{
		testCase{
			name:     "Invalid query",
			query:    "?error=%%",
			response: `{"Error":"invalid URL escape \"%%\""}`,
			status:   http.StatusBadRequest,
		},
		testCase{
			name:     "Missing name",
			query:    "?test",
			response: `{"Error":"missing provider name in the request"}`,
			status:   http.StatusBadRequest,
		},
		testCase{
			name:     "Known name",
			query:    "?name=ShuftiPro",
			response: "{\"Implemented\":true}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Unknown name",
			query:    "?name=fake",
			response: "{\"Implemented\":false}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Non provider",
			query:    "?name=CipherTrace",
			response: "{\"Implemented\":true}\n",
			status:   http.StatusOK,
		},
		testCase{
			name:     "Provider list",
			response: string(list),
			status:   http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, endpoint+tc.query, nil)
			w := httptest.NewRecorder()
			IsProviderImplemented(w, req)
			assert.Equal(t, tc.response, w.Body.String())
			assert.Equal(t, tc.status, w.Code)
		})
	}
}
