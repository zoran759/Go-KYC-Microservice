package complyadvantage

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"modulus/kyc/common"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var approvedResp = []byte(`
{
    "code": 200,
    "status": "success",
    "content": {
        "data": {
            "id": 93650279,
            "ref": "1544882792-BH4jVx8b",
            "searcher_id": 3970,
            "assignee_id": 3970,
            "filters": {
                "birth_year": 1896,
                "country_codes": [],
                "remove_deceased": 0,
                "types": [
                    "sanction",
                    "warning",
                    "fitness-probity"
                ],
                "exact_match": true,
                "fuzziness": 0
            },
            "match_status": "no_match",
            "risk_level": "unknown",
            "search_term": "Georgy Konstantinovich Zhukov",
            "submitted_term": "Georgy Konstantinovich Zhukov",
            "client_ref": null,
            "total_hits": 0,
            "updated_at": "2018-12-15 14:06:32",
            "created_at": "2018-12-15 14:06:32",
            "tags": [],
            "limit": 100,
            "offset": 0,
            "searcher": {
                "id": 3970,
                "email": "john@capdax.com",
                "name": "John",
                "phone": "",
                "created_at": "2018-05-30 21:50:05"
            },
            "assignee": {
                "id": 3970,
                "email": "john@capdax.com",
                "name": "John",
                "phone": "",
                "created_at": "2018-05-30 21:50:05"
            },
            "hits": []
        }
    }
}`)

