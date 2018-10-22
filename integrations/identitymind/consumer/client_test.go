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
    "state": "R",
    "tid": "26860023",
    "upr": "UNKNOWN",
    "user": "UNKNOWN"
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

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("should construct proper client instance", func() {
			config := Config{
				Host:     "fake_host",
				Username: "fake_name",
				Password: "fake_password",
			}

			testclient := &Client{
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

			result, err := client.CheckCustomer(&common.UserData{AccountName: "john_doe"})

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("during sending request: Post host/account/consumer: no responder found"))
		})

		It("should fail with fake error response", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, underReviewResponse))
			httpmock.RegisterResponder(http.MethodGet, client.host+stateRetrievalEndpoint+"26860023", httpmock.NewStringResponder(http.StatusOK, `{"error_message":"failed"}`))

			result, err := client.CheckCustomer(&common.UserData{AccountName: "john_doe"})

			Expect(result.StatusCheck).NotTo(BeNil())
			Expect(result.StatusCheck.Provider).To(Equal(common.IdentityMind))
			Expect(result.StatusCheck.ReferenceID).To(Equal("26860023"))

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).NotTo(HaveOccurred())

			result, err = client.CheckStatus(result.StatusCheck.ReferenceID)

			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("failed"))
		})

		It("should success with approved status", func() {
			Expect(client).ToNot(BeNil())

			httpmock.RegisterResponder(http.MethodPost, client.host+consumerEndpoint, httpmock.NewStringResponder(http.StatusOK, acceptedResponse))

			result, err := client.CheckCustomer(&common.UserData{AccountName: "john_doe"})

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

			result, err := client.CheckCustomer(&common.UserData{AccountName: "john_doe"})

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
