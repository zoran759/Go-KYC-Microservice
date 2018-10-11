package main

import (
	"modulus/kyc/main/handlers"
	"os"

	"fmt"
	"log"
	"net/http"
)

func main() {
	// Set up the CheckCustomer API handler.
	http.HandleFunc("/CheckCustomer", handlers.CheckCustomerHandler)

	// Set up the CheckStatus API handler.
	http.HandleFunc("/CheckStatus", handlers.CheckStatusHandler)

	// If the service port isn't set in the envvar use the default value instead.
	var port string
	if port = os.Getenv("KYC_SERVER_PORT"); len(port) == 0 {
		port = "80"
	}

	// Start(Blocking) the server
	log.Printf("Listen on :%v", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}
