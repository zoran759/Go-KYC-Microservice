package trulioo

import (
	"fmt"
	"modulus/kyc/common"
	"modulus/kyc/integrations/trulioo/configuration"
	"modulus/kyc/integrations/trulioo/verification"

	"github.com/pkg/errors"
)

// Trulioo defines the verification service.
type Trulioo struct {
	configuration configuration.Configuration
	verification  verification.Verification
}

// New constructs a new service object.
func New(config Config) Trulioo {

	return Trulioo{
		configuration: configuration.NewService(config.ToConfigurationConfig()),
		verification:  verification.NewService(config.ToVerificationConfig()),
	}
}

// CheckCustomer implements CustomerChecker interface for Trulioo.
func (service Trulioo) CheckCustomer(customer *common.UserData) (res common.KYCResult, err error) {
	if customer == nil {
		err = errors.New("No customer supplied")
		return
	}

	consents, err := service.configuration.Consents(customer.CountryAlpha2)
	if err != nil {
		return
	}

	dataFields := verification.MapCustomerToDataFields(customer)

	response, err := service.verification.Verify(customer.CountryAlpha2, consents, dataFields)
	if err != nil {
		return
	}

	if len(response.Errors) > 0 {
		err = response.Errors
		return
	}

	if response.Record.RecordStatus == Match {
		res.Status = common.Approved
		return
	}
	if len(response.Record.Errors) > 0 {
		err = response.Record.Errors
		return
	}

	res.Details = &common.KYCDetails{}

	for _, result := range response.Record.DatasourceResults {
		status := ""

		if result.DatasourceStatus != "" {
			status += fmt.Sprintf("status: %s; ", result.DatasourceStatus)
		}

		fieldsStatuses := ""
		for _, field := range result.DatasourceFields {
			if field.Status != "" {
				fieldsStatuses += fmt.Sprintf("%s : %s; ", field.FieldName, field.Status)
			}
		}

		if fieldsStatuses != "" {
			status += fmt.Sprintf("field statuses: %s", fieldsStatuses)
		}

		if result.Errors != nil && len(result.Errors) != 0 {
			status += fmt.Sprintf(
				"error: %s",
				result.Errors.Error(),
			)
		}

		if status != "" {
			res.Details.Reasons = append(
				res.Details.Reasons,
				fmt.Sprintf(
					"Datasource %s has %s",
					result.DatasourceName,
					status,
				),
			)
		}

	}

	if response.Record.RecordStatus == NoMatch {
		res.Status = common.Denied
	} else {
		res.Status = common.Unclear
	}

	return
}
