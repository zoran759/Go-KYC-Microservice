package consumer

import (
	"encoding/base64"
	"net/http"

	"modulus/kyc/common"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var underReviewResponse = `
{
    "ednaScoreCard": {
        "sc": [],
        "etr": [
            {
                "test": "ed:23",
                "fired": false,
                "details": "Associated Devices: 0",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:37",
                "fired": false,
                "details": "ed:37(false) = true",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:59",
                "fired": false,
                "details": "The MAN count is 1",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:36",
                "fired": false,
                "details": "ed:36(false) = true",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:25",
                "fired": false,
                "details": "Associated Billing Addresses: 0",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ss:0",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "mm:0",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "dv:1",
                "fired": true,
                "details": "[Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "dv:1",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "au:1",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "dv:0",
                "fired": true,
                "details": "[Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "dv:0",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:32",
                "fired": false,
                "details": "ed:32(false) = true",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:503",
                "fired": false,
                "details": "1 hour account account creation velocity = 1.  ",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:514",
                "fired": false,
                "details": "ed:514(0) > 5",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:513",
                "fired": false,
                "details": "28 day account account creation velocity = 1.  ",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:504",
                "fired": false,
                "details": "ed:504(0) > 2",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ra:0",
                "details": "false",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:509",
                "fired": false,
                "details": "ed:509(0) > 3",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:508",
                "fired": false,
                "details": "24 hour account account creation velocity = 1.  ",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:2",
                "fired": false,
                "details": "No transaction parameter on watchlist",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:1",
                "fired": false,
                "details": "ed:1(false) = true",
                "ts": 1549976042000,
                "stage": "1"
            },
            {
                "test": "ed:19",
                "fired": false,
                "details": "User Account reputation is Bad or Blacklisted",
                "ts": 1549976042000,
                "stage": "1"
            }
        ],
        "ar": {
            "result": "DISABLED"
        },
        "er": {
            "profile": "DEFAULT",
            "reportedRule": {
                "description": "ID document validation",
                "resultCode": "MANUAL_REVIEW",
                "details": "[Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'; [Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'",
                "ruleId": 30060,
                "testResults": [
                    {
                        "test": "dv:1",
                        "fired": true,
                        "details": "[Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'",
                        "condition": {
                            "left": "dv:1",
                            "right": false,
                            "operator": "eq",
                            "type": "info"
                        },
                        "ts": 1549976042000,
                        "stage": "1"
                    },
                    {
                        "test": "dv:0",
                        "fired": true,
                        "details": "[Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'",
                        "condition": {
                            "left": "dv:0",
                            "right": false,
                            "operator": "eq",
                            "type": "info"
                        },
                        "ts": 1549976042000,
                        "stage": "1"
                    }
                ],
                "name": "Document Validation"
            }
        }
    },
    "mtid": "86c5468b323346378083d571a5dc480a",
    "state": "R",
    "tid": "86c5468b323346378083d571a5dc480a",
    "rcd": ""
}`

var acceptedResponse = `
{
    "ednaScoreCard": {
        "er": {
            "reportedRule": {
                "description": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
                "details": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
                "name": "Unknown Fallthrough",
                "resultCode": "ACCEPT",
                "ruleId": 1002,
                "testResults": []
            }
        },
        "sc": []
    },
    "erd": "Unknown User",
    "frd": "Fallthrough for transaction with an unknown entity. No other rules triggered.",
    "frn": "Unknown Fallthrough",
    "frp": "ACCEPT",
    "mtid": "26860023",
    "rcd": "1002,101,202,111,131,50005,150",
    "res": "ACCEPT",
    "state": "A",
    "tid": "26860023",
    "upr": "UNKNOWN",
    "user": "UNKNOWN"
}`

var rejectedResponse = `
{
    "ednaScoreCard": {
        "er": {
            "reportedRule": {
                "description": "Fake fraud attempt simulation.",
                "details": "Fake transaction with the suspicious entity. No other rules triggered.",
                "name": "Fraud Attempt",
                "resultCode": "DENY",
                "ruleId": 1002,
                "testResults": []
            }
        },
        "sc": []
    },
    "erd": "Fraudster",
    "frd": ".",
    "frn": "Fraud Attempt",
    "frp": "DENY",
    "mtid": "26860023",
    "rcd": "1002,101,202,111,131,50005,150",
    "res": "DENY",
    "state": "D",
    "tid": "26860023",
    "upr": "SUSPICIOUS",
    "user": "BAD"
}`

// "mtid" key type intentionally changed from string to integer.
var malformedResponse = `
{
	"erd": "Fraudster",
    "frd": ".",
    "frn": "Fraud Attempt",
    "frp": "DENY",
    "mtid": 26860023,
    "rcd": "1002,101,202,111,131,50005,150",
    "res": "DENY",
    "state": "D",
    "tid": "26860023",
    "upr": "SUSPICIOUS",
    "user": "Unknown"
}`

// "state" value intentionally changed to unknown value.
var unknownStateResponse = `
{
	"erd": "Fraudster",
    "frd": ".",
    "frn": "Fraud Attempt",
    "frp": "DENY",
    "mtid": "26860023",
    "rcd": "1002,101,202,111,131,50005,150",
    "res": "DENY",
    "state": "U",
    "tid": "26860023",
    "upr": "UNKNOWN",
    "user": "Unknown"
}`

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("should construct proper client instance", func() {
			config := Config{
				Host:     "fake_host",
				Username: "fake_name",
				Password: "fake_password",
			}

			testclient := Client{
				host:        config.Host,
				credentials: "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Username+":"+config.Password)),
			}

			client := NewClient(config)

			Expect(client).NotTo(BeNil())
			Expect(client).To(Equal(testclient))
		})
	})

	Describe("CheckCustomer", func() {
		var client = NewClient(Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		})

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail with nil customer", func() {
			Expect(client).ToNot(BeNil())

			result, err := client.CheckCustomer(nil)

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no customer supplied"))
		})

		It("should fail with malformed customer", func() {
			Expect(client).ToNot(BeNil())

			malformedCustomer := &common.UserData{
				AccountName: "very long account name that is exceeding the limit of 60 symbols",
			}

			result, err := client.CheckCustomer(malformedCustomer)

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("invalid customer data: account length 64 exceeded limit of 60 symbols"))
		})

		It("should fail with missing httpmock responder", func() {
			Expect(client).ToNot(BeNil())

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: Post host/account/consumer: no responder found"))
		})

		It("should fail with fake error response", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, `{"error_message":"failed"}`))

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(result.StatusCheck).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: failed"))
		})

		It("should success with approved status", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, acceptedResponse))

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: UNKNOWN"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: ACCEPT"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Unknown User"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: ACCEPT"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should success with rejected status", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, rejectedResponse))

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Denied))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: BAD"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: DENY"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Fraudster"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: DENY"))
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with unknown status", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, unknownStateResponse))

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("Customer reputation: Unknown"))
			Expect(result.Details.Reasons[1]).To(Equal("Fraud policy evaluation result: DENY"))
			Expect(result.Details.Reasons[2]).To(Equal("Customer reputation reason: Fraudster"))
			Expect(result.Details.Reasons[3]).To(Equal("Combined fraud and automated review evaluations result: DENY"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown state of the verification from the API: U"))
		})

		It("should success with 'under review' status", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, underReviewResponse))

			result, err := client.CheckCustomer(&common.UserData{})

			Expect(result.Status).To(Equal(common.Denied))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(5))
			Expect(result.Details.Reasons[0]).To(Equal("Some of checks triggered 'MANUAL REVIEW' status"))
			Expect(result.Details.Reasons[1]).To(Equal("Profile: DEFAULT"))
			Expect(result.Details.Reasons[2]).To(Equal("Rule: id 30060 | ID document validation"))
			Expect(result.Details.Reasons[3]).To(Equal("Test: 'dv:1' | [Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'"))
			Expect(result.Details.Reasons[4]).To(Equal("Test: 'dv:0' | [Fired] No remaining queries for this third party service. Please increase limit in the Admin tab of the UI. Select the Merchant Preferences section and 'Third Party Overview'"))
			Expect(result.StatusCheck).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("sendRequest", func() {
		var client = NewClient(Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		})

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail with error message", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, malformedResponse))

			resp, errorCode, err := client.sendRequest([]byte{})

			Expect(resp).ToNot(BeNil())
			Expect(errorCode).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("should success", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, acceptedResponse))

			resp, errorCode, err := client.sendRequest([]byte{})

			Expect(resp).ToNot(BeNil())
			Expect(errorCode).To(BeNil())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("CheckStatus", func() {
		var client = NewClient(Config{
			Host:     "host",
			Username: "test",
			Password: "test",
		})

		var referenceID = "26860023"

		BeforeEach(func() {
			httpmock.Activate()
		})

		AfterEach(func() {
			httpmock.DeactivateAndReset()
		})

		It("should fail with error message", func() {
			Expect(client).ToNot(BeNil())

			resp, err := client.CheckStatus(referenceID)

			Expect(resp.Details).To(BeNil())
			Expect(resp.ErrorCode).To(BeEmpty())
			Expect(resp.StatusCheck).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: Get host/account/consumer/v2/26860023: no responder found"))
		})

		It("should fail with error message", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodGet, client.host+stateRetrievalEndpoint+"26860023", httpmock.NewStringResponder(http.StatusOK, `{"error_message":"failed"}`))

			resp, err := client.CheckStatus(referenceID)

			Expect(resp).ToNot(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed"))
		})

		It("should fail with error message", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodGet, client.host+stateRetrievalEndpoint+"26860023", httpmock.NewStringResponder(http.StatusOK, malformedResponse))

			resp, err := client.CheckStatus(referenceID)

			Expect(resp).ToNot(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("should success", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodGet, client.host+stateRetrievalEndpoint+"26860023", httpmock.NewStringResponder(http.StatusOK, acceptedResponse))

			resp, err := client.CheckStatus(referenceID)

			Expect(resp).ToNot(BeNil())
			Expect(resp.Status).To(Equal(common.Approved))
			Expect(resp.ErrorCode).To(BeEmpty())
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
