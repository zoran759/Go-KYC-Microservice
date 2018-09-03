package consumer

import (
	"fmt"

	"gitlab.com/lambospeed/kyc/common"
)

// ApplicationResponseData defines the model for Response Data for a Consumer or Merchant KYC.
type ApplicationResponseData struct {
	// The current reputation of the user involved in this transaction.
	CurrentUserReputation EDNAPolicyResult `json:"user"`
	// The previous reputation of the User, that is, the reputation of the user the last time that it was evaluated.
	PreviousUserReputation EDNAPolicyResult `json:"upr"`
	// FIXME: do we really need this?
	// ednaScoreCard:	ExternalizedTransactionScorecard{...}
	// The name of the fraud rule that fired.
	FraudRuleName string `json:"frn"`
	// Result of fraud evaluation.
	FraudPolicyResult FraudPolicyResult `json:"frp"`
	// The description of the fraud rule that fired.
	FraudRuleDescription string `json:"frd"`
	// Result of automated review evaluation.
	ARPResult AutomatedReviewPolicyResult `json:"arpr"`
	// The description, if any, of the automated review rule that fired.
	ARPDescription string `json:"arpd"`
	// The id, if any, of the automated review rule that fired.
	ARPID string `json:"arpid"`
	// FIXME: do we really need this?
	// graphScore:	GraphRiskScore{...}
	// The transaction id for this KYC. This id should be provided on subsequent updates to the KYC.
	KYCTxID string `json:"mtid"`
	// The current state of the KYC. A - Accepted; R - Under Review; D - Rejected.
	State KYCState `json:"state"`
	// FIXME: do we really need this?
	// oowQuestions:	QuestionsWrapper{...}
	ACVerification string `json:"acVerification"`
	// FIXME: do we really need this?
	DocVerification     *DocumentVerification `json:"docVerification"`
	SMSVerification     string                `json:"smsVerification"`
	OwnerApplicationIDs []string              `json:"ownerApplicationIds"`
	ParentMerchant      string                `json:"parentMerchant"`
	MerchantAPIName     string                `json:"merchantAPIName"`
	// A description of the reason for the User’s reputation.
	ReputationReasonDescription string `json:"erd"`
	// Result of policy evaluation. Combines the result of fraud and automated review evaluations.
	Result FraudPolicyResult `json:"res"`
	// The set of result codes from the evaluation of the current transaction.
	TxResultCodes string `json:"rcd"`
	// The transaction id of the current transaction. If no “tid” data was provided in the request data then a unique id will be generated. No assumptions should be made about the format of the generated id and it will be a maximimum length of 64 alphanumeric characters.
	TxID string `json:"tid"`
}

// DocumentVerification is the part of ApplicationResponseData.
type DocumentVerification struct {
	RedirectURL string `json:"redirectURL"`
	RequestID   string `json:"requestId"`
}

// toResult processes the response and generates the verification result.
func (r *ApplicationResponseData) toResult() (result common.KYCResult, details *common.DetailedKYCResult, err error) {
	switch r.State {
	case Accepted:
		result = common.Approved
	case Rejected:
		result = common.Denied
	}

	details = &common.DetailedKYCResult{
		Finality: common.Unknown,
	}
	if len(r.CurrentUserReputation) > 0 {
		details.Reasons = append(details.Reasons, fmt.Sprintf("Customer reputation: %s", r.CurrentUserReputation))
	}
	if len(r.FraudPolicyResult) > 0 {
		details.Reasons = append(details.Reasons, fmt.Sprintf("Fraud policy evaluation result: %s", r.FraudPolicyResult))
	}
	if len(r.ReputationReasonDescription) > 0 {
		details.Reasons = append(details.Reasons, fmt.Sprintf("Customer reputation reason: %s", r.ReputationReasonDescription))
	}
	if len(r.Result) > 0 {
		details.Reasons = append(details.Reasons, fmt.Sprintf("Combined fraud and automated review evaluations result: %s", r.Result))
	}

	return
}
