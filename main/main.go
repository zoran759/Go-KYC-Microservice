package main

import (
	"flag"
	"log"
	"net/http"

	client "modulus/common/licensing-client"
	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"
)

const (
	prodCfgFile = "kyc.cfg"
	devCfgFile  = "kyc_dev.cfg"
)

// DevEnv is the flag to manage prod/dev builds.
// For a production build, this flag value should be set to "false" upon compilation time using: [-ldflags "-X main.DevEnv=false"]
var DevEnv = "true"

var (
	cfgFile = flag.String("config", "", "Load the service configuration from the file specified")
	port    = flag.String("port", "", "Listen on the port specified")
)

func main() {

	// Validate license in production environment.
	// Do os.Exit if failed.
	if DevEnv == "false" {
		client.ValidateLicenseOrFail()
	}

	flag.Parse()

	// Load config from the file.
	// If the command line flag is set its value will be used as the config file name otherwise
	// a predefined file name will be used depending on the value of DevEnv variable.
	if len(*cfgFile) == 0 {
		switch DevEnv {
		case "true":
			*cfgFile = devCfgFile
		default:
			*cfgFile = prodCfgFile
		}
	}
	if err := config.FromFile(*cfgFile); err != nil {
		log.Fatalf("Loading configuration from %s: %s\n", *cfgFile, err)
	}

	createHandlers()

	// Set the listening port for the service.
	// If the command line flag is set its value will be used for the listening port
	// otherwise the option from the service config will be used.
	if len(*port) == 0 {
		*port = config.Cfg.ServicePort()
	}

	log.Printf("Listen on :%v", *port)

	if err := http.ListenAndServe(":"+*port, nil); err != nil {
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
