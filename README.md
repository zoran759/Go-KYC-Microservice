# KYC Package

* **[KYC service configuration](#kyc-service-configuration)**
  * **[Command line options](#command-line-options)**
  * **[Configuration file options](#configuration-file-options)**
* **[KYC providers configuration options](#kyc-providers-configuration-options)**
  * **[Coinfirm](#coinfirm-configuration-options)**
  * **[ComplyAdvantage](#complyadvantage-configuration-options)**
  * **[IdentityMind](#identitymind-configuration-options)**
  * **[IDology](#idology-configuration-options)**
  * **[Jumio](#jumio-configuration-options)**
  * **[Shufti Pro](#shufti-pro-configuration-options)**
  * **[Sum&Substance](#sum&substance-configuration-options)**
  * **[SynapseFI](#synapsefi-configuration-options)**
  * **[Thomson Reuters](#thomson-reuters-configuration-options)**
  * **[Trulioo](#trulioo-configuration-options)**
* **[REST API](#rest-api)**
  * **[Endpoints](#endpoints)**
  * **[Customer verification request](#checkcustomer-request-fields-description)**
  * **[Customer verification status request](#checkstatus-request-fields-description)**
  * **[API response](#api-response-fields-description)**
  * **[API Error response](#api-error-response-fields-description)**
  * **[HTTP response codes](#http-response-codes)**
  * **[Checking if a KYC provider is implemented in the API](#checking-if-a-kyc-provider-is-implemented-in-the-api)**
* **[For developers](#for-developers)**

## **KYC service configuration**

### **Command line options**

The service supports the following command line options:

| **Name** | **Description**                                                                              |
| -------- | -------------------------------------------------------------------------------------------- |
| `help`   | Prints info about supported command-line options and exits.                                  |
| `config` | Specifies the file to use for configuration.                                                 |
| `port`   | Specifies the port for the service to listen for incoming requests. The default port is 8080 |

### **Configuration file options**

All options must be placed under the **`Config`** section of the configuration file. The service supports the following options in the configuration file:

| **Name** | **Description**                                            |
| -------- | ---------------------------------------------------------- |
| `Port`   | Has the same meaning as the command-line **`port`** option |

> **WARNING!** If a command line option is specified its value overrides the configuration file value for that option.

## **KYC providers configuration options**

Current implementation of the service allows configuration for the supported KYC providers through the configuration file **[kyc.cfg](main/kyc.cfg)**. We provide the sample file without credentials. It includes all providers supported by the service. Feel free to modify it to suit your needs.

Below are the options required for each provider.

### **Coinfirm configuration options**

| **Name** | **Description**                            |
| -------- | ------------------------------------------ |
| Host     | Coinfirm API url without trailing slash    |
| Email    | Email of the Coinfirm user                 |
| Password | Password of the Coinfirm user              |
| Company  | Name of the project registered in Coinfirm |

### **ComplyAdvantage configuration options**

| **Name**  | **Description**                                                                                              |
| --------- | ------------------------------------------------------------------------------------------------------------ |
| Host      | ComplyAdvantage API url without trailing slash                                                               |
| APIkey    | ComplyAdvantage API key which can be generated within ComplyAdvantage web platform                           |
| Fuzziness | Determines how closely the returned results must match the supplied name. Float value in range [0.0 ... 1.0] |

### **IdentityMind configuration options**

| **Name** | **Description**                                                               |
| -------- | ----------------------------------------------------------------------------- |
| Host     | IdentityMind API url without trailing slash                                   |
| Username | IdentityMind API username that is supplied during registration in the service |
| Password | IdentityMind API key that is supplied during registration in the service      |

### **IDology configuration options**

| **Name**         | **Description**                                                                                                           |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Host             | IDology API url without trailing slash                                                                                    |
| Username         | IDology API username supplied by the service                                                                              |
| Password         | IDology API password supplied by the service                                                                              |
| UseSummaryResult | If Summary Results are enabled in the Enterprise Configuration, this must be set to true. Boolean value `true` or `false` |

### **Jumio configuration options**

| **Name** | **Description**                         |
| -------- | --------------------------------------- |
| BaseURL  | Jumio API url without trailing slash    |
| Token    | Jumio API token supplied by the service |
| Secret   | Jumio API secret supplied by the service. You can view and manage your API token and secret in the Customer Portal under Settings > API credentials |

### **Shufti Pro configuration options**

| **Name**    | **Description**                                                                                        |
| ----------- | ------------------------------------------------------------------------------------------------------ |
| Host        | Shufti Pro API url without trailing slash                                                              |
| ClientID    | Shufti Pro Client ID supplied by the service                                                           |
| SecretKey   | Shufti Pro Secret Key supplied by the service                                                          |
| RedirectURL | The url to redirect your customer after the verification process completes. Currently, leave it intact |

### **Sum&Substance configuration options**

| **Name** | **Description**                               |
| -------- | --------------------------------------------- |
| Host     | Sum&Substance API url without trailing slash  |
| APIKey   | Sum&Substance API key supplied by the service |

### **SynapseFI configuration options**

| **Name**     | **Description**                                  |
| ------------ | ------------------------------------------------ |
| Host         | Synapse FI API url **with** trailing slash       |
| ClientID     | Synapse FI client id supplied by the service     |
| ClientSecret | Synapse FI client secret supplied by the service |

### **Thomson Reuters configuration options**

| **Name**  | **Description**                                                                                                |
| --------- | -------------------------------------------------------------------------------------------------------------- |
| Host      | Thomson Reuters World-Check One API (WC1 API) url **with** trailing slash                                      |
| APIkey    | WC1 API key generated by WC1 and made available to WC1 administrators via the user administration interface    |
| APIsecret | WC1 API secret generated by WC1 and made available to WC1 administrators via the user administration interface |

### **Trulioo configuration options**

| **Name**     | **Description**                                                        |
| ------------ | ---------------------------------------------------------------------- |
| Host         | Trulioo GlobalGateway Normalized API (NAPI) url without trailing slash |
| NAPILogin    | The NAPI username supplied by the service                              |
| NAPIPassword | The NAPI password supplied by the service                              |

## **REST API**

The KYC service exposes REST API for interaction. The data payload of requests should be JSON encoded. The API responds with JSON-encoded payload as well.

### **Endpoints**

Our API makes available the following Endpoints:

| **Method** | **Route**        |  **Description**                                       |
| ---------- | ---------------- | ------------------------------------------------------ |
| GET        | `/`              | Answers with the welcome message in plain text format  |
| GET        | `/Ping`          | Answers with the "Pong!" response in plain text format |
| GET        | `/Provider`      | Check whether a specified provider is implemented      |
| POST       | `/CheckCustomer` | Send KYC verification requests                         |
| POST       | `/CheckStatus`   | Send KYC verification current status check requests    |

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

### **[API response](common/rest.go#L23) fields description**

| **Name**   | **Type**     | **Description**                                                             |
| ---------- | ------------ | --------------------------------------------------------------------------- |
| **Result** | _***[Result](#commonresult-fields-description)**_ | A result of the KYC verification       |
| **Error**  | _**string**_ | A text of an error message if the error has occured during the verification |

If a **KYC provider** doesn't support the instant result response then check and use the [**Result.StatusCheck**](#kycstatuscheck-fields-description) field for the info required for the KYC verification status check requests.

### **[API Error response](common/rest.go#L18) fields description**

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

### **Checking if a KYC provider is implemented in the API**

The service accepts only one param for this kind of check - **`name`**. Other params are ignored. If the request is valid then the JSON response will be returned showing whether the specified provider is implemented:

```json
{
    "Implemented": true
}
```

This example represents the positive response.

If the request performed without params then the sorted list of implemented KYC providers will be returned in the response:

```json
[
    "ComplyAdvantage",
    "IDology",
    ...
    "ThomsonReuters",
    "Trulioo"
]
```

## **FOR DEVELOPERS**

> **This part may be of interest mainly to developers.**

* **[Integration interface](#integration-interface)**
* **[KYC request](#kyc-request)**
* **[KYC response](#kyc-response)**
* **[Applicable fields grouped per provider](#applicable-fields-grouped-per-provider)**
  * **[Coinfirm](#coinfirm)**
  * **[ComplyAdvantage](#complyadvantage)**
  * **[IdentityMind](#identitymind)**
  * **[IDology](#idology)**
  * **[Jumio](#jumio)**
  * **[Shufti Pro](#shufti-pro)**
  * **[Sum&Substance](#sum&substance)**
  * **[SynapseFI](#synapsefi)**
  * **[Thomson Reuters](#thomson-reuters)**
  * **[Trulioo](#trulioo)**
* **[The countries supported by KYC providers and the fields variability](#the-countries-supported-by-kyc-providers-and-the-fields-variability)**
  * **[Coinfirm](#coinfirm-covered-countries)**
  * **[ComplyAdvantage](#complyadvantage-covered-countries)**
  * **[IdentityMind](#identitymind-covered-countries)**
  * **[IDology](#idology-covered-countries)**
  * **[Jumio](#jumio-covered-countries)**
  * **[Shufti Pro](#shufti-pro-covered-countries)**
  * **[Sum&Substance](#sum&substance-covered-countries)**
  * **[SynapseFI](#synapsefi-covered-countries)**
  * **[Thomson Reuters](#thomson-reuters-covered-countries)**
  * **[Trulioo](#trulioo-covered-countries)**

## **Integration interface**

All KYC providers implement [**common.KYCPlatform**](common/contract.go#L3) interface for the verification process:

```go
type KYCPlatform interface {
    CheckCustomer(customer *UserData) (KYCResult, error)
    CheckStatus(referenceID string) (KYCResult, error)
}
```

KYC providers handle KYC process differently. Some return KYC result instantly in the response. Some require to poll the customer verification status to check if the process is completed. For this purpose the __*common.KYCResponse.Result.StatusCheck__ field is provided. If a polling is required and no error has occured then this field will be non-nil.

The rest required for interaction with KYC providers is in the **`common`** package including request and response structures.

## **KYC request**

For the verification request use a request of the [**common.UserData**](#userdata-fields-description) type.

### **[UserData](common/model.go#L8) fields description**

| **Name**                     | **Type**                           | **Description**                                                       |
| ---------------------------- | ---------------------------------- | --------------------------------------------------------------------- |
| **FirstName**                | _**string**_                       | _Required_. First name (FN) of the customer, for ex. "John"           |
| **LastName**                 | _**string**_                       | _Required_. Last name (LN) of the customer, for ex. "Doe"             |
| **MaternalLastName**         | _**string**_                       | Maternal (second) last name (MatN) of the customer                    |
| **MiddleName**               | _**string**_                       | Middle name (MidN) of the customer, for ex. "Benedikt"                |
| **FullName**                 | _**string**_                       | Full name as found on identification documents (if it isn't composed as FN+MidN+LN+MatN) |
| **LegalName**                | _**string**_                       | Legal name of the customer, for ex. "Astrid Lindgren"                 |
| **LatinISO1Name**            | _**string**_                       | Latin ISO1 name of the customer, for ex. russian "Иван Сидоров" be "IVAN SIDOROV" |
| **AccountName**              | _**string**_                       | Account name for the customer, for ex. "john_doe"                     |
| **Email**                    | _**string**_                       | Email of the customer                                                 |
| **IPaddress**                | _**string**_                       | Customer’s IP address                                                 |
| **Gender**                   | [_**Gender**_](common/enum.go#L27) | Gender of the customer                                                |
| **DateOfBirth**              | [_**Time**_](common/time.go#L5)    | Date of birth of the customer in RFC3339 format, for ex. "2006-01-02T15:04:05Z07:00" |
| **PlaceOfBirth**             | _**string**_                       | Place of birth of the customer                                        |
| **CountryOfBirthAlpha2**     | _**string**_                       | Country of birth of the customer in ISO 3166-1 alpha-2 format, for ex. "US" |
| **StateOfBirth**             | _**string**_                       | State of birth of the customer, for ex. "GA"                          |
| **CountryAlpha2**            | _**string**_                       | Country of the customer in ISO 3166-1 alpha-2 format, for ex. "DE"    |
| **Nationality**              | _**string**_                       | Citizenship of the customer. ISO 3166-1 alpha-2 format, for ex. "TH"  |
| **Phone**                    | _**string**_                       | Primary phone of the customer. It isn't the mobile phone!             |
| **MobilePhone**              | _**string**_                       | Mobile phone of the customer                                          |
| **BankAccountNumber**        | _**string**_                       | Chinese bank account number                                           |
| **VehicleRegistrationPlate** | _**string**_                       | New Zealand vehicle registration plate                                |
| **CurrentAddress**           | [_**Address**_](#address-fields-description) | Current address of the customer                             |
| **SupplementalAddresses**    | _**[]Address**_                    | List of supplemental addresses of the customer                        |
| **Location**                 | _***[Location](#location-fields-description)**_ | Geopositional data of the customer                       |
| **Business**                 | _***[Business](#business-fields-description)**_ | The business which the customer is linked to or is one of the owners |
| **Passport**                 | _***[Passport](#passport-fields-description)**_ | Passport of the customer                                 |
| **IDCard**                   | _***[IDCard](#idcard-fields-description)**_     | National ID Number (Malaysia NRIC Number, Swedish PIN, etc...)       |
| **SNILS**                    | _***[SNILS](#snils-fields-description)**_       | SNILS (Russian insurance number of individual ledger account) of the customer |
| **HealthID**                 | _***[HealthID](#healthid-fields-description)**_               | National Health Service Identification     |
| **SocialServiceID**          | _***[SocialServiceID](#socialserviceid-fields-description)**_ | National Social Service Identification (Social Security Number, Social Insurance Number, National Insurance Number) |
| **TaxID**                    | _***[TaxID](#taxid-fields-description)**_                     | National Taxpayer Personal Identification Number       |
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
| **VideoAuth**                | _***[VideoAuth](#videoauth-fields-description)**_             | Short authorization video of the customer (up to 5 seconds)    |
| **IsCompany**                | _**bool**_                         | Indicates when a Company data is provided for KYC onboarding          |
| **CompanyName**              | _**string**_                       | Company full name                                                     |
| **Website**                  | _**string**_                       | Company's website URL                                                 |
| **CompanyBoard**             | _***CompanyBoard**_                | A certified document containing a list of members of company's board of directors (e.g. an extract from company register or an officially certified document) |
| **CompanyRegistration**      | _***CompanyRegistration**_         | A certificate of company registration                                 |

### **[Address](common/address.go#L5) fields description**

| **Name**              | **Type**     | **Description**                                                  |
| --------------------- | ------------ | ---------------------------------------------------------------- |
| **CountryAlpha2**     | _**string**_ | Country in ISO 3166-1 alpha-2 format, for ex. "FR"               |
| **County**            | _**string**_ | County if applicable                                             |
| **State**             | _**string**_ | Name of the state, for ex. "Alabama"                             |
| **Town**              | _**string**_ | City or town name                                                |
| **Suburb**            | _**string**_ | Suburb if applicable                                             |
| **Street**            | _**string**_ | Street name, for ex. "PeachTree Place", "7th street"             |
| **StreetType**        | _**string**_ | Street type, for ex. "Avenue"                                    |
| **SubStreet**         | _**string**_ | Substreet if applicable                                          |
| **BuildingName**      | _**string**_ | Building or house name                                           |
| **BuildingNumber**    | _**string**_ | Building or house number                                         |
| **FlatNumber**        | _**string**_ | Apartment number                                                 |
| **PostOfficeBox**     | _**string**_ | Post office box                                                  |
| **PostCode**          | _**string**_ | Zip or postal code                                               |
| **StateProvinceCode** | _**string**_ | Abbreviated name of the state or province, for ex. "CA"          |
| **StartDate**         | _**Time**_   | When the customer settled at this address, in RFC3339 format     |
| **EndDate**           | _**Time**_   | When the customer moved out from this address, in RFC3339 format |

### **[Business](common/model.go#L63) fields description**

| **Name**                      | **Type**     | **Description**                                         |
| ----------------------------- | ------------ | ------------------------------------------------------- |
| **Name**                      | _**string**_ | Name of the Enterprise the customer relates to          |
| **RegistrationNumber**        | _**string**_ | Registration number of the Enterprise                   |
| **IncorporationDate**         | _**Time**_   | Incorporation date of the Enterprise, in RFC3339 format |
| **IncorporationJurisdiction** | _**string**_ | Incorporation jurisdiction of the Enterprise            |

### **[Passport](common/model.go#L104) fields description**

| **Name**          | **Type**                                                | **Description**                                         |
| ----------------- | ------------------------------------------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_                                            | Passport number without whitespaces and dashes          |
| **Mrz1**          | _**string**_                                            | First line of the Machine Readable Zone (MRZ) of passport, 44 letters and digits, i.e. "P<CZESPECIMEN<<VZOR<<<<<<<<<<<<<<<<<<<<<<<<<" |
| **Mrz2**          | _**string**_                                            | Second line of the Machine Readable Zone (MRZ) of passport, 44 letters and digits, i.e. "99003853<1CZE1101018M1207046110101111<<<<<94" |
| **CountryAlpha2** | _**string**_                                            | Country in ISO 3166-1 alpha-2 format, for ex. "SG"      |
| **State**         | _**string**_                                            | Abbreviated name of the state or province, for ex. "TX" |
| **IssuedDate**    | _**Time**_                                              | Issued date, in RFC3339 format                          |
| **ValidUntil**    | _**Time**_                                              | Valid until date, in RFC3339 format                     |
| **Image**         | _***[DocumentFile](#documentfile-fields-description)**_ | Scan or photo of the passport                           |

### **[IDCard](common/model.go#L116) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **Number**        | _**string**_        | Id card number without whitespaces and dashes      |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "CN" |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                     |
| **Image**         | _***DocumentFile**_ | Scan or photo of the card                          |

### **[SNILS](common/model.go#L124) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **Number**     | _**string**_        | SNILS number without whitespaces and dashes |
| **IssuedDate** | _**Time**_          | Issued date, in RFC3339 format              |
| **Image**      | _***DocumentFile**_ | Scan or photo of the SNILS                  |

### **[HealthID](common/model.go#L223) fields description**

| **Name**          | **Type**            | **Description**                               |
| ----------------- | ------------------- | --------------------------------------------- |
| **Number**        | _**string**_        | Number of the document                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                 |

### **[SocialServiceID](common/model.go#L229) fields description**

| **Name**          | **Type**            | **Description**                               |
| ----------------- | ------------------- | --------------------------------------------- |
| **Number**        | _**string**_        | Number of the document                        |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                 |

### **[TaxID](common/model.go#L237) fields description**

| **Name**          | **Type**            | **Description**                               |
| ----------------- | ------------------- | --------------------------------------------- |
| **Number**        | _**string**_        | Number of the document                        |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                 |

### **[DriverLicense](common/model.go#L131) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license number                                   |
| **Version**       | _**string**_        | New Zealand driver license version number (this number changes each time a new card is issued) |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "DE"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "KY" |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                          |
| **ValidUntil**    | _**Time**_          | Valid until date, in RFC3339 format                     |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license   |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license    |

### **[DriverLicenseTranslation](common/model.go#L143) fields description**

| **Name**          | **Type**            | **Description**                                                   |
| ----------------- | ------------------- | ----------------------------------------------------------------- |
| **Number**        | _**string**_        | Driver license translation number                                 |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "LV"                |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "MT"           |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                                    |
| **ValidUntil**    | _**Time**_          | Valid until date, in RFC3339 format                               |
| **FrontImage**    | _***DocumentFile**_ | Scan or photo of the front side of the driver license translation |
| **BackImage**     | _***DocumentFile**_ | Scan or photo of the back side of the driver license translation  |

### **[CreditCard](common/model.go#L154) fields description**

| **Name**       | **Type**            | **Description**                                   |
| -------------- | ------------------- | ------------------------------------------------- |
| **Number**     | _**string**_        | Credit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date, in RFC3339 format               |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the credit card |

### **[DebitCard](common/model.go#L161) fields description**

| **Name**       | **Type**            | **Description**                                  |
| -------------- | ------------------- | ------------------------------------------------ |
| **Number**     | _**string**_        | Debit card number without whitespaces and dashes |
| **ValidUntil** | _**Time**_          | Valid until date, in RFC3339 format              |
| **Image**      | _***DocumentFile**_ | Scan or photo of the face side of the debit card |

### **[UtilityBill](common/model.go#L168) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ID" |
| **Image**         | _***DocumentFile**_ | Scan or photo of the utility bill                  |

### **[ResidencePermit](common/model.go#L174) fields description**

| **Name**          | **Type**            | **Description**                                    |
| ----------------- | ------------------- | -------------------------------------------------- |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "GB" |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                     |
| **ValidUntil**    | _**Time**_          | Valid until date, in RFC3339 format                |
| **Image**         | _***DocumentFile**_ | Scan or photo of the residence permit              |

### **[Agreement](common/model.go#L182) fields description**

| **Name**  | **Type**            | **Description**                |
| --------- | ------------------- | ------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the agreement |

### **[EmploymentCertificate](common/model.go#L192) fields description**

| **Name**       | **Type**            | **Description**                             |
| -------------- | ------------------- | ------------------------------------------- |
| **IssuedDate** | _**Time**_          | Issued date, in RFC3339 format              |
| **Image**      | _***DocumentFile**_ | Scan or photo of the employment certificate |

### **[Contract](common/model.go#L187) fields description**

| **Name**  | **Type**            | **Description**               |
| --------- | ------------------- | ----------------------------- |
| **Image** | _***DocumentFile**_ | Scan or photo of the contract |

### **[DocumentPhoto](common/model.go#L208) fields description**

| **Name**  | **Type**            | **Description**                            |
| --------- | ------------------- | ------------------------------------------ |
| **Image** | _***DocumentFile**_ | Scan or photo of the photo from a document |

### **[Selfie](common/model.go#L198) fields description**

| **Name**  | **Type**            | **Description** |
| --------- | ------------------- | --------------- |
| **Image** | _***DocumentFile**_ | Selfie image    |

### **[Avatar](common/model.go#L203) fields description**

| **Name**  | **Type**            | **Description**          |
| --------- | ------------------- | ------------------------ |
| **Image** | _***DocumentFile**_ | Profile image aka avatar |

### **[Other](common/model.go#L213) fields description**

| **Name**          | **Type**            | **Description**                                         |
| ----------------- | ------------------- | ------------------------------------------------------- |
| **Number**        | _**string**_        | Document number without whitespaces and dashes          |
| **CountryAlpha2** | _**string**_        | Country in ISO 3166-1 alpha-2 format, for ex. "ES"      |
| **State**         | _**string**_        | Abbreviated name of the state or province, for ex. "PA" |
| **IssuedDate**    | _**Time**_          | Issued date, in RFC3339 format                          |
| **ValidUntil**    | _**Time**_          | Valid until date, in RFC3339 format                     |
| **Image**         | _***DocumentFile**_ | Scan or photo of the document                           |

### **[DocumentFile](common/model.go#L71) fields description**

| **Name**        | **Type**     | **Description**                                               |
| --------------- | ------------ | ------------------------------------------------------------- |
| **Filename**    | _**string**_ | File name of the document image, for ex. "passport_front.jpg" |
| **ContentType** | _**string**_ | MIME type of the content, for ex. "image/jpeg"                |
| **Data**        | _**[]byte**_ | Raw content of the document image file                        |

### **[VideoAuth](common/model.go#L243) fields description**

| **Name**        | **Type**     | **Description**                                  |
| --------------- | ------------ | ------------------------------------------------ |
| **Filename**    | _**string**_ | Name of the video file, for ex. "auth_video.mp4" |
| **ContentType** | _**string**_ | MIME type of the content, for ex. "video/mp4"    |
| **Data**        | _**[]byte**_ | Raw content of the video file                    |

### **[Location](common/model.go#L57) fields description**

| **Name**      | **Type**     | **Description**                             |
| ------------- | ------------ | ------------------------------------------- |
| **Latitude**  | _**string**_ | The location latitude, for ex. "55.678849"  |
| **Longitude** | _**string**_ | The location longitude, for ex. "52.327662" |

## **KYC response**

The verification response consist of two elements: a result and an error if occurred. The result is of the type [**common.Result**](#commonresult-fields-description).

> Some KYC providers might require to poll the customer verification status to check if the process is completed. For this purpose the __*StatusCheck__ field is provided. If a polling is required and no error has occured then this field will be non-nil.

### **[common.Result](common/rest.go#L29) fields description**

| **Name**        | **Type**                                                  | **Description**                                                               |
| --------------- | --------------------------------------------------------- | ----------------------------------------------------------------------------- |
| **Status**      | _**[string](#status-possible-values-description)**_       | Status of the verification                                                    |
| **Details**     | _***[Details](#details-fields-description)**_             | Details of the verification if provided                                       |
| **ErrorCode**   | _**string**_                                              | Error code returned by a KYC provider if the provider support error codes     |
| **StatusCheck** | _***[KYCStatusCheck](#kycstatuscheck-fields-description)**_ | Data required to do the customer verification status check requests if needed |

### **[Status](common/mapping.go#L3) possible values description**

| **Value**    | **Description**                                                                                                                |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| **Error**    | Verification has failed. Probably, some error has occurred. Returned error value must be non-nil and **`common.Result.ErrorCode`** may contain error code value |
| **Approved** | Successful verification with approved result. The details maybe non-nil and contain additional info about the verification     |
| **Denied**   | Successful verification with rejected result. The details should be non-nil and contain additional info about the verification |
| **Unclear**  | Needs subsequent status polling or the verification completed with an indefinite result. That might mean that some additional info is required. The details should be non-nil and contain additional info. If status polling is required then **`common.Result.StatusCheck`** must be non-nil |

### **[Details](common/rest.go#L37) fields description**

| **Name**     | **Type**                                              | **Description**                                                          |
| ------------ | ----------------------------------------------------- | ------------------------------------------------------------------------ |
| **Finality** | [_**string**_](#finality-possible-values-description) | Rejection type of the result (if the negative answer is given)           |
| **Reasons**  | _**[]string**_                                        | List of additional response info describing result-related circumstances |

### **[Finality](common/mapping.go#L11) possible values description**

| **Value**    | **Description**                                                                                                             |
| ------------ | --------------------------------------------------------------------------------------------------------------------------- |
| **Final**    | Final reject, e.g. when a person is a fraudster, or a client does not want to accept such kind of clients in his/her system |
| **NonFinal** | A reject that can be fixed, e.g. by uploading an image of better quality                                                    |
| **Unknown**  | The provider doesn't support **`Finality`** feature                                                                         |

### **[KYCStatusCheck](common/model.go#L92) fields description**

| **Name**        | **Type**                                | **Description**                                                |
| --------------- | --------------------------------------- | -------------------------------------------------------------- |
| **Provider**    | _**[KYCProvider](common/enum.go#L36)**_ | An identificator for the KYC provider name                     |
| **ReferenceID** | _**string**_                            | An identificator that references to this verification submission. It mention in docs as applicantId/mtid/jumioIdScanReference/etc. Its value is specific for a provider |
| **LastCheck**   | _**time.Time**_                         | Last time a verification status was checked, in RFC3339 format |

## **Applicable fields grouped per provider**

[**common.UserData**](#userdata-fields-description) provides a wide range of possible data that might require the verification. However, not every KYC provider will surely use all available fields of the model. Therefore, to ease the process of integration for administrators, here you'll find the grouping of applicable fields per provider.

### **Coinfirm**

#### **Individual KYC onboarding**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**                     | **Type**         | **Required**        | **Comment**                                |
| ---------------------------- | ---------------- | :-----------------: | ------------------------------------------ |
| **FirstName**                | _string_         | **Yes**             |                                            |
| **LastName**                 | _string_         | **Yes**             |                                            |
| MiddleName                   | _string_         |                     |                                            |
| **Email**                    | _string_         | **Yes**             |                                            |
| IPaddress                    | _string_         |                     |                                            |
| **DateOfBirth**              | _Time_           | **Yes**             |                                            |
| **CountryAlpha2**            | _string_         | **Yes**             |                                            |
| **Nationality**              | _string_         | **Yes**             |                                            |
| Phone                        | _string_         |                     |                                            |
| MobilePhone                  | _string_         |                     |                                            |
| **CurrentAddress**           | _Address_        | **Yes**             |                                            |
| **Passport**                 | _*Passport_      | __*__ (see comment) | __*__ Provide anyone of required documents |
| **IDCard**                   | _*IDCard_        | __*__ (see comment) |                                            |
| **SNILS**                    | _*SNILS_         | __*__ (see comment) |                                            |
| **DriverLicense**            | _*DriverLicense_ | __*__ (see comment) |                                            |
| **DriverLicenseTranslation** | _*DriverLicenseTranslation_ | __*__ (see comment)  |                                |
| UtilityBill                  | _*UtilityBill_   |                     |                                            |

[**Address**](#address-fields-description) mandatory fields:

| **Name**              | **Type** |
| --------------------- | -------- |
| **Town**              | _string_ |
| **Street**            | _string_ |
| **PostCode**          | _string_ |

> **DOCUMENTS NOTE:** Supported extensions for document files: **"jpg", "jpeg", "png", "gif", "bmp", "svg", "psd", "tif", "tiff", "webp", "pdf"**.

#### **Company KYC onboarding**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**                | **Type**               | **Required**         | **Comment**                                 |
| ----------------------- | ---------------------- | :------------------: | ------------------------------------------- |
| **Email**               | _string_               | **Yes**              |                                             |
| IPaddress               | _string_               |                      |                                             |
| **CountryAlpha2**       | _string_               | **Yes**              |                                             |
| **CurrentAddress**      | _Address_              | **Yes**              | As for Individual                           |
| **Passport**            | _*Passport_            | __*__ (see comment)  | __*__ Provide anyone of required documents  |
| **IDCard**              | _*IDCard_              | __*__ (see comment)  |                                             |
| **DriverLicense**       | _*DriverLicense_       | __*__ (see comment)  |                                             |
| **IsCompany**           | _bool_                 | **Yes**              | Must be set to "true"                       |
| **CompanyName**         | _string_               | **Yes**              |                                             |
| **Website**             | _string_               |                      |                                             |
| **CompanyBoard**        | _*CompanyBoard_        | __**__ (see comment) | __**__ Provide anyone of required documents |
| **CompanyRegistration** | _*CompanyRegistration_ | __**__ (see comment) |                                             |

> **DOCUMENTS NOTE:** Supported extensions for document files: **"jpg", "jpeg", "png", "gif", "bmp", "svg", "psd", "tif", "tiff", "webp", "pdf"**.

### **ComplyAdvantage**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**      | **Type**     | **Required** | **Comment**                                            |
| ------------- | ------------ | :----------: | ------------------------------------------------------ |
| **FirstName** | _**string**_ | **(*)**      | __*__ Either provide first and last names or full name |
| **LastName**  | _**string**_ | **(*)**      |                                                        |
| MiddleName    | _string_     |              |                                                        |
| **FullName**  | _**string**_ | **(*)**      | __*__ Either provide this or first and last names      |
| DateOfBirth   | _Time_       |              | Recommend for better results                           |

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

### **Shufti Pro**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**             | **Type**         | **Required**       | **Comment**                                 |
| -------------------- | ---------------- | :----------------: | ------------------------------------------- |
| **FirstName**        | _string_         | **Yes**            |                                             |
| **LastName**         | _string_         | **Yes**            |                                             |
| MiddleName           | _string_         |                    |                                             |
| **Email**            | _string_         |                    |                                             |
| **DateOfBirth**      | _Time_           |                    |                                             |
| **CountryAlpha2**    | _string_         | **Yes**            |                                             |
| CurrentAddress       | _Address_        |                    |                                             |
| **Passport**(*)      | _*Passport_      | **See comment(*)** | (*)Anyone of documents marked with asterisk |
| **IDCard**(*)        | _*IDCard_        | **(*)**            |                                             |
| SNILS                | _*SNILS_         |                    |                                             |
| **DriverLicense**(*) | _*DriverLicense_ | **(*)**            |                                             |
| **CreditCard**(*)    | _*CreditCard_    | **(*)**            |                                             |
| **DebitCard**(*)     | _*DebitCard_     |                    |                                             |
| UtilityBill          | _*UtilityBill_   |                    | Can be used in case of address verification |
| **Selfie**           | _*Selfie_        | **Yes**            |                                             |

> **DOCUMENTS NOTE:** Include image file(s) for the document used for the verification.

### **Sum&Substance**

According to the [API Reference](https://developers.sumsub.com) all fields of [**UserData**](#userdata-fields-description) are applicable except the following:

* MaternalLastName
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

### **SynapseFI**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**             | **Type**         | **FileType**         | **Required** | **Comment**                                         |
| -------------------- | ---------------- | -------------------- | :----------: | --------------------------------------------------- |
| **FirstName**        | _string_         |                      | **Yes**      |                                                     |
| **LastName**         | _string_         |                      | **Yes**      |                                                     |
| MiddleName           | _string_         |                      |              |                                                     |
| Gender               | _Gender_         |                      |              |                                                     |
| **Email**            | _string_         |                      | **Yes**      |                                                     |
| **DateOfBirth**      | _Time_           |                      | **Yes**      | Required for documents only                         |
| **CountryAlpha2**    | _string_         |                      | **Yes**      |                                                     |
| **Phone**            | _string_         |                      | **(**)**     | (**)Anyone of documents marked with double asterisk |
| **Mobile phone**     | _string_         |                      | **(**)**     |                                                     |
| **CurrentAddress**   | _Address_        |                      | **Yes**      | Required for documents only                         |
| **Passport**(*)      | _*Passport_      | .png/.jpg/.jpeg      | **(*)**      | (*)Anyone of documents marked with asterisk         |
| **IDCard**(*)        | _*IDCard_        | .png/.jpg/.jpeg      | **(*)**      |                                                     |
| **DriverLicense**(*) | _*DriverLicense_ | .png/.jpg/.jpeg      | **(*)**      |                                                     |
| UtilityBill          | _*UtilityBill_   | .png/.jpg/.jpeg/.pdf |              |                                                     |
| **Selfie**           | _*Selfie_        | .png/.jpg/.jpeg      | **Yes**      | Deprecated in favor of video authorization          |
| **VideoAuth**        | _*VideoAuth_     | .mov/.mp4/.webm      | **Yes**      | 5 second authorization video of the customer. Requires document with customer photo (passport, etc.) |

> **DOCUMENTS NOTE:** Include image file(s) for the document used for the verification.

### **Thomson Reuters**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**             | **Type**     | **Required**        | **Comment**                                                           |
| -------------------- | ------------ | :-----------------: | --------------------------------------------------------------------- |
| **FirstName**        | _**string**_ | __*__ (see comment) | __*__ Either use first name and last name or full name                |
| **LastName**         | _**string**_ | __*__ (see comment) | if the order of the names is different from western usual composition |
| MiddleName           | _string_     |                     | or those names are difficult to represent separately                  |
| **FullName**         | _**string**_ | __*__ (see comment) |                                                                       |
| Gender               | _Gender_     |                     |                                                                       |
| DateOfBirth          | _Time_       |                     |                                                                       |
| CountryOfBirthAlpha2 | _string_     |                     |                                                                       |
| CountryAlpha2        | _string_     |                     |                                                                       |
| Nationality          | _string_     |                     |                                                                       |

> **NOTE:** For better result, please, fill as much fields as possible.

### **Trulioo**

[**UserData**](#userdata-fields-description) applicable fields:

| **Name**                 | **Type**           | **Required** | **Comment** |
| ------------------------ | ------------------ | :----------: | ----------- |
| FirstName                | _string_           |              |             |
| MaternalLastName         | _string_           |              |             |
| LastName                 | _string_           |              |             |
| MiddleName               | _string_           |              |             |
| LatinISO1Name            | _string_           |              |             |
| Email                    | _string_           |              |             |
| Gender                   | _Gender_           |              |             |
| DateOfBirth              | _Time_             |              |             |
| **CountryAlpha2**        | _string_           | **Yes**      |             |
| Phone                    | _string_           |              |             |
| MobilePhone              | _string_           |              |             |
| BankAccountNumber        | _string_           |              |             |
| VehicleRegistrationPlate | _string_           |              |             |
| CurrentAddress           | _Address_          |              |             |
| Business                 | _*Business_        |              |             |
| Passport                 | _*Passport_        |              |             |
| IDCard                   | _*IDCard_          |              |             |
| HealthID                 | _*HealthID_        |              |             |
| SocialServiceID          | _*SocialServiceID_ |              |             |
| TaxID                    | _*TaxID_           |              |             |
| DriverLicense            | _*DriverLicense_   |              |             |
| ResidencePermit          | _*ResidencePermit_ |              |             |
| Selfie                   | _*Selfie_          |              |             |

> **DOCUMENTS NOTE:** Include image file(s) for a document used in the verification (some documents haven't physical form only a number, for ex. UK NI and NHS Numbers).

### **The countries supported by KYC providers and the fields variability**

KYC providers may require various set of `common.UserData` fields depending on the customer country. Also, they may service to the limited number of countries and this number of countries might configurable in a web-interface of the provider.

### **Coinfirm covered countries**

* International
* No fields variations found in the docs

### **ComplyAdvantage covered countries**

* International (no list of supported countries)
* No fields variations found in the docs

### **IdentityMind covered countries**

* International
* No fields variations found in the docs

### **IDology covered countries**

* USA and Canada
* No fields variations found in the docs

### **Jumio covered countries**

* International
* No fields variations found in the docs

### **Shufti Pro covered countries**

* International (["Shufti Pro provides support for all countries"](https://github.com/shuftipro/RESTful-API-v1.3/blob/master/off-site_without_ocr/countries.md#supported-countries))
* No fields variations found in the docs

### **Sum&Substance covered countries**

* International
* No fields variations found in the docs

### **SynapseFI covered countries**

* International (no list of supported countries)
* No fields variations found in the docs

### **Thomson Reuters covered countries**

* International
* No fields variations found in the docs

### **Trulioo covered countries**

These are the countries that supported since last check.

| **Country code** | **Country Name**                                     |
| :--------------: | ---------------------------------------------------- |
| AE               | United Arab Emirates                                 |
| AR               | Argentina                                            |
| AT               | Austria                                              |
| AU               | Australia                                            |
| BE               | Belgium                                              |
| BR               | Brazil                                               |
| CA               | Canada                                               |
| CH               | Switzerland                                          |
| CL               | Chile                                                |
| CN               | China                                                |
| CO               | Colombia                                             |
| CR               | Costa Rica                                           |
| DE               | Germany                                              |
| DK               | Denmark                                              |
| EC               | Ecuador                                              |
| EG               | Egypt                                                |
| ES               | Spain                                                |
| FR               | France                                               |
| GB               | United Kingdom of Great Britain and Northern Ireland |
| HK               | Hong Kong                                            |
| IE               | Ireland                                              |
| IN               | India                                                |
| IT               | Italy                                                |
| JP               | Japan                                                |
| KR               | Republic of Korea                                    |
| KW               | Kuwait                                               |
| LB               | Lebanon                                              |
| MX               | Mexico                                               |
| MY               | Malaysia                                             |
| NL               | Netherlands                                          |
| NZ               | New Zealand                                          |
| OM               | Oman                                                 |
| PE               | Peru                                                 |
| PT               | Portugal                                             |
| RO               | Romania                                              |
| RU               | Russian Federation                                   |
| SA               | Saudi Arabia                                         |
| SE               | Sweden                                               |
| SG               | Singapore                                            |
| SV               | El Salvador                                          |
| TH               | Thailand                                             |
| UA               | Ukraine                                              |
| US               | United States of America                             |
| ZA               | South Africa                                         |

* API provides the group of methods for retrieving the lists of:
  * Consents
  * Supported countries
  * Available fields dynamically based on a country
  * Document Types available for a country
  * Test Entities configured for a country
  * Datasource groups configured for a country

[**UserData**](#userdata-fields-description) applicable fields for all supported countries:

| **Name**           | **Type**             | **Required** | **Countries for which the field is required** |
| -----------------  | -------------------- | :----------: | --------------------------------------------- |
| **FirstName**      | **_string_**         |              | AE, AR, AU, BR, CA, CL, CN, CO, CR, DE, DK, EC, EG, ES, FR, GB, HK, IE, IT, JP, KW, LB, MX, MY, NL, NZ, OM, PE, RU, SA, SE, SG, SV, TH, UA, US, ZA |
| **LastName**       | **_string_**         |              | AE, AR, AT, AU, BE, BR, CA, CH, CL, CN, CO, CR, DE, DK, EC, EG, ES, FR, GB, HK, IE, IN, IT, JP, KW, LB, MX, MY, NL, NZ, OM, PE, PE, PT, RO, RU, SA, SE, SG, SV, TH, UA, US, ZA |
| MaternalLastName   | _string_             |              |                                               |
| **MiddleName**     | _string_             |              | RU                                            |
| LatinISO1Name      | _string_             |              |                                               |
| Email              | _string_             |              |                                               |
| **Gender**         | _Gender_             |              | GB, MX, MY                                    |
| **DateOfBirth**    | **_Time_**           |              | AE, AU, CR, DE, DK, EG, FR, GB, IT, KR, KW, LB, MX, MY, NL, NZ, OM, RU, SA, SE, SV, ZA |
| **CountryAlpha2**  | **_string_**         | **Yes**      |                                               |
| **Phone**          | **_string_**         |              | CR                                            |
| MobilePhone        | _string_             |              |                                               |
| **CurrentAddress** | **_Address_**        |              | AR AU BE CA CH CR DE DK ES FR GB IE IT JP MX NL NZ PE PT SE US ZA |
| Business           | _*Business_          |              |                                               |
| **Passport**       | **_*Passport_**      |              | AE, AR, AT, AU, BE, BR, CA, CH, CL, CN, CO, CR, DE, DK, EC, EG, ES, FR, GB, HK, IE, IT, JP, KR, KW, LB, MX, MY, NL, NZ, OM, PE, PT, RU, SA, SE, SG, SV, TH, UA, US, ZA |
| IDCard             | _*IDCard_            |              | AE, AR, BR, CN, CO, CR, DK, EC, EG, FR, HK, KR, KW, LB, MX, MY, NL, OM, RO, SA, SE, SG, SV, TH, ZA |
| SocialServiceID    | _*SocialServiceID_   |              | CA, GB, IE, IT, UA                            |
| TaxID              | _*TaxID_             |              |                                               |
| **DriverLicense**  | **_*DriverLicense_** |              | GB, KR, NZ, US                                |
| ResidencePermit    | _*ResidencePermit_   |              |                                               |
| Selfie             | _*Selfie_            |              |                                               |

[**UserData**](#userdata-fields-description) required fields for the specific countries:

| **Name**                 | **Type**    | **Countries for which the field is required** |
| ------------------------ | ----------- | --------------------------------------------- |
| MaternalLastName         | _string_    | CO, MX, PE                                    |
| FullName                 | _string_    | MY, SG                                        |
| CountryOfBirthAlpha2     | _string_    | MY                                            |
| StateOfBirth             | _string_    | MX, MY                                        |
| BankAccountNumber        | _string_    | CN                                            |
| VehicleRegistrationPlate | _string_    | NZ                                            |
| HealthID                 | _*HealthID_ | GB                                            |

[**Address**](#address-fields-description) required fields for the specific countries:

| **Name**          | **Type** | **Countries for which the field is required**                          |
| ----------------- | -------- | ---------------------------------------------------------------------- |
| County            | _string_ | IE                                                                     |
| Town              | _string_ | CA, CH, CR, DE, IE, JP, NZ, PT, SE, US, ZA                             |
| Suburb            | _string_ | JP, ZA                                                                 |
| Street            | _string_ | AU, CA, CH, DE, DK, ES, FR, GB, IE, MX, NZ, PT, US                     |
| BuildingNumber    | _string_ | AU, CA, DE, ES, JP, MX, NL, US                                         |
| PostCode          | _string_ | AR, AU, BE, CA, CH, DE, ES, FR, GB, IT, JP, MX, NL, NZ, PE, PT, US, ZA |
| StateProvinceCode | _string_ | CA, CR, JP, US                                                         |

[**Passport**](#passport-fields-description) required fields for the specific countries:

| **Name**      | **Type**        | **Countries for which the field is required**                                                                                                                      |
| ------------- | --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Number        | _string_        | AU, RU                                                                                                                                                             |
| Mrz1          | _string_        | AE, AR, AT, AU, BE, BR, CA, CH, CL, CN, CO, CR, DE, DK, EC, EG, ES, FR, GB, HK, IE, IT, JP, KR, KW, LB, MX, MY, NL, NZ, OM, PE, PT, SA, SE, SG, SV, TH, UA, US, ZA |
| Mrz2          | _string_        | AE, AR, AT, AU, BE, BR, CA, CH, CL, CN, CO, CR, DE, DK, EC, EG, ES, FR, GB, HK, IE, IT, JP, KR, KW, LB, MX, MY, NL, NZ, OM, PE, PT, SA, SE, SG, SV, TH, UA, US, ZA |
| IssuedDate    | _Time_          | RU                                                                                                                                                                 |

[**IDCard**](#idcard-fields-description) required fields for the specific countries:

| **Name** | **Type** | **Countries for which the field is required**                                                                  |
| -------- | -------- | -------------------------------------------------------------------------------------------------------------- |
| Number   | _string_ | AE, AR, BR, CA, CN, CO, CR, DK, EC, EG, FR, HK, IE, IT, KW, LB, MX, MY, NL, OM, RO, SA, SE, SG, SV, TH, UA, ZA |

[**HealthID**](#healthid-fields-description) required fields for the specific countries:

| **Name** | **Type** | **Countries for which the field is required** |
| -------- | -------- | --------------------------------------------- |
| Number   | _string_ | GB                                            |

[**SocialServiceID**](#socialserviceid-fields-description) required fields for the specific countries:

| **Name** | **Type** | **Countries for which the field is required** |
| -------- | -------- | --------------------------------------------- |
| Number   | _string_ | CA, GB, IE, IT, UA                            |

[**DriverLicense**](#driverlicense-fields-description) required fields for the specific countries:

| **Name** | **Type** | **Countries for which the field is required** |
| -------- | -------- | --------------------------------------------- |
| Number   | _string_ | GB, KR, NZ, US                                |
| Version  | _string_ | NZ                                            |
| State    | _string_ | US                                            |
