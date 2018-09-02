package verification

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNewService(t *testing.T) {
	config := Config{
		Host:         "host",
		ClientID:     "client_id",
		ClientSecret: "secret",
	}

	testService := service{config: config}

	service := NewService(config)

	assert.Equal(t, testService, service)
}

func Test_service_CreateUser(t *testing.T) {
	service := NewService(Config{
		Host:         "https://uat-api.synapsefi.com/v3.1/users",
		ClientID:     "client_id",
		ClientSecret: "secret",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://uat-api.synapsefi.com/v3.1/users",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
    "_id": "594e0fa2838454002ea317a0",
    "_links": {
        "self": {
            "href": "https://uat-api.synapsefi.com/v3.1/users/594e0fa2838454002ea317a0"
        }
    },
    "client": {
        "id": "589acd9ecb3cd400fa75ac06",
        "name": "SynapseFI"
    },
    "doc_status": {
        "physical_doc": "MISSING|INVALID",
        "virtual_doc": "MISSING|INVALID"
    },
    "documents": [
        {
            "id": "2a4a5957a3a62aaac1a0dd0edcae96ea2cdee688ec6337b20745eed8869e3ac8",
            "name": "Test User",
            "permission_scope": "UNVERIFIED",
            "physical_docs": [
                {
                    "document_type": "GOVT_ID",
                    "id": "c486c2cb8c1bce695fcfae3197e14aa5b8ddec184c2779d00d581abee5d9a04c",
                    "last_updated": 1498288031319,
                    "status": "SUBMITTED|REVIEWING"
                },
{
                    "document_type": "SELFIE",
                    "id": "c486c2cb8c1bce695fcfae3197e14aa5b8ddec184c2779d00d581abee5d9a04c",
                    "last_updated": 1498288031319,
                    "status": "SUBMITTED"
                }
            ]
        }
    ],
    "emails": [],
    "extra": {
        "cip_tag": 1,
        "date_joined": 1498288029784,
        "extra_security": false,
        "is_business": false,
        "last_updated": 1498288029784,
        "public_note": null,
        "supp_id": "122eddfgbeafrfvbbb"
    },
    "is_hidden": false,
    "legal_names": [
        "Test User"
    ],
    "logins": [
        {
            "email": "test@synapsefi.com",
            "scope": "READ_AND_WRITE"
        }
    ],
    "permission": "UNVERIFIED",
    "phone_numbers": [
        "901.111.1111",
        "test@synapsefi.com"
    ],
    "photos": [],
    "refresh_token": "refresh_ehG7YBS8ZiD0sLa6PQHMUxryovVkJzElC5gWROXq"
}`), nil
		},
	)

	response, err := service.CreateUser(CreateUserRequest{})
	if assert.NoError(t, err) {
		assert.Equal(t, "594e0fa2838454002ea317a0", response.ID)
		assert.Equal(t, "MISSING|INVALID", response.DocumentStatus.PhysicalDoc)
		assert.Len(t, response.Documents, 1)
		document := response.Documents[0]
		assert.Len(t, document.PhysicalDocs, 2)
		assert.Equal(t, "GOVT_ID", document.PhysicalDocs[0].DocumentType)
		assert.Equal(t, "SUBMITTED|REVIEWING", document.PhysicalDocs[0].Status)
		assert.Equal(t, "SELFIE", document.PhysicalDocs[1].DocumentType)
		assert.Equal(t, "SUBMITTED", document.PhysicalDocs[1].Status)
	}
}

func Test_service_CreateUserError(t *testing.T) {
	service := NewService(Config{
		Host:         "https://uat-api.synapsefi.com/v3.1/users",
		ClientID:     "client_id",
		ClientSecret: "secret",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://uat-api.synapsefi.com/v3.1/users",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`), nil
		},
	)

	response, err := service.CreateUser(CreateUserRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodPost,
		"https://uat-api.synapsefi.com/v3.1/users",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = service.CreateUser(CreateUserRequest{})
	assert.Error(t, err)
	assert.Nil(t, response)
}

