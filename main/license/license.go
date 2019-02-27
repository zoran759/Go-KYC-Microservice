package license

import (
	"errors"
	"log"
	"sync"

	client "modulus/common/licensing-client"
)

// devmode turns off license check when its value is true.
var devmode bool

// lic holds the current license.
var lic license

type license struct {
	sync.RWMutex
	data  string
	valid bool
}

// ErrNoValidLicense represents an error of invalid or missing license.
var ErrNoValidLicense = errors.New("missing or invalid license for the KYC service")

// Update updates the license.
func Update(newlic string) error {
	if devmode {
		return nil
	}

	err := client.ValidateClientLicense(newlic)
	if err != nil {
		log.Println("The license is invalid:", err)
	}

	lic.Lock()
	lic.data = newlic
	lic.valid = err == nil
	lic.Unlock()

	return err
}

// Valid returns current state of the license.
func Valid() bool {
	return devmode || lic.valid
}

// SetDevMode turns off license check.
func SetDevMode() {
	devmode = true
}
