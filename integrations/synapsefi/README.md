# SynapseFI integration

This instruction describes how to use **`synapsefi`** package.

Current API version: **3.1**

## Configuration options description

| **Name**         | **Type**     | **Description**                                            |
| ---------------- | ------------ | ---------------------------------------------------------- |
| **Host**         | _**string**_ | SynapseFI API host. `https://api.synapsefi.com/v3.1/` for live mode, `https://uat-api.synapsefi.com/v3.1/` for the sandbox |
| **ClientID**     | _**string**_ | Client’s ID provided by SynapseFI                          |
| **ClientSecret** | _**string**_ | Secret Key provided by SynapseFI                           |

## How to use the package

1) Create new [**config**](verification/contract.go#L3) for SynapseFI API usage.

2) Obtain a new service object by calling the [**New()**](synapsefi.go#L17) constructor.

3) Use service's verifier [**CheckCustomer**](synapsefi.go#L24) for the customer verification.

4) Use service's KYC status checker [**CheckStatus**](synapsefi.go#L71) for checking of the current state of the customer verification process.

## Sample code

```go
customer := &common.UserData{
    ...
}

...

config = synapsefi.Config{
    Host:         "https://uat-api.synapsefi.com/v3.1/",
    ClientID:     "client_id_9tCAZNrlj3gOeUYGIKucqEb0pQmx6zy1W2VBasX7",
    ClientSecret: "client_secret_UWVu3Y5E82Amai9q1Tk0PGlXLhyZHf0rNCneSpos",
}

service := synapsefi.New(config)

result, err := service.CheckCustomer(customer)

...

referenceID := result.StatusCheck.ReferenceID

...

result, err := service.CheckStatus(referenceID)
...
```

### Testing notes

Even in sandbox mode, SynapseFI uses USPS to verify addresses. It means that address like "123 test st Chikago TX" will be rejected with message "Supplied address is invalid / Unable to verify address".

Use [LOB](https://lob.com/products/address-verification) to check if an address is deliverable in the testing mode.
