package trulioo

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/configuration"
	"gitlab.com/lambospeed/kyc/integrations/trulioo/verification"
	"fmt"
	"github.com/pkg/errors"
)

type Trulioo struct {
	configuration configuration.Configuration
	verification  verification.Verification
}

func New(config Config) Trulioo {

	return Trulioo{
		configuration: configuration.NewService(config.ToConfigurationConfig()),
		verification:  verification.NewService(config.ToVerificationConfig()),
	}
}

func (service Trulioo) CheckCustomer(customer *common.UserData) (common.KYCResult, *common.DetailedKYCResult, error) {
	if customer == nil {
		return common.Error, nil, errors.New("No customer supplied")
	}

	consents, err := service.configuration.Consents(customer.CountryAlpha2)
	if err != nil {
		return common.Error, nil, err
	}

	dataFields := verification.MapCustomerToDataFields(*customer)

	response, err := service.verification.Verify(customer.CountryAlpha2, consents, dataFields)
	if err != nil {
		return common.Error, nil, err
	}

	if response.Errors != nil && len(response.Errors) > 0 {
		return common.Error, nil, response.Errors
	}

	if response.Record.RecordStatus == Match {
		return common.Approved, nil, nil
	} else {
		if response.Record.Errors != nil && len(response.Record.Errors) > 0 {
			return common.Error, nil, response.Record.Errors
		}

		detailedResult := common.DetailedKYCResult{
			Finality: common.Unknown,
		}

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
				detailedResult.Reasons = append(
					detailedResult.Reasons,
					fmt.Sprintf(
						"Datasource %s has %s",
						result.DatasourceName,
						status,
					),
				)
			}

		}

		if response.Record.RecordStatus == NoMatch {
			return common.Denied, &detailedResult, nil
		}

		return common.Unclear, &detailedResult, nil
	}
}
