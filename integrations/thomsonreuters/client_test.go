package thomsonreuters

import (
	"fmt"
	"net/http"
	"testing"

	"modulus/kyc/integrations/thomsonreuters/model"

	"gopkg.in/jarcoal/httpmock.v1"
	"github.com/stretchr/testify/assert"
)

var tr = ThomsonReuters{
	scheme: "https",
	host:   "rms-world-check-one-api-pilot.thomsonreuters.com",
	path:   "/v1/",
	key:    "c7863652-3d05-4f02-8bf7-40ebb70fe17b",
	secret: "KXT8Pkj5n0Ttm4OSfD31x3Au4zf+2QqSbZIXBFoWq1oi7eGWh0k0dkqSdXmSmy15QcWyob7S/ENIdviedBCLRA==",
}

var groupsResponse = `
[
    {
        "id": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
        "name": "CriptoHub S.A. - API (P)",
        "parentId": null,
        "hasChildren": true,
        "status": "ACTIVE",
        "children": [
            {
                "id": "0a3687cf-65b4-1aaa-9975-f229000006ba",
                "name": "CriptoHub S.A. - Screening",
                "parentId": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
                "hasChildren": false,
                "status": "ACTIVE",
                "children": []
            }
        ]
    }
]`

var groupResponse = `
{
    "id": "0a3687cf-65b4-1aaa-9975-f229000006ba",
    "name": "CriptoHub S.A. - Screening",
    "parentId": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
    "hasChildren": false,
    "status": "ACTIVE",
    "children": []
}`

var caseTemplateResponse = `
{
    "groupId": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
    "groupScreeningType": "CASE_MANAGEMENT_AUDIT",
    "customFields": [],
    "secondaryFieldsByProvider": {
        "watchlist": {
            "secondaryFieldsByEntity": {
                "individual": [
                    {
                        "typeId": "SFCT_1",
                        "fieldValueType": "GENDER",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "GENDER"
                    },
                    {
                        "typeId": "SFCT_2",
                        "fieldValueType": "DATE",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "DATE_OF_BIRTH"
                    },
                    {
                        "typeId": "SFCT_3",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "COUNTRY_LOCATION"
                    },
                    {
                        "typeId": "SFCT_4",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "PLACE_OF_BIRTH"
                    },
                    {
                        "typeId": "SFCT_5",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "NATIONALITY"
                    }
                ],
                "vessel": [
                    {
                        "typeId": "SFCT_7",
                        "fieldValueType": "TEXT",
                        "regExp": "[0-9]{7}",
                        "fieldRequired": false,
                        "label": "IMO_NUMBER"
                    }
                ],
                "organisation": [
                    {
                        "typeId": "SFCT_6",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "REGISTERED_COUNTRY"
                    }
                ]
            }
        },
        "clientWatchlist": {
            "secondaryFieldsByEntity": {
                "individual": [
                    {
                        "typeId": "SFCT_1",
                        "fieldValueType": "GENDER",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "GENDER"
                    },
                    {
                        "typeId": "SFCT_2",
                        "fieldValueType": "DATE",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "DATE_OF_BIRTH"
                    },
                    {
                        "typeId": "SFCT_3",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "COUNTRY_LOCATION"
                    },
                    {
                        "typeId": "SFCT_4",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "PLACE_OF_BIRTH"
                    },
                    {
                        "typeId": "SFCT_5",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "NATIONALITY"
                    }
                ],
                "vessel": [
                    {
                        "typeId": "SFCT_7",
                        "fieldValueType": "TEXT",
                        "regExp": "[0-9]{7}",
                        "fieldRequired": false,
                        "label": "IMO_NUMBER"
                    }
                ],
                "organisation": [
                    {
                        "typeId": "SFCT_6",
                        "fieldValueType": "COUNTRY",
                        "regExp": null,
                        "fieldRequired": false,
                        "label": "REGISTERED_COUNTRY"
                    }
                ]
            }
        },
        "passportCheck": {
            "secondaryFieldsByEntity": {
                "individual": [
                    {
                        "typeId": "SFCT_8",
                        "fieldValueType": "TEXT",
                        "regExp": ".{0,1000}",
                        "fieldRequired": true,
                        "label": "PASSPORT_GIVEN_NAMES"
                    },
                    {
                        "typeId": "SFCT_9",
                        "fieldValueType": "TEXT",
                        "regExp": ".{0,1000}",
                        "fieldRequired": true,
                        "label": "PASSPORT_LAST_NAME"
                    },
                    {
                        "typeId": "SFCT_10",
                        "fieldValueType": "GENDER",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_GENDER"
                    },
                    {
                        "typeId": "SFCT_11",
                        "fieldValueType": "STATE",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_ISSUING_STATE"
                    },
                    {
                        "typeId": "SFCT_12",
                        "fieldValueType": "STATE",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_NATIONALITY"
                    },
                    {
                        "typeId": "SFCT_13",
                        "fieldValueType": "DATE",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_DATE_OF_BIRTH"
                    },
                    {
                        "typeId": "SFCT_14",
                        "fieldValueType": "PASSPORT_DOCUMENT_TYPE",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_DOCUMENT_TYPE"
                    },
                    {
                        "typeId": "SFCT_15",
                        "fieldValueType": "TEXT",
                        "regExp": "^[a-zA-Z0-9<]{0,1000}$",
                        "fieldRequired": true,
                        "label": "PASSPORT_ID_NUMBER"
                    },
                    {
                        "typeId": "SFCT_16",
                        "fieldValueType": "DATE",
                        "regExp": null,
                        "fieldRequired": true,
                        "label": "PASSPORT_DATE_OF_EXPIRY"
                    }
                ]
            }
        }
    },
    "mandatoryProviderTypes": [
        "WATCHLIST"
    ]
}`

