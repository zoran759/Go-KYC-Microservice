# Jumio integration

This instruction describes how to use **`jumio`** package.

## Configuration options description

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Host** | _**string**_ | Endpoint URL of the Jumio API. It looks like `https://netverify.com/api/netverify/v2/performNetverify` |
| **Token** | _**string**_ | Jumio API token |
| **Secret** | _**string**_ | Jumio API secret |

## How to use the package

1) Create new [**config**](contract.go#L9) for Jumio API usage.

2) Obtain the new service object by calling the [**New()**](service.go#L15) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use the service object for the customer verification.

4) For the convenience, the package contains [**USendpoint**](contract.go#L5) and [**EUendpoint**](contract.go#L6) constants for Jumio API Endpoints URLs.

## Sample code

```go
customer := &common.UserData{
    ...
}

...

config := jumio.Config{
    Host: jumio.USendpoint,
    Token: "token",
    Secret: "secret",
}

service := jumio.New(config)

result, details, err := service.CheckCustomer(customer)
...
```
