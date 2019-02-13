package consumer

import (
	"fmt"

	"modulus/kyc/common"
)

// ApplicationResponseData defines the model for Response Data for a Consumer or Merchant KYC.
type ApplicationResponseData struct {
	// The current reputation of the user involved in this transaction.
	CurrentUserReputation EDNAPolicyResult `json:"user"`
	// The previous reputation of the User, that is, the reputation of the user the last time that it was evaluated.
	PreviousUserReputation EDNAPolicyResult `json:"upr"`
	// The score card for the current transaction.
	EdnaScoreCard ExternalizedTransactionScorecard `json:"ednaScoreCard"`
	// The name of the fraud rule that fired.
	FraudRuleName string `json:"frn"`
	// Result of fraud evaluation.
	FraudPolicyResult FraudPolicyResult `json:"frp"`
	// The description of the fraud rule that fired.
	FraudRuleDescription string `json:"frd"`
	// Result of automated review evaluation.
	ARPResult ReviewResult `json:"arpr"`
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
	// In the case of error the error message.
	ErrorMessage string `json:"error_message"`
}

// ExternalizedTransactionScorecard represents the score card for the current transaction.
type ExternalizedTransactionScorecard struct {
	// The test results for this transaction.
	ScoreCard []ConditionResult `json:"sc"`
	// The evaluated test results for this transaction.
	EvaluatedTestResults []ConditionResult            `json:"etr"`
	AutomatedResult      AutomatedReviewEngineResult  `json:"ar"`
	EvaluationResult     ExternalizedEvaluationResult `json:"er"`
}

// ConditionResult represents the result of the evaluation of a condition in a rule or security test.
type ConditionResult struct {
	// The id of security test or the key of the transaction data to which the condition applied.
	Test string `json:"test"`
	// Whether the condition fired or not.
	Fired bool `json:"fired"`
	// Textual result of the condition.
	Details string `json:"details"`
	// Indicates that the result is waiting for an asynchronous response from the customer and/or a third party service.
	WaitingForData bool `json:"waitingForData"`
	// The time in milliseconds UTC at which the result was created. Is only present in the result of Consumer and Merchant KYC.
	Timestamp int64 `json:"ts"`
	// The stage during which this result was created. Is only present in the result of Consumer and Merchant KYC.
	Stage string `json:"stage"`
}

// AutomatedReviewEngineResult represents the result of the automated review policy for this transaction.
type AutomatedReviewEngineResult struct {
	// Result of rule.
	Result ReviewResult `json:"result"`
	// The unique rule identifier.
	RuleID string `json:"ruleId"`
	// The rule name.
	RuleName string `json:"ruleName"`
	// The rule description.
	RuleDescription string `json:"ruleDescription"`
}

// ExternalizedEvaluationResult represents the result of the fraud policy evaluation for this transaction.
type ExternalizedEvaluationResult struct {
	// If multiple rules fired during evaluation then this is complete set of rules that fired. Otherwise it is absent.
	FiredRules []Rule `json:"firedRules"`
	// A rule that fired for the current transaction.
	ReportedRule Rule `json:"reportedRule"`
	// The name of the profile used for evaluation.
	Profile string `json:"profile"`
}

// Rule represents a rule that fired for the current transaction.
type Rule struct {
	// The rule description.
	Description string `json:"description"`
	// The unique rule identifier.
	RuleID int `json:"ruleId"`
	// Result of rule.
	ResultCode FraudPolicyResult `json:"resultCode"`
	// Details of the evaluation of this rule for the current transaction.
	Details string `json:"details"`
	// The results of the individual assertions of the rule.
	TestResults []ConditionResult `json:"testResults"`
	// The rule name.
	Name string `json:"name"`
}

// DocumentVerification is the part of ApplicationResponseData.
type DocumentVerification struct {
	RedirectURL string `json:"redirectURL"`
	RequestID   string `json:"requestId"`
}

// toResult processes the response and generates the verification result.
func (r *ApplicationResponseData) toResult() (result common.KYCResult, err error) {
	switch r.State {
	case UnderReview:
		result.Status = common.Denied
		reasons := []string{}
		reasons = append(reasons, "Some of checks triggered 'MANUAL REVIEW' status")
		reasons = append(reasons, "Profile: "+r.EdnaScoreCard.EvaluationResult.Profile)
		reasons = append(reasons, fmt.Sprintf("Rule: id %d | %s", r.EdnaScoreCard.EvaluationResult.ReportedRule.RuleID, r.EdnaScoreCard.EvaluationResult.ReportedRule.Description))
		for _, tr := range r.EdnaScoreCard.EvaluationResult.ReportedRule.TestResults {
			reasons = append(reasons, fmt.Sprintf("Test: '%s' | %s", tr.Test, tr.Details))
		}
		result.Details = &common.KYCDetails{Reasons: reasons}
		return
	case Accepted:
		result.Status = common.Approved
	case Rejected:
		result.Status = common.Denied
	default:
		err = fmt.Errorf("unknown state of the verification from the API: %s", r.State)
	}

	details := &common.KYCDetails{}

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

	if len(details.Reasons) > 0 {
		result.Details = details
	}

	return
}