var resolutionToolkitsResponse = `
{
    "WATCHLIST": {
        "groupId": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
        "providerType": "WATCHLIST",
        "resolutionFields": {
            "statuses": [
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000696",
                    "label": "POSITIVE",
                    "type": "POSITIVE"
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c0000069b",
                    "label": "POSSIBLE",
                    "type": "POSSIBLE"
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c000006a1",
                    "label": "FALSE",
                    "type": "FALSE"
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c000006a4",
                    "label": "UNSPECIFIED",
                    "type": "UNSPECIFIED"
                }
            ],
            "risks": [
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000695",
                    "label": "UNKNOWN",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000692",
                    "label": "HIGH",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000693",
                    "label": "MEDIUM",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000694",
                    "label": "LOW",
                    "type": null
                }
            ],
            "reasons": [
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000690",
                    "label": "No Match",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c0000068e",
                    "label": "Full Match",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c0000068f",
                    "label": "Partial Match",
                    "type": null
                },
                {
                    "id": "0a3687d0-65b4-1cc3-9975-f20c00000691",
                    "label": "Unknown",
                    "type": null
                }
            ]
        },
        "resolutionRules": {
            "0a3687d0-65b4-1cc3-9975-f20c0000069b": {
                "reasons": [
                    "0a3687d0-65b4-1cc3-9975-f20c0000068f"
                ],
                "remarkRequired": false,
                "reasonRequired": true,
                "risks": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000693",
                    "0a3687d0-65b4-1cc3-9975-f20c00000692",
                    "0a3687d0-65b4-1cc3-9975-f20c00000695",
                    "0a3687d0-65b4-1cc3-9975-f20c00000694"
                ]
            },
            "0a3687d0-65b4-1cc3-9975-f20c00000696": {
                "reasons": [
                    "0a3687d0-65b4-1cc3-9975-f20c0000068e"
                ],
                "remarkRequired": false,
                "reasonRequired": true,
                "risks": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000693",
                    "0a3687d0-65b4-1cc3-9975-f20c00000692",
                    "0a3687d0-65b4-1cc3-9975-f20c00000694"
                ]
            },
            "0a3687d0-65b4-1cc3-9975-f20c000006a1": {
                "reasons": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000690"
                ],
                "remarkRequired": false,
                "reasonRequired": true,
                "risks": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000695"
                ]
            },
            "0a3687d0-65b4-1cc3-9975-f20c000006a4": {
                "reasons": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000691"
                ],
                "remarkRequired": false,
                "reasonRequired": true,
                "risks": [
                    "0a3687d0-65b4-1cc3-9975-f20c00000695"
                ]
            }
        }
    },
    "CLIENT_WATCHLIST": {
        "groupId": "0a3687d0-65b4-1cc3-9975-f20b0000066f",
        "providerType": "CLIENT_WATCHLIST",
        "resolutionFields": {
            "statuses": [
                {
                    "id": "DEFAULT_STATUS_POSITIVE",
                    "label": "Positive",
                    "type": "POSITIVE"
                },
                {
                    "id": "DEFAULT_STATUS_POSSIBLE",
                    "label": "Possible",
                    "type": "POSSIBLE"
                },
                {
                    "id": "DEFAULT_STATUS_FALSE",
                    "label": "False",
                    "type": "FALSE"
                },
                {
                    "id": "DEFAULT_STATUS_UNSPECIFIED",
                    "label": "Unspecified",
                    "type": "UNSPECIFIED"
                }
            ],
            "risks": [],
            "reasons": [
                {
                    "id": "DEFAULT_REASON_FALSE",
                    "label": "Default Reason for False Status",
                    "type": null
                },
                {
                    "id": "DEFAULT_REASON_POSITIVE",
                    "label": "Default Reason for Positive Status",
                    "type": null
                },
                {
                    "id": "DEFAULT_REASON_POSSIBLE",
                    "label": "Default Reason for Possible Status",
                    "type": null
                },
                {
                    "id": "DEFAULT_REASON_UNSPECIFIED",
                    "label": "Default Reason for Unspecified Status",
                    "type": null
                }
            ]
        },
        "resolutionRules": {
            "DEFAULT_STATUS_UNSPECIFIED": {
                "reasons": [
                    "DEFAULT_REASON_UNSPECIFIED"
                ],
                "remarkRequired": true,
                "reasonRequired": true,
                "risks": []
            },
            "DEFAULT_STATUS_FALSE": {
                "reasons": [
                    "DEFAULT_REASON_FALSE"
                ],
                "remarkRequired": true,
                "reasonRequired": true,
                "risks": []
            },
            "DEFAULT_STATUS_POSSIBLE": {
                "reasons": [
                    "DEFAULT_REASON_POSSIBLE"
                ],
                "remarkRequired": true,
                "reasonRequired": true,
                "risks": []
            },
            "DEFAULT_STATUS_POSITIVE": {
                "reasons": [
                    "DEFAULT_REASON_POSITIVE"
                ],
                "remarkRequired": true,
                "reasonRequired": true,
                "risks": []
            }
        }
    }
}`

