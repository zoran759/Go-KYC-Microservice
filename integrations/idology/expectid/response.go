package expectid

import (
	"encoding/xml"
	"fmt"

	"gitlab.com/lambospeed/kyc/common"
)

// SummaryResult defines "summary-result" part in the response.
type SummaryResult struct {
	XMLName xml.Name         `xml:"summary-result"`
	Key     SummaryResultKey `xml:"key"`
	Message string           `xml:"message"`
}

// Results defines "results" part in the response.
type Results struct {
	XMLName xml.Name  `xml:"results"`
	Key     ResultKey `xml:"key"`
	Message string    `xml:"message"`
}

// Qualifiers defines "qualifiers" part in the response.
type Qualifiers struct {
	XMLName    xml.Name    `xml:"qualifiers"`
	Qualifiers []Qualifier `xml:"qualifier"`
}

// Qualifier defines "qualifier" part in the response.
type Qualifier struct {
	XMLName xml.Name `xml:"qualifier"`
	Key     string   `xml:"key"`
	Message string   `xml:"message"`
}

// PatriotAct defines "pa" part in the response.
type PatriotAct struct {
	XMLName     xml.Name `xml:"pa"`
	List        string   `xml:"list"`
	Score       int      `xml:"score"`
	DateOfBirth string   `xml:"dob"`
}

// Restriction defines "restriction" part in the response.
type Restriction struct {
	XMLName    xml.Name   `xml:"restriction"`
	Key        string     `xml:"key"`
	Message    string     `xml:"message"`
	PatriotAct PatriotAct `xml:"pa"`
}

// Response defines a response from IDology ExpectID® API.
type Response struct {
	XMLName       xml.Name      `xml:"response"`
	IDNumber      int           `xml:"id-number"`
	SummaryResult SummaryResult `xml:"summary-result"`
	Results       Results       `xml:"results"`
	Restriction   *Restriction  `xml:"restriction"`
	Qualifiers    *Qualifiers   `xml:"qualifiers"`
	Error         *string       `xml:"error"`
}

// toResult processes the response and generates the verification result.
func (r *Response) toResult(useSummaryResult bool) (result common.KYCResult, err error) {
	detailsCreateIfNil := func(details **common.KYCDetails) {
		if *details == nil {
			*details = &common.KYCDetails{}
		}
	}

	switch useSummaryResult {
	case true:
		switch r.SummaryResult.Key {
		case Success:
			result.Status = common.Approved
		case Failure:
			result.Status = common.Denied
		case Partial:
			result.Status = common.Unclear
		}
	case false:
		switch r.Results.Key {
		case Match:
			result.Status = common.Approved
		case NoMatch, MatchRestricted:
			result.Status = common.Denied
		}
	}

	if r.Restriction != nil {
		detailsCreateIfNil(&result.Details)
		result.Details.Reasons = []string{
			r.Restriction.Message,
			r.Restriction.PatriotAct.List,
			fmt.Sprintf("Patriot Act score: %d", r.Restriction.PatriotAct.Score),
		}
	}

	if r.Qualifiers != nil {
		detailsCreateIfNil(&result.Details)
		for _, q := range r.Qualifiers.Qualifiers {
			result.Details.Reasons = append(result.Details.Reasons, q.Message)
		}
	}

	return
}
