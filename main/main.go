package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"
)

var configFile = flag.String("config", "kyc.cfg", "Configuration file for KYC providers")

func main() {
	flag.Parse()

	if err := config.FromFile(*configFile); err != nil {
		log.Fatalln("Loading configuration:", err)
	}

	createHandlers()

	var port string

	if port = os.Getenv("KYC_PORT"); len(port) == 0 {
		port = "8080"
	}

	log.Printf("Listen on :%v", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}
}

func createHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the KYC service. Have a nice day!\n"))
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong!"))
	})
	http.HandleFunc("/CheckCustomer", handlers.CheckCustomerHandler)
	http.HandleFunc("/CheckStatus", handlers.CheckStatusHandler)
}