var syncScreeningResponseDenied = `
{
    "caseId": "24da33ec-9ad9-463c-9ef7-9e0dce1bfcbb",
    "results": [
        {
            "resultId": "0a3687d0-673a-15cf-9a06-ae7c00d3929c",
            "referenceId": "e_tr_wci_846144",
            "matchStrength": "EXACT",
            "matchedTerm": "Сергей Владимирович Железняк",
            "submittedTerm": "Сергей Владимирович Железняк",
            "matchedNameType": "NATIVE_AKA",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": "1970-07-30"
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1970-07-30",
                    "matchedValue": null,
                    "matchedDateTimeValue": "1970-07-30",
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                }
            ],
            "sources": [
                "b_trwc_67",
                "b_trwc_CANSEMRUS",
                "b_trwc_14",
                "b_trwc_15",
                "b_trwc_PEP N",
                "b_trwc_UANSDC",
                "b_trwc_OFAC-UKR2",
                "b_trwc_21",
                "b_trwc_190",
                "b_trwc_SECO-UKR",
                "b_trwc_160",
                "b_trwc_151",
                "b_trwc_162",
                "b_trwc_141",
                "b_trwc_EU-UKR2",
                "b_trwc_110",
                "b_trwc_188",
                "b_trwc_386",
                "b_trwc_155",
                "b_trwc_112",
                "b_trwc_156",
                "b_trwc_312",
                "b_trwc_445"
            ],
            "categories": [
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "PEP",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions",
                "Sanctions"
            ],
            "creationDate": "2019-01-04T21:17:00.013Z",
            "modificationDate": "2019-01-04T21:17:00.013Z",
            "primaryName": "Sergei Vladimirovich ZHELEZNYAK",
            "events": [
                {
                    "day": 30,
                    "month": 7,
                    "year": 1970,
                    "address": null,
                    "fullDate": "1970-07-30",
                    "allegedAddresses": [],
                    "type": "BIRTH"
                }
            ],
            "countryLinks": [
                {
                    "countryText": "UNKNOWN",
                    "country": {
                        "code": "ZZZ",
                        "name": "UNKNOWN"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "POLITICAL INDIVIDUAL",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        },
        {
            "resultId": "0a3687d0-673a-15cf-9a06-ae7c00d3923a",
            "referenceId": "e_tr_wci_2878909",
            "matchStrength": "STRONG",
            "matchedTerm": "Sergey ZHELEZNYAK",
            "submittedTerm": "Сергей Владимирович Железняк",
            "matchedNameType": "PRIMARY",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1970-07-30",
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                }
            ],
            "sources": [
                "b_trwc_4"
            ],
            "categories": [
                "Other Bodies"
            ],
            "creationDate": "2019-01-04T21:17:00.013Z",
            "modificationDate": "2019-01-04T21:17:00.013Z",
            "primaryName": "Sergey ZHELEZNYAK",
            "events": [],
            "countryLinks": [
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "CRIME - FINANCIAL",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        },
        {
            "resultId": "0a3687d0-673a-15cf-9a06-ae7c00d39192",
            "referenceId": "e_tr_wci_3692518",
            "matchStrength": "MEDIUM",
            "matchedTerm": "ZHELEZNYAKOV,Vladimir Igorevich",
            "submittedTerm": "Сергей Владимирович Железняк",
            "matchedNameType": "AKA",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": "1988-05-21"
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1970-07-30",
                    "matchedValue": null,
                    "matchedDateTimeValue": "1988-05-21",
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                }
            ],
            "sources": [
                "b_trwc_UAMVS"
            ],
            "categories": [
                "Law Enforcement"
            ],
            "creationDate": "2019-01-04T21:17:00.013Z",
            "modificationDate": "2019-01-04T21:17:00.013Z",
            "primaryName": "Vladimir ZHELEZNYAKOV",
            "events": [
                {
                    "day": 21,
                    "month": 5,
                    "year": 1988,
                    "address": null,
                    "fullDate": "1988-05-21",
                    "allegedAddresses": [],
                    "type": "BIRTH"
                }
            ],
            "countryLinks": [
                {
                    "countryText": "UKRAINE",
                    "country": {
                        "code": "UKR",
                        "name": "UKRAINE"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "UKRAINE",
                    "country": {
                        "code": "UKR",
                        "name": "UKRAINE"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "INDIVIDUAL",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        },
        {
            "resultId": "0a3687d0-673a-15cf-9a06-ae7c00d39272",
            "referenceId": "e_tr_wci_4239867",
            "matchStrength": "WEAK",
            "matchedTerm": "КОЛЕСНИК,Сергей Владимирович",
            "submittedTerm": "Сергей Владимирович Железняк",
            "matchedNameType": "NATIVE_AKA",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": "1989-05-20"
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1970-07-30",
                    "matchedValue": null,
                    "matchedDateTimeValue": "1989-05-20",
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                }
            ],
            "sources": [
                "b_trwc_UAMVS"
            ],
            "categories": [
                "Law Enforcement"
            ],
            "creationDate": "2019-01-04T21:17:00.013Z",
            "modificationDate": "2019-01-04T21:17:00.013Z",
            "primaryName": "Sergey KOLESNIK",
            "events": [
                {
                    "day": 20,
                    "month": 5,
                    "year": 1989,
                    "address": null,
                    "fullDate": "1989-05-20",
                    "allegedAddresses": [],
                    "type": "BIRTH"
                }
            ],
            "countryLinks": [
                {
                    "countryText": "UKRAINE",
                    "country": {
                        "code": "UKR",
                        "name": "UKRAINE"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "UKRAINE",
                    "country": {
                        "code": "UKR",
                        "name": "UKRAINE"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "INDIVIDUAL",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        }
    ]
}`

