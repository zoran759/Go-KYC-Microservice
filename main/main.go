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
	log.Printf("KYC http server started on :%v", configs.GetPort(configs.KycServer))
	err := http.ListenAndServe(fmt.Sprintf(":%v", configs.GetPort(configs.KycServer)), apiHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
