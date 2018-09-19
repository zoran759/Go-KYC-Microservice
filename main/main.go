package main

import (
	"gitlab.com/lambospeed/kyc/main/handlers"
	"os"

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
	log.Printf("KYC http server started on :%v", os.Getenv("KYC_SERVER_PORT"))
	err := http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("KYC_SERVER_PORT")), apiHandler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
