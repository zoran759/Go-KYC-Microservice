# Trulioo integration

This manual describes how to use `trulioo` package.

## Configuration options description

* **`Host`** - type             _`string`_ - Trulioo base API host. It looks like `https://api.globaldatacompany.com`.
* **`NAPILogin`** - type         _`string`_ - Trulioo Normalized API login.
* **`NAPIPassword`** - type         _`string`_ - Trulioo Normalized API password.

## How to use the package

1) Create new [config](contract.go#L8) for Trulioo API usage.

2) Obtain a new service object by calling the [New()](service.go#L17) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [CheckCustomer](service.go#L12) for the customer verification.

Sample code:

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

result, details, err := service.CheckCustomer(customer)
...
```
