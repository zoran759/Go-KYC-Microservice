package jumio

import "gitlab.com/lambospeed/kyc/common"

// The mapping from the common document types to the API-acceptable values.
var documentTypeToIDType = map[common.DocumentType]IDType{
	// FIXME: does "VISA" idType have the counterpart in common.DocumentType?
	common.IDCard:   IDCard,
	common.Passport: Passport,
	common.Drivers:  DrivingLicense,
}

// The mapping from the common document types to the usual names.
var docTypeToName = map[common.DocumentType]string{
	common.IDCard:   "id card",
	common.Passport: "passport",
	common.Drivers:  "driving license",
}