var syncScreeningResponseApproved = `
{
    "caseId": "32737c50-0058-4f28-a0fa-01776aba71e4",
    "results": [
        {
            "resultId": "0a3687cf-673a-1553-9a06-c6dc00d5378b",
            "referenceId": "e_tr_wci_1229873",
            "matchStrength": "EXACT",
            "matchedTerm": "САРБАШ,Сергей Васильевич",
            "submittedTerm": "Сергей Васильевич Сарбаш",
            "matchedNameType": "NATIVE_AKA",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": "1967-06-12"
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1975-09-21",
                    "matchedValue": null,
                    "matchedDateTimeValue": "1967-06-12",
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "RUS",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "RUS",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                }
            ],
            "sources": [
                "b_trwc_PEP N"
            ],
            "categories": [
                "PEP"
            ],
            "creationDate": "2019-01-04T23:03:29.980Z",
            "modificationDate": "2019-01-04T23:03:29.980Z",
            "primaryName": "Sergey SARBASH",
            "events": [
                {
                    "day": 12,
                    "month": 6,
                    "year": 1967,
                    "address": null,
                    "fullDate": "1967-06-12",
                    "allegedAddresses": [],
                    "type": "BIRTH"
                }
            ],
            "countryLinks": [
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "RUSSIAN FEDERATION",
                    "country": {
                        "code": "RUS",
                        "name": "RUSSIAN FEDERATION"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "LEGAL",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        },
        {
            "resultId": "0a3687cf-673a-1553-9a06-c6dc00d538b1",
            "referenceId": "e_tr_wci_1512563",
            "matchStrength": "STRONG",
            "matchedTerm": "МІШУСТІН,Сергій Васильович",
            "submittedTerm": "Сергей Васильевич Сарбаш",
            "matchedNameType": "NATIVE_AKA",
            "secondaryFieldResults": [
                {
                    "field": {
                        "typeId": null,
                        "value": "MALE",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "MALE",
                    "submittedDateTimeValue": null,
                    "matchedValue": "MALE",
                    "matchedDateTimeValue": null,
                    "fieldResult": "MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": "1946-08-29"
                    },
                    "typeId": null,
                    "submittedValue": null,
                    "submittedDateTimeValue": "1975-09-21",
                    "matchedValue": null,
                    "matchedDateTimeValue": "1946-08-29",
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "GAB",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "GAB",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": null,
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": null,
                    "matchedDateTimeValue": null,
                    "fieldResult": "UNKNOWN"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                },
                {
                    "field": {
                        "typeId": null,
                        "value": "UKR",
                        "dateTimeValue": null
                    },
                    "typeId": null,
                    "submittedValue": "RUS",
                    "submittedDateTimeValue": null,
                    "matchedValue": "UKR",
                    "matchedDateTimeValue": null,
                    "fieldResult": "NOT_MATCHED"
                }
            ],
            "sources": [
                "b_trwc_PEP N"
            ],
            "categories": [
                "PEP"
            ],
            "creationDate": "2019-01-04T23:03:29.980Z",
            "modificationDate": "2019-01-04T23:03:29.980Z",
            "primaryName": "Serhiy MISHUSTIN",
            "events": [
                {
                    "day": 29,
                    "month": 8,
                    "year": 1946,
                    "address": null,
                    "fullDate": "1946-08-29",
                    "allegedAddresses": [],
                    "type": "BIRTH"
                }
            ],
            "countryLinks": [
                {
                    "countryText": "UKRAINE",
                    "country": {
                        "code": "UKR",
                        "name": "UKRAINE"
                    },
                    "type": "NATIONALITY"
                },
                {
                    "countryText": "CONGO DEM. REP.",
                    "country": {
                        "code": "COD",
                        "name": "CONGO, THE DEMOCRATIC REPUBLIC OF THE"
                    },
                    "type": "LOCATION"
                },
                {
                    "countryText": "GABON",
                    "country": {
                        "code": "GAB",
                        "name": "GABON"
                    },
                    "type": "LOCATION"
                }
            ],
            "identityDocuments": [],
            "category": "DIPLOMAT",
            "providerType": "WATCHLIST",
            "gender": "MALE"
        }
    ]
}`

