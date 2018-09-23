package jumio

import "gitlab.com/lambospeed/kyc/common"

// IDType represents document types acceptable by the API.
type IDType string

// The mapping from the common document types to the API-acceptable values.
var documentTypeMap = map[common.DocumentType]IDType{
	// FIXME: does "VISA" idType have the counterpart in common.DocumentType?
	common.IDCard:   "ID_CARD",
	common.Passport: "PASSPORT",
	common.Drivers:  "DRIVING_LICENSE",
}
