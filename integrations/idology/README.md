# IDology integration

This manual describes how to use `idology` package.

## Configuration options description

* **`Host`** - type             _`string`_ - a full URL of the IDology ExpectID® API Endpoint. It looks like `https://web.idologylive.com/api/idiq.svc`.
* **`Username`** - type         _`string`_ - an IDology API username (128 bytes limit).
* **`Password`** - type         _`string`_ - an IDology API password (255 bytes limit).
* **`UseSummaryResult`**, type _`bool`_ - use Summary Results instead of ExpectID results. This depends on the Enterprise Configuration in the web portal (IDCenter).

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