func TestGetRootGroups(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))

	groups, status, err := tr.getRootGroups()

	assert.NoError(err)
	assert.Nil(status)
	assert.Len(groups, 1)

	group := groups[0]
	assert.Equal("0a3687d0-65b4-1cc3-9975-f20b0000066f", group.ID)
	assert.Equal("CriptoHub S.A. - API (P)", group.Name)
	assert.Empty(group.ParentID)
	assert.True(group.HasChildren)
	assert.Equal(model.ActiveStatus, group.Status)
	assert.Len(group.Children, 1)
}

func TestGetGroup(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))
	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups/0a3687cf-65b4-1aaa-9975-f229000006ba", httpmock.NewStringResponder(http.StatusOK, groupResponse))

	groups, status, err := tr.getRootGroups()

	assert.NoError(err)
	assert.Nil(status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.HasChildren {
			assert.NotEmpty(g.Children)
			gID = g.Children[0].ID
			break
		}
	}

	assert.NotEmpty(gID)

	group, status, err := tr.getGroup(gID)

	assert.NoError(err)
	assert.Nil(status)

	assert.Equal("0a3687cf-65b4-1aaa-9975-f229000006ba", group.ID)
	assert.Equal("CriptoHub S.A. - Screening", group.Name)
	assert.Equal("0a3687d0-65b4-1cc3-9975-f20b0000066f", group.ParentID)
	assert.Equal(model.ActiveStatus, group.Status)
	assert.False(group.HasChildren)
	assert.Empty(group.Children)
}

