# IDology integration

This manual describes how to use `idology` package.

## How to use the package

1) Create new [config](contract.go#L9) for Idology API usage.

2) Obtain verifier's object by calling the [New()](contract.go#L24) constructor. As a parameter, pass it the configuration you created in step 1.

3) Use verifier's checker [ExpectID](contract.go#L19) for the verification.

4) The method for the customer check is [CheckCustomer](contract.go#L32).

Sample code:

```go
customer := &common.UserData{
    ...
}

...

config := idology.Config{
    Host: "host",
    Username: "username",
    Password: "password",
}

verifier := idology.New(config)

result, details, err := verifier.ExpectID.CheckCustomer(customer)
...
```
