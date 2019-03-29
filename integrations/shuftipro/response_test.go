package shuftipro

import (
	"fmt"
	"testing"

	"modulus/kyc/common"

	"github.com/stretchr/testify/assert"
)

func TestToKYCResult(t *testing.T) {
	declinedReason := "Document originality is not verified."

	type testCase struct {
		name     string
		response Response
		result   common.KYCResult
	}

	testCases := []testCase{
		testCase{
			name: "Verification accepted",
			response: Response{
				Event: Accepted,
			},
			result: common.KYCResult{
				Status: common.Approved,
			},
		},
		testCase{
			name: "Verification declined",
			response: Response{
				Event:          Declined,
				DeclinedReason: declinedReason,
			},
			result: common.KYCResult{
				Status: common.Denied,
				Details: &common.KYCDetails{
					Reasons: []string{declinedReason},
				},
			},
		},
		testCase{
			name: "Verification cancelled",
			response: Response{
				Event: Cancelled,
			},
			result: common.KYCResult{
				Status: common.Denied,
				Details: &common.KYCDetails{
					Reasons: []string{fmt.Sprintf("Returned event cannot be processed: '%s'", Cancelled)},
				},
			},
		},
		testCase{
			name: "Empty event - pending verification",
			response: Response{
				Reference: "777",
			},
			result: common.KYCResult{
				Status: common.Unclear,
				StatusCheck: &common.KYCStatusCheck{
					Provider:    common.ShuftiPro,
					ReferenceID: "777",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res := tc.response.ToKYCResult()
			if tc.result.StatusCheck != nil {
				tc.result.StatusCheck.LastCheck = res.StatusCheck.LastCheck
			}
			assert.Equal(t, tc.result, res)
		})
	}
}
