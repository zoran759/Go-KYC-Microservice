package license

import (
	"log"
	"sync"

	client "modulus/common/licensing-client"
)

// lic holds the current license.
var lic license

type license struct {
	sync.RWMutex
	data  string
	valid bool
}

// UpdateLicense updates the license.
func UpdateLicense(newlic string) {
	err := client.ValidateClientLicense(newlic)
	if err != nil {
		log.Println("The license is invalid:", err)
	}

	lic.Lock()
	lic.data = newlic
	lic.valid = err == nil
	lic.Unlock()
}

// Valid returns current state of the license.
func Valid() bool {
	lic.Lock()
	valid := lic.valid
	lic.Unlock()
	return valid
}
