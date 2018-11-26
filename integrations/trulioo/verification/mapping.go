package verification

import (
	"encoding/base64"
	"strconv"
	"time"

	"modulus/kyc/common"
)

// MapCustomerToDataFields converts input customer data to the format acceptable by the API.
func MapCustomerToDataFields(customer *common.UserData) DataFields {
	return DataFields{
		PersonInfo:      mapCustomerToPersonalInfo(customer),
		Location:        mapCustomerAddressToLocation(customer.CurrentAddress),
		Communication:   mapCustomerToCommunication(customer),
		Document:        mapCustomerDocument(customer),
		Business:        mapCustomerBusiness(customer.Business),
		Passport:        mapCustomerPassport(customer.Passport),
		DriverLicence:   mapCustomerDriverLicence(customer.DriverLicense),
		NationalIds:     mapCustomerToNationalIds(customer),
		CountrySpecific: mapCustomerToCountrySpecific(customer),
	}
}

func mapCustomerToPersonalInfo(customer *common.UserData) *PersonInfo {
	dateOfBirth := time.Time(customer.DateOfBirth)

	pi := &PersonInfo{
		FirstGivenName: customer.FirstName,
		MiddleName:     customer.MiddleName,
		FirstSurName:   customer.LastName,
		SecondSurname:  customer.MaternalLastName,
		ISOLatin1Name:  customer.LatinISO1Name,
		DayOfBirth:     dateOfBirth.Day(),
		MonthOfBirth:   int(dateOfBirth.Month()),
		YearOfBirth:    dateOfBirth.Year(),
		Gender:         mapGender(customer.Gender),
	}

	switch customer.CountryAlpha2 {
	case "MY", "SG":
		pi.AdditionalFields = &PIAdditionalFields{
			FullName: customer.FullName(),
		}
	}

	return pi
}

func mapCustomerAddressToLocation(address common.Address) *Location {
	if address == (common.Address{}) {
		return nil
	}

	l := &Location{
		BuildingNumber:    address.BuildingNumber,
		BuildingName:      address.BuildingName,
		UnitNumber:        address.FlatNumber,
		StreetName:        address.Street,
		StreetType:        address.StreetType,
		City:              address.Town,
		Suburb:            address.Suburb,
		County:            address.County,
		Country:           address.CountryAlpha2,
		StateProvinceCode: address.StateProvinceCode,
		PostalCode:        address.PostCode,
		POBox:             address.PostOfficeBox,
	}

	if address.CountryAlpha2 == "ZA" {
		l.AdditionalFields = &AdditionalFields{
			Address1: address.String(),
		}
	}

	return l
}

func mapCustomerToCommunication(customer *common.UserData) *Communication {
	if customer.Email == "" && customer.Phone == "" && customer.MobilePhone == "" {
		return nil
	}
	return &Communication{
		MobileNumber: customer.MobilePhone,
		Telephone:    customer.Phone,
		EmailAddress: customer.Email,
	}
}

func mapCustomerDocument(customer *common.UserData) (document *Document) {
	document = &Document{}

	if customer.Selfie != nil && customer.Selfie.Image != nil {
		document.LivePhoto = base64.StdEncoding.EncodeToString(customer.Selfie.Image.Data)
	}

	if customer.Passport != nil && customer.Passport.Image != nil {
		document.DocumentType = "Passport"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.Passport.Image.Data)
		return
	}

	if customer.DriverLicense != nil && customer.DriverLicense.FrontImage != nil {
		document.DocumentType = "DrivingLicence"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.FrontImage.Data)
		if customer.DriverLicense.BackImage != nil {
			document.DocumentBackImage = base64.StdEncoding.EncodeToString(customer.DriverLicense.BackImage.Data)
		}
		return
	}

	if customer.IDCard != nil && customer.IDCard.Image != nil {
		document.DocumentType = "IdentityCard"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.IDCard.Image.Data)
		return
	}

	if customer.ResidencePermit != nil && customer.ResidencePermit.Image != nil {
		document.DocumentType = "ResidencePermit"
		document.DocumentFrontImage = base64.StdEncoding.EncodeToString(customer.ResidencePermit.Image.Data)
		return
	}

	return nil
}

func mapCustomerBusiness(business *common.Business) *Business {
	if business == nil {
		return nil
	}

	incorporationDate := time.Time(business.IncorporationDate)

	return &Business{
		BusinessName:                business.Name,
		BusinessRegistrationNumber:  business.RegistrationNumber,
		DayOfIncorporation:          incorporationDate.Day(),
		MonthOfIncorporation:        int(incorporationDate.Month()),
		YearOfIncorporation:         incorporationDate.Year(),
		JurisdictionOfIncorporation: business.IncorporationJurisdiction,
	}
}

