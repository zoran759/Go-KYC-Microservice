package consumer

import "gitlab.com/lambospeed/kyc/common"

// The mapping from the common gender values to the API-acceptable values.
var genderMap = map[common.Gender]string{
	common.Male:   "M",
	common.Female: "F",
}

// The mapping from the common document types to the API-acceptable values.
var documentTypeMap = map[common.DocumentType]DocumentType{
	common.Passport:        Passport,
	common.Drivers:         DriverLicence,
	common.IDCard:          GovernmentIssuedIDCard,
	common.ResidencePermit: ResidencePermit,
	common.UtilityBill:     UtilityBill,
}
