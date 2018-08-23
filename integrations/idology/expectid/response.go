package expectid

import (
	"github.com/achiku/xml"
)

// SummaryResult defines "summary-result" part of the response.
type SummaryResult struct {
	XMLName xml.Name             `xml:"summary-result"`
	Key     SummaryResultKey     `xml:"key"`
	Message SummaryResultMessage `xml:"message"`
}

// Results defines "results" part of the response.
type Results struct {
	XMLName xml.Name  `xml:"results"`
	Key     ResultKey `xml:"key"`
	Message string    `xml:"message"`
}

// Qualifier defines "qualifiers" part of the response.
type Qualifier struct {
	XMLName xml.Name `xml:"qualifier"`
	Key     string   `xml:"key"`
	Message string   `xml:"message"`
}

// Response defines a response from IDology ExpectID® API.
type Response struct {
	XMLName       xml.Name      `xml:"response"`
	IDNumber      int           `xml:"id-number"`
	SummaryResult SummaryResult `xml:"summary-result"`
	Results       Results       `xml:"results"`
	Error         *string       `xml:"error"`
	Qualifiers    []Qualifier   `xml:"qualifiers"`
}
