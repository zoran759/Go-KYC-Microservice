# Sum&Substance integration

This instruction describes how to use **`sumsub`** package.

## Configuration options description

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Host** | _**string**_ | Sum&Substance base API host. It looks like `https://test-api.sumsub.com` |
| **APIKey** | _**string**_ | API key provided by Sum&Substance to you |
| **TimeoutTreshold** | _**int64**_ | A threshold, in seconds, after which polling stops and returns timeout error |

## How to use the package

1) Create new [**config**](contract.go#L3) for Sumsub API usage.

2) Obtain a new service object by calling the [**New()**](sumsub.go#L23) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [**CheckCustomer**](sumsub.go#L42) for the customer verification.

### Sample code

```go
customer := &common.UserData{
    ...
}

...

config := sumsub.Config{
    Host: "host",
    APIKey: "APIKey",
}

service := sumsub.New(config)

result, details, err := service.CheckCustomer(customer)
...
```
