# KYC Package

> _**This is the initial specification.**_

## Integrations

The main package will initiate a KYC check with a chosen KYC provider.

The KYC provider is changed by changing the package being imported to the _kyc_provider_ named import.

To create a new integration, simply add a new subpackage under integrations.

An example integration package is provided. As you can see, all an integration package must do is provide a single exported method:

```go
func CheckCustomer(customer *common.UserData) bool
```

The main package will call this method in a goroutine to perform a check. The method should block until a response is has been received from the KYC provider.

---

## **Specification of the current implementation**

### **Common part**

TODO: write this.

### **Specific KYC providers**

KYC providers have different configuration options so twas inevitable to implement a specific config for each one of them. But mostly they are identical.

For instructions on integration of a specific KYC provider, please, refer this list:

* [IDology](integrations/idology/README.md)
* [Sum&Substance](integrations/sumsub/README.md)
* [Trulioo](integrations/trulioo/README.md)
* [Shufti Pro](integrations/shuftipro/README.md)

### **Supported fields and format**

TODO: write this.