func Test_service_GetUser(t *testing.T) {
	service := NewService(Config{
		Host:         "https://uat-api.synapsefi.com/v3.1/users",
		ClientID:     "client_id",
		ClientSecret: "secret",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://uat-api.synapsefi.com/v3.1/users/id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{
    "_id": "594e0fa2838454002ea317a0",
    "_links": {
        "self": {
            "href": "https://uat-api.synapsefi.com/v3.1/users/594e0fa2838454002ea317a0"
        }
    },
    "client": {
        "id": "589acd9ecb3cd400fa75ac06",
        "name": "SynapseFI"
    },
    "doc_status": {
        "physical_doc": "SUBMITTED|VALID",
        "virtual_doc": "SUBMITTED|VALID"
    },
    "documents": [
        {
            "id": "2a4a5957a3a62aaac1a0dd0edcae96ea2cdee688ec6337b20745eed8869e3ac8",
            "name": "Test User",
            "permission_scope": "SEND|RECEIVE|1000|DAILY",
            "physical_docs": [
                {
                    "document_type": "GOVT_ID",
                    "id": "c486c2cb8c1bce695fcfae3197e14aa5b8ddec184c2779d00d581abee5d9a04c",
                    "last_updated": 1498288034877,
                    "status": "SUBMITTED|VALID"
                },
				{
                    "document_type": "SELFIE",
                    "id": "c486c2cb8c1bce695fcfae3197e14aa5b8ddec184c2779d00d581abee5d9a04c",
                    "last_updated": 1498288034877,
                    "status": "SUBMITTED|UNVALID"
                }
            ]
        }
    ],
    "emails": [],
    "extra": {
        "cip_tag": 1,
        "date_joined": 1498288029784,
        "extra_security": false,
        "is_business": false,
        "last_updated": 1498288034864,
        "public_note": null,
        "supp_id": "122eddfgbeafrfvbbb"
    },
    "is_hidden": false,
    "legal_names": [
        "Test User"
    ],
    "logins": [
        {
            "email": "test@synapsefi.com",
            "scope": "READ_AND_WRITE"
        }
    ],
    "permission": "SEND-AND-RECEIVE",
    "phone_numbers": [
        "test@synapsefi.com",
        "901.111.1111"
    ],
    "photos": [],
    "refresh_token": "refresh_ehG7YBS8ZiD0sLa6PQHMUxryovVkJzElC5gWROXq"
}`), nil
		},
	)

	response, err := service.GetUser("id")
	if assert.NoError(t, err) {
		assert.Equal(t, "594e0fa2838454002ea317a0", response.ID)
		assert.Equal(t, "SUBMITTED|VALID", response.DocumentStatus.PhysicalDoc)

		document := response.Documents[0]

		if assert.Len(t, document.PhysicalDocs, 2) {
			assert.Equal(t, "GOVT_ID", document.PhysicalDocs[0].DocumentType)
			assert.Equal(t, "SUBMITTED|VALID", document.PhysicalDocs[0].Status)
			assert.Equal(t, "SELFIE", document.PhysicalDocs[1].DocumentType)
			assert.Equal(t, "SUBMITTED|UNVALID", document.PhysicalDocs[1].Status)
		}
	}
}

func Test_service_GetUserError(t *testing.T) {
	service := NewService(Config{
		Host:         "https://uat-api.synapsefi.com/v3.1/users",
		ClientID:     "client_id",
		ClientSecret: "secret",
	})

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodGet,
		"https://uat-api.synapsefi.com/v3.1/users/id",
		func(request *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(
				http.StatusOK,
				`{`), nil
		},
	)

	response, err := service.GetUser("id")
	assert.Error(t, err)
	assert.Nil(t, response)

	httpmock.Reset()
	httpmock.RegisterResponder(
		http.MethodGet,
		"https://uat-api.synapsefi.com/v3.1/users/id",
		func(request *http.Request) (*http.Response, error) {
			return nil, errors.New("test_error")
		},
	)

	response, err = service.GetUser("id")
	assert.Error(t, err)
	assert.Nil(t, response)
}
