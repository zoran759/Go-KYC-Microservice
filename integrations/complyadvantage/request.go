package complyadvantage

import (
	"time"

	"modulus/kyc/common"
)

// Filters represents filters within the search to narrow down the results.
type Filters struct {
	Types          []string `json:"types,omitempty"`
	BirthYear      int      `json:"birth_year,omitempty"`
	RemoveDeceased int      `json:"remove_deceased,omitempty"`
}

// Request represents search request.
type Request struct {
	SearchTerm    string            `json:"search_term"`
	ClientRef     string            `json:"client_ref,omitempty"`
	SearchProfile string            `json:"search_profile,omitempty"`
	Fuzziness     float32           `json:"fuzziness"`
	Offset        int               `json:"offset,omitempty"`
	Limit         int               `json:"limit,omitempty"`
	Filters       Filters           `json:"filters,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	ShareURL      int               `json:"share_url,omitempty"`
}

// Note that search_profile and types are mutually exclusive, and only one of these two options should be provided.

// newRequest constructs new Request object from the customer data.
func (s service) newRequest(customer *common.UserData) Request {
	r := Request{
		SearchTerm: customer.Fullname(),
		Fuzziness:  s.fuzziness,
	}

	if !time.Time(customer.DateOfBirth).IsZero() {
		r.Filters.BirthYear = time.Time(customer.DateOfBirth).Year()
	}

	return r
}
