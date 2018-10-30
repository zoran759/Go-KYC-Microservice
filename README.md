# KYC Package

## **Table of contents**

* **[REST API](#rest-api)**
* **[Integration interface](#integration-interface)**
* **[KYC request](#kyc-request)**
* **[KYC response](#kyc-response)**
* **[Specific KYC providers](#specific-kyc-providers)**
* **[Applicable fields grouped per provider](#applicable-fields-grouped-per-provider)**
  * **[IDology](#idology)**
  * **[Sum&Substance](#sum&substance)**
  * **[Trulioo](#trulioo)**
  * **[Shufti Pro](#shufti-pro)**
  * **[IdentityMind](#identitymind)**
  * **[Jumio](#jumio)**
* **[The countries supported by KYC providers and the fields variability](#the-countries-supported-by-kyc-providers-and-the-fields-variability)**

## **REST API**

The KYC service provides REST API to interact with other components of the application. The data payload of requests should be JSON encoded. The API responds with JSON payload as well.

| **Route**          |  **Description**                                                    |
| ------------------ | ------------------------------------------------------------------- |
| **/**              | Answers with the welcome message in plain text format               |
| **/ping**          | Answers with the "pong!" response in plain text format              |
| **/CheckCustomer** | The endpoint to send KYC verification requests                      |
| **/CheckStatus**   | The endpoint to send KYC verification current status check requests |

The models for requests and responses are provided.

### **[CheckCustomer request](common/rest.go#L6) fields description**

| **Name**     | **Type**                                       | **Description**                             |
| ------------ | ---------------------------------------------- | ------------------------------------------- |
| **Provider** | _**[KYCProvider](common/enum.go#L36)**_        | The identificator for the KYC provider name |
| **UserData** | _**[UserData](#userdata-fields-description)**_ | A verification data of the customer         |

### **[CheckStatus request](common/rest.go#L12) fields description**

| **Name**        | **Type**                                | **Description**                                                                        |
| --------------- | --------------------------------------- | -------------------------------------------------------------------------------------- |
| **Provider**    | _**[KYCProvider](common/enum.go#L36)**_ | The identificator for the KYC provider name                                            |
| **ReferenceID** | _**string**_                            | The identificator of the verification submission. Its value is specific for a provider |

### **[API response](common/rest.go#L18) fields description**

| **Name**   | **Type**     | **Description**                                                             |
| ---------- | ------------ | --------------------------------------------------------------------------- |
| **Result** | _***[KYCResult](#commonkycresult-fields-description)**_ | A result of the KYC verification |
| **Error**  | _**string**_ | A text of an error message if the error has occured during the verification |

If a **KYC provider** doesn't support the instant result response then check and use the [**KYCResult.StatusPolling**](#statuspolling-fields-description) field for the info required for the KYC verification status check requests.

### **[API Error response](common/rest.go#L24) fields description**

| **Name**  | **Type**     | **Description**    |
| --------- | ------------ | ------------------ |
| **Error** | _**string**_ | A text of an error |

### **HTTP response codes**

| **Code** | **Description**                                                                                                  |
| -------- | ---------------------------------------------------------------------------------------------------------------- |
| **200**  | A request has been successfully processed. The response should be inspected for possible KYC verification errors |
| **400**  | It happens when something wrong with the request. If the request is somehow malformed or missed a required param |
| **404**  | It happens when a KYC provider in the request is unknown for the API                                             |
| **422**  | It happens when a KYC provider doesn't support requested method or it isn't implemented yet                      |
| **500**  | It happens when something goes wrong in the server (serialization errors, KYC config's errors, etc...)           |

## **Integration interface**

All KYC providers implement [**common.CustomerChecker**](common/contract.go#L3) interface for the verification process:

```go
type CustomerChecker interface {
    CheckCustomer(customer *UserData) (KYCResult, error)
}
```

Some KYC providers might require to poll the customer verification status to check if the process is completed. For this purpose the [**common.StatusChecker**](common/contract.go#L10) interface is provided:

```go
type StatusChecker interface {
    CheckStatus(string) (KYCResult, error)
}
```

Those KYC providers that are using some kind of polling mechanism implement this interface.

Providers are configurable by their configs. Configuration options for each provider are described in the respective integration instructions in [Specific KYC providers](#specific-kyc-providers).

The rest required for interaction with KYC providers is in the **`common`** package including request and response structures.

## **KYC request**

For the verification request use a request of the [**common.UserData**](#userdata-fields-description) type.

### **[UserData](common/model.go#L8) fields description**

| **Name**                     | **Type**                           | **Description**                                                       |
| ---------------------------- | ---------------------------------- | --------------------------------------------------------------------- |
| **FirstName**                | _**string**_                       | _Required_. First name of the customer, for ex. "John"                |
| **PaternalLastName**         | _**string**_                       | Paternal last name of the customer                                    |
| **LastName**                 | _**string**_                       | _Required_. Last name of the customer, for ex. "Doe"                  |
| **MiddleName**               | _**string**_                       | Middle name of the customer, for ex. "Benedikt"                       |
| **LegalName**                | _**string**_                       | Legal name of the customer, for ex. "Astrid Lindgren"                 |
| **LatinISO1Name**            | _**string**_                       | Latin ISO1 name of the customer, for ex. russian "Иван" be IVAN       |
| **AccountName**              | _**string**_                       | Account name for the customer, for ex. "john_doe"                     |
| **Email**                    | _**string**_                       | Email of the customer                                                 |
| **IPaddress**                | _**string**_                       | Customer’s IP address                                                 |
| **Gender**                   | [_**Gender**_](common/enum.go#L27) | Gender of the customer                                                |
| **DateOfBirth**              | [_**Time**_](common/time.go#L5)    | Date of birth of the customer                                         |
| **PlaceOfBirth**             | _**string**_                       | Place of birth of the customer                                        |
| **CountryOfBirthAlpha2**     | _**string**_                       | Country of birth of the customer in ISO 3166-1 alpha-2 format, for ex. "US" |
| **StateOfBirth**             | _**string**_                       | State of birth of the customer, for ex. "GA"                          |
| **CountryAlpha2**            | _**string**_                       | Country of the customer in ISO 3166-1 alpha-2 format, for ex. "DE"    |
| **Nationality**              | _**string**_                       | Citizenship of the customer. Perhaps, it should be country's name, for ex. "Italy" |
| **Phone**                    | _**string**_                       | Primary phone of the customer. It isn't the mobile phone!             |
| **MobilePhone**              | _**string**_                       | Mobile phone of the customer                                          |
| **CurrentAddress**           | [_**Address**_](#address-fields-description) | Current address of the customer                             |
| **SupplementalAddresses**    | _**[]Address**_                    | List of supplemental addresses of the customer                        |
| **Location**                 | _***[Location](#location-fields-description)**_ | Geopositional data of the customer                       |
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
| **DocumentPhoto**            | _***[DocumentPhoto](#documentphoto-fields-description)**_     | Document photo of the customer (like a photo from the passport) |
| **Selfie**                   | _***[Selfie](#selfie-fields-description)**_                   | Selfie image of the customer               |
| **Avatar**                   | _***[Avatar](#avatar-fields-description)**_                   | A profile image aka avatar of the customer |
| **Other**                    | _***[Other](#other-fields-description)**_                     | Other document (should be used only when nothing else applies) |

### **[Address](common/model.go#L50) fields description**

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

### **[Business](common/model.go#L164) fields description**

| **Name**                      | **Type**     | **Description**                                |
| ----------------------------- | ------------ | ---------------------------------------------- |
| **Name**                      | _**string**_ | Name of the Enterprise the customer relates to |
| **RegistrationNumber**        | _**string**_ | Registration number of the Enterprise          |
| **IncorporationDate**         | _**Time**_   | Incorporation date of the Enterprise           |
| **IncorporationJurisdiction** | _**string**_ | Incorporation jurisdiction of the Enterprise   |

### **[Passport](common/model.go#L205) fields description**

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

### **[IDCard](common/model.go#L217) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **Number**        | _**string**_        | Id card number without whitespaces and dashes      |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "CN" |
| **IssuedDate**    | _**Time**_          | Issued date                                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the card                          |

### **[SNILS](common/model.go#L225) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **Number**     | _**string**_        | SNILS number without whitespaces and dashes |
| **IssuedDate** | _**Time**_          | Issued date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the SNILS                  |

### **[DriverLicense](common/model.go#L232) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license number                                   |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "DE"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "KY" |
| **IssuedDate**    | _**Time**_          | Issued date                                             |
| **ValidUntil**    | _**Time**_          | Valid until date                                        |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license   |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license    |

### **[DriverLicenseTranslation](common/model.go#L243) fields description**

| **Name**          | **Type**            | **Description**                                                   |
| ----------------- | ------------------- | ----------------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license translation number                                 |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "LV"                |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "MT"           |
| **IssuedDate**    | _**Time**_          | Issued date                                                       |
| **ValidUntil**    | _**Time**_          | Valid until date                                                  |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license translation |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license translation  |

### **[CreditCard](common/model.go#L254) fields description**

| **Name**       | **Type**            | **Description**                                   |
| -------------- | ------------------- | ------------------------------------------------- |
| **Number**     | _**string**_        | Credit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date                                  |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the credit card |

### **[DebitCard](common/model.go#L261) fields description**

| **Name**       | **Type**            | **Description**                                  |
| -------------- | ------------------- | ------------------------------------------------ |
| **Number**     | _**string**_        | Debit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the debit card |

### **[UtilityBill](common/model.go#L268) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ID" |
| **Image**         | _***DocumentFile**_ | Scan or photo of the utility bill                  |

### **[ResidencePermit](common/model.go#L274) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "GB" |
| **IssuedDate**    | _**Time**_          | Issued date                                        |
| **ValidUntil**    | _**Time**_          | Valid until date                                   |
| **Image**         | _***DocumentFile**_ | Scan or photo of the residence permit              |

### **[Agreement](common/model.go#L282) fields description**

| **Name**  | **Type**            | **Description**                |
| --------- | ------------------- | ------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the agreement |

### **[EmploymentCertificate](common/model.go#L292) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **IssuedDate** | _**Time**_          | Issued date                                 |
| **Image**      | _***DocumentFile**_ | Scan or photo of the employment certificate |

### **[Contract](common/model.go#L287) fields description**

| **Name**  | **Type**            | **Description**               |
| --------- | ------------------- | ----------------------------- |
| **Image** | _***DocumentFile**_ | Scan or photo of the contract |

### **[DocumentPhoto](common/model.go#L308) fields description**

| **Name**  | **Type**            | **Description**                            |
| --------- | ------------------- | ------------------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the photo from a document |

### **[Selfie](common/model.go#L298) fields description**

| **Name**  | **Type**            | **Description** |
| --------- | ------------------- | --------------- |
| **Image** | _***DocumentFile**_ | Selfie image    |

### **[Avatar](common/model.go#L303) fields description**

| **Name**  | **Type**            | **Description**          |
| --------- | ------------------- | ------------------------ |
| **Image** | _***DocumentFile**_ | Profile image aka avatar |

### **[Other](common/model.go#L313) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Document number without whitespaces and dashes          |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ES"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "PA" |
| **IssuedDate**    | _**Time**_          | Issued date                                             |
| **ValidUntil**    | _**Time**_          | Valid until date                                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                           |

### **[DocumentFile](common/model.go#L172) fields description**

| **Name**        | **Type**     | **Description**                                               |
| --------------- | ------------ | ------------------------------------------------------------- |
| **Filename**    | _**string**_ | File name of the document image, for ex. "passport_front.jpg" |
| **ContentType** | _**string**_ | MIME type of the content, for ex. "image/jpeg"                |
| **Data**        | _**[]byte**_ | Raw content of the document image file                        |

### **[Location](common/model.go#L158) fields description**

| **Name**      | **Type**     | **Description**                             |
| ------------- | ------------ | ------------------------------------------- |
| **Latitude**  | _**string**_ | The location latitude, for ex. "55.678849"  |
| **Longitude** | _**string**_ | The location longitude, for ex. "52.327662" |

## **KYC response**

The verification response consist of two elements: a result and an error if occurred. The result is of the type [**common.KYCResult**](#commonkycresult-fields-description).

> Some KYC providers might require to poll the customer verification status to check if the process is completed. For this purpose the __*KYCStatusCheck__ field is provided. If a polling is required and no error has occured then this field will be non-nil.

### **[common.KYCResult](common/model.go#L185) fields description**

| **Name**        | **Type**                                                    | **Description**                                                           |
| --------------- | ----------------------------------------------------------- | ------------------------------------------------------------------------- |
| **Status**      | _**[KYCStatus](#kycstatus-possible-values-description)**_   | Status of the verification                                                |
| **Details**     | _***[KYCDetails](#kycdetails-fields-description)**_         | Details of the verification if provided                                   |
| **ErrorCode**   | _**string**_                                                | Error code returned by a KYC provider if the provider support error codes |
| **StatusCheck** | _***[KYCStatusCheck](#kycstatuscheck-fields-description)**_ | Data required to do the KYC verification status check requests if needed  |

### **[KYCStatus](common/enum.go#L6) possible values description**

| **Value**    | **Description**                                                                                                                |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| **Error**    | Verification has failed. Probably, some error has occurred. Returned error value must be non-nil and **`common.KYCResult.ErrorCode`** may contain error code value |
| **Approved** | Successful verification with approved result. The details maybe non-nil and contain additional info about the verification     |
| **Denied**   | Successful verification with rejected result. The details should be non-nil and contain additional info about the verification |
| **Unclear**  | Needs subsequent status polling or the verification completed with an indefinite result. That might mean that some additional info is required. The details should be non-nil and contain additional info. If status polling is required then **`common.KYCResult.StatusCheck`** must be non-nil |

### **[KYCDetails](common/model.go#L179) fields description**

| **Name**     | **Type**                                                      | **Description**                                                          |
| ------------ | ------------------------------------------------------------- | ------------------------------------------------------------------------ |
| **Finality** | [_**KYCFinality**_](#kycfinality-possible-values-description) | Rejection type of the result (if the negative answer is given)           |
| **Reasons**  | _**[]string**_                                                | List of additional response info describing result-related circumstances |

### **[KYCFinality](common/enum.go#L17) possible values description**

| **Value**    | **Description**                                                                                                             |
| ------------ | --------------------------------------------------------------------------------------------------------------------------- |
| **Final**    | Final reject, e.g. when a person is a fraudster, or a client does not want to accept such kind of clients in his/her system |
| **NonFinal** | A reject that can be fixed, e.g. by uploading an image of better quality                                                    |
| **Unknown**  | The provider doesn't support **`Finality`** feature                                                                         |

### **[KYCStatusCheck](common/model.go#L193) fields description**

| **Name**        | **Type**                                | **Description**                             |
| --------------- | --------------------------------------- | ------------------------------------------- |
| **Provider**    | _**[KYCProvider](common/enum.go#L36)**_ | An identificator for the KYC provider name  |
| **ReferenceID** | _**string**_                            | An identificator that references to this verification submission. It mention in docs as applicantId/mtid/jumioIdScanReference/etc. Its value is specific for a provider |
| **LastCheck**   | _**time.Time**_                         | Last time a verification status was checked |

## **Specific KYC providers**

KYC providers have different configuration options and that was implemented as a specific config for each one of them. But mostly they are identical.

For instructions on integration of a specific KYC provider, please, refer this list:

* [**IDology**](integrations/idology/README.md)
* [**Sum&Substance**](integrations/sumsub/README.md)
* [**Trulioo**](integrations/trulioo/README.md)
* [**Shufti Pro**](integrations/shuftipro/README.md)
* [**IdentityMind**](integrations/identitymind/README.md)
* [**Jumio**](integrations/jumio/README.md)

## **Applicable fields grouped per provider**

[**common.UserData**](#userdata-fields-description) provides a wide range of possible data that might require the verification. However, not every KYC provider will surely use all available fields of the model. Therefore, to ease the process of integration for administrators, here you'll find the grouping of applicable fields per provider.

### **IDology**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**                  | **Type**    | **Required** | **Comment**                                                          |
| ------------------------- | ----------- | :----------: | -------------------------------------------------------------------- |
| **FirstName**             | _string_    | **Yes**      |                                                                      |
| **LastName**              | _string_    | **Yes**      |                                                                      |
| Email                     | _string_    |              |                                                                      |
| DateOfBirth               | _Time_      |              |                                                                      |
| Phone                     | _string_    |              | It will be used if non-empty and has length of 10                    |
| MobilePhone               | _string_    |              | It will be used if has lenght of 10 and the **Phone** field is empty |
| **CurrentAddress**        | _Address_   | **Yes**      |                                                                      |
| SupplementalAddresses     | _[]Address_ |              | It might be a shipping address                                       |
| **IDCard**                | _*IDCard_   | **Yes**      |                                                                      |

[**Address**](#address-fields-description) mandatory fields:

| **Name**              | **Type** |
| --------------------- | -------- |
| **Town**              | _string_ |
| **StateProvinceCode** | _string_ |
| **PostCode**          | _string_ |

### **Sum&Substance**

According to the [API Reference](https://developers.sumsub.com) all fields of [**UserData**](#userdata-fields-description) are applicable except the following:

* PaternalLastName
* LatinISO1Name
* AccountName
* IPaddress
* Location
* Business

Sum&Substance requires at least one document to be present to start the verification process, so anyone of the following fields should present:

| **Name**                     | **Type**                    |
| ---------------------------- | --------------------------- |
| **Passport**                 | _*Passport_                 |
| **IDCard**                   | _*IDCard_                   |
| **SNILS**                    | _*SNILS_                    |
| **DriverLicense**            | _*DriverLicense_            |
| **DriverLicenseTranslation** | _*DriverLicenseTranslation_ |
| **CreditCard**               | _*CreditCard_               |
| **DebitCard**                | _*DebitCard_                |
| **UtilityBill**              | _*UtilityBill_              |
| **ResidencePermit**          | _*ResidencePermit_          |
| **Agreement**                | _*Agreement_                |
| **EmploymentCertificate**    | _*EmploymentCertificate_    |
| **Contract**                 | _*Contract_                 |
| **DocumentPhoto**            | _*DocumentPhoto_            |
| **Selfie**                   | _*Selfie_                   |
| **Avatar**                   | _*Avatar_                   |
| **Other**                    | _*Other_                    |

All fields in the Reference are marked as optional but at least first name and last name should be provided in addition to a doc.

### **Trulioo**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**          | **Type**           | **Required** | **Comment** |
| ----------------- | ------------------ | :----------: | ----------- |
| FirstName         | _string_           |              |             |
| PaternalLastName  | _string_           |              |             |
| LastName          | _string_           |              |             |
| MiddleName        | _string_           |              |             |
| LatinISO1Name     | _string_           |              |             |
| Email             | _string_           |              |             |
| Gender            | _Gender_           |              |             |
| DateOfBirth       | _Time_             |              |             |
| **CountryAlpha2** | _string_           | **Yes**      |             |
| Phone             | _string_           |              |             |
| MobilePhone       | _string_           |              |             |
| CurrentAddress    | _Address_          |              |             |
| Business          | _*Business_        |              |             |
| Passport          | _*Passport_        |              |             |
| IDCard            | _*IDCard_          |              |             |
| DriverLicense     | _*DriverLicense_   |              |             |
| ResidencePermit   | _*ResidencePermit_ |              |             |
| Selfie            | _*Selfie_          |              |             |

> **DOCUMENTS NOTE:** Include image file(s) for a document used in the verification.

It's unclear from the [API Reference](https://developer.trulioo.com/v1.0/reference) what fields are mandatory so, it's better to provide as much info as possible.

### **Shufti Pro**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**             | **Type**         | **Required**       | **Comment**                                                       |
| -------------------- | ---------------- | :----------------: | ----------------------------------------------------------------- |
| **FirstName**        | _string_         | **Yes**            |                                                                   |
| **LastName**         | _string_         | **Yes**            |                                                                   |
| MiddleName           | _string_         |                    |                                                                   |
| Email                | _string_         |                    |                                                                   |
| DateOfBirth          | _Time_           |                    |                                                                   |
| **CountryAlpha2**    | _string_         | **Yes**            |                                                                   |
| **Phone**            | _string_         | **Yes**            | Customer’s phone number with country code. Example: +440000000000 |
| CurrentAddress       | _Address_        |                    |                                                                   |
| **Passport**(*)      | _*Passport_      | **See comment(*)** | (*)Anyone of documents marked with asterisk                       |
| **IDCard**(*)        | _*IDCard_        | **(*)**            |                                                                   |
| SNILS                | _*SNILS_         |                    |                                                                   |
| **DriverLicense**(*) | _*DriverLicense_ | **(*)**            |                                                                   |
| **CreditCard**(*)    | _*CreditCard_    | **(*)**            |                                                                   |
| UtilityBill          | _*UtilityBill_   |                    |                                                                   |
| **Selfie**           | _*Selfie_        | **Yes**            |                                                                   |

> **DOCUMENTS NOTE:** Include image file(s) for the document used for the verification.

### **IdentityMind**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**              | **Type**           | **Required**       | **Comment**                                  |
| --------------------- | ------------------ | :----------------: | -------------------------------------------- |
| FirstName             | _string_           |                    |                                              |
| LastName              | _string_           |                    |                                              |
| MiddleName            | _string_           |                    |                                              |
| **AccountName**       | _string_           | **Yes**            |                                              |
| Email                 | _string_           |                    |                                              |
| IPaddress             | _string_           |                    |                                              |
| Gender                | _Gender_           |                    |                                              |
| DateOfBirth           | _Time_             |                    |                                              |
| CountryAlpha2         | _string_           |                    |                                              |
| Phone                 | _string_           |                    |                                              |
| MobilePhone           | _string_           |                    |                                              |
| CurrentAddress        | _Address_          |                    |                                              |
| Location              | _*Location_        |                    |                                              |
| Passport              | _*Passport_        |                    |                                              |
| IDCard                | _*IDCard_          |                    |                                              |
| SNILS                 | _*SNILS_           |                    |                                              |
| DriverLicense         | _*DriverLicense_   |                    |                                              |
| UtilityBill           | _*UtilityBill_     |                    |                                              |
| ResidencePermit       | _*ResidencePermit_ |                    |                                              |
| Selfie(*)             | _*Selfie_          | **See comment(*)** | (*)Provide it if using Document Verification |

> **DOCUMENTS NOTE:** Include image file(s) for the document used for the verification.

### **Jumio**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**             | **Type**         | **Required**       | **Comment**                                 |
| -------------------- | ---------------- | :----------------: | ------------------------------------------- |
| FirstName            | _string_         |                    |                                             |
| LastName             | _string_         |                    |                                             |
| DateOfBirth          | _Time_           |                    |                                             |
| **Passport**(*)      | _*Passport_      | **See comment(*)** | (*)Anyone of documents marked with asterisk |
| **IDCard**(*)        | _*IDCard_        | **(*)**            |                                             |
| **SNILS**(*)         | _*SNILS_         | **(*)**            |                                             |
| **DriverLicense**(*) | _*DriverLicense_ | **(*)**            |                                             |
| **Selfie**           | _*Selfie_        | **See comment(*)** | (*)Mandatory if Face match enabled          |

> **DOCUMENTS NOTE:** Include image file(s) for the document used for the verification.

## **The countries supported by KYC providers and the fields variability**

KYC providers may require various set of `common.UserData` fields depending on the customer country. Also, they may service to the limited number of countries and this number of countries might configurable in a web-interface of the provider.

### **IDology covered countries**

* USA and Canada
* No fields variations found in the docs

### **Sum&Substance covered countries**

* International
* No fields variations found in the docs

### **Trulioo covered countries**

* International
* API provides the group of methods for retrieving the lists of:
  * Consents
  * Supported countries
  * Available fields dynamically based on a country
  * Document Types available for a country
  * Test Entities configured for a country
  * Datasource groups configured for a country

### **Shufti Pro covered countries**

* International (the list of supported country codes is similar to ISO 3166-1 alpha-2)
* No fields variations found in the docs

### **IdentityMind covered countries**

* International
* No fields variations found in the docs

### **Jumio covered countries**

* International
* No fields variations found in the docs
