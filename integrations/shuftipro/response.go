package shuftipro

import (
	"fmt"
	"modulus/kyc/common"
)

// List of ResultValue values.
const (
	AcceptedValue ResultValue = 1
	DeclinedValue ResultValue = 0
)

// List of Event values.
const (
	Accepted        Event = "verification.accepted"
	Declined        Event = "verification.declined"
	Cancelled       Event = "verification.cancelled"
	StatusChanged   Event = "verification.status.changed"
	ReqPending      Event = "request.pending"
	ReqInvalid      Event = "request.invalid"
	ReqUnauthorized Event = "request.unauthorized"
	ReqDeleted      Event = "request.deleted"
)

var event2description = map[Event]string{
	ReqInvalid:      "request parameters provided in the request are invalid; ",
	ReqUnauthorized: "the information provided in authorization header is invalid; ",
}

// Response represents a response of the Shufti Pro Verification API.
type Response struct {
	Reference      string      `json:"reference"`
	Event          Event       `json:"event"`
	Error          interface{} `json:"error"`
	Token          string      `json:"token"`
	Result         *Result     `json:"verification_result"`
	DeclinedReason string      `json:"declined_reason"`
}

// Event represents the status of request.
type Event string

// ResultValue represents verification result value.
type ResultValue int

// Result represents verification result.
type Result struct {
	Face             *ResultValue            `json:"face"`
	Document         *DocumentResult         `json:"document"`
	Address          *AddressResult          `json:"address"`
	BackgroundChecks *BackgroundChecksResult `json:"background_checks"`
}

// DocumentResult represents document verification result.
type DocumentResult struct {
	SelectedType              *ResultValue `json:"selected_type"`
	Name                      *ResultValue `json:"name"`
	BirthDate                 *ResultValue `json:"dob"`
	Number                    *ResultValue `json:"document_number"`
	IssueDate                 *ResultValue `json:"issue_date"`
	ExpiryDate                *ResultValue `json:"expiry_date"`
	Document                  *ResultValue `json:"document"`
	CustomerLooksLikeXYearOld *ResultValue `json:"customer_looks_like_x_year_old"`
	FaceOnDocumentMatched     *ResultValue `json:"face_on_document_matched"`
	Visibility                *ResultValue `json:"document_visibility"`
	MustNotBeExpired          *ResultValue `json:"document_must_not_be_expired"`
	Country                   *ResultValue `json:"document_country"`
}

// AddressResult represents address verification result.
type AddressResult struct {
	SelectedType             *ResultValue `json:"selected_type"`
	FullAddress              *ResultValue `json:"full_address"`
	Name                     *ResultValue `json:"name"`
	DocumentVisibility       *ResultValue `json:"address_document_visibility"`
	DocumentMustNotBeExpired *ResultValue `json:"address_document_must_not_be_expired"`
	DocumentCountry          *ResultValue `json:"address_document_country"`
	Document                 *ResultValue `json:"address_document"`
}

// BackgroundChecksResult represents background checks result.
type BackgroundChecksResult struct {
	Name        *ResultValue `json:"name"`
	DateOfBirth *ResultValue `json:"dob"`
}

// ToKYCResult converts Shufti Pro API response to the KYC result.
func (r Response) ToKYCResult() common.KYCResult {
	res := common.KYCResult{}

	switch r.Event {
	case Accepted:
		res.Status = common.Approved
	case Declined:
		res.Status = common.Denied
		if len(r.DeclinedReason) > 0 {
			res.Details = &common.KYCDetails{
				Reasons: []string{r.DeclinedReason},
			}
		}
	default:
		res.Status = common.Denied
		res.Details = &common.KYCDetails{
			Reasons: []string{fmt.Sprintf("Returned event cannot be processed: '%s'", r.Event)},
		}
	}

	return res
}