func mapCustomerPassport(passport *common.Passport) *Passport {
	if passport == nil {
		return nil
	}

	p := &Passport{
		Number: passport.Number,
		Mrz1:   passport.Mrz1,
		Mrz2:   passport.Mrz2,
	}

	if !time.Time(passport.ValidUntil).IsZero() {
		p.YearOfExpiry = time.Time(passport.ValidUntil).Year()
		p.MonthOfExpiry = int(time.Time(passport.ValidUntil).Month())
		p.DayOfExpiry = time.Time(passport.ValidUntil).Day()
	}

	return p
}

func mapCustomerDriverLicence(drivers *common.DriverLicense) *DriverLicence {
	if drivers == nil {
		return nil
	}

	return &DriverLicence{
		Number: drivers.Number,
	}
}

func mapCustomerToNationalIds(customer *common.UserData) (nIDs []NationalID) {
	switch customer.CountryAlpha2 {
	case "GB":
		if len(customer.UKNHSNumber) > 0 {
			nIDs = append(nIDs, NationalID{
				Number: customer.UKNHSNumber,
				Type:   "health",
			})
		}
		if len(customer.UKNINumber) > 0 {
			nIDs = append(nIDs, NationalID{
				Number: customer.UKNINumber,
				Type:   "socialservice",
			})
		}
	case "AE", "AR", "BR", "CN", "CO", "CR", "DK", "EC", "EG", "FR", "HK", "KW",
		"LB", "MX", "MY", "NL", "OM", "RO", "SA", "SE", "SG", "SV", "TH", "ZA":
		if customer.IDCard != nil {
			nIDs = append(nIDs, NationalID{
				Number: customer.IDCard.Number,
				Type:   "nationalid",
			})
		}
	case "CA", "IE", "IT", "UA":
		if customer.IDCard != nil {
			nIDs = append(nIDs, NationalID{
				Number: customer.IDCard.Number,
				Type:   "socialservice",
			})
		}
	}

	return
}

func mapCustomerToCountrySpecific(customer *common.UserData) map[CountryCode]CountrySpecific {
	cspec := CountrySpecific{}

	switch customer.CountryAlpha2 {
	case "AE", "AR", "AT", "BE", "BR", "CA", "CH", "CL", "CO", "CR", "DE", "DK",
		"EC", "EG", "ES", "FR", "HK", "IE", "IT", "JP", "KW", "LB", "NL", "OM",
		"PE", "PT", "SA", "SE", "SG", "SV", "TH", "UA", "ZA":
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
	case "AU":
		if customer.Passport != nil {
			cspec.PassportCountry = customer.Passport.CountryAlpha2
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
			cspec.PassportNumber = customer.Passport.Number
		}
	case "CN":
		cspec.BankAccountNumber = customer.BankAccountNumber
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
	case "GB":
		if customer.DriverLicense != nil {
			cspec.DriverLicenceNumber = customer.DriverLicense.Number
		}
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
	case "KR":
		if customer.DriverLicense != nil {
			cspec.DriverLicenceNumber = customer.DriverLicense.Number
		}
		if customer.IDCard != nil {
			cspec.NameOnCard = customer.FirstName
			cspec.SerialNumber = customer.IDCard.Number
		}
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
	case "MX":
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
		cspec.StateOfBirth = customer.StateOfBirth
	case "MY":
		cspec.CountryOfBirth = customer.CountryOfBirthAlpha2
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
		cspec.StateOfBirth = customer.StateOfBirth
	case "NZ":
		if customer.DriverLicense != nil {
			cspec.DriverLicenceNumber = customer.DriverLicense.Number
			cspec.DriverLicenceVerNumber = customer.DriverLicense.Version
		}
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
		cspec.VehicleRegistrationPlate = customer.VehicleRegistrationPlate
	case "RU":
		if customer.Passport != nil {
			cspec.DayOfIssue = strconv.Itoa(time.Time(customer.Passport.IssuedDate).Day())
			cspec.MonthOfIssue = strconv.Itoa(int(time.Time(customer.Passport.IssuedDate).Month()))
			cspec.YearOfIssue = strconv.Itoa(time.Time(customer.Passport.IssuedDate).Year())
			if len(customer.Passport.Number) > 4 {
				cspec.PassportSerie = customer.Passport.Number[:4]
				cspec.InternalPassportNumber = customer.Passport.Number[4:]
			}
		}
	case "US":
		if customer.DriverLicense != nil {
			cspec.DriverLicenceNumber = customer.DriverLicense.Number
			cspec.DriverLicenceState = customer.DriverLicense.State
		}
		if customer.Passport != nil {
			cspec.PassportMRZLine1 = customer.Passport.Mrz1
			cspec.PassportMRZLine2 = customer.Passport.Mrz2
		}
	default:
		return nil
	}

	if cspec == (CountrySpecific{}) {
		return nil
	}

	return map[CountryCode]CountrySpecific{
		customer.CountryAlpha2: cspec,
	}
}

func mapGender(gender common.Gender) string {
	switch gender {
	case common.Male:
		return "M"
	case common.Female:
		return "F"
	default:
		return ""
	}
}
