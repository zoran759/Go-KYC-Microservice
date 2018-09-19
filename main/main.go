package main

import (
	"github.com/grigored/cex/src/cecommon/configs"
	"gitlab.com/lambospeed/kyc/main/handlers"

	"fmt"
	"log"
	"net/http"
)

func main() {

	// Set up the CheckCustomer API handler. Strip prefix of the API version.
	apiHandler := http.StripPrefix(
		"/api/v1",
		http.HandlerFunc(handlers.CheckCustomerHandler),
	)
	http.Handle("/api/v1/", apiHandler)

	// Start(Blocking) the server
	log.Printf("KYC http server started on :%v", configs.GetPort(configs.KYCServer))
	err := http.ListenAndServe(fmt.Sprintf(":%v", configs.GetPort(configs.KYCServer)), apiHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
