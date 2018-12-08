package verification

import (
	"crypto/sha256"
	"fmt"
)

// Config represents service config.
type Config struct {
	Host         string
	ClientID     string
	ClientSecret string
	fingerprint  string
}

// Verification describes the verification interface.
type Verification interface {
	CreateUser(User) (*Response, *string, error)
	AddPhysicalDocs(string, string, string, []SubDocument) (*string, error)
	GetUser(string) (*Response, *string, error)
}

func (c Config) calcFingerprint() string {
	if len(c.fingerprint) > 0 {
		return c.fingerprint
	}

	return fmt.Sprintf("%x", sha256.Sum256([]byte(c.ClientID+c.ClientSecret)))
}
