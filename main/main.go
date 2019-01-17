package main

import (
	client "modulus/common/licensing-client"
	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"

	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	prodCfgFile = "kyc.cfg"
	devCfgFile  = "kyc_dev.cfg"
)

// DevEnv is the flag to manage prod/dev builds.
// For a production build, this flag value should be set to "false" upon compilation time using: [-ldflags "-X main.DevEnv=false"]
var DevEnv = "true"

func main() {

	// Validate license in production environment.
	// Do os.Exit if failed.
	if DevEnv == "false" {
		client.ValidateLicenseOrFail()
	}

	// Load config from the file.
	cfgFile := prodCfgFile
	if DevEnv == "true" {
		cfgFile = devCfgFile
	}
	if err := config.FromFile(cfgFile); err != nil {
		log.Fatalf("Loading configuration from %s: %s\n", cfgFile, err)
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

// createHandlers registers the API handlers in the DefaultServeMux.
func createHandlers() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the KYC service. Have a nice day!\n"))
	})
	http.HandleFunc("/Ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Pong!"))
	})
	http.HandleFunc("/CheckCustomer", handlers.CheckCustomer)
	http.HandleFunc("/CheckStatus", handlers.CheckStatus)
	http.HandleFunc("/Provider", handlers.IsProviderImplemented)
}
