# Jumio integration

This instruction describes how to use **`jumio`** package.

## Configuration options description

| **Name**    | **Type**     | **Description**  |
| ----------- | ------------ | ---------------- |
| **BaseURL** | _**string**_ | Base URL for requests to the Jumio performNetverify API. It looks like `https://netverify.com/api/netverify/v2` |
| **Token**   | _**string**_ | Jumio API token  |
| **Secret**  | _**string**_ | Jumio API secret |

## How to use the package

1) Create new [**config**](contract.go#L9) for Jumio API usage.

2) Obtain the new service object by calling the [**New()**](service.go#L20) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use the service object for the customer verification.

4) For the convenience, the package contains [**USbaseURL**](contract.go#L5) and [**EUbaseURL**](contract.go#L6) constants for Jumio performNetverify API base URLs.

## Sample code

```go
customer := &common.UserData{
    ...
}

...

config := jumio.Config{
    BaseURL: jumio.USbaseURL,
    Token: "token",
    Secret: "secret",
}

service := jumio.New(config)

// Check the customer.
result, err := service.CheckCustomer(customer)
...

referenceID := result.StatusCheck.ReferenceID

...

// Check the current state of the customer verification.
result, err := service.CheckStatus(referenceID)
...

```
