package main

import (
	"fmt"
	"log"
	"modulus/common/configs"
	"modulus/kyc/main/handlers"
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
	log.Printf("KYC http server started on :%v", configs.KycServerPort)
	err := http.ListenAndServe(fmt.Sprintf(":%v", configs.KycServerPort), apiHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
