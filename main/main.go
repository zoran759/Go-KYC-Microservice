package main

import (
	"modulus-project-troon/common/licensing-client"
	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"

	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Indicates if it is Development Environment.
// The flag value is set to "false" upon compilation time.
var DevEnv = "true"

var configFile = flag.String("config", "kyc.cfg", "Configuration file for KYC providers")

func main() {

	// Validate license in production environment
	if DevEnv == "false" {
		validateLicense()
	}

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
}

// Validates a license against modulus servers
func validateLicense() {

	// Read license key
	licenseData, err := ioutil.ReadFile("LicenseKey.txt")
	if err != nil {
		log.Fatalln("Error upon reading of license file:", err)
	}
	licenseKey := strings.TrimSpace(string(licenseData))

	// Validate the license
	log.Println("Validating the license: ", licenseKey)
	res, err := client.ValidateLicense(licenseKey)
	if err != nil {
		log.Fatalln("Error upon license validation:", err)
	}
	if res == nil || !res.Pass {
		log.Fatalln("License validation has failed.")
	} else {
		log.Println("The license is validated.")
	}
}
