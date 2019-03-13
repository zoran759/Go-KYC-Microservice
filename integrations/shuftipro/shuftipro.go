package shuftipro

import (
	"errors"

	"modulus/kyc/common"
)

var _ common.KYCPlatform = ShuftiPro{}

// ShuftiPro represents the verification service.
type ShuftiPro struct {
	client Client
}

// New constructs a new verification service object.
func New(config Config) ShuftiPro {
	return ShuftiPro{
		client: NewClient(config),
	}
}

// CheckCustomer implements KYCPlatform interface for ShuftiPro.
func (s ShuftiPro) CheckCustomer(customer *common.UserData) (result common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	result, err = s.client.CheckCustomer(customer)

	return
}

// CheckStatus implements KYCPlatform interface for the ShuftiPro.
func (ShuftiPro) CheckStatus(referenceID string) (res common.KYCResult, err error) {
	err = errors.New("Shufti Pro doesn't support a verification status check")
	return
}