var deniedResp = []byte(`
{
    "code": 200,
    "status": "success",
    "content": {
        "data": {
            "id": 93797112,
            "ref": "1544997039-RsF9dii1",
            "searcher_id": 3970,
            "assignee_id": 3970,
            "filters": {
                "country_codes": [],
                "remove_deceased": 1,
                "types": [
                    "sanction",
                    "warning",
                    "fitness-probity"
                ],
                "exact_match": true,
                "fuzziness": 0
            },
            "match_status": "potential_match",
            "risk_level": "unknown",
            "search_term": "Alexey Miller",
            "submitted_term": "Alexey Miller",
            "client_ref": null,
            "total_hits": 1,
            "updated_at": "2018-12-16 21:50:39",
            "created_at": "2018-12-16 21:50:39",
            "tags": [],
            "limit": 100,
            "offset": 0,
            "searcher": {
                "id": 3970,
                "email": "john@capdax.com",
                "name": "John",
                "phone": "",
                "created_at": "2018-05-30 21:50:05"
            },
            "assignee": {
                "id": 3970,
                "email": "john@capdax.com",
                "name": "John",
                "phone": "",
                "created_at": "2018-05-30 21:50:05"
            },
            "hits": [
                {
                    "doc": {
                        "aka": [
                            {
                                "name": "A Lie Ke Xie *Bao Li Suo Wei Qi *Mi Le"
                            },
                            {
                                "name": "Aleksey Miler"
                            },
                            {
                                "name": "Aleksei Borisovich Miller"
                            },
                            {
                                "name": "Aleksei Miller"
                            },
                            {
                                "name": "Aleksej Borisovic Miller"
                            },
                            {
                                "name": "Aleksej Borisovič Miller"
                            },
                            {
                                "name": "Aleksej Miler"
                            },
                            {
                                "name": "Aleksej Miller"
                            },
                            {
                                "name": "Aleksiej Miller"
                            },
                            {
                                "name": "Alexei Borissowitsch Miller"
                            },
                            {
                                "name": "Alexei Miller"
                            },
                            {
                                "name": "Alexey Miller"
                            },
                            {
                                "name": "Alexeï Miller"
                            },
                            {
                                "name": "Miller Oleksii Borisovich"
                            },
                            {
                                "name": "alragsei milrereu"
                            },
                            {
                                "name": "arekuseimireru"
                            },
                            {
                                "name": "lkhsy mylr"
                            },
                            {
                                "name": "lksyy mylr"
                            },
                            {
                                "name": "Алексей Борисович Миллер"
                            },
                            {
                                "name": "Алексей Миллер"
                            },
                            {
                                "name": "Алексеј Милер"
                            },
                            {
                                "name": "Міллер Олексій Борисович"
                            },
                            {
                                "name": "Ալեքսեյ Միլեր"
                            },
                            {
                                "name": "אלכסיי מילר"
                            },
                            {
                                "name": "الکسی میلر"
                            },
                            {
                                "name": "アレクセイ・ミレル"
                            },
                            {
                                "name": "阿列克谢·鲍里索维奇·米勒"
                            },
                            {
                                "name": "알락세이 밀레르"
                            },
                            {
                                "name": "Miller Alexey Borisovich"
                            }
                        ],
                        "entity_type": "person",
                        "fields": [
                            {
                                "name": "Country",
                                "source": "complyadvantage",
                                "value": "Russian Federation"
                            },
                            {
                                "name": "Country",
                                "source": "ofac-sdn-list",
                                "value": "Russian Federation"
                            },
                            {
                                "name": "Nationality",
                                "source": "complyadvantage",
                                "value": "Russian Federation"
                            },
                            {
                                "name": "Countries",
                                "tag": "country_names",
                                "value": "Russian Federation"
                            },
                            {
                                "name": "Date of Birth",
                                "source": "complyadvantage",
                                "tag": "date_of_birth",
                                "value": "1962-01-31"
                            },
                            {
                                "name": "Date of Birth",
                                "source": "ofac-sdn-list",
                                "tag": "date_of_birth",
                                "value": "1962-01-31"
                            },
                            {
                                "name": "Political Position",
                                "source": "complyadvantage",
                                "tag": "political_position",
                                "value": "Businessperson"
                            },
                            {
                                "name": "Political Position",
                                "source": "complyadvantage",
                                "tag": "political_position",
                                "value": "CEO of Gazprom"
                            },
                            {
                                "name": "Political Position",
                                "source": "complyadvantage",
                                "tag": "political_position",
                                "value": "Economist"
                            },
                            {
                                "name": "Gender",
                                "source": "complyadvantage",
                                "value": "male"
                            },
                            {
                                "name": "Gender",
                                "source": "ofac-sdn-list",
                                "value": "male"
                            },
                            {
                                "name": "Address",
                                "source": "ofac-sdn-list",
                                "value": "Moscow, Russia"
                            },
                            {
                                "name": "ofac Id",
                                "source": "ofac-sdn-list",
                                "value": "OFAC-24083"
                            },
                            {
                                "name": "Programs",
                                "source": "ofac-sdn-list",
                                "value": "* UKRAINE-EO13661: Executive Order 13661"
                            }
                        ],
                        "id": "TPD7IM2E6LW91J3",
                        "keywords": [],
                        "last_updated_utc": "2018-12-13T09:54:28Z",
                        "name": "Miller Alexey Borisovich",
                        "source_notes": {
                            "complyadvantage": {
                                "aml_types": [
                                    "pep-class-2"
                                ],
                                "country_codes": [
                                    "RU"
                                ],
                                "name": "DBPedia",
                                "url": "http://complyadvantage.com"
                            },
                            "ofac-sdn-list": {
                                "aml_types": [
                                    "sanction"
                                ],
                                "country_codes": [
                                    "RU"
                                ],
                                "listing_started_utc": "2018-04-09T00:00:00Z",
                                "name": "OFAC SDN List",
                                "url": "http://www.treasury.gov/resource-center/sanctions/SDN-List/Pages/default.aspx"
                            }
                        },
                        "sources": [
                            "complyadvantage",
                            "ofac-sdn-list"
                        ],
                        "types": [
                            "pep",
                            "pep-class-2",
                            "sanction"
                        ]
                    },
                    "match_types": [
                        "name_exact"
                    ],
                    "match_types_details": {
                        "Aleksei Borisovich Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksei Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksej Borisovic Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksej Borisovič Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksej Miler": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_fuzzy",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksej Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Aleksiej Miller": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Alexei Borissowitsch Miller": {
                            "match_types": {
                                "alexey": [
                                    "name_fuzzy",
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Alexei Miller": {
                            "match_types": {
                                "alexey": [
                                    "name_fuzzy",
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Alexey Miller": {
                            "match_types": {
                                "alexey": [
                                    "name_exact",
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Alexeï Miller": {
                            "match_types": {
                                "alexey": [
                                    "name_fuzzy",
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Alksyy mylr": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Miller Alexey Borisovich": {
                            "match_types": {
                                "alexey": [
                                    "name_exact",
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "name"
                        },
                        "Miller Oleksii Borisovich": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "name_exact",
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Алексей Борисович Миллер": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Алексей Миллер": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Алексеј Милер": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        },
                        "Міллер Олексій Борисович": {
                            "match_types": {
                                "alexey": [
                                    "phonetic_name"
                                ],
                                "miller": [
                                    "phonetic_name"
                                ]
                            },
                            "type": "aka"
                        }
                    },
                    "score": 35.348774
                }
            ]
        }
    }
}`)

