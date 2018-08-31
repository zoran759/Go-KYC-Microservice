# Shufti Pro integration

This manual describes how to use `shuftipro` package.

## Configuration options description

* **`Host`** - type             _`string`_ - Shufti Pro base API host. It looks like `https://api.shuftipro.com/ `.
* **`ClientID`** - type         _`string`_ - Client’s ID provided by Shufti Pro to you.
* **`SecretKey`** - type         _`string`_ - Secret Key provided by Shufti Pro to you.
* **`RedirectURL`**, type _`string`_ - A redirect URL, to which the user will be redirected after the verification (Stated as required in the integration guide, but not used in the actual process)

## How to use the package

1) Create new [config](contract.go#L8) for Shufti Pro API usage.

2) Obtain a new service object by calling the [New()](service.go#L17) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [CheckCustomer](service.go#L12) for the customer verification.

Sample code:

```go
customer := &common.UserData{
    ...
}

...

config := shuftipro.Config{
    Host: "host",
    ClientID: "ClientID",
    SecretKey: "SecretKey",
    RedirectURL: "RedirectURL",
}

service := shuftipro.New(config)

result, details, err := service.CheckCustomer(customer)
...
```
