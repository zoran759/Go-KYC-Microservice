# Trulioo integration

This instruction describes how to use **`trulioo`** package.

## Configuration options description

| **Name**         | **Type**     | **Description**                                                          |
| ---------------- | ------------ | ------------------------------------------------------------------------ |
| **Host**         | _**string**_ | Trulioo base API host. It looks like `https://api.globaldatacompany.com` |
| **NAPILogin**    | _**string**_ | Trulioo Normalized API login                                             |
| **NAPIPassword** | _**string**_ | Trulioo Normalized API password                                          |

## How to use the package

1) Create new [**config**](contract.go#L9) for Trulioo API usage.

2) Obtain a new service object by calling the [**New()**](trulioo.go#L18) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [**CheckCustomer**](trulioo.go#L27) for the customer verification.

## Sample code

```go
customer := &common.UserData{
    ...
}

...

config := trulioo.Config{
    Host: "host",
    NAPILogin: "NAPILogin",
    NAPIPassword: "NAPIPassword",
}

service := trulioo.New(config)

result, err := service.CheckCustomer(customer)
...
```
