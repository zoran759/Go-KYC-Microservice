# Sum&Substance integration

This instruction describes how to use **`sumsub`** package.

## Configuration options description

| **Name**   | **Type**     | **Description**                                                          |
| ---------- | ------------ | ------------------------------------------------------------------------ |
| **Host**   | _**string**_ | Sum&Substance base API host. It looks like `https://test-api.sumsub.com` |
| **APIKey** | _**string**_ | API key provided by Sum&Substance to you                                 |

## How to use the package

1) Create new [**config**](contract.go#L3) for Sumsub API usage.

2) Obtain a new service object by calling the [**New()**](sumsub.go#L25) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [**CheckCustomer**](sumsub.go#L44) to submit the customer data for the verification.

4) Use service's status checker [**CheckStatus**](sumsub.go#L94) for checking of the current state of the customer verification process.

### Sample code

```go
customer := &common.UserData{
    ...
}

...

// Create the config for the service.
config := sumsub.Config{
    Host: "host",
    APIKey: "APIKey",
}

// Obtain new service object for the verification.
service := sumsub.New(config)

// Submit the customer data to the service.
result, err := service.CheckCustomer(customer)

...

referenceID := result.StatusCheck.ReferenceID

...

// Check the current state of the customer verification.
result, err := service.CheckStatus(referenceID)
...
```
