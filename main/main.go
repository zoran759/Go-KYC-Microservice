package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"

	"modulus/kyc/main/config"
	"modulus/kyc/main/handlers"
	"modulus/kyc/main/license"
)

// devbuild holds the hardcoded value of the DevEnv flag for dev builds.
const devbuild = "D"

// DevEnv is the flag to manage prod/dev builds.
// For a production build, set its value to something other than its default upon compilation time using: -ldflags "-X main.DevEnv="
var DevEnv = devbuild

// Command line flags supported by the service.
var (
	cfgFile = flag.String("config", "", "Use the specified file for the service configuration")
	port    = flag.String("port", "", "Listen on the port specified")
)

func main() {
	file, err := os.OpenFile("logs.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Printf("error opening log file %s\n", err)
	}
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	// Turn the license control off in development environment.
	if DevEnv == devbuild {
		license.SetDevMode()
		log.Println("WARNING! Development mode is ON")
	}

	flag.Parse()

	// Load config from the file.
	// If the command line flag is set its value will be used as the config file name otherwise
	// a predefined file name will be used depending on the value of DevEnv variable.
	if len(*cfgFile) == 0 {
		switch DevEnv {
		case devbuild:
			*cfgFile = config.DefaultDevFile
		default:
			*cfgFile = config.DefaultFile
		}
	}
	config.Load(*cfgFile)

	// watch config changes.
	go watchConfigs()

	createHandlers()

	// Set the listening port for the service.
	// If the command line flag is set its value will be used for the listening port
	// otherwise the option from the service config will be used.
	if len(*port) == 0 {
		*port = config.ServicePort()
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
	http.HandleFunc("/Config", handlers.UpdateConfig)
	http.HandleFunc("/cipherTrace", handlers.CipherTraceCheck)
}

// Watch config file and update configs when events happen.
func watchConfigs() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("file watcher can not be created")
	}
	err = watcher.Add(*cfgFile)
	if err != nil {
		log.Fatalf("Watching configuration from %s: %s\n", *cfgFile, err)
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if ok && event.Op.String() == "WRITE" {
				config.Load(*cfgFile)
				log.Printf("Reloading configuration from %s\n", *cfgFile)
			}
		case err, _ := <-watcher.Errors:
			log.Println("error watching config file:", err)
		}
	}
}