var errorResp = []byte(`
{
    "code": 400,
    "status": "failure",
    "message": "Malformed Request",
    "errors": {
        "tags": "tags are invalid: {\"name\":\"value\"}"
    }
}`)

func TestNew(t *testing.T) {
	cfg := Config{
		Host:   "test",
		APIkey: "fake",
	}

	svc := New(cfg)
	svc2 := ComplyAdvantage{
		config: cfg,
	}

	assert := assert.New(t)

	assert.Equal(reflect.TypeOf(ComplyAdvantage{}), reflect.TypeOf(svc))
	assert.Equal(svc2, svc)
}

func TestCheckCustomerApproved(t *testing.T) {
	c := &common.UserData{
		FirstName:   "Georgy",
		LastName:    "Zhukov",
		MiddleName:  "Konstantinovich",
		DateOfBirth: common.Time(time.Date(1896, 11, 19, 0, 0, 0, 0, time.UTC)),
	}

	s := New(Config{
		Host:   "host",
		APIkey: "key",
	})

	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponder(http.MethodPost, "host/searches", httpmock.NewBytesResponder(http.StatusOK, approvedResp))

	res, err := s.CheckCustomer(c)

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(common.Approved, res.Status)
	assert.Nil(res.Details)
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCheckCustomerDenied(t *testing.T) {
	c := &common.UserData{
		FirstName: "Alexey",
		LastName:  "Miller",
	}

	s := New(Config{
		Host:   "host",
		APIkey: "key",
	})

	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponder(http.MethodPost, "host/searches", httpmock.NewBytesResponder(http.StatusOK, deniedResp))

	res, err := s.CheckCustomer(c)

	assert := assert.New(t)

	assert.Nil(err)
	assert.Equal(common.Denied, res.Status)
	assert.NotNil(res.Details)
	assert.Equal(common.Unknown, res.Details.Finality)
	assert.Len(res.Details.Reasons, 2)
	assert.Equal("Search ID: 93797112", res.Details.Reasons[0])
	assert.Equal("[Name: Miller Alexey Borisovich] Match types: name_exact", res.Details.Reasons[1])
	assert.Empty(res.ErrorCode)
	assert.Nil(res.StatusCheck)
}

func TestCheckCustomerError(t *testing.T) {
	s := New(Config{
		Host:   "host",
		APIkey: "key",
	})

	httpmock.Activate()
	defer httpmock.Deactivate()

	httpmock.RegisterResponder(http.MethodPost, "host/searches", httpmock.NewBytesResponder(http.StatusBadRequest, errorResp))

	res, err := s.CheckCustomer(&common.UserData{})

	assert := assert.New(t)

	assert.Equal(common.Error, res.Status)
	assert.Nil(res.Details)
	assert.Equal("400", res.ErrorCode)
	assert.Nil(res.StatusCheck)
	assert.NotNil(err)
	assert.Equal("400 Malformed Request | {\"tags\":\"tags are invalid: {\\\"name\\\":\\\"value\\\"}\"}", err.Error())
}

func TestCheckCustomerOtherErrors(t *testing.T) {
	s := New(Config{
		Host:   "host2",
		APIkey: "key",
	})

	httpmock.Activate()
	defer httpmock.Deactivate()

	_, err := s.CheckCustomer(&common.UserData{})

	assert.Error(t, err)

	httpmock.RegisterResponder(http.MethodPost, "host2/searches", httpmock.NewBytesResponder(http.StatusBadRequest, []byte(`unexpected response format`)))

	_, err = s.CheckCustomer(&common.UserData{})

	assert.Error(t, err)
}

func TestCheckStatus(t *testing.T) {
	assert := assert.New(t)

	s := ComplyAdvantage{}

	res, err := s.CheckStatus("")

	assert.Equal(common.Error, res.Status)
	assert.Error(err)
	assert.Equal("ComplyAdvantage doesn't support this method", err.Error())
}
