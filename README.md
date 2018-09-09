# KYC Package

> _**This is the initial specification. For the current implementation see [SPECIFICATION OF THE CURRENT IMPLEMENTATION](#specification-of-the-current-implementation)**_

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

## **SPECIFICATION OF THE CURRENT IMPLEMENTATION**

### **Table of contents**

* **[Integration interface](#integration-interface)**
* **[KYC request](#kyc-request)**
* **[KYC response](#kyc-response)**
* **[Specific KYC providers](#specific-kyc-providers)**

### **Integration interface**

All KYC providers implement [**common.CustomerChecker**](common/contract.go#L3) interface for the verification process:

```go
type CustomerChecker interface {
    CheckCustomer(customer *UserData) (KYCResult, *DetailedKYCResult, error)
}
```

Providers are configurable by their configs. Configuration options for each provider are described in the respective integration instructions in [Specific KYC providers](#specific-kyc-providers).

The rest required for interaction with KYC providers is in the **`common`** package including request and response structures.

### **KYC request**

For the verification request use [**common.UserData**](common/model.go#L8) structure.

#### **[UserData](common/model.go#L8) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **FirstName** | _**string**_ | _required_. First name of the customer, for ex. "John" |
| **PaternalLastName** | _**string**_ | paternal last name of the customer |
| **LastName** | _**string**_ | _required_. Last name of the customer, for ex. "Doe" |
| **MiddleName** | _**string**_ | middle name of the customer, for ex. "Benedikt" |
| **LegalName** | _**string**_ | legal name of the customer, for ex. "Foobar Co." |
| **LatinISO1Name** | _**string**_ | latin ISO1 name of the customer |
| **Email** | _**string**_ | email of the customer |
| **Gender** | [_**Gender**_](common/enum.go#L27) | gender of the customer |
| **DateOfBirth** | [_**Time**_](common/model.go#L65) | date of birth of the customer |
| **PlaceOfBirth** | _**string**_ | place of birth of the customer |
| **CountryOfBirthAlpha2** | _**string**_ | country of birth of the customer in ISO 3166-1 alpha-2 format, for ex. "US" |
| **StateOfBirth** | _**string**_ | state of birth of the customer, for ex. "GA" |
| **CountryAlpha2** | _**string**_ | country of the customer in ISO 3166-1 alpha-2 format, for ex. "DE" |
| **Nationality** | _**string**_ | citizenship of the customer. Perhaps, it should be country's name, for ex. "Italy" |
| **Phone** | _**string**_ | primary phone of the customer. It isn't a mobile phone! |
| **MobilePhone** | _**string**_ | mobile phone of the customer |
| **AddressString** | _**string**_ | current address of the customer in the string format. It should conform to a standard form of postal address |
| **CurrentAddress** | [_**Address**_](#address-fields-description) | current address of the customer |
| **SupplementalAddresses** | _**[][Address](#address-fields-description)**_ | array of supplemental addresses of the customer |
| **Documents** | _**[][Document](#document-fields-description)**_ | array of documents of the customer |
| **Business** | [_**Business**_](#business-fields-description) | the business which the customer is linked to or is one of the owners |

#### **[Address](common/model.go#L33) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **CountryAlpha2** | _**string**_ | country in ISO 3166-1 alpha-2 format, for ex. "FR" |
| **County** | _**string**_ | county if applicable |
| **State** | _**string**_ | name of the state, for ex. "Alabama" |
| **Town** | _**string**_ | city or town name |
| **Suburb** | _**string**_ | suburb if applicable |
| **Street** | _**string**_ | street name, for ex. "PeachTree Place", "7th street" |
| **StreetType** | _**string**_ | street type, for ex. "Avenue" |
| **SubStreet** | _**string**_ | substreet if applicable |
| **BuildingName** | _**string**_ | building or house name |
| **BuildingNumber** | _**string**_ | building or house number |
| **FlatNumber** | _**string**_ | flat or apartment number |
| **PostOfficeBox** | _**string**_ | post office box |
| **PostCode** | _**string**_ | zip or postal code |
| **StateProvinceCode** | _**string**_ | abbreviated name of the state, for ex. "CA" |
| **StartDate** | [_**Time**_](common/model.go#L65) | when the customer settled into this address |
| **EndDate** | [_**Time**_](common/model.go#L65) | when the customer moved out from this address |

#### **[Document](common/model.go#L84) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Metadata** | [_**DocumentMetadata**_](#documentmetadata-fields-description) | document info |
| **Front** | _**[*DocumentFile](#documentfile-fields-description)**_ | front-side document image |
| **Back** | _**[*DocumentFile](#documentfile-fields-description)**_ | back-side document image |

#### **[DocumentMetadata](common/model.go#L91) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Type** | [_**DocumentType**_](common/enum.go#L36) | the document type |
| **Country** | _**string**_ | country name where the document was issued, for ex. "JAPAN" |
| **DateIssued** | [_**Time**_](common/model.go#L65) | the date when the document was issued |
| **ValidUntil** | [_**Time**_](common/model.go#L65) | the date to which the document is valid |
| **Number** | _**string**_ | the document number |
| **CardFirst6Digits** | _**string**_ | first six digits of the document number if applicable (SSN, SNILS, banking card, etc.) |
| **CardLast4Digits** | _**string**_ | last four digits of the document number if applicable (SSN, SNILS, banking card, etc.) |

#### **[DocumentFile](common/model.go#L102) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Filename** | _**string**_ | file name of the document image, for ex. "passport_front.jpg" |
| **ContentType** | _**string**_ | mime type of the content, for ex. "image/jpeg" |
| **Data** | _**[]byte**_ | raw content of the document image file |

#### **[Business](common/model.go#L76) fields description**

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Name** | _**string**_ | name of the Enterprise the customer relates to |
| **RegistrationNumber** | _**string**_ | registration number of the Enterprise |
| **IncorporationDate** | [_**Time**_](common/model.go#L65) | incorporation date of the Enterprise |
| **IncorporationJurisdiction** | _**string**_ | incorporation jurisdiction of the Enterprise |

### **KYC response**

The verification response consist of three elements: a result, a detailed result and an error if occurred.

The result is of type [**common.KYCResult**](common/enum.go#L3) and may hold following values:

| **Value** | **Description** |
| --------- | --------------- |
| **Error** | the verification has failed. That must mean that some error has occurred. Returned error value must be non-nil |
| **Approved** | successful verification with approved result. The detailed result maybe non-nil and contain additional info about the verification |
| **Denied** | successful verification with rejected result. The detailed result must be non-nil and contain additional info about the verification |
| **Unclear** | the verification completed with an indefinite result. That might mean that some additional info is required. The detailed result must be non-nil and contain additional info  |

The detailed result is of type [***common.DetailedKYCResult**](common/model.go#L109) and consist of the following fields:

| **Name** | **Type** | **Description** |
| -------- | -------- | --------------- |
| **Finality** | [_**KYCFinality**_](common/enum.go#L17) | finality of the result. Possible values are `Final`, `NonFinal` and `Unknown`. Not all providers support "finality" property hence tristate value |
| **Reasons** | _**[]string**_ | array of additional service responses describing result-related circumstances |

### **Specific KYC providers**

KYC providers have different configuration options so twas inevitable to implement a specific config for each one of them. But mostly they are identical.

For instructions on integration of a specific KYC provider, please, refer this list:

* [**IDology**](integrations/idology/README.md)
* [**Sum&Substance**](integrations/sumsub/README.md)
* [**Trulioo**](integrations/trulioo/README.md)
* [**Shufti Pro**](integrations/shuftipro/README.md)

---
