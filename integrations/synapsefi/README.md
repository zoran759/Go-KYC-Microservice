# SynapseFI integration

This instruction describes how to use **`synapsefi`** package.

Current API version: **3.1**

## Configuration options description

| **Name**          | **Type**     | **Description**                                            |
| ----------------- | ------------ | ---------------------------------------------------------- |
| **Host***         | _**string**_ | SynapseFI API host. `https://api.synapsefi.com/v3.1/` for live mode, `https://uat-api.synapsefi.com/v3.1/` for the sandbox |
| **ClientID***     | _**string**_ | Client’s ID provided by SynapseFI                          |
| **ClientSecret*** | _**string**_ | Secret Key provided by SynapseFI                           |
| TimeoutThreshold  | _**int64**_  | Timeout for polling requests, optional. Default: 5 minutes |
| KYCFlow           | _**string**_ | User creating flow (see below), optional. Default: simple  |

## How to use the package

1) Create new [**config**](verification/contract.go#L3) for SynapseFI API usage.

2) Obtain a new service object by calling the [**New()**](synapsefi.go#L18) constructor.

3) Use service's verifier [**CheckCustomer**](synapsefi.go#L26) for the customer verification.

### Sample code

```go
customer := &common.UserData{
    ...
}

...

config = synapsefi.Config{
    Connection: synapsefi.Connection{
        Host:         "https://uat-api.synapsefi.com/v3.1/",
        ClientID:     "client_id_9tCAZNrlj3gOeUYGIKucqEb0pQmx6zy1W2VBasX7",
        ClientSecret: "client_secret_UWVu3Y5E82Amai9q1Tk0PGlXLhyZHf0rNCneSpos",
    },
    TimeoutThreshold: int64(time.Hour.Seconds()),   // optional, might be skipped
}

service := synapsefi.New(config)

result, err := service.CheckCustomer(customer)
...
```

### KYC flow

[Documentation](https://docs.synapsefi.com/docs/adding-documents) says the following

> Submitting more than one physical document at a time might result in a delay in response time. It is recommended that you perform separate calls to add each physical document if you have more than one.

Usually this is not a problem and the mentioned delay is barely more than 20 secs, but sometimes it becomes up to 5+ minutes (very rarely, but it happens).

In case when the such delay isn't critical for the application, you can use the default **"simple"** user registration flow, that issues the only one step and API call -- [CreateUser](https://docs.synapsefi.com/docs/create-a-user).

Otherwise, you can switch the **KYCFlow** configuration option to any other **string** value -- let's say, **"complex"**, i.e.:

```go
config = synapsefi.Config{
    Connection: synapsefi.Connection{
        Host:         "https://uat-api.synapsefi.com/v3.1/",
        ClientID:     "client_id_9tCAZNrlj3gOeUYGIKucqEb0pQmx6zy1W2VBasX7",
        ClientSecret: "client_secret_UWVu3Y5E82Amai9q1Tk0PGlXLhyZHf0rNCneSpos",
    },
    TimeoutThreshold: int64(time.Hour.Seconds()),   // optional, might be skipped
    KYCFlow: "complex",                             // optional, might be skipped
}
```

This flow implies a bit more complex sequence:
1) [CreateUser](https://docs.synapsefi.com/docs/create-a-user) without documents
2) [OAuth user](https://docs.synapsefi.com/docs/get-oauth_key-refresh-token) to obtain write permissions
3) [Upload user documents](https://docs.synapsefi.com/docs/adding-documents) using OAuth key

### Testing notes

Even in sandbox mode, SynapseFI uses USPS to verify addresses. It means that address like "123 test st Chikago TX" will be rejected with message "Supplied address is invalid / Unable to verify address".

Use [LOB](https://lob.com/products/address-verification) to check if an address is deliverable in the testing mode.
