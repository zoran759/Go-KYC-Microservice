package common

// KYCStatus defines the verification status of the KYC check.
type KYCStatus int

// Possible KYCStatus values.
const (
	Error KYCStatus = iota
	Approved
	Denied
	Pending
	Queued
	Unclear
)

// KYCFinality defines the finality of the result of the KYC check.
type KYCFinality int

// Possible KYCFinality values.
const (
	Unknown KYCFinality = iota
	Final
	NonFinal
)

// Gender defines user's gender.
type Gender int

// Gender values.
const (
	Male   Gender = 1
	Female Gender = 2
)

// KYCProvider represents a KYC Provider.
type KYCProvider string

// Possible values of KYCProvider.
const (
	IdentityMind KYCProvider = "IdentityMind"
	IDology      KYCProvider = "IDology"
	ShuftiPro    KYCProvider = "ShuftiPro"
	SumSub       KYCProvider = "Sum&Substance"
	Trulioo      KYCProvider = "Trulioo"
)

// KYCProviders enumerates the implemented KYC providers.
var KYCProviders = map[KYCProvider]bool{
	IdentityMind: true,
	IDology:      true,
	ShuftiPro:    true,
	SumSub:       true,
	Trulioo:      true,
}
