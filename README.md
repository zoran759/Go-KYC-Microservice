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
* **[Required fields](#required-fields)**
  * **[IDology](#idology)**
  * **[Sum&Substance](#sum&substance)**
  * **[Trulioo](#trulioo)**
  * **[Shufti Pro](#shufti-pro)**
* **[Applicable fields grouped per provider](#applicable-fields-grouped-per-provider)**
  * **[IDology](#fields-applicable-for-idology)**
  * **[Sum&Substance](#fields-applicable-for-sum&substance)**
  * **[Trulioo](#fields-applicable-for-trulioo)**
  * **[Shufti Pro](#fields-applicable-for-shufti-pro)**
* **[The countries supported by KYC providers and the fields variability](#the-countries-supported-by-kyc-providers-and-the-fields-variability)**

### **Integration interface**

All KYC providers implement [**common.CustomerChecker**](common/contract.go#L3) interface for the verification process:

```go
type CustomerChecker interface {
    CheckCustomer(customer *UserData) (KYCResult, error)
}
```

Providers are configurable by their configs. Configuration options for each provider are described in the respective integration instructions in [Specific KYC providers](#specific-kyc-providers).

The rest required for interaction with KYC providers is in the **`common`** package including request and response structures.

### **KYC request**

For the verification request use a request of the [**common.UserData**](common/model.go#L8) type.

#### **[UserData](common/model.go#L8) fields description**

| **Name**                     | **Type**                           | **Description**                                                       |
| ---------------------------- | ---------------------------------- | --------------------------------------------------------------------- |
| **FirstName**                | _**string**_                       | _Required_. First name of the customer, for ex. "John"                |
| **PaternalLastName**         | _**string**_                       | Paternal last name of the customer                                    |
| **LastName**                 | _**string**_                       | _Required_. Last name of the customer, for ex. "Doe"                  |
| **MiddleName**               | _**string**_                       | Middle name of the customer, for ex. "Benedikt"                       |
| **LegalName**                | _**string**_                       | Legal name of the customer, for ex. "Foobar Co."                      |
| **LatinISO1Name**            | _**string**_                       | Latin ISO1 name of the customer                                       |
| **Email**                    | _**string**_                       | Email of the customer                                                 |
| **Gender**                   | [_**Gender**_](common/enum.go#L27) | Gender of the customer                                                |
| **DateOfBirth**              | [_**Time**_](common/model.go#L132) | Date of birth of the customer                                         |
| **PlaceOfBirth**             | _**string**_                       | Place of birth of the customer                                        |
| **CountryOfBirthAlpha2**     | _**string**_                       | Country of birth of the customer in ISO 3166-1 alpha-2 format, for ex. "US" |
| **StateOfBirth**             | _**string**_                       | State of birth of the customer, for ex. "GA"                          |
| **CountryAlpha2**            | _**string**_                       | Country of the customer in ISO 3166-1 alpha-2 format, for ex. "DE"    |
| **Nationality**              | _**string**_                       | Citizenship of the customer. Perhaps, it should be country's name, for ex. "Italy" |
| **Phone**                    | _**string**_                       | Primary phone of the customer. It isn't the mobile phone!             |
| **MobilePhone**              | _**string**_                       | Mobile phone of the customer                                          |
| **CurrentAddress**           | [_**Address**_](#address-fields-description) | Current address of the customer                             |
| **SupplementalAddresses**    | _**[]Address**_                    | List of supplemental addresses of the customer                        |
| **Business**                 | _***[Business](#business-fields-description)**_ | The business which the customer is linked to or is one of the owners |
| **Passport**                 | _***[Passport](#passport-fields-description)**_               | Passport of the customer                   |
| **IDCard**                   | _***[IDCard](#idcard-fields-description)**_                   | Id card of the customer, for ex. US SSN    |
| **SNILS**                    | _***[SNILS](#snils-fields-description)**_                     | SNILS (Russian insurance number of individual ledger account) of the customer |
| **DriverLicense**            | _***[DriverLicense](#driverlicense-fields-description)**_     | Driver license of the customer             |
| **DriverLicenseTranslation** | _***[DriverLicenseTranslation](#driverlicensetranslation-fields-description)**_ | Driver license translation of the customer (translation of the driving license required in the target country) |
| **CreditCard**               | _***[CreditCard](#creditcard-fields-description)**_           | Banking credit card of the customer        |
| **DebitCard**                | _***[DebitCard](#debitcard-fields-description)**_             | Banking debit card of the customer         |
| **UtilityBill**              | _***[UtilityBill](#utilitybill-fields-description)**_         | Utility bill                               |
| **ResidencePermit**          | _***[ResidencePermit](#residencepermit-fields-description)**_ | Residence permit of the customer           |
| **Agreement**                | _***[Agreement](#agreement-fields-description)**_             | Agreement of some sort, e.g. for processing personal info      |
| **EmploymentCertificate**    | _***[EmploymentCertificate](#employmentcertificate-fields-description)**_ | Employment certificate of the customer (a document from an employer, e.g. proof that a user works there) |
| **Contract**                 | _***[Contract](#contract-fields-description)**_               | Some sort of contract                      |
| **DocumentPhoto**            | _***[DocumentPhoto](#documentphoto-fields-description)**_     | Document photo of the customer (like a photo from a passport)  |
| **Selfie**                   | _***[Selfie](#selfie-fields-description)**_                   | Selfie image of the customer               |
| **Avatar**                   | _***[Avatar](#avatar-fields-description)**_                   | A profile image aka avatar of the customer |
| **Other**                    | _***[Other](#other-fields-description)**_                     | Other document (should be used only when nothing else applies) |

#### **[Address](common/model.go#L47) fields description**

| **Name**              | **Type**     | **Description**                                         |
| --------------------- | ------------ | ------------------------------------------------------- |
| **CountryAlpha2**     | _**string**_ | Country in ISO 3166-1 alpha-2 format, for ex. "FR"      |
| **County**            | _**string**_ | County if applicable                                    |
| **State**             | _**string**_ | Name of the state, for ex. "Alabama"                    |
| **Town**              | _**string**_ | City or town name                                       |
| **Suburb**            | _**string**_ | Suburb if applicable                                    |
| **Street**            | _**string**_ | Street name, for ex. "PeachTree Place", "7th street"    |
| **StreetType**        | _**string**_ | Street type, for ex. "Avenue"                           |
| **SubStreet**         | _**string**_ | Substreet if applicable                                 |
| **BuildingName**      | _**string**_ | Building or house name                                  |
| **BuildingNumber**    | _**string**_ | Building or house number                                |
| **FlatNumber**        | _**string**_ | Apartment number                                        |
| **PostOfficeBox**     | _**string**_ | Post office box                                         |
| **PostCode**          | _**string**_ | Zip or postal code                                      |
| **StateProvinceCode** | _**string**_ | Abbreviated name of the state or province, for ex. "CA" |
| **StartDate**         | _**Time**_   | When the customer settled at this address               |
| **EndDate**           | _**Time**_   | When the customer moved out from this address           |

#### **[Business](common/model.go#L143) fields description**

| **Name**                      | **Type**     | **Description**                                |
| ----------------------------- | ------------ | ---------------------------------------------- |
| **Name**                      | _**string**_ | Name of the Enterprise the customer relates to |
| **RegistrationNumber**        | _**string**_ | Registration number of the Enterprise          |
| **IncorporationDate**         | _**Time**_   | Incorporation date of the Enterprise           |
| **IncorporationJurisdiction** | _**string**_ | Incorporation jurisdiction of the Enterprise   |

#### **[Passport](common/model.go#L193) fields description**

| **Name**          | **Type**                                                | **Description**                                         |
| ----------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_                                            | Passport number without whitespaces and dashes          |
| **Mrz1**          | _**string**_                                            | First line of the Machine Readable Zone (MRZ) of passport, 44 letters and digits, i.e. "P<CZESPECIMEN<<VZOR<<<<<<<<<<<<<<<<<<<<<<<<<" |
| **Mrz2**          | _**string**_                                            | Second line of the Machine Readable Zone (MRZ) of passport, 44 letters and digits, i.e. "99003853<1CZE1101018M1207046110101111<<<<<94" |
| **CountryAlpha2** | _**string**_                                            | Country in ISO 3166-1 alpha-2 format, for ex. "SG"      |
| **State**         | _**string**_                                            | Abbreviated name of the state or province, for ex. "TX" |
| **IssuedDate**    | _**Time**_                                              | Issued date                                             |
| **ValidUntil**    | _**Time**_                                              | Valid until date                                        |
| **Image**         | _***[DocumentFile](#documentfile-fields-description)**_ | Scan or photo of the passport                           |

#### **[IDCard](common/model.go#L205) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **Number**        | _**string**_        | Id card number without whitespaces and dashes      |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "CN" |
| **IssuedDate**    | _**Time**_          | Issued date                                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the card                          |

#### **[SNILS](common/model.go#L213) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **Number**     | _**string**_        | SNILS number without whitespaces and dashes |
| **IssuedDate** | _**Time**_          | Issued date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the SNILS                  |

#### **[DriverLicense](common/model.go#L220) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license number                                   |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "DE"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "KY" |
| **IssuedDate**    | _**Time**_          | Issued date                                             |
| **ValidUntil**    | _**Time**_          | Valid until date                                        |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license   |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license    |

#### **[DriverLicenseTranslation](common/model.go#L231) fields description**

| **Name**          | **Type**            | **Description**                                                   |
| ----------------- | ------------------- | ----------------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license translation number                                 |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "LV"                |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "MT"           |
| **IssuedDate**    | _**Time**_          | Issued date                                                       |
| **ValidUntil**    | _**Time**_          | Valid until date                                                  |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license translation |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license translation  |

#### **[CreditCard](common/model.go#L242) fields description**

| **Name**       | **Type**            | **Description**                                   |
| -------------- | ------------------- | ------------------------------------------------- |
| **Number**     | _**string**_        | Credit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date                                  |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the credit card |

#### **[DebitCard](common/model.go#L249) fields description**

| **Name**       | **Type**            | **Description**                                  |
| -------------- | ------------------- | ------------------------------------------------ |
| **Number**     | _**string**_        | Debit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the debit card |

#### **[UtilityBill](common/model.go#L256) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ID" |
| **Image**         | _***DocumentFile**_ | Scan or photo of the utility bill                  |

#### **[ResidencePermit](common/model.go#L262) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "GB" |
| **IssuedDate**    | _**Time**_          | Issued date                                        |
| **ValidUntil**    | _**Time**_          | Valid until date                                   |
| **Image**         | _***DocumentFile**_ | Scan or photo of the residence permit              |

#### **[Agreement](common/model.go#L270) fields description**

| **Name**  | **Type**            | **Description**                |
| --------- | ------------------- | ------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the agreement |

#### **[EmploymentCertificate](common/model.go#L280) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **IssuedDate** | _**Time**_          | Issued date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the employment certificate |

#### **[Contract](common/model.go#L275) fields description**

| **Name**  | **Type**            | **Description**               |
| --------- | ------------------- | ----------------------------- |
| **Image** | _***DocumentFile**_ | Scan or photo of the contract |

#### **[DocumentPhoto](common/model.go#L296) fields description**

| **Name**  | **Type**            | **Description**                            |
| --------- | ------------------- | ------------------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the photo from a document |

#### **[Selfie](common/model.go#L286) fields description**

| **Name**  | **Type**            | **Description** |
| --------- | ------------------- | --------------- |
| **Image** | _***DocumentFile**_ | Selfie image    |

#### **[Avatar](common/model.go#L291) fields description**

| **Name**  | **Type**            | **Description**          |
| --------- | ------------------- | ------------------------ |
| **Image** | _***DocumentFile**_ | Profile image aka avatar |

#### **[Other](common/model.go#L301) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Document number without whitespaces and dashes          |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ES"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "PA" |
| **IssuedDate**    | _**Time**_          | Issued date                                             |
| **ValidUntil**    | _**Time**_          | Valid until date                                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                           |

#### **[DocumentFile](common/model.go#L151) fields description**

| **Name**        | **Type**     | **Description**                                               |
| --------------- | ------------ | ------------------------------------------------------------- |
| **Filename**    | _**string**_ | File name of the document image, for ex. "passport_front.jpg" |
| **ContentType** | _**string**_ | MIME type of the content, for ex. "image/jpeg"                |
| **Data**        | _**[]byte**_ | Raw content of the document image file                        |

### **KYC response**

The verification response consist of two elements: a result and an error if occurred. The result is of the type [**common.KYCResult**](common/model.go#L164).

#### **[common.KYCResult](common/model.go#L164) fields description**

| **Name**      | **Type**                                                  | **Description**                                                           |
| ------------- | --------------------------------------------------------- | ------------------------------------------------------------------------- |
| **Status**    | [_**KYCStatus**_](#kycstatus-possible-values-description) | Status of the verification                                                |
| **Details**   | _***[KYCDetails](#kycdetails-fields-description)**_       | Details of the verification if provided                                   |
| **ErrorCode** | _**string**_                                              | Error code returned by a KYC provider if the provider support error codes |

#### **[KYCStatus](common/enum.go#L6) possible values description**

| **Value**    | **Description** |
| ------------ | --------------- |
| **Error**    | Verification has failed. Probably, some error has occurred. Returned error value must be non-nil and **`common.KYCResult.ErrorCode`** may contain error code value |
| **Approved** | Successful verification with approved result. The details maybe non-nil and contain additional info about the verification     |
| **Denied**   | Successful verification with rejected result. The details should be non-nil and contain additional info about the verification |
| **Unclear**  | Verification completed with an indefinite result. That might mean that some additional info is required. The details should be non-nil and contain additional info |

#### **[KYCDetails](common/model.go#L158) fields description**

| **Name**     | **Type** | **Description** |
| ------------ | -------- | --------------- |
| **Finality** | [_**KYCFinality**_](#kycfinality-possible-values-description) | Rejection type of the result (if the negative answer is given)           |
| **Reasons**  | _**[]string**_                                                | List of additional response info describing result-related circumstances |

#### **[KYCFinality](common/enum.go#L17) possible values description**

| **Value**    | **Description** |
| ------------ | --------------- |
| **Final**    | Final reject, e.g. when a person is a fraudster, or a client does not want to accept such kind of clients in his/her system |
| **NonFinal** | A reject that can be fixed, e.g. by uploading an image of better quality |
| **Unknown**  | The provider doesn't support **`Finality`** feature |

### **Specific KYC providers**

KYC providers have different configuration options so twas inevitable to implement a specific config for each one of them. But mostly they are identical.

For instructions on integration of a specific KYC provider, please, refer this list:

* [**IDology**](integrations/idology/README.md)
* [**Sum&Substance**](integrations/sumsub/README.md)
* [**Trulioo**](integrations/trulioo/README.md)
* [**Shufti Pro**](integrations/shuftipro/README.md)

### **Required fields**

Each KYC provider has its own subset of minimum required info of the customer. Use this as a reference when integrating with a specific provider what fields of [**common.UserData**](common/model.go#L8) it requires.

> Of course, independently of that the sane minimum of data must always be present.
> Also, the more data you provide to the service the more accurate will be the result.

---

#### **IDology**

[common.UserData](common/model.go#L8) required fields:

| **Name** | **Type** |
| -------- | -------- |
| [FirstName](common/model.go#L10) | _string_ |
| [LastName](common/model.go#L12) | _string_ |
| [CurrentAddress](common/model.go#L26) | [_Address_](common/model.go#L32) |

[common.Address](common/model.go#L32) required fields:

| **Name** | **Type** |
| -------- | -------- |
| [Town](common/model.go#L37) | _string_ |
| [StateProvinceCode](common/model.go#L47) | _string_ |
| [PostCode](common/model.go#L46) | _string_ |

---

#### **Sum&Substance**

According to the [Sum&Substance API Reference](https://developers.sumsub.com/#applicants-api) there are no explicitly required fields from the customer data so providing as much info as possible is the rule.

---

#### **Trulioo**

From the [Trulioo API Reference](https://api.globaldatacompany.com/docs) it is unclear what fields are mandatory so providing as much info as possible is the rule.

---

#### **Shufti Pro**

[common.UserData](common/model.go#L8) required fields:

| **Name** | **Type** |
| -------- | -------- |
| [CountryAlpha2](common/model.go#L22) | _string_ |
| [Documents](common/model.go#L28) | _[][Document](common/model.go#L136)_ |

[common.Document](common/model.go#L136) required fields of [Documents](common/model.go#L28):

| **Name** | **Type** |
| -------- | -------- |
| [Metadata](common/model.go#L138) | [DocumentMetadata](common/model.go#L143) |
| [Front](common/model.go#L139) | [*DocumentFile](common/model.go#L154) |

[common.DocumentMetadata](common/model.go#L143) required fields of [Documents](common/model.go#L28):

| **Name** | **Type** |
| -------- | -------- |
| [Type](common/model.go#L145) | [DocumentType](common/enum.go#L36) |

> **Please, consult [Fields applicable for Shufti Pro](#fields-applicable-for-shufti-pro) for the details about required Documents.**

### **Applicable fields grouped per provider**

[**common.UserData**](common/model.go#L8) provides a wide range of possible data that might require the verification. However, not every KYC provider will surely use all available fields of the model. Therefore, to ease the process of integration for administrators, here you'll find the grouping of applicable fields per provider.

#### **Fields applicable for IDology**

[common.UserData](common/model.go#L8) applicable fields:

| **Name** | **Type** | **Comment** |
| -------- | -------- | ----------- |
| **FirstName** | _string_ | |
| **LastName** | _string_ | |
| **DateOfBirth** | _Time_ | |
| **Email** | _string_ | |
| **Phone** | _string_ | it will be used if non-empty and has length of 10 |
| **MobilePhone** | _string_ | it will be used if has lenght of 10 and the **Phone** field is empty |
| **CurrentAddress** | _Address_ | |
| **SupplementalAddresses** | _[]Address_ | for ex. it might be a shipping address |
| **Documents** | _[]Document_ | `common.IDCard` document type (**SSN**) |

#### **Fields applicable for Sum&Substance**

[common.UserData](common/model.go#L8) applicable fields:

| **Name** | **Type** | **Comment** |
| -------- | -------- | ----------- |
| **FirstName** | _string_ | |
| **LastName** | _string_ | |
| **MiddleName** | _string_ | |
| **LegalName** | _string_ | |
| **Gender** | _Gender_ | |
| **DateOfBirth** | _Time_ | |
| **PlaceOfBirth** | _string_ | |
| **CountryOfBirthAlpha2** | _string_ | |
| **StateOfBirth** | _string_ | |
| **CountryAlpha2** | _string_ | |
| **Nationality** | _string_ | |
| **Phone** | _string_ | |
| **MobilePhone** | _string_ | |
| **CurrentAddress** | _Address_ | |
| **SupplementalAddresses** | _[]Address_ | |
| **Documents** | _[]Document_ | |

#### **Fields applicable for Trulioo**

[common.UserData](common/model.go#L8) applicable fields:

| **Name** | **Type** | **Comment** |
| -------- | -------- | ----------- |
| **FirstName** | _string_ | |
| **PaternalLastName** | _string_ | |
| **LastName** | _string_ | |
| **MiddleName** | _string_ | |
| **LatinISO1Name** | _string_ | |
| **CountryAlpha2** | _string_ | |
| **DateOfBirth** | _Time_ | |
| **Gender** | _Gender_ | |
| **Email** | _string_ | |
| **Phone** | _string_ | |
| **MobilePhone** | _string_ | |
| **CurrentAddress** | _Address_ | |
| **Documents** | _[]Document_ |  |
| **Business** | _Business_ | |

#### **Fields applicable for Shufti Pro**

[common.UserData](common/model.go#L8) applicable fields:

| **Name** | **Type** | **Comment** |
| -------- | -------- | ----------- |
| **FirstName** | _string_ | |
| **LastName** | _string_ | |
| **MiddleName** | _string_ | |
| **DateOfBirth** | _Time_ | |
| **Email** | _string_ | |
| **CountryAlpha2** | _string_ | |
| **CurrentAddress** | _Address_ | |
| **Documents** | _[]Document_ | There are different services which require different documents. For face: **`common.Selfie`**. For documents, anyone of: **`common.Passport`**, **`common.IDCard`**, **`common.Drivers`**, **`common.BankCard`**. For addresses, anyone of: **`common.IDCard`**, **`common.UtilityBill`**. With image data included |

### **The countries supported by KYC providers and the fields variability**

KYC providers may require various set of `common.UserData` fields depending on the customer country. Also, they may service to the limited number of countries.

#### **The countries supported by IDology**

* USA and Canada
* No fields variations found in the docs

#### **The countries supported by Sum&Substance**

* International
* No fields variations found in the docs

#### **The countries supported by Trulioo**

* International
* API provides the group of methods for retrieving the lists of:
  * Consents
  * Supported countries
  * Available fields dynamically based on a country
  * Document Types available for a country
  * Test Entities configured for a country
  * Datasource groups configured for a country

#### **The countries supported by Shufti Pro**

* International (the list of supported country codes is similar to ISO 3166-1 alpha-2)
* No fields variations found in the docs
