# IDology integration

This manual describes how to use `idology` package.

## How to use the package

1) Create new [config](contract.go#L8) for Idology API usage.

2) Obtain a new service object by calling the [New()](service.go#L17) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [ExpectID](contract.go#L12) for the customer verification.

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

service := idology.New(config)

result, details, err := service.ExpectID.CheckCustomer(customer)
...
```
