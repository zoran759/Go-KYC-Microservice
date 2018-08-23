# IDology integration

This manual describes how to use `idology` package.

## How to use the package

1) Create new [config](../contract.go#9) for Idology API usage.

2) Obtain verifier's object by calling the [New()](../contract.go#25) constructor. As a parameter, pass it the configuration you created in step 1.

3) Use verifier's checkers for the verification. You may use either `.ExpectID` for identity verification or `.AlertList` for check against the alert list of bad players or fraudsters or them both.

4) The method for the customer check is `CheckCustomer()`.

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

result, details, err = verifier.AlertList.CheckCustomer(customer)
...
```
