package handlers

import (
	"gitlab.com/lambospeed/kyc/common"
	"gitlab.com/lambospeed/kyc/common/kycErrors"
	"gitlab.com/lambospeed/kyc/integrations/shuftipro"

	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Represents response for the CheckCustomer request
type checkCustomerResponse struct {

	// KYC result
	KYCResult common.KYCResult

	// KYC detailed result
	DetailedKYCResult *common.DetailedKYCResult
}

// Handler for the CustomerHandler function
func CheckCustomerHandler(w http.ResponseWriter, r *http.Request) {

	// Parse request parameters
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	// Read photo of user's ID
	id, err := ioutil.ReadFile(r.PostFormValue("idPhoto"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	// Read photo of user's face
	face, err := ioutil.ReadFile(r.PostFormValue("facePhoto"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
		return
	}

	// Assemble customer data
	customer := &common.UserData{
		FirstName:        r.PostFormValue("firstName"),
		PaternalLastName: r.PostFormValue("paternalLastName"),
		LastName:         r.PostFormValue("lastName"),
		MiddleName:       "M",
		CountryAlpha2:    "US",
		Phone:            "+3221-214-4456",
		Email:            "jsmith@yahoo.com",
		DateOfBirth:      common.Time(time.Date(1982, 4, 3, 0, 0, 0, 0, time.UTC)),
		SupplementalAddresses: []common.Address{
			{},
		},
		CurrentAddress: common.Address{
			PostCode:          "90010",
			Town:              "Chicago",
			BuildingNumber:    "452",
			FlatNumber:        "2",
			State:             "MI",
			StateProvinceCode: "MI",
			StreetType:        "Avenue",
			Street:            "Michigan",
		},
		Documents: []common.Document{
			{
				Metadata: common.DocumentMetadata{
					Type:    common.IDCard,
					Country: "RUS",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        id,
				},
				Back: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        id,
				},
			},
			{
				Metadata: common.DocumentMetadata{
					Type:    common.Selfie,
					Country: "RUS",
				},
				Front: &common.DocumentFile{
					Filename:    "passport.png",
					ContentType: "image/png",
					Data:        face,
				},
			},
		},
	}

	log.Printf("CustomerHandler request received Customer:%#v\nProvider:%v\n", customer, r.PostFormValue("provider"))

	// Prepare result variables
	var kycRes common.KYCResult
	var detailedKYCRes *common.DetailedKYCResult

	switch r.PostFormValue("provider") {

	case "shuftipro":
		{
			//Example Shufti Pro integration
			service := shuftipro.New(shuftipro.Config{
				Host:        "https://api.shuftipro.com",
				ClientID:    "ac93f3a0fee5afa2d9399d5d0f257dc92bbde89b1e48452e1bfac3c5c1dc99db",
				SecretKey:   "lq34eOTxDe1e6G8a1P7Igqo5YK3ABCDF",
				RedirectURL: "https://api.shuftipro.com",
			})

			// Make a request to the KYC provider
			kycRes, detailedKYCRes, err = service.CheckCustomer(customer)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: err.Error()})
				return
			}

			log.Printf("Res: %#v\n", kycRes)
			log.Printf("detailedRes: %#v\n", detailedKYCRes)
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(kycErrors.ErrorResponse{Error: kycErrors.InvalidKYCProvider.Error()})
	}

	// Assemble response
	response := checkCustomerResponse{
		KYCResult:         kycRes,
		DetailedKYCResult: detailedKYCRes,
	}

	// Send the response over HTTP
	json.NewEncoder(w).Encode(response)
}
