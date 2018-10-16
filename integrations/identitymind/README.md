# IdentityMind integration

This instruction describes how to use **`identitymind`** package.

## Configuration options description

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Host** | _**string**_ | Endpoint URL of the IdentityMind API. It looks like `https://edna.identitymind.com/im` |
| **Username** | _**string**_ | IdentityMind API username |
| **Password** | _**string**_ | IdentityMind API password/license key |

## How to use the package

1) Create new [**config**](contract.go#L12) for IdentityMind API usage.

2) Obtain a new service object by calling the [**New()**](service.go#L21) constructor. As the parameter, pass it the configuration you created in step 1.

3) Use service's verifier [**CheckCustomer**](service.go#L28) for the customer verification.

4) For the convenience, the package contains [**ProductionBaseURL**](contract.go#L9) constant for IdentityMind API Endpoint URL.

5) If the service returned the data for KYC status checking use service's status checker [**CheckStatus**](service.go#L35) for checking the current state of the customer verification process.

## Sample code

```go
customer := &common.UserData{
    ...
}

...

config := identitymind.Config{
    Host: identitymind.ProductionBaseURL,
    Username: "username",
    Password: "password",
}

service := identitymind.New(config)

result, err := service.CheckCustomer(customer)
...

customerID := result.StatusPolling.CustomerID

...

// Check the current state of the customer verification.
result, err := service.CheckStatus(customerID)
...

```
