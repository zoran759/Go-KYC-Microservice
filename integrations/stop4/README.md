# 4Stop integration

This instruction describes how to use **`stop4`** package.

## Configuration options description

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Host*** | _**string**_ | 4Stop API host, `https://private-f1649f-coreservices2.apiary-mock.com`
| **MerchantID*** | _**string**_ | Your merchant ID provided by 4Stop |
| **Password*** | _**string**_ | Your merchant password provided by 4Stop |

## How to use the package

1) Create new [**config**](verification/contract.go#L4) for 4Stop API usage.

2) Obtain a new service object by calling the [**New()**](service.go#L17) constructor.

3) Use service's verifier [**CheckCustomer**](service.go#L24) for the customer verification.

### Sample code

```go
customer := &common.UserData{
    ...
}

...

config = stop4.Config{
    Connection: stop4.Connection{
        Host:         "https://private-f1649f-coreservices2.apiary-mock.com",
        MerchantID:   "testmerchant123",
        Password:     "testpassword123",
    },
}

service := service.New(config)

result, err := service.CheckCustomer(customer)
...
```

### Testing notes
In the API [documentation](https://coreservices2.docs.apiary.io/#reference/0/test-cases) you can find a list of test cases and possible values

 

   