func TestGetCaseTemplate(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))
	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups/0a3687d0-65b4-1cc3-9975-f20b0000066f/caseTemplate", httpmock.NewStringResponder(http.StatusOK, caseTemplateResponse))

	groups, status, err := tr.getRootGroups()

	assert.NoError(err)
	assert.Nil(status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		gID = g.ID
		break
	}

	assert.NotEmpty(gID)

	ctr, status, err := tr.getCaseTemplate(gID)

	// FIXME: perhaps, can add more checks.

	assert.NoError(err)
	assert.Nil(status)

	assert.Equal("0a3687d0-65b4-1cc3-9975-f20b0000066f", ctr.GroupID)
	assert.Equal(model.CaseManagementAudit, ctr.GroupScreeningType)
	assert.Len(ctr.MandatoryProviderTypes, 1)
	assert.Equal(model.WatchList, ctr.MandatoryProviderTypes[0])
	assert.Empty(ctr.CustomFields)
	assert.Contains(ctr.SecondaryFieldsByProvider, "watchlist")
	assert.Contains(ctr.SecondaryFieldsByProvider, "clientWatchlist")
	assert.Contains(ctr.SecondaryFieldsByProvider, "passportCheck")

	watchlist := ctr.SecondaryFieldsByProvider["watchlist"]

	assert.Contains(watchlist.SecondaryFieldsByEntity, "individual")
	assert.Contains(watchlist.SecondaryFieldsByEntity, "vessel")
	assert.Contains(watchlist.SecondaryFieldsByEntity, "organisation")
	assert.Len(watchlist.SecondaryFieldsByEntity["individual"], 5)

	clientWatchlist := ctr.SecondaryFieldsByProvider["clientWatchlist"]

	assert.Contains(clientWatchlist.SecondaryFieldsByEntity, "individual")
	assert.Contains(clientWatchlist.SecondaryFieldsByEntity, "vessel")
	assert.Contains(clientWatchlist.SecondaryFieldsByEntity, "organisation")
	assert.Len(clientWatchlist.SecondaryFieldsByEntity["individual"], 5)

	passportCheck := ctr.SecondaryFieldsByProvider["passportCheck"]

	assert.Contains(passportCheck.SecondaryFieldsByEntity, "individual")
	assert.NotContains(passportCheck.SecondaryFieldsByEntity, "vessel")
	assert.NotContains(passportCheck.SecondaryFieldsByEntity, "organisation")
	assert.Len(passportCheck.SecondaryFieldsByEntity["individual"], 9)
}

