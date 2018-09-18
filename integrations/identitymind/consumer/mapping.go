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

// The mapping from the common document types to the usual names.
var docTypeToNameMap = map[common.DocumentType]string{
	common.IDCard:          "id card",
	common.Passport:        "passport",
	common.Drivers:         "driving license",
	common.UtilityBill:     "utility bill",
	common.SNILS:           "SNILS",
	common.ResidencePermit: "residence permit",
}
