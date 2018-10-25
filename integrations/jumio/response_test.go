package jumio

import (
	"encoding/json"
	"modulus/kyc/common"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var approvedResponse = []byte(`
{
	"timestamp": "2014-08-14T08:16:20.845Z",
	"scanReference": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	"document": {
		"type": "PASSPORT",
		"dob": "1990-01-01",
		"expiry": "2022-12-31",
		"firstName": "FIRSTNAME",
		"issuingCountry": "USA",
		"lastName": "LASTNAME",
		"number": "P1234",
		"status": "APPROVED_VERIFIED"
	},
	"transaction": {
		"clientIp": "xxx.xxx.xxx.xxx",
		"customerId": "CUSTOMERID",
		"date": "2014-08-10T10:27:29.494Z",
		"source": "API",
		"status": "DONE"
	},
	"verification": {
		"mrzCheck": "OK"
	}
}`)

var deniedResponse = []byte(`
{
	"timestamp": "2014-08-14T08:16:20.845Z",
	"scanReference": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	"document": {
		"type": "PASSPORT",
		"dob": "1990-01-01",
		"expiry": "2022-12-31",
		"firstName": "FIRSTNAME",
		"issuingCountry": "USA",
		"lastName": "LASTNAME",
		"number": "P1234",
		"status": "DENIED_FRAUD"
	},
	"transaction": {
		"clientIp": "xxx.xxx.xxx.xxx",
		"customerId": "CUSTOMERID",
		"date": "2014-08-10T10:27:29.494Z",
		"source": "API",
		"status": "DONE"
	},
	"verification": {
		"mrzCheck": "OK",
		"rejectReason": {
			"rejectReasonCode": "100",
			"rejectReasonDescription": "MANIPULATED_DOCUMENT",
			"rejectReasonDetails": [
				{
					"detailsCode": "1002",
					"detailsDescription": "DOCUMENT_NUMBER"
				},
				{
					"detailsCode": "1007",
					"detailsDescription": "SECURITY_CHECKS"
				}
			]
		},
		"identityVerification": {
			"similarity": "NO_MATCH",
			"validity": "FALSE",
			"reason": "SELFIE_MANIPULATED"
		}
	}
}`)

var failedResponse = []byte(`
{
	"timestamp": "2014-08-14T08:16:20.845Z",
	"scanReference": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	"document": {
		"type": "PASSPORT",
		"dob": "1990-01-01",
		"expiry": "2022-12-31",
		"firstName": "FIRSTNAME",
		"issuingCountry": "USA",
		"lastName": "LASTNAME",
		"number": "P1234",
		"status": "NO_ID_UPLOADED"
	},
	"transaction": {
		"clientIp": "xxx.xxx.xxx.xxx",
		"customerId": "CUSTOMERID",
		"date": "2014-08-10T10:27:29.494Z",
		"source": "API",
		"status": "FAILED"
	},
	"verification": {
		"mrzCheck": "OK",
		"identityVerification": {
			"similarity": "MATCH",
			"validity": "TRUE"
		}
	}
}`)

var _ = Describe("Response", func() {
	Describe("toResult", func() {
		It("should success with approved result", func() {
			r := &DetailsResponse{}

			err := json.Unmarshal(approvedResponse, r)

			Expect(err).NotTo(HaveOccurred())

			result, err := r.toResult()

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal(common.Approved))
			Expect(result.Details).To(BeNil())
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})

		It("should success with denied result", func() {
			r := &DetailsResponse{}

			err := json.Unmarshal(deniedResponse, r)

			Expect(err).NotTo(HaveOccurred())

			result, err := r.toResult()

			Expect(err).NotTo(HaveOccurred())
			Expect(result.Status).To(Equal(common.Denied))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(4))
			Expect(result.Details.Reasons[0]).To(Equal("MANIPULATED_DOCUMENT"))
			Expect(result.Details.Reasons[1]).To(Equal("1002 DOCUMENT_NUMBER"))
			Expect(result.Details.Reasons[2]).To(Equal("1007 SECURITY_CHECKS"))
			Expect(result.Details.Reasons[3]).To(Equal("Identity Verification: similarity = NO_MATCH | validity = FALSE | reason = SELFIE_MANIPULATED"))
			Expect(result.ErrorCode).To(Equal("100"))
			Expect(result.StatusCheck).To(BeNil())
		})

		It("should fail with error message", func() {
			r := &DetailsResponse{}

			err := json.Unmarshal(failedResponse, r)

			Expect(err).NotTo(HaveOccurred())

			result, err := r.toResult()

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("for some reason Jumio returned the 'FAILED' status for the verification transaction"))
			Expect(result.Status).To(Equal(common.Error))
			Expect(result.Details).NotTo(BeNil())
			Expect(result.Details.Finality).To(Equal(common.Unknown))
			Expect(result.Details.Reasons).To(HaveLen(1))
			Expect(result.Details.Reasons[0]).To(Equal("Identity Verification: similarity = MATCH | validity = TRUE"))
			Expect(result.ErrorCode).To(BeEmpty())
			Expect(result.StatusCheck).To(BeNil())
		})
	})
})