func TestPerformSynchronousScreening(t *testing.T) {
	assert := assert.New(t)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, tr.scheme+"://"+tr.host+tr.path+"groups", httpmock.NewStringResponder(http.StatusOK, groupsResponse))
	httpmock.RegisterResponder(http.MethodPost, tr.scheme+"://"+tr.host+tr.path+"cases/screeningRequest", httpmock.NewStringResponder(http.StatusOK, syncScreeningResponseDenied))

	groups, status, err := tr.getRootGroups()

	assert.NoError(err)
	assert.Nil(status)
	assert.NotEmpty(groups)

	gID := ""
	for _, g := range groups {
		if g.Status != model.ActiveStatus {
			continue
		}

		gID = g.ID
		break
	}

	assert.NotEmpty(gID)

	newcase := model.NewCase{
		GroupID:    gID,
		EntityType: model.IndividualCET,
		Name:       "Сергей Владимирович Железняк",
		ProviderTypes: []model.ProviderType{
			model.WatchList,
		},
		SecondaryFields: []model.Field{
			model.Field{
				TypeID: "SFCT_1",
				Value:  model.Male,
			},
			model.Field{
				TypeID:        "SFCT_2",
				DateTimeValue: "1970-07-30",
			},
			model.Field{
				TypeID: "SFCT_3",
				Value:  "RUS",
			},
			model.Field{
				TypeID: "SFCT_4",
				Value:  "RUS",
			},
			model.Field{
				TypeID: "SFCT_5",
				Value:  "RUS",
			},
		},
	}

	src, status, err := tr.performSynchronousScreening(newcase)
	if status != nil {
		fmt.Println(*status)
	}

	assert.NoError(err)
	assert.Nil(status)

	assert.Equal("24da33ec-9ad9-463c-9ef7-9e0dce1bfcbb", src.CaseID)
	assert.Len(src.Results, 4)
	assert.Equal(model.Exact, src.Results[0].MatchStrength)
	assert.Equal(newcase.Name, src.Results[0].MatchedTerm)
	assert.Equal(model.Strong, src.Results[1].MatchStrength)
	assert.NotEqual(newcase.Name, src.Results[1].MatchedTerm)
	assert.Equal(model.Medium, src.Results[2].MatchStrength)
	assert.NotEqual(newcase.Name, src.Results[2].MatchedTerm)
	assert.Equal(model.Weak, src.Results[3].MatchStrength)
	assert.NotEqual(newcase.Name, src.Results[3].MatchedTerm)
}
