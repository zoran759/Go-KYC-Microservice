package consumer

import "modulus/kyc/common"

// The mapping from the common gender values to the API-acceptable values.
var gender2API = map[common.Gender]string{
	common.Male:   "M",
	common.Female: "F",
}